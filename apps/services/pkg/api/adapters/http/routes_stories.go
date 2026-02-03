package http

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

// Slug validation regex for stories.
var storySlugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

func RegisterHTTPRoutesForStories( //nolint:funlen
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	storyService *stories.Service,
) {
	routes.
		Route("GET /{locale}/stories", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			records, err := storyService.List(ctx.Request.Context(), localeParam, cursor)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage(err.Error()),
				)
			}

			return ctx.Results.JSON(records)
		}).
		HasSummary("List stories").
		HasDescription("List stories.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/stories/{slug}", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")

			// Try to get viewer's user ID from session (optional - not required)
			var viewerUserID *string

			sessionID, err := GetSessionIDFromCookie(ctx.Request, authService.Config)
			if err == nil && sessionID != "" {
				session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
				if sessionErr == nil && session != nil && session.LoggedInUserID != nil {
					viewerUserID = session.LoggedInUserID
				}
			}

			record, err := storyService.GetBySlugForViewer(
				ctx.Request.Context(),
				localeParam,
				slugParam,
				viewerUserID,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage(err.Error()),
				)
			}

			if record == nil {
				return ctx.Results.NotFound(httpfx.WithErrorMessage("story not found"))
			}

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get story by slug").
		HasDescription("Get story by slug.").
		HasResponse(http.StatusOK)

	// Check story slug availability
	routes.
		Route("GET /{locale}/stories/{slug}/_check", func(ctx *httpfx.Context) httpfx.Result {
			slugParam := ctx.Request.PathValue("slug")
			excludeIDParam := ctx.Request.URL.Query().Get("exclude_id")
			storyIDParam := ctx.Request.URL.Query().Get("story_id")
			publishedAtParam := ctx.Request.URL.Query().Get("published_at")
			includeDeletedParam := ctx.Request.URL.Query().Get("include_deleted")

			var excludeID *string
			if excludeIDParam != "" {
				excludeID = &excludeIDParam
			}

			var storyID *string
			if storyIDParam != "" {
				storyID = &storyIDParam
			}

			var publishedAt *time.Time
			if publishedAtParam != "" {
				if parsed, parseErr := time.Parse(time.RFC3339, publishedAtParam); parseErr == nil {
					publishedAt = &parsed
				}
			}

			includeDeleted := includeDeletedParam == "true"

			availability, err := storyService.CheckSlugAvailability(
				ctx.Request.Context(),
				slugParam,
				excludeID,
				storyID,
				publishedAt,
				includeDeleted,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage(err.Error()),
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
		HasSummary("Check story slug availability").
		HasDescription("Check if a story slug is available (not taken) and validates date prefix for published stories.").
		HasResponse(http.StatusOK)

	// Story CRUD routes (protected, requires authentication)

	// Create story
	routes.Route(
		"POST /{locale}/profiles/{slug}/_stories",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			localeParam := ctx.Request.PathValue("locale")
			profileSlugParam := ctx.Request.PathValue("slug")

			var requestBody struct {
				Slug              string   `json:"slug"`
				Kind              string   `json:"kind"`
				Title             string   `json:"title"`
				Summary           string   `json:"summary"`
				Content           string   `json:"content"`
				StoryPictureURI   *string  `json:"story_picture_uri"`
				PublishToProfiles []string `json:"publish_to_profiles"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Sanitize inputs
			requestBody.Slug = lib.SanitizeSlug(strings.TrimSpace(requestBody.Slug))
			requestBody.Kind = strings.TrimSpace(requestBody.Kind)
			requestBody.Title = strings.TrimSpace(requestBody.Title)
			requestBody.Summary = strings.TrimSpace(requestBody.Summary)

			// Validate kind
			validKinds := map[string]bool{
				"article": true, "announcement": true, "news": true,
				"status": true, "content": true, "presentation": true,
			}
			if !validKinds[requestBody.Kind] {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid story kind"))
			}

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
			if !storySlugRegex.MatchString(requestBody.Slug) {
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

			story, err := storyService.Create(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				profileSlugParam,
				localeParam,
				requestBody.Slug,
				requestBody.Kind,
				requestBody.Title,
				requestBody.Summary,
				requestBody.Content,
				requestBody.StoryPictureURI,
				requestBody.PublishToProfiles,
			)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage(
							"You do not have permission to create stories for this profile",
						),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Story creation failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("profile_slug", profileSlugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to create story"),
				)
			}

			wrappedResponse := map[string]any{
				"data":  story,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Create Story").
		HasDescription("Create a new story for the profile.").
		HasResponse(http.StatusOK)

	// Get story for editing
	routes.Route(
		"GET /{locale}/profiles/{slug}/_stories/{storyId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			localeParam := ctx.Request.PathValue("locale")
			storyIDParam := ctx.Request.PathValue("storyId")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			canEdit, err := storyService.CanUserEditStory(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				storyIDParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Permission check failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("story_id", storyIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to check permissions"),
				)
			}

			if !canEdit {
				return ctx.Results.Error(
					http.StatusForbidden,
					httpfx.WithErrorMessage("You do not have permission to edit this story"),
				)
			}

			story, err := storyService.GetForEdit(ctx.Request.Context(), localeParam, storyIDParam)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Get story for edit failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("story_id", storyIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get story"),
				)
			}

			if story == nil {
				return ctx.Results.NotFound(httpfx.WithErrorMessage("Story not found"))
			}

			wrappedResponse := map[string]any{
				"data":  story,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get Story for Editing").
		HasDescription("Get story raw content for editing by authorized users.").
		HasResponse(http.StatusOK)

	// Update story
	routes.Route(
		"PATCH /{locale}/profiles/{slug}/_stories/{storyId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			storyIDParam := ctx.Request.PathValue("storyId")

			var requestBody struct {
				Slug            string  `json:"slug"`
				StoryPictureURI *string `json:"story_picture_uri"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
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
			if !storySlugRegex.MatchString(requestBody.Slug) {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage(
						"Slug can only contain lowercase letters, numbers, and hyphens",
					),
				)
			}

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

			story, err := storyService.Update(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				storyIDParam,
				requestBody.Slug,
				requestBody.StoryPictureURI,
			)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this story"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Story update failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("story_id", storyIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to update story"),
				)
			}

			wrappedResponse := map[string]any{
				"data":  story,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Update Story").
		HasDescription("Update story main fields (slug, picture).").
		HasResponse(http.StatusOK)

	// Update story translation
	routes.Route(
		"PATCH /{locale}/profiles/{slug}/_stories/{storyId}/translations/{translationLocale}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			storyIDParam := ctx.Request.PathValue("storyId")
			translationLocaleParam := ctx.Request.PathValue("translationLocale")

			var requestBody struct {
				Title   string `json:"title"`
				Summary string `json:"summary"`
				Content string `json:"content"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if requestBody.Title == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Title is required"))
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			err := storyService.UpdateTranslation(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				storyIDParam,
				translationLocaleParam,
				requestBody.Title,
				requestBody.Summary,
				requestBody.Content,
			)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to edit this story"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Story translation update failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("story_id", storyIDParam),
					slog.String("locale", translationLocaleParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to update story translation"),
				)
			}

			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
					"message": "Translation updated successfully",
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Update Story Translation").
		HasDescription("Update story translation for a specific locale.").
		HasResponse(http.StatusOK)

	// Delete story
	routes.Route(
		"DELETE /{locale}/profiles/{slug}/_stories/{storyId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			storyIDParam := ctx.Request.PathValue("storyId")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			err := storyService.Delete(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				storyIDParam,
			)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to delete this story"),
					)
				}

				if errors.Is(err, stories.ErrHasActivePublications) {
					return ctx.Results.Error(
						http.StatusConflict,
						httpfx.WithErrorMessage(
							"Cannot delete story with active publications. Unpublish from all profiles first.",
						),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Story deletion failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("story_id", storyIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to delete story"),
				)
			}

			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
					"message": "Story deleted successfully",
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Delete Story").
		HasDescription("Delete a story (soft delete).").
		HasResponse(http.StatusOK)

	// Check story permissions
	routes.Route(
		"GET /{locale}/profiles/{slug}/_stories/{storyId}/_permissions",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			storyIDParam := ctx.Request.PathValue("storyId")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			canEdit, err := storyService.CanUserEditStory(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				storyIDParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Permission check failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("story_id", storyIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to check permissions"),
				)
			}

			result := map[string]any{
				"can_edit": canEdit,
			}

			wrappedResponse := map[string]any{
				"data":  result,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Check Story Edit Permissions").
		HasDescription("Check if the authenticated user can edit the specified story.").
		HasResponse(http.StatusOK)

	// --- Story Publication CRUD routes ---

	// List publications for a story
	routes.Route(
		"GET /{locale}/profiles/{slug}/_stories/{storyId}/publications",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam := ctx.Request.PathValue("locale")
			storyIDParam := ctx.Request.PathValue("storyId")

			publications, err := storyService.ListPublications(
				ctx.Request.Context(),
				localeParam,
				storyIDParam,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to list publications"),
				)
			}

			wrappedResponse := map[string]any{
				"data":  publications,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("List Story Publications").
		HasDescription("List all publications for a story with profile info.").
		HasResponse(http.StatusOK)

	// Add publication
	routes.Route(
		"POST /{locale}/profiles/{slug}/_stories/{storyId}/publications",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			localeParam := ctx.Request.PathValue("locale")
			storyIDParam := ctx.Request.PathValue("storyId")

			var requestBody struct {
				ProfileID  string `json:"profile_id"`
				IsFeatured bool   `json:"is_featured"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if requestBody.ProfileID == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("profile_id is required"))
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			publication, err := storyService.AddPublication(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				storyIDParam,
				requestBody.ProfileID,
				localeParam,
				requestBody.IsFeatured,
			)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage(
							"You do not have permission to publish this story",
						),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Add publication failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("story_id", storyIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to add publication"),
				)
			}

			wrappedResponse := map[string]any{
				"data":  publication,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Add Story Publication").
		HasDescription("Publish a story to a profile.").
		HasResponse(http.StatusOK)

	// Update publication (is_featured)
	routes.Route(
		"PATCH /{locale}/profiles/{slug}/_stories/{storyId}/publications/{publicationId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			storyIDParam := ctx.Request.PathValue("storyId")
			publicationIDParam := ctx.Request.PathValue("publicationId")

			var requestBody struct {
				IsFeatured bool `json:"is_featured"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			err := storyService.UpdatePublication(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				storyIDParam,
				publicationIDParam,
				requestBody.IsFeatured,
			)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage(
							"You do not have permission to update this publication",
						),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Update publication failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("story_id", storyIDParam),
					slog.String("publication_id", publicationIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to update publication"),
				)
			}

			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
					"message": "Publication updated successfully",
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Update Story Publication").
		HasDescription("Update a story publication (e.g. is_featured).").
		HasResponse(http.StatusOK)

	// Remove publication
	routes.Route(
		"DELETE /{locale}/profiles/{slug}/_stories/{storyId}/publications/{publicationId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			localeParam := ctx.Request.PathValue("locale")
			storyIDParam := ctx.Request.PathValue("storyId")
			publicationIDParam := ctx.Request.PathValue("publicationId")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			err := storyService.RemovePublication(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				storyIDParam,
				publicationIDParam,
				localeParam,
			)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage(
							"You do not have permission to remove this publication",
						),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Remove publication failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("story_id", storyIDParam),
					slog.String("publication_id", publicationIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to remove publication"),
				)
			}

			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
					"message": "Publication removed successfully",
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Remove Story Publication").
		HasDescription("Remove a story publication (unpublish from a profile).").
		HasResponse(http.StatusOK)
}
