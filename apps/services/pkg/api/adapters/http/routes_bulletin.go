package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	bulletinbiz "github.com/eser/aya.is/services/pkg/api/business/bulletin"
	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// RegisterHTTPRoutesForBulletin registers routes for bulletin subscription management.
func RegisterHTTPRoutesForBulletin( //nolint:funlen
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	bulletinService *bulletinbiz.Service,
	telegramService *telegrambiz.Service,
) {
	// List current user's bulletin subscriptions
	routes.Route(
		"GET /{locale}/bulletin/subscriptions",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			if user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("No individual profile found"),
				)
			}

			subs, err := bulletinService.GetSubscriptions(
				ctx.Request.Context(),
				*user.IndividualProfileID,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to list bulletin subscriptions",
					slog.String("error", err.Error()))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			if subs == nil {
				subs = []*bulletinbiz.Subscription{}
			}

			return ctx.Results.JSON(map[string]any{
				"data":  subs,
				"error": nil,
			})
		},
	).HasDescription("List bulletin subscriptions for the current user")

	// Subscribe to bulletin
	routes.Route(
		"POST /{locale}/bulletin/subscribe",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			if user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("No individual profile found"),
				)
			}

			var input struct {
				Channel       string `json:"channel"`
				Frequency     string `json:"frequency"`
				PreferredTime int    `json:"preferred_time"`
			}

			decodeErr := json.NewDecoder(ctx.Request.Body).Decode(&input)
			if decodeErr != nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			if input.Channel == "" {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("channel is required"),
				)
			}

			if input.PreferredTime < 0 || input.PreferredTime > 23 {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("preferred_time must be between 0 and 23 (UTC hour)"),
				)
			}

			frequency := bulletinbiz.DigestFrequency(input.Frequency)
			if frequency == "" {
				frequency = bulletinbiz.FrequencyDaily
			}

			sub, err := bulletinService.Subscribe(
				ctx.Request.Context(),
				*user.IndividualProfileID,
				bulletinbiz.ChannelKind(input.Channel),
				frequency,
				input.PreferredTime,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to create bulletin subscription",
					slog.String("error", err.Error()))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  sub,
				"error": nil,
			})
		},
	).HasDescription("Subscribe to bulletin digest")

	// Update subscription preferences
	routes.Route(
		"PUT /{locale}/bulletin/subscriptions/{id}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			if user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("No individual profile found"),
				)
			}

			idParam := ctx.Request.PathValue("id")

			var input struct {
				Frequency     string `json:"frequency"`
				PreferredTime int    `json:"preferred_time"`
			}

			decodeErr := json.NewDecoder(ctx.Request.Body).Decode(&input)
			if decodeErr != nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			if input.PreferredTime < 0 || input.PreferredTime > 23 {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("preferred_time must be between 0 and 23 (UTC hour)"),
				)
			}

			frequency := bulletinbiz.DigestFrequency(input.Frequency)
			if frequency == "" {
				frequency = bulletinbiz.FrequencyDaily
			}

			err = bulletinService.UpdatePreferences(
				ctx.Request.Context(),
				idParam,
				frequency,
				input.PreferredTime,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to update bulletin subscription",
					slog.String("error", err.Error()),
					slog.String("subscription_id", idParam))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, bulletinbiz.ErrSubscriptionNotFound) {
					statusCode = http.StatusNotFound
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]any{"id": idParam, "updated": true},
				"error": nil,
			})
		},
	).HasDescription("Update bulletin subscription preferences")

	// Unsubscribe
	routes.Route(
		"DELETE /{locale}/bulletin/subscriptions/{id}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			if user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("No individual profile found"),
				)
			}

			idParam := ctx.Request.PathValue("id")

			err = bulletinService.Unsubscribe(ctx.Request.Context(), idParam)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to unsubscribe from bulletin",
					slog.String("error", err.Error()),
					slog.String("subscription_id", idParam))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, bulletinbiz.ErrSubscriptionNotFound) {
					statusCode = http.StatusNotFound
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]any{"id": idParam, "deleted": true},
				"error": nil,
			})
		},
	).HasDescription("Unsubscribe from bulletin digest")

	// Get unified bulletin preferences for the settings UI
	routes.Route(
		"GET /{locale}/bulletin/preferences",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			if user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("No individual profile found"),
				)
			}

			subs, err := bulletinService.GetSubscriptions(
				ctx.Request.Context(),
				*user.IndividualProfileID,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to get bulletin preferences",
					slog.String("error", err.Error()))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			// Derive unified preferences from active subscriptions
			var frequency *string

			var preferredTime *int

			var lastBulletinAt *string

			activeChannels := make([]string, 0, len(subs))

			for _, sub := range subs {
				activeChannels = append(activeChannels, string(sub.Channel))

				// All subscriptions share the same frequency + preferredTime
				if frequency == nil {
					f := string(sub.Frequency)
					frequency = &f

					pt := sub.PreferredTime
					preferredTime = &pt
				}

				if sub.LastBulletinAt != nil {
					formatted := sub.LastBulletinAt.UTC().Format("2006-01-02T15:04:05Z")
					lastBulletinAt = &formatted
				}
			}

			// Resolve email from the authenticated user
			var email *string
			if user.Email != nil && *user.Email != "" {
				email = user.Email
			}

			// Check Telegram connection
			telegramConnected := false
			if telegramService != nil {
				_, linkErr := telegramService.GetProfileTelegramLink(
					ctx.Request.Context(), *user.IndividualProfileID,
				)
				if linkErr == nil {
					telegramConnected = true
				}
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]any{
					"frequency":          frequency,
					"preferred_time":     preferredTime,
					"channels":           activeChannels,
					"email":              email,
					"telegram_connected": telegramConnected,
					"last_bulletin_at":   lastBulletinAt,
				},
				"error": nil,
			})
		},
	).HasDescription("Get unified bulletin preferences for settings UI")

	// Update bulletin preferences atomically
	routes.Route(
		"PUT /{locale}/bulletin/preferences",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			if user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("No individual profile found"),
				)
			}

			var input struct {
				Frequency     string   `json:"frequency"`
				PreferredTime int      `json:"preferred_time"`
				Channels      []string `json:"channels"`
			}

			decodeErr := json.NewDecoder(ctx.Request.Body).Decode(&input)
			if decodeErr != nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			if input.PreferredTime < 0 || input.PreferredTime > 23 {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("preferred_time must be between 0 and 23 (UTC hour)"),
				)
			}

			// Validate frequency
			frequency := bulletinbiz.DigestFrequency(input.Frequency)
			if frequency == "" {
				frequency = bulletinbiz.FrequencyDaily
			}

			validFrequencies := map[bulletinbiz.DigestFrequency]bool{
				bulletinbiz.FrequencyDaily:   true,
				bulletinbiz.FrequencyBiDaily: true,
				bulletinbiz.FrequencyWeekly:  true,
			}
			if !validFrequencies[frequency] {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("frequency must be one of: daily, bidaily, weekly"),
				)
			}

			// Validate channels
			validChannels := map[string]bool{
				string(bulletinbiz.ChannelEmail):    true,
				string(bulletinbiz.ChannelTelegram): true,
			}

			channels := make([]bulletinbiz.ChannelKind, 0, len(input.Channels))
			for _, ch := range input.Channels {
				if !validChannels[ch] {
					return ctx.Results.Error(
						http.StatusBadRequest,
						httpfx.WithErrorMessage("channels must be a subset of: email, telegram"),
					)
				}

				channels = append(channels, bulletinbiz.ChannelKind(ch))
			}

			err = bulletinService.UpdateBulletinPreferences(
				ctx.Request.Context(),
				*user.IndividualProfileID,
				frequency,
				input.PreferredTime,
				channels,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to update bulletin preferences",
					slog.String("error", err.Error()))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]any{"updated": true},
				"error": nil,
			})
		},
	).HasDescription("Update bulletin preferences atomically")

	// Admin: manually trigger bulletin processing
	routes.Route(
		"POST /admin/bulletin/trigger",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			if user.Kind != "admin" {
				return ctx.Results.Error(
					http.StatusForbidden,
					httpfx.WithErrorMessage("Admin access required"),
				)
			}

			processErr := bulletinService.ProcessDigestWindow(ctx.Request.Context())
			if processErr != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to trigger bulletin processing",
					slog.String("error", processErr.Error()))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(processErr),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]any{"triggered": true},
				"error": nil,
			})
		},
	).HasDescription("Manually trigger bulletin digest processing (admin only)")
}
