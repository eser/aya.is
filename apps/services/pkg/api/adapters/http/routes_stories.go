package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

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
					httpfx.WithPlainText(err.Error()),
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

			record, err := storyService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			// if record == nil {
			// 	return ctx.Results.NotFound(httpfx.WithPlainText("story not found"))
			// }

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get story by slug").
		HasDescription("Get story by slug.").
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
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			localeParam := ctx.Request.PathValue("locale")
			profileSlugParam := ctx.Request.PathValue("slug")

			var requestBody struct {
				Slug            string  `json:"slug"`
				Kind            string  `json:"kind"`
				Title           string  `json:"title"`
				Summary         string  `json:"summary"`
				Content         string  `json:"content"`
				StoryPictureURI *string `json:"story_picture_uri"`
				Status          string  `json:"status"`
				IsFeatured      bool    `json:"is_featured"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithPlainText("Invalid request body"))
			}

			if requestBody.Slug == "" || requestBody.Kind == "" ||
				requestBody.Title == "" || requestBody.Status == "" {
				return ctx.Results.BadRequest(
					httpfx.WithPlainText("Slug, kind, title, and status are required"),
				)
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			story, err := storyService.Create(
				ctx.Request.Context(),
				*session.LoggedInUserID,
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
			)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithPlainText(
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
					httpfx.WithPlainText("Failed to create story"),
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
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			localeParam := ctx.Request.PathValue("locale")
			storyIDParam := ctx.Request.PathValue("storyId")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
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
					httpfx.WithPlainText("Failed to check permissions"),
				)
			}

			if !canEdit {
				return ctx.Results.Error(
					http.StatusForbidden,
					httpfx.WithPlainText("You do not have permission to edit this story"),
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
					httpfx.WithPlainText("Failed to get story"),
				)
			}

			if story == nil {
				return ctx.Results.NotFound(httpfx.WithPlainText("Story not found"))
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
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			storyIDParam := ctx.Request.PathValue("storyId")

			var requestBody struct {
				Slug            string  `json:"slug"`
				Status          string  `json:"status"`
				IsFeatured      bool    `json:"is_featured"`
				StoryPictureURI *string `json:"story_picture_uri"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithPlainText("Invalid request body"))
			}

			if requestBody.Slug == "" || requestBody.Status == "" {
				return ctx.Results.BadRequest(
					httpfx.WithPlainText("Slug and status are required"),
				)
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			story, err := storyService.Update(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				storyIDParam,
				requestBody.Slug,
				requestBody.Status,
				requestBody.IsFeatured,
				requestBody.StoryPictureURI,
			)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithPlainText("You do not have permission to edit this story"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Story update failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("story_id", storyIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to update story"),
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
					httpfx.WithPlainText("Session ID not found in context"),
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
				return ctx.Results.BadRequest(httpfx.WithPlainText("Invalid request body"))
			}

			if requestBody.Title == "" {
				return ctx.Results.BadRequest(httpfx.WithPlainText("Title is required"))
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
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
						httpfx.WithPlainText("You do not have permission to edit this story"),
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
					httpfx.WithPlainText("Failed to update story translation"),
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
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			storyIDParam := ctx.Request.PathValue("storyId")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
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
						httpfx.WithPlainText("You do not have permission to delete this story"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Story deletion failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("story_id", storyIDParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to delete story"),
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
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			storyIDParam := ctx.Request.PathValue("storyId")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
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
					httpfx.WithPlainText("Failed to check permissions"),
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
