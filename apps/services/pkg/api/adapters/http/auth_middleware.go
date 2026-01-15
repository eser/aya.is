package http

import (
	"context"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthHeader                            = "Authorization"
	ContextKeySessionID httpfx.ContextKey = "session_id"
)

func AuthMiddleware(authService *auth.Service, userService *users.Service) httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.Result {
		// FIXME(@eser) no need to check if the header is specified
		auth := ctx.Request.Header.Get(AuthHeader)

		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return ctx.Results.Unauthorized(httpfx.WithPlainText("Unauthorized"))
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			return []byte(authService.Config.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			return ctx.Results.Unauthorized(httpfx.WithPlainText("Invalid token"))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Results.Unauthorized(httpfx.WithPlainText("Invalid claims"))
		}

		sessionID, _ := claims["session_id"].(string)
		if sessionID == "" {
			return ctx.Results.Unauthorized(httpfx.WithPlainText("No session"))
		}

		// Load session from repository
		session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
		if err != nil || session.Status != "active" {
			return ctx.Results.Unauthorized(httpfx.WithPlainText("Session invalid"))
		}

		// Update logged_in_at
		_ = userService.UpdateSessionLoggedInAt(ctx.Request.Context(), sessionID, time.Now())

		// Store session ID in context for route handlers
		newContext := context.WithValue(
			ctx.Request.Context(),
			ContextKeySessionID,
			sessionID,
		)
		ctx.UpdateContext(newContext)

		result := ctx.Next()

		return result
	}
}
