package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	telegramadapter "github.com/eser/aya.is/services/pkg/api/adapters/telegram"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// RegisterHTTPRoutesForTelegram registers the Telegram webhook and token generation endpoints.
func RegisterHTTPRoutesForTelegram( //nolint:cyclop,funlen
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	telegram *TelegramProviders,
) {
	// POST /telegram/webhook — receives updates from Telegram Bot API
	routes.Route(
		"POST /telegram/webhook",
		func(ctx *httpfx.Context) httpfx.Result {
			// Verify the webhook secret header
			secretHeader := ctx.Request.Header.Get("X-Telegram-Bot-Api-Secret-Token")
			if telegram.WebhookSecret != "" && secretHeader != telegram.WebhookSecret {
				logger.WarnContext(ctx.Request.Context(), "Telegram webhook: invalid secret header")

				return ctx.Results.Error(
					http.StatusForbidden,
					httpfx.WithErrorMessage("Invalid webhook secret"),
				)
			}

			// Parse the update
			var update telegramadapter.Update
			err := json.NewDecoder(ctx.Request.Body).Decode(&update)
			if err != nil {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Invalid update payload"),
				)
			}

			// Process asynchronously to not block Telegram
			go telegram.Bot.HandleUpdate(context.Background(), &update)

			return ctx.Results.Ok()
		},
	).
		HasSummary("Telegram Webhook").
		HasDescription("Receives updates from Telegram Bot API.")

	// POST /{locale}/profiles/{slug}/_links/telegram/generate-token
	// Authenticated endpoint: generates a link token and returns the deep link URL
	routes.Route(
		"POST /{locale}/profiles/{slug}/_links/telegram/generate-token",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")

			// Get session
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil || session.LoggedInUserID == nil {
				return ctx.Results.Error(
					http.StatusUnauthorized,
					httpfx.WithErrorMessage("Not authenticated"),
				)
			}

			// Check permissions (maintainer or above)
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

			// Get profile ID
			profile, err := profileService.GetBySlug(ctx.Request.Context(), "en", slugParam)
			if err != nil || profile == nil {
				return ctx.Results.Error(
					http.StatusNotFound,
					httpfx.WithErrorMessage("Profile not found"),
				)
			}

			// Generate link token
			token, tokenErr := telegram.Service.GenerateLinkToken(
				ctx.Request.Context(),
				profile.ID,
				profile.Slug,
				*session.LoggedInUserID,
			)
			if tokenErr != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to generate Telegram link token",
					slog.String("error", tokenErr.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusConflict,
					httpfx.WithErrorMessage(tokenErr.Error()),
				)
			}

			// Build deep link using the bot client
			deepLink := telegram.Bot.Client().DeepLink(token)

			return ctx.Results.JSON(map[string]any{
				"data": map[string]string{
					"token":     token,
					"deep_link": deepLink,
				},
				"error": nil,
			})
		},
	).
		HasSummary("Generate Telegram Link Token").
		HasDescription("Generates a link token for connecting a Telegram account to a profile.")
}
