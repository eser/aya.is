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
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// RegisterHTTPRoutesForBulletin registers routes for bulletin subscription management.
func RegisterHTTPRoutesForBulletin(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	bulletinService *bulletinbiz.Service,
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

			sub, err := bulletinService.Subscribe(
				ctx.Request.Context(),
				*user.IndividualProfileID,
				bulletinbiz.ChannelKind(input.Channel),
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
				PreferredTime int `json:"preferred_time"`
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

			err = bulletinService.UpdatePreferences(
				ctx.Request.Context(),
				idParam,
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
