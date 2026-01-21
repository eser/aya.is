package http

import (
	"log/slog"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/uploads"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForSite(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	uploadService *uploads.Service,
) {
	routes.
		Route(
			"GET /{locale}/site/custom-domains/{domain}",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam := ctx.Request.PathValue("locale")
				domainParam := ctx.Request.PathValue("domain")

				records, err := profileService.GetByCustomDomain(
					ctx.Request.Context(),
					localeParam,
					domainParam,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText(err.Error()),
					)
				}

				wrappedResponse := cursors.WrapResponseWithCursor(records, nil)

				return ctx.Results.JSON(wrappedResponse)
			},
		).
		HasSummary("Get profile by a custom domain").
		HasDescription("Get profile by a custom domain.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/site/spotlight", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			// localeParam := ctx.Request.PathValue("locale")

			// Static spotlight items
			spotlightItems := []profiles.SpotlightItem{
				// {
				// 	Icon:  "Users",
				// 	To:    "/" + localeParam + "/aya",
				// 	Title: "AYA",
				// },
				// {
				// 	Icon:  "User",
				// 	To:    "/" + localeParam + "/eser",
				// 	Title: "Eser Özvataf | SW³",
				// },
			}

			wrappedResponse := cursors.WrapResponseWithCursor(spotlightItems, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Gets spotlight metadata").
		HasDescription("Gets spotlight metadata.").
		HasResponse(http.StatusOK)

	// Upload routes (protected, requires authentication)

	// Generate presigned URL for upload
	routes.Route(
		"POST /{locale}/site/uploads/presign",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			var requestBody struct {
				Filename    string `json:"filename"`
				ContentType string `json:"content_type"`
				Purpose     string `json:"purpose"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithPlainText("Invalid request body"))
			}

			if requestBody.Filename == "" || requestBody.ContentType == "" ||
				requestBody.Purpose == "" {
				return ctx.Results.BadRequest(
					httpfx.WithPlainText("Filename, content_type, and purpose are required"),
				)
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			resp, err := uploadService.GenerateUploadURL(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				uploads.PresignedURLRequest{
					Filename:    requestBody.Filename,
					ContentType: requestBody.ContentType,
					Purpose:     uploads.Purpose(requestBody.Purpose),
				},
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to generate presigned URL",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("filename", requestBody.Filename))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to generate upload URL"),
				)
			}

			wrappedResponse := map[string]any{
				"data":  resp,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Generate Presigned Upload URL").
		HasDescription("Generate a presigned URL for uploading a file to S3-compatible storage.").
		HasResponse(http.StatusOK)

	// Remove uploaded file
	routes.Route(
		"DELETE /{locale}/site/uploads/{key...}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			keyParam := ctx.Request.PathValue("key")
			if keyParam == "" {
				return ctx.Results.BadRequest(httpfx.WithPlainText("Key is required"))
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			err := uploadService.RemoveObject(ctx.Request.Context(), keyParam)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to remove uploaded file",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("key", keyParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to remove file"),
				)
			}

			wrappedResponse := map[string]any{
				"data": map[string]any{
					"success": true,
					"message": "File removed successfully",
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Remove Uploaded File").
		HasDescription("Remove a previously uploaded file from S3-compatible storage.").
		HasResponse(http.StatusOK)
}
