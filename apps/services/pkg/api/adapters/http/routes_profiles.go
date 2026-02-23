package http

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/aifx"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

// Slug validation regex for pages.
var pageSlugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

func RegisterHTTPRoutesForProfiles( //nolint:funlen,cyclop,maintidx
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	storyService *stories.Service,
	profilePointsService *profile_points.Service,
	aiModels *aifx.Registry,
) {
	routes.
		Route("GET /{locale}/profiles", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			filterKind, filterKindOk := cursor.Filters["kind"]
			if !filterKindOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("filter_kind is required"))
			}

			kinds := strings.SplitSeq(filterKind, ",")
			for kind := range kinds {
				if kind != "individual" && kind != "organization" && kind != "product" {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("filter_kind is invalid"))
				}
			}

			records, err := profileService.List(ctx.Request.Context(), localeParam, cursor)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  records,
				"error": nil,
			})
		}).
		HasSummary("List profiles").
		HasDescription("List profiles.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			if slugParam == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("slug parameter is required"))
			}

			viewerUserID := GetViewerUserID(ctx.Request, authService, userService)

			record, err := profileService.GetBySlugExWithViewerUser(
				ctx.Request.Context(),
				localeParam,
				slugParam,
				viewerUserID,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get profile by slug").
		HasDescription("Get profile by slug.").
		HasResponse(http.StatusOK)

	routes.Route(
		"GET /{locale}/profiles/{slug}/_check",
		func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			slugParam := ctx.Request.PathValue("slug")
			includeDeletedParam := ctx.Request.URL.Query().Get("include_deleted")

			if slugParam == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("slug parameter is required"))
			}

			includeDeleted := includeDeletedParam == "true"

			availability, err := profileService.CheckSlugAvailability(
				ctx.Request.Context(),
				slugParam,
				includeDeleted,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			result := map[string]any{
				"available": availability.Available,
				"message":   availability.Message,
				"severity":  availability.Severity,
			}

			wrappedResponse := cursors.WrapResponseWithCursor(result, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Check profile slug availability").
		HasDescription("Check if a profile slug is available (not taken or reserved).").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/pages", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			viewerUserID := GetViewerUserID(ctx.Request, authService, userService)

			records, err := profileService.ListPagesBySlugForViewer(
				ctx.Request.Context(),
				localeParam,
				slugParam,
				viewerUserID,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(records, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("List profile pages by profile slug").
		HasDescription("List profile pages by profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/pages/{pageSlug}",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}
				slugParam := ctx.Request.PathValue("slug")
				pageSlugParam := ctx.Request.PathValue("pageSlug")

				viewerUserID := GetViewerUserID(ctx.Request, authService, userService)

				records, err := profileService.GetPageBySlugForViewer(
					ctx.Request.Context(),
					localeParam,
					slugParam,
					pageSlugParam,
					viewerUserID,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				wrappedResponse := cursors.WrapResponseWithCursor(records, nil)

				return ctx.Results.JSON(wrappedResponse)
			},
		).
		HasSummary("List profile page by profile slug and page slug").
		HasDescription("List profile page by profile slug and page slug.").
		HasResponse(http.StatusOK)

	// Check page slug availability
	routes.
		Route(
			"GET /{locale}/profiles/{slug}/pages/{pageSlug}/_check",
			func(ctx *httpfx.Context) httpfx.Result {
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}
				profileSlugParam := ctx.Request.PathValue("slug")
				pageSlugParam := ctx.Request.PathValue("pageSlug")
				excludeIDParam := ctx.Request.URL.Query().Get("exclude_id")
				includeDeletedParam := ctx.Request.URL.Query().Get("include_deleted")

				var excludeID *string
				if excludeIDParam != "" {
					excludeID = &excludeIDParam
				}

				includeDeleted := includeDeletedParam == "true"

				availability, err := profileService.CheckPageSlugAvailability(
					ctx.Request.Context(),
					localeParam,
					profileSlugParam,
					pageSlugParam,
					excludeID,
					includeDeleted,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				result := map[string]any{
					"available": availability.Available,
					"message":   availability.Message,
					"severity":  availability.Severity,
				}

				wrappedResponse := cursors.WrapResponseWithCursor(result, nil)

				return ctx.Results.JSON(wrappedResponse)
			},
		).
		HasSummary("Check page slug availability").
		HasDescription("Check if a page slug is available within a profile (not taken).").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/links", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			// Return all links (including non-featured) for the dedicated links page.
			// The empty viewerProfileID means only public links are returned.
			// NOTE: The service has visibility filtering logic (CanViewLink, FilterVisibleLinks)
			// that can show different links based on viewer's membership level (follower,
			// sponsor, etc.). To enable this, we would need to optionally detect the session
			// and pass the viewer's profile ID here. Currently, all public endpoints only
			// show visibility=public links.
			records, err := profileService.ListAllLinksBySlug(
				ctx.Request.Context(),
				localeParam,
				slugParam,
				"", // Empty = anonymous viewer, only public links visible
			)
			if err != nil {
				if errors.Is(err, profiles.ErrLinksNotEnabled) {
					return ctx.Results.Error(
						http.StatusNotFound,
						httpfx.WithErrorMessage("links feature is not enabled for this profile"),
					)
				}

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(records, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("List all profile links by profile slug").
		HasDescription("List all profile links by profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/stories", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			viewerUserID := GetViewerUserID(ctx.Request, authService, userService)

			records, err := storyService.ListByPublicationProfileSlugForViewer(
				ctx.Request.Context(),
				localeParam,
				slugParam,
				cursor,
				viewerUserID,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  records,
				"error": nil,
			})
		}).
		HasSummary("List stories published to profile slug").
		HasDescription("List stories published to profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/stories-authored",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}
				slugParam := ctx.Request.PathValue("slug")
				cursor := cursors.NewCursorFromRequest(ctx.Request)

				viewerUserID := GetViewerUserID(ctx.Request, authService, userService)

				records, err := storyService.ListByAuthorProfileSlugForViewer(
					ctx.Request.Context(),
					localeParam,
					slugParam,
					cursor,
					viewerUserID,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(map[string]any{
					"data":  records,
					"error": nil,
				})
			},
		).
		HasSummary("List stories authored by profile slug").
		HasDescription("List stories authored by profile slug, including unpublished.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/stories/{storySlug}",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}
				// slugParam := ctx.Request.PathValue("slug")
				storySlugParam := ctx.Request.PathValue("storySlug")

				// Try to get viewer's user ID from session (optional - not required)
				var viewerUserID *string

				sessionID, err := GetSessionIDFromCookie(ctx.Request, authService.Config)
				if err == nil && sessionID != "" {
					session, sessionErr := userService.GetSessionByID(
						ctx.Request.Context(),
						sessionID,
					)
					if sessionErr == nil && session != nil && session.LoggedInUserID != nil {
						viewerUserID = session.LoggedInUserID
					}
				}

				// TODO(@eser) pass profile slug too for getting story by profile slug and story slug
				record, err := storyService.GetBySlugForViewer(
					ctx.Request.Context(),
					localeParam,
					storySlugParam,
					viewerUserID,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				if record == nil {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("story not found"))
				}

				wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

				return ctx.Results.JSON(wrappedResponse)
			},
		).
		HasSummary("Get story by profile slug and story slug").
		HasDescription("Get story by profile slug and story slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/contributions",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}
				slugParam := ctx.Request.PathValue("slug")
				cursor := cursors.NewCursorFromRequest(ctx.Request)

				records, err := profileService.ListProfileContributionsBySlug(
					ctx.Request.Context(),
					localeParam,
					slugParam,
					cursor,
				)
				if err != nil {
					if errors.Is(err, profiles.ErrRelationsNotEnabled) {
						return ctx.Results.Error(
							http.StatusNotFound,
							httpfx.WithErrorMessage(
								"relations feature is not enabled for this profile",
							),
						)
					}

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(map[string]any{
					"data":  records,
					"error": nil,
				})
			},
		).
		HasSummary("List profile contributions by profile slug").
		HasDescription("List profile contributions by profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/members",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}
				slugParam := ctx.Request.PathValue("slug")
				cursor := cursors.NewCursorFromRequest(ctx.Request)

				records, err := profileService.ListProfileMembersBySlug(
					ctx.Request.Context(),
					localeParam,
					slugParam,
					cursor,
				)
				if err != nil {
					if errors.Is(err, profiles.ErrRelationsNotEnabled) {
						return ctx.Results.Error(
							http.StatusNotFound,
							httpfx.WithErrorMessage(
								"relations feature is not enabled for this profile",
							),
						)
					}

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(map[string]any{
					"data":  records,
					"error": nil,
				})
			},
		).
		HasSummary("List profile members by profile slug").
		HasDescription("List profile members by profile slug.").
		HasResponse(http.StatusOK)

	// Register profile creation route (protected, requires authentication)
	routes.Route(
		"POST /{locale}/profiles/_create",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get locale from path
			locale := ctx.Request.PathValue("locale")

			// Parse request body
			var requestBody struct {
				Kind        string `json:"kind"`
				Slug        string `json:"slug"`
				Title       string `json:"title"`
				Description string `json:"description"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Sanitize inputs - trim whitespace
			requestBody.Kind = strings.TrimSpace(requestBody.Kind)
			requestBody.Slug = strings.TrimSpace(strings.ToLower(requestBody.Slug))
			requestBody.Title = strings.TrimSpace(requestBody.Title)
			requestBody.Description = strings.TrimSpace(requestBody.Description)

			// Validate kind
			if requestBody.Kind != "individual" && requestBody.Kind != "organization" &&
				requestBody.Kind != "product" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid profile kind"))
			}

			// Validate slug
			if len(requestBody.Slug) < 2 {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Slug must be at least 2 characters"),
				)
			}
			if len(requestBody.Slug) > 50 {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Slug must be at most 50 characters"),
				)
			}
			slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
			if !slugRegex.MatchString(requestBody.Slug) {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage(
						"Slug can only contain lowercase letters, numbers, and hyphens",
					),
				)
			}

			// Validate title
			if len(requestBody.Title) == 0 {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Title is required"))
			}
			if len(requestBody.Title) > 100 {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Title is too long"))
			}

			// Validate description (optional but has max length)
			if len(requestBody.Description) > 500 {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Description is too long"))
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get current user to check individual profile restriction
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Check individual profile restriction
			if requestBody.Kind == "individual" && user.IndividualProfileID != nil {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("User already has an individual profile"),
				)
			}

			// Check if slug already exists
			exists, err := profileService.CheckSlugExists(ctx.Request.Context(), requestBody.Slug)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to check slug availability"),
				)
			}

			if exists {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Slug is already taken"))
			}

			// Create the profile
			profile, err := profileService.Create(
				ctx.Request.Context(),
				user.ID,
				locale,
				requestBody.Slug,
				requestBody.Kind,
				requestBody.Title,
				requestBody.Description,
				nil, // profilePictureURI
				nil, // pronouns
				nil, // properties
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Profile creation failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", user.ID),
					slog.String("slug", requestBody.Slug))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to create profile"),
				)
			}

			// If creating an individual profile, set it as the user's individual profile
			if requestBody.Kind == "individual" {
				err = userService.SetIndividualProfileID(ctx.Request.Context(), user.ID, profile.ID)
				if err != nil {
					logger.ErrorContext(
						ctx.Request.Context(),
						"Failed to set user individual profile ID",
						slog.String("error", err.Error()),
						slog.String("session_id", sessionID),
						slog.String("user_id", user.ID),
						slog.String("profile_id", profile.ID),
					)

					// Note: We don't return an error here because the profile was created successfully
					// The user can still access their profile, they just need to manually link it
					logger.WarnContext(
						ctx.Request.Context(),
						"Profile created but failed to link to user",
						slog.String("session_id", sessionID),
						slog.String("user_id", user.ID),
						slog.String("profile_id", profile.ID),
					)
				}

				// Create self-membership for individual profile (profile is its own owner)
				membershipErr := profileService.CreateProfileMembership(
					ctx.Request.Context(),
					user.ID,
					profile.ID,
					&profile.ID, // Self-reference: member_profile_id = profile_id
					"owner",
				)
				if membershipErr != nil {
					logger.ErrorContext(
						ctx.Request.Context(),
						"Failed to create profile membership for individual profile",
						slog.String("error", membershipErr.Error()),
						slog.String("session_id", sessionID),
						slog.String("user_id", user.ID),
						slog.String("profile_id", profile.ID),
					)
				}

				// Auto-create managed GitHub link from session's stored OAuth tokens
				if session.OAuthProvider != nil &&
					*session.OAuthProvider == "github" &&
					session.OAuthAccessToken != nil &&
					user.GithubRemoteID != nil {
					githubHandle := ""
					if user.GithubHandle != nil {
						githubHandle = *user.GithubHandle
					}

					tokenScope := ""
					if session.OAuthTokenScope != nil {
						tokenScope = *session.OAuthTokenScope
					}

					upsertManagedGitHubLink(
						ctx.Request.Context(),
						logger,
						profileService,
						profile.ID,
						*user.GithubRemoteID,
						githubHandle,
						*session.OAuthAccessToken,
						tokenScope,
					)
				}
			} else {
				// For org/product profiles, create membership with creator's individual profile as owner
				if user.IndividualProfileID != nil {
					membershipErr := profileService.CreateProfileMembership(
						ctx.Request.Context(),
						user.ID,
						profile.ID,
						user.IndividualProfileID,
						"owner",
					)
					if membershipErr != nil {
						logger.ErrorContext(
							ctx.Request.Context(),
							"Failed to create profile membership for org/product profile",
							slog.String("error", membershipErr.Error()),
							slog.String("session_id", sessionID),
							slog.String("user_id", user.ID),
							slog.String("profile_id", profile.ID),
						)
					}
				} else {
					logger.WarnContext(
						ctx.Request.Context(),
						"User has no individual profile, cannot create membership for org/product profile",
						slog.String("session_id", sessionID),
						slog.String("user_id", user.ID),
						slog.String("profile_id", profile.ID),
					)
				}
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  profile,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Create Profile").
		HasDescription("Create a new profile for the authenticated user.").
		HasResponse(http.StatusOK)

	// Profile editing routes (protected, requires authentication)
	routes.Route(
		"PATCH /{locale}/profiles/{slug}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			// Parse request body
			var requestBody struct {
				ProfilePictureURI               *string        `json:"profile_picture_uri"`
				Pronouns                        *string        `json:"pronouns"`
				Properties                      map[string]any `json:"properties"`
				FeatureRelations                *string        `json:"feature_relations"`
				FeatureLinks                    *string        `json:"feature_links"`
				FeatureQA                       *string        `json:"feature_qa"`
				FeatureDiscussions              *string        `json:"feature_discussions"`
				OptionStoryDiscussionsByDefault *bool          `json:"option_story_discussions_by_default"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get user to determine kind for URI prefix validation
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Update the profile
			profile, err := profileService.Update(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				requestBody.ProfilePictureURI,
				requestBody.Pronouns,
				requestBody.Properties,
				requestBody.FeatureRelations,
				requestBody.FeatureLinks,
				requestBody.FeatureQA,
				requestBody.FeatureDiscussions,
				requestBody.OptionStoryDiscussionsByDefault,
			)
			if err != nil {
				if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Profile update failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to update profile"),
				)
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  profile,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Update Profile").
		HasDescription("Update profile main fields (profile picture, pronouns, properties).").
		HasResponse(http.StatusOK)

	routes.Route(
		"PATCH /{locale}/profiles/{slug}/translations/{translationLocale}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			slugParam := ctx.Request.PathValue("slug")
			translationLocaleParam := ctx.Request.PathValue("translationLocale")

			// Parse request body
			var requestBody struct {
				Title       string         `json:"title"`
				Description string         `json:"description"`
				Properties  map[string]any `json:"properties"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Sanitize inputs - trim whitespace
			requestBody.Title = strings.TrimSpace(requestBody.Title)
			requestBody.Description = strings.TrimSpace(requestBody.Description)

			// Validate title (required)
			if len(requestBody.Title) == 0 {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Title is required"))
			}
			if len(requestBody.Title) > 100 {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Title is too long"))
			}

			// Validate description (optional but has max length)
			if len(requestBody.Description) > 500 {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Description is too long"))
			}

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get user to determine kind
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Update the translation
			err := profileService.UpdateTranslation(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				translationLocaleParam,
				requestBody.Title,
				requestBody.Description,
				requestBody.Properties,
			)
			if err != nil {
				if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Profile translation update failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam),
					slog.String("locale", translationLocaleParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to update profile translation"),
				)
			}

			// Return success response
			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
					"message": "Translation updated successfully",
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Update Profile Translation").
		HasDescription("Update profile translation for a specific locale (title, description).").
		HasResponse(http.StatusOK)

	routes.Route(
		"GET /{locale}/profiles/{slug}/_permissions",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			slugParam := ctx.Request.PathValue("slug")

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			canEdit, viewerMembershipKind, permErr := profileService.GetProfilePermissions(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
			)
			if permErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(permErr),
				)
			}

			// Return permissions
			result := map[string]any{
				"can_edit":               canEdit,
				"viewer_membership_kind": viewerMembershipKind,
			}

			wrappedResponse := map[string]any{
				"data":  result,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Check Profile Edit Permissions").
		HasDescription("Check if the authenticated user can edit the specified profile.").
		HasResponse(http.StatusOK)

	routes.Route(
		"GET /{locale}/profiles/{slug}/_tx",
		func(ctx *httpfx.Context) httpfx.Result {
			// Get variables from path
			slugParam := ctx.Request.PathValue("slug")

			// Get all profile translations
			translations, err := profileService.GetProfileTranslations(
				ctx.Request.Context(),
				slugParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Profile translations retrieval failed",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to retrieve profile translations"),
				)
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  translations,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get Profile Translations").
		HasDescription("Get all profile translations for editing by authorized users.").
		HasResponse(http.StatusOK)

	// Profile Links management routes
	routes.Route(
		"GET /{locale}/profiles/{slug}/_links",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get user to determine kind
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Get all profile links for editing
			links, err := profileService.ListProfileLinksBySlug(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
			)
			if err != nil {
				if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Profile links retrieval failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to retrieve profile links"),
				)
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  links,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("List Profile Links for Editing").
		HasDescription("List all profile links (including hidden) for editing by authorized users.").
		HasResponse(http.StatusOK)

	routes.Route(
		"POST /{locale}/profiles/{slug}/_links",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			// Parse request body
			var requestBody struct {
				Kind        string  `json:"kind"`
				URI         *string `json:"uri"`
				Title       string  `json:"title"`
				Icon        *string `json:"icon"`
				Group       *string `json:"group"`
				Description *string `json:"description"`
				IsFeatured  bool    `json:"is_featured"`
				Visibility  string  `json:"visibility"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Validate required fields
			if requestBody.Kind == "" || requestBody.Title == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Kind and title are required"),
				)
			}

			// Validate kind against allowed values
			validKinds := map[string]bool{
				"website":     true,
				"github":      true,
				"x":           true,
				"linkedin":    true,
				"instagram":   true,
				"youtube":     true,
				"speakerdeck": true,
				"bsky":        true,
				"discord":     true,
				"telegram":    true,
			}
			if !validKinds[requestBody.Kind] {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid link kind"))
			}

			// Default visibility to public if not specified
			if requestBody.Visibility == "" {
				requestBody.Visibility = string(profiles.LinkVisibilityPublic)
			}

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get user to determine kind
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Create the profile link
			link, err := profileService.CreateProfileLink(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				requestBody.Kind,
				requestBody.URI,
				requestBody.Title,
				requestBody.Icon,
				requestBody.Group,
				requestBody.Description,
				requestBody.IsFeatured,
				profiles.LinkVisibility(requestBody.Visibility),
			)
			if err != nil {
				if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Profile link creation failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to create profile link"),
				)
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  link,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Create Profile Link").
		HasDescription("Create a new social media link or external link for the profile.").
		HasResponse(http.StatusOK)

	routes.Route(
		"PATCH /{locale}/profiles/{slug}/_links/{linkId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")
			linkIDParam := ctx.Request.PathValue("linkId")

			// Parse request body
			var requestBody struct {
				Kind        string  `json:"kind"`
				Order       int     `json:"order"`
				URI         *string `json:"uri"`
				Title       string  `json:"title"`
				Icon        *string `json:"icon"`
				Group       *string `json:"group"`
				Description *string `json:"description"`
				IsFeatured  bool    `json:"is_featured"`
				Visibility  string  `json:"visibility"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Validate required fields
			if requestBody.Kind == "" || requestBody.Title == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Kind and title are required"),
				)
			}

			// Validate kind against allowed values
			validKinds := map[string]bool{
				"website":     true,
				"github":      true,
				"x":           true,
				"linkedin":    true,
				"instagram":   true,
				"youtube":     true,
				"speakerdeck": true,
				"bsky":        true,
				"discord":     true,
				"telegram":    true,
			}
			if !validKinds[requestBody.Kind] {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid link kind"))
			}

			// Default visibility to public if not specified
			if requestBody.Visibility == "" {
				requestBody.Visibility = string(profiles.LinkVisibilityPublic)
			}

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get user to determine kind
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Update the profile link
			link, err := profileService.UpdateProfileLink(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				linkIDParam,
				requestBody.Kind,
				requestBody.Order,
				requestBody.URI,
				requestBody.Title,
				requestBody.Icon,
				requestBody.Group,
				requestBody.Description,
				requestBody.IsFeatured,
				profiles.LinkVisibility(requestBody.Visibility),
			)
			if err != nil {
				if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Profile link update failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam),
					slog.String("link_id", linkIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to update profile link"),
				)
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  link,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Update Profile Link").
		HasDescription("Update an existing profile link.").
		HasResponse(http.StatusOK)

	routes.Route(
		"DELETE /{locale}/profiles/{slug}/_links/{linkId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			slugParam := ctx.Request.PathValue("slug")
			linkIDParam := ctx.Request.PathValue("linkId")

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get user to determine kind
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Delete the profile link
			err := profileService.DeleteProfileLink(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				linkIDParam,
			)
			if err != nil {
				if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Profile link deletion failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam),
					slog.String("link_id", linkIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to delete profile link"),
				)
			}

			// Return success response
			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
					"message": "Profile link deleted successfully",
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Delete Profile Link").
		HasDescription("Delete a profile link.").
		HasResponse(http.StatusOK)

	// Profile Pages management routes
	routes.Route(
		"GET /{locale}/profiles/{slug}/_pages",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			canEdit, permErr := profileService.HasUserAccessToProfile(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				profiles.MembershipKindMaintainer,
			)
			if permErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(permErr),
				)
			}

			if !canEdit {
				return ctx.Results.Error(
					http.StatusForbidden,
					httpfx.WithErrorMessage("You do not have permission to edit this profile"),
				)
			}

			// Get all profile pages for editing
			pages, err := profileService.ListPagesBySlug(
				ctx.Request.Context(),
				localeParam,
				slugParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Profile pages retrieval failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to retrieve profile pages"),
				)
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  pages,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("List Profile Pages for Editing").
		HasDescription("List all profile pages for editing by authorized users.").
		HasResponse(http.StatusOK)

	routes.Route(
		"POST /{locale}/profiles/{slug}/_pages",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			// Parse request body
			var requestBody struct {
				Slug            string  `json:"slug"`
				Title           string  `json:"title"`
				Summary         string  `json:"summary"`
				Content         string  `json:"content"`
				CoverPictureURI *string `json:"cover_picture_uri"`
				PublishedAt     *string `json:"published_at"`
				Visibility      string  `json:"visibility"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if requestBody.Visibility == "" {
				requestBody.Visibility = string(profiles.PageVisibilityPublic)
			}

			// Sanitize inputs
			requestBody.Slug = lib.SanitizeSlug(strings.TrimSpace(requestBody.Slug))
			requestBody.Title = strings.TrimSpace(requestBody.Title)
			requestBody.Summary = strings.TrimSpace(requestBody.Summary)

			// Validate slug
			if len(requestBody.Slug) < 2 {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Slug must be at least 2 characters"),
				)
			}
			if len(requestBody.Slug) > 100 {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Slug must be at most 100 characters"),
				)
			}
			if !pageSlugRegex.MatchString(requestBody.Slug) {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage(
						"Slug can only contain lowercase letters, numbers, and hyphens",
					),
				)
			}

			// Validate title
			if len(requestBody.Title) == 0 {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Title is required"))
			}
			if len(requestBody.Title) > 200 {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Title is too long"))
			}

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get user to determine kind for URI prefix validation
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Create the profile page
			page, err := profileService.CreateProfilePage(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				requestBody.Slug,
				localeParam,
				requestBody.Title,
				requestBody.Summary,
				requestBody.Content,
				requestBody.CoverPictureURI,
				requestBody.PublishedAt,
				requestBody.Visibility,
			)
			if err != nil {
				if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Profile page creation failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to create profile page"),
				)
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  page,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Create Profile Page").
		HasDescription("Create a new custom page for the profile.").
		HasResponse(http.StatusOK)

	routes.Route(
		"PATCH /{locale}/profiles/{slug}/_pages/{pageId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			slugParam := ctx.Request.PathValue("slug")
			pageIDParam := ctx.Request.PathValue("pageId")

			// Parse request body
			var requestBody struct {
				Slug            string  `json:"slug"`
				Order           int     `json:"order"`
				CoverPictureURI *string `json:"cover_picture_uri"`
				PublishedAt     *string `json:"published_at"`
				Visibility      string  `json:"visibility"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if requestBody.Visibility == "" {
				requestBody.Visibility = string(profiles.PageVisibilityPublic)
			}

			// Sanitize inputs
			requestBody.Slug = lib.SanitizeSlug(strings.TrimSpace(requestBody.Slug))

			// Validate slug
			if len(requestBody.Slug) < 2 {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Slug must be at least 2 characters"),
				)
			}
			if len(requestBody.Slug) > 100 {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Slug must be at most 100 characters"),
				)
			}
			if !pageSlugRegex.MatchString(requestBody.Slug) {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage(
						"Slug can only contain lowercase letters, numbers, and hyphens",
					),
				)
			}

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get user to determine kind for URI prefix validation
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Update the profile page
			page, err := profileService.UpdateProfilePage(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				pageIDParam,
				requestBody.Slug,
				requestBody.Order,
				requestBody.CoverPictureURI,
				requestBody.PublishedAt,
				requestBody.Visibility,
			)
			if err != nil {
				if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Profile page update failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam),
					slog.String("page_id", pageIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to update profile page"),
				)
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  page,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Update Profile Page").
		HasDescription("Update an existing profile page.").
		HasResponse(http.StatusOK)

	routes.Route(
		"PATCH /{locale}/profiles/{slug}/_pages/{pageId}/translations/{translationLocale}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			slugParam := ctx.Request.PathValue("slug")
			pageIDParam := ctx.Request.PathValue("pageId")
			translationLocaleParam := ctx.Request.PathValue("translationLocale")

			// Parse request body
			var requestBody struct {
				Title   string `json:"title"`
				Summary string `json:"summary"`
				Content string `json:"content"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Validate required fields
			if requestBody.Title == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Title is required"),
				)
			}

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get user to determine kind
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Update the translation
			err := profileService.UpdateProfilePageTranslation(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				pageIDParam,
				translationLocaleParam,
				requestBody.Title,
				requestBody.Summary,
				requestBody.Content,
			)
			if err != nil {
				if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Profile page translation update failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam),
					slog.String("page_id", pageIDParam),
					slog.String("locale", translationLocaleParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to update profile page translation"),
				)
			}

			// Return success response
			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
					"message": "Page translation updated successfully",
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Update Profile Page Translation").
		HasDescription("Update profile page translation for a specific locale.").
		HasResponse(http.StatusOK)

	routes.Route(
		"DELETE /{locale}/profiles/{slug}/_pages/{pageId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			slugParam := ctx.Request.PathValue("slug")
			pageIDParam := ctx.Request.PathValue("pageId")

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Get user to determine kind
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Delete the profile page
			err := profileService.DeleteProfilePage(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				pageIDParam,
			)
			if err != nil {
				if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Profile page deletion failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam),
					slog.String("page_id", pageIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to delete profile page"),
				)
			}

			// Return success response
			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
					"message": "Profile page deleted successfully",
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Delete Profile Page").
		HasDescription("Delete a profile page.").
		HasResponse(http.StatusOK)

	// --- Profile Page Translation Management routes ---

	// List profile page translation locales
	routes.Route(
		"GET /{locale}/profiles/{slug}/_pages/{pageId}/_tx",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			slugParam := ctx.Request.PathValue("slug")
			pageIDParam := ctx.Request.PathValue("pageId")

			locales, err := profileService.ListProfilePageTranslationLocales(
				ctx.Request.Context(),
				slugParam,
				pageIDParam,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to list page translation locales"),
				)
			}

			wrappedResponse := map[string]any{
				"data":  locales,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("List Profile Page Translation Locales").
		HasDescription("List all locales that have translations for a profile page.").
		HasResponse(http.StatusOK)

	// Delete profile page translation
	routes.Route(
		"DELETE /{locale}/profiles/{slug}/_pages/{pageId}/translations/{translationLocale}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			pageIDParam := ctx.Request.PathValue("pageId")
			translationLocaleParam := ctx.Request.PathValue("translationLocale")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			err := profileService.DeleteProfilePageTranslation(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				pageIDParam,
				translationLocaleParam,
			)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage(
							"You do not have permission to edit this profile",
						),
					)
				}

				logger.ErrorContext(
					ctx.Request.Context(),
					"Profile page translation deletion failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam),
					slog.String("page_id", pageIDParam),
					slog.String("locale", translationLocaleParam),
				)

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to delete page translation"),
				)
			}

			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Delete Profile Page Translation").
		HasDescription("Delete a profile page translation for a specific locale.").
		HasResponse(http.StatusOK)

	// Generate CV page from profile data using AI
	pageGenerator := NewAIContentGenerator(aiModels)

	routes.Route(
		"POST /{locale}/profiles/{slug}/_pages/generate-cv",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			if user.IndividualProfileID == nil {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("User has no individual profile"),
				)
			}

			// Extend HTTP write deadline for long-running AI call
			rc := http.NewResponseController(ctx.ResponseWriter)
			_ = rc.SetWriteDeadline(time.Now().Add(90 * time.Second))

			page, err := profileService.GenerateCVPage(
				ctx.Request.Context(),
				profiles.GenerateCVPageParams{
					UserID:              *session.LoggedInUserID,
					UserKind:            user.Kind,
					IndividualProfileID: *user.IndividualProfileID,
					ProfileSlug:         slugParam,
					Locale:              localeParam,
				},
				pageGenerator,
				profilePointsService,
			)
			if err != nil {
				if errors.Is(err, profiles.ErrUnauthorized) {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				if errors.Is(err, profile_points.ErrInsufficientPoints) {
					return ctx.Results.Error(
						http.StatusPaymentRequired,
						httpfx.WithErrorMessage(
							"Insufficient points for content generation (requires 5 points)",
						),
					)
				}

				if errors.Is(err, ErrAITranslationNotAvailable) {
					return ctx.Results.Error(
						http.StatusServiceUnavailable,
						httpfx.WithErrorMessage("AI content generation not available"),
					)
				}

				if errors.Is(err, profiles.ErrNoLinkedInLinkFound) {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("No LinkedIn link found on this profile"),
					)
				}

				logger.Error(
					"Failed to generate CV page",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
				)

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			wrappedResponse := map[string]any{
				"data":  page,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Generate CV Page").
		HasDescription("Generate a CV page from profile data using AI.")

	// Auto-translate profile page
	pageTranslator := NewAIContentTranslator(aiModels)

	routes.Route(
		"POST /{locale}/profiles/{slug}/_pages/{pageId}/translations/{targetLocale}/auto-translate",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			pageIDParam := ctx.Request.PathValue("pageId")
			targetLocaleParam := ctx.Request.PathValue("targetLocale")

			var requestBody struct {
				SourceLocale string `json:"source_locale"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if requestBody.SourceLocale == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("source_locale is required"),
				)
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			if user.IndividualProfileID == nil {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("User has no individual profile"),
				)
			}

			// Extend HTTP write deadline for long-running AI call
			rc := http.NewResponseController(ctx.ResponseWriter)
			_ = rc.SetWriteDeadline(time.Now().Add(90 * time.Second))

			err := profileService.AutoTranslateProfilePage(
				ctx.Request.Context(),
				profiles.AutoTranslatePageParams{
					UserID:              *session.LoggedInUserID,
					UserKind:            user.Kind,
					IndividualProfileID: *user.IndividualProfileID,
					ProfileSlug:         slugParam,
					PageID:              pageIDParam,
					SourceLocale:        requestBody.SourceLocale,
					TargetLocale:        targetLocaleParam,
				},
				pageTranslator,
				profilePointsService,
			)
			if err != nil {
				if errors.Is(err, profiles.ErrUnauthorized) {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this profile"),
					)
				}

				if errors.Is(err, profile_points.ErrInsufficientPoints) {
					return ctx.Results.Error(
						http.StatusPaymentRequired,
						httpfx.WithErrorMessage(
							"Insufficient points for auto-translation (requires 5 points)",
						),
					)
				}

				if errors.Is(err, ErrAITranslationNotAvailable) {
					return ctx.Results.Error(
						http.StatusServiceUnavailable,
						httpfx.WithErrorMessage("AI translation not available"),
					)
				}

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Auto-translate Profile Page").
		HasDescription("Auto-translate profile page content from source locale to target locale using AI.").
		HasResponse(http.StatusOK)
}
