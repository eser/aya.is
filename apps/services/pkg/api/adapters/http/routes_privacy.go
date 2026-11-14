package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/privacy"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

type deleteAccountRequest struct {
	ConfirmSlug   string `json:"confirm_slug"`
	DeleteContent bool   `json:"delete_content"`
}

// privacyUserID returns the id of the signed-in user. Every privacy route acts on the caller's own account,
// never on an id taken from the request, so there is nothing to forge.
func privacyUserID(ctx *httpfx.Context, userService *users.Service) (string, bool) {
	sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
	if !ok {
		return "", false
	}

	session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
	if err != nil || session == nil || session.LoggedInUserID == nil {
		return "", false
	}

	return *session.LoggedInUserID, true
}

// privacyError maps a privacy service error to an HTTP response.
//
//nolint:cyclop // one flat switch that maps each domain error to a status code
func privacyError(
	ctx *httpfx.Context,
	logger *logfx.Logger,
	err error,
) httpfx.Result {
	var cooldown *privacy.CooldownError

	var owned *privacy.OwnedProfilesError

	switch {
	case errors.As(err, &cooldown):
		return ctx.Results.Error(http.StatusTooManyRequests, httpfx.WithJSON(map[string]any{
			"error":           "a new export can not be requested yet",
			"code":            "export_cooldown",
			"next_request_at": cooldown.NextRequestAt,
		}))
	case errors.As(err, &owned):
		return ctx.Results.Error(http.StatusConflict, httpfx.WithJSON(map[string]any{
			"error":    "transfer ownership of these profiles first",
			"code":     "owns_profiles",
			"profiles": owned.Profiles,
		}))
	case errors.Is(err, privacy.ErrExportInProgress):
		return ctx.Results.Error(http.StatusConflict, httpfx.WithJSON(map[string]string{
			"error": "an export is already in progress", "code": "export_in_progress",
		}))
	case errors.Is(err, privacy.ErrExportUnavailable):
		return ctx.Results.Error(http.StatusServiceUnavailable, httpfx.WithJSON(map[string]string{
			"error": "personal data export is not available", "code": "export_unavailable",
		}))
	case errors.Is(err, privacy.ErrConfirmationMismatch):
		return ctx.Results.BadRequest(httpfx.WithJSON(map[string]string{
			"error": "the confirmation does not match your username", "code": "confirmation_mismatch",
		}))
	case errors.Is(err, privacy.ErrExportNotFound), errors.Is(err, privacy.ErrNoAccount):
		return ctx.Results.NotFound(httpfx.WithErrorMessage("not found"))
	case errors.Is(err, privacy.ErrExportNotReady):
		return ctx.Results.Error(http.StatusConflict, httpfx.WithJSON(map[string]string{
			"error": "the export is not ready", "code": "export_not_ready",
		}))
	case errors.Is(err, privacy.ErrExportExpired):
		return ctx.Results.Error(http.StatusGone, httpfx.WithJSON(map[string]string{
			"error": "the export has expired", "code": "export_expired",
		}))
	case errors.Is(err, privacy.ErrAccountDeleted):
		return ctx.Results.Error(
			http.StatusGone,
			httpfx.WithErrorMessage("account has been deleted"),
		)
	default:
		logger.Error("privacy request failed", "error", err)

		return ctx.Results.Error(http.StatusInternalServerError, httpfx.WithSanitizedError(err))
	}
}

func RegisterHTTPRoutesForPrivacy( //nolint:cyclop,funlen
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	privacyService *privacy.Service,
) {
	routes.
		Route(
			"GET /{locale}/privacy",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				userID, ok := privacyUserID(ctx, userService)
				if !ok {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Invalid session"))
				}

				state, err := privacyService.GetState(ctx.Request.Context(), userID)
				if err != nil {
					return privacyError(ctx, logger, err)
				}

				return ctx.Results.JSON(cursors.WrapResponseWithCursor(state, nil))
			},
		).
		HasSummary("Get the signed-in user's privacy state").
		HasDescription("Hidden state and personal data export status of the signed-in user.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"POST /{locale}/privacy/exports",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				userID, ok := privacyUserID(ctx, userService)
				if !ok {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Invalid session"))
				}

				export, err := privacyService.RequestExport(ctx.Request.Context(), userID)
				if err != nil {
					return privacyError(ctx, logger, err)
				}

				return ctx.Results.Accepted(
					httpfx.WithJSON(cursors.WrapResponseWithCursor(map[string]any{
						"id":           export.ID,
						"status":       export.Status,
						"requested_at": export.RequestedAt,
					}, nil)),
				)
			},
		).
		HasSummary("Request a personal data export").
		HasDescription("Queues the collection of the signed-in user's personal data. The user is e-mailed when it is ready.").
		HasResponse(http.StatusAccepted)

	routes.
		Route(
			"GET /{locale}/privacy/exports/{id}/download",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				userID, ok := privacyUserID(ctx, userService)
				if !ok {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Invalid session"))
				}

				url, err := privacyService.GetExportDownloadURL(
					ctx.Request.Context(),
					userID,
					ctx.Request.PathValue("id"),
				)
				if err != nil {
					return privacyError(ctx, logger, err)
				}

				return ctx.Results.JSON(
					cursors.WrapResponseWithCursor(map[string]string{"url": url}, nil),
				)
			},
		).
		HasSummary("Get a download link for a personal data export").
		HasDescription("Returns a short-lived signed URL of a finished export owned by the signed-in user.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"POST /{locale}/privacy/hide",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				userID, ok := privacyUserID(ctx, userService)
				if !ok {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Invalid session"))
				}

				err := privacyService.HideAccount(ctx.Request.Context(), userID)
				if err != nil {
					return privacyError(ctx, logger, err)
				}

				// Every session of the user has just been ended; drop the cookie of this one as well.
				ClearSessionCookie(ctx.ResponseWriter, authService.Config)

				return ctx.Results.JSON(
					cursors.WrapResponseWithCursor(map[string]bool{"hidden": true}, nil),
				)
			},
		).
		HasSummary("Hide the signed-in user's profile").
		HasDescription("Hides the profile and signs the user out. It reappears when they log in again.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"POST /{locale}/privacy/delete",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				userID, ok := privacyUserID(ctx, userService)
				if !ok {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Invalid session"))
				}

				var request deleteAccountRequest

				decodeErr := json.NewDecoder(ctx.Request.Body).Decode(&request)
				if decodeErr != nil {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
				}

				err := privacyService.DeleteAccount(
					ctx.Request.Context(),
					userID,
					request.ConfirmSlug,
					privacy.DeleteOptions{DeleteContent: request.DeleteContent},
				)
				if err != nil {
					return privacyError(ctx, logger, err)
				}

				ClearSessionCookie(ctx.ResponseWriter, authService.Config)

				return ctx.Results.JSON(
					cursors.WrapResponseWithCursor(map[string]bool{"deleted": true}, nil),
				)
			},
		).
		HasSummary("Delete the signed-in user's account").
		HasDescription("Irreversibly erases personal data and frees the username. Requires the username as confirmation.").
		HasResponse(http.StatusOK)
}
