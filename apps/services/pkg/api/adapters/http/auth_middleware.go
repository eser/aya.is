package http

import (
	"context"
	"log/slog"
	"strings"

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
		authHeader := ctx.Request.Header.Get(AuthHeader)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return ctx.Results.Unauthorized(httpfx.WithPlainText("Unauthorized"))
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Debug: Log JWT secret length (not the actual secret for security)
		secretLen := len(authService.Config.JwtSecret)
		slog.Info("AuthMiddleware: Validating token",
			slog.Int("jwt_secret_length", secretLen),
			slog.Int("token_length", len(tokenStr)),
		)

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				slog.Warn("AuthMiddleware: Unexpected signing method",
					slog.String("method", token.Method.Alg()),
				)

				return nil, auth.ErrInvalidSigningMethod
			}

			return []byte(authService.Config.JwtSecret), nil
		})
		if err != nil {
			slog.Warn("AuthMiddleware: Token parse failed",
				slog.String("error", err.Error()),
			)

			return ctx.Results.Unauthorized(httpfx.WithPlainText("Invalid token"))
		}

		if !token.Valid {
			slog.Warn("AuthMiddleware: Token invalid after parse")

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
		if err != nil || session == nil || session.Status != users.SessionStatusActive {
			return ctx.Results.Unauthorized(httpfx.WithPlainText("Session invalid"))
		}

		// Update last activity and user agent
		userAgent := ctx.Request.Header.Get("User-Agent")
		_ = userService.UpdateSessionActivity(ctx.Request.Context(), sessionID, &userAgent)

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
