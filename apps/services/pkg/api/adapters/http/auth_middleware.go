package http

import (
	"context"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

const ContextKeySessionID httpfx.ContextKey = "session_id"

// AuthMiddleware resolves the caller's session from the session cookie.
// The cookie carries the session ID directly. The session is loaded from
// the database and must be active with a logged-in user.
// The session ID is stored in context for downstream handlers.
func AuthMiddleware(authService *auth.Service, userService *users.Service) httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.Result {
		sessionID, err := GetSessionIDFromCookie(ctx.Request, authService.Config)
		if err != nil || sessionID == "" {
			return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Unauthorized"))
		}

		// Load session from repository and verify it is active
		session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
		if err != nil || session == nil || session.Status != users.SessionStatusActive {
			return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Session invalid"))
		}

		// Protected routes require an authenticated (non-anonymous) session
		if session.LoggedInUserID == nil {
			return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Unauthorized"))
		}

		// Update last activity and user agent
		userAgent := ctx.Request.Header.Get("User-Agent")

		if err := userService.UpdateSessionActivity(ctx.Request.Context(), sessionID, &userAgent); err != nil {
			slog.Warn(
				"failed to update session activity",
				slog.String("error", err.Error()),
				slog.String("session_id", sessionID),
			)
		}

		// Store session ID in context for route handlers
		newContext := context.WithValue(
			ctx.Request.Context(),
			ContextKeySessionID,
			sessionID,
		)
		ctx.UpdateContext(newContext)

		return ctx.Next()
	}
}
