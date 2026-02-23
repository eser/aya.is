package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	telegramadapter "github.com/eser/aya.is/services/pkg/api/adapters/telegram"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// RegisterHTTPRoutesForTelegram registers the Telegram webhook and code verification endpoints.
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

	// POST /{locale}/profiles/{slug}/_links/telegram/verify-code
	// Authenticated endpoint: verifies a code from the bot and creates the managed link
	routes.Route(
		"POST /{locale}/profiles/{slug}/_links/telegram/verify-code",
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

			// Parse request body
			var body struct {
				Code string `json:"code"`
			}

			err := json.NewDecoder(ctx.Request.Body).Decode(&body)
			if err != nil || strings.TrimSpace(body.Code) == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Missing or invalid verification code"),
				)
			}

			// Normalize code to uppercase
			code := strings.ToUpper(strings.TrimSpace(body.Code))

			// Get profile
			profile, profileErr := profileService.GetBySlug(ctx.Request.Context(), "en", slugParam)
			if profileErr != nil || profile == nil {
				return ctx.Results.Error(
					http.StatusNotFound,
					httpfx.WithErrorMessage("Profile not found"),
				)
			}

			// Verify code and create link
			result, verifyErr := telegram.Service.VerifyCodeAndLink(
				ctx.Request.Context(),
				code,
				profile.ID,
				profile.Slug,
				*session.LoggedInUserID,
			)
			if verifyErr != nil {
				logger.WarnContext(ctx.Request.Context(), "Telegram code verification failed",
					slog.String("error", verifyErr.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusBadRequest

				switch {
				case errors.Is(verifyErr, telegrambiz.ErrCodeNotFound):
					statusCode = http.StatusNotFound
				case errors.Is(verifyErr, telegrambiz.ErrAlreadyLinked),
					errors.Is(verifyErr, telegrambiz.ErrProfileAlreadyHasTelegram):
					statusCode = http.StatusConflict
				}

				return ctx.Results.Error(
					statusCode,
					httpfx.WithErrorMessage(verifyErr.Error()),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]any{
					"profile_id":        result.ProfileID,
					"profile_slug":      result.ProfileSlug,
					"telegram_user_id":  result.TelegramUserID,
					"telegram_username": result.TelegramUsername,
				},
				"error": nil,
			})
		},
	).
		HasSummary("Verify Telegram Code").
		HasDescription("Verifies a code from the Telegram bot and links the account to the profile.")

	// POST /{locale}/profiles/{slug}/_resources/telegram/verify-code
	// Verifies a registration code from the bot and creates a Telegram group resource
	routes.Route(
		"POST /{locale}/profiles/{slug}/_resources/telegram/verify-code",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result { //nolint:cyclop,funlen
			localeParam, localeOk := validateLocale(ctx)
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

			// Parse request body
			var body struct {
				Code string `json:"code"`
			}

			err := json.NewDecoder(ctx.Request.Body).Decode(&body)
			if err != nil || strings.TrimSpace(body.Code) == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Missing or invalid registration code"),
				)
			}

			code := strings.ToUpper(strings.TrimSpace(body.Code))

			// Resolve the registration code
			codeData, resolveErr := telegram.Service.ResolveGroupRegisterCode(
				ctx.Request.Context(), code,
			)
			if resolveErr != nil {
				statusCode := http.StatusNotFound
				if errors.Is(resolveErr, telegrambiz.ErrCodeNotFound) {
					statusCode = http.StatusNotFound
				}

				return ctx.Results.Error(
					statusCode,
					httpfx.WithErrorMessage("Invalid or expired registration code"),
				)
			}

			// Get profile
			profile, profileErr := profileService.GetBySlug(ctx.Request.Context(), "en", slugParam)
			if profileErr != nil || profile == nil {
				return ctx.Results.Error(
					http.StatusNotFound,
					httpfx.WithErrorMessage("Profile not found"),
				)
			}

			// Get user info for resource creation
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			// Create the Telegram group as a profile resource
			chatIDStr := strconv.FormatInt(codeData.ChatID, 10)

			resource, createErr := profileService.CreateProfileResource(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				"telegram_group",
				true,
				&chatIDStr,
				nil,
				nil,
				codeData.ChatTitle,
				nil,
				map[string]any{
					"telegram_chat_id": codeData.ChatID,
				},
			)
			if createErr != nil {
				logger.ErrorContext(
					ctx.Request.Context(),
					"Failed to create Telegram group resource",
					slog.String("error", createErr.Error()),
					slog.String("slug", slugParam),
				)

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(createErr),
				)
			}

			// Consume the code
			consumeErr := telegram.Service.ConsumeGroupRegisterCode(ctx.Request.Context(), code)
			if consumeErr != nil {
				logger.WarnContext(ctx.Request.Context(),
					"Failed to consume register code (resource was already created)",
					slog.String("code", code),
					slog.String("error", consumeErr.Error()))
			}

			// Send confirmation message to the Telegram group (fire-and-forget)
			go func() {
				confirmText := fmt.Sprintf(
					"This group has been registered with <b>%s</b> (https://aya.is/%s).\n"+
						"Eligible members can now receive invites to this group.",
					profile.Title,
					profile.Slug,
				)
				_ = telegram.Bot.Client().SendMessage(
					context.Background(), codeData.ChatID, confirmText,
				)
			}()

			return ctx.Results.JSON(map[string]any{
				"data":  resource,
				"error": nil,
			})
		},
	).
		HasSummary("Verify Telegram Group Registration Code").
		HasDescription("Verifies a registration code and creates a Telegram group resource on the profile.")
}
