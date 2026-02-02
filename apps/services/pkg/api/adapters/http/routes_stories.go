package http

import (
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
			statusParam := ctx.Request.URL.Query().Get("status")
			publishedAtParam := ctx.Request.URL.Query().Get("published_at")
			includeDeletedParam := ctx.Request.URL.Query().Get("include_deleted")

			var excludeID *string
			if excludeIDParam != "" {
				excludeID = &excludeIDParam
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
				statusParam,
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
				Slug            string     `json:"slug"`
				Kind            string     `json:"kind"`
				Title           string     `json:"title"`
				Summary         string     `json:"summary"`
				Content         string     `json:"content"`
				StoryPictureURI *string    `json:"story_picture_uri"`
				Status          string     `json:"status"`
				IsFeatured      bool       `json:"is_featured"`
				PublishedAt     *time.Time `json:"published_at"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Sanitize inputs
			requestBody.Slug = lib.SanitizeSlug(strings.TrimSpace(requestBody.Slug))
			requestBody.Kind = strings.TrimSpace(requestBody.Kind)
			requestBody.Title = strings.TrimSpace(requestBody.Title)
			requestBody.Summary = strings.TrimSpace(requestBody.Summary)
			requestBody.Status = strings.TrimSpace(requestBody.Status)

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

			// Validate status
			validStatuses := map[string]bool{"draft": true, "published": true}
			if !validStatuses[requestBody.Status] {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid status"))
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
				profileSlugParam,
				localeParam,
				requestBody.Slug,
				requestBody.Kind,
				requestBody.Status,
				requestBody.Title,
				requestBody.Summary,
				requestBody.Content,
				requestBody.StoryPictureURI,
				requestBody.IsFeatured,
				requestBody.PublishedAt,
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
				Slug            string     `json:"slug"`
				Status          string     `json:"status"`
				IsFeatured      bool       `json:"is_featured"`
				StoryPictureURI *string    `json:"story_picture_uri"`
				PublishedAt     *time.Time `json:"published_at"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Sanitize inputs
			requestBody.Slug = lib.SanitizeSlug(strings.TrimSpace(requestBody.Slug))
			requestBody.Status = strings.TrimSpace(requestBody.Status)

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

			// Validate status
			validStatuses := map[string]bool{"draft": true, "published": true}
			if !validStatuses[requestBody.Status] {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid status"))
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
				requestBody.Status,
				requestBody.IsFeatured,
				requestBody.StoryPictureURI,
				requestBody.PublishedAt,
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
		HasDescription("Update story main fields (slug, status, featured, picture).").
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
}
