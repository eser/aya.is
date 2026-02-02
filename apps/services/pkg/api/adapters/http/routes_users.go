package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/sessions"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForUsers( //nolint:funlen,cyclop
	baseURI string,
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	sessionService *sessions.Service,
) {
	routes.
		Route("GET /{locale}/users", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			records, err := userService.List(ctx.Request.Context(), cursor)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage(err.Error()),
				)
			}

			return ctx.Results.JSON(records)
		}).
		HasSummary("List users").
		HasDescription("List users.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/users/{id}", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			idParam := ctx.Request.PathValue("id")

			record, err := userService.GetByID(ctx.Request.Context(), idParam)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage(err.Error()),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get user by ID").
		HasDescription("Get user by ID.").
		HasResponse(http.StatusOK)

	// --- Auth endpoints ---
	routes.
		Route("GET /{locale}/auth/{authProvider}/login", func(ctx *httpfx.Context) httpfx.Result {
			// get auth provider from path
			localeParam := ctx.Request.PathValue("locale")
			authProvider := ctx.Request.PathValue("authProvider")

			url := ctx.Request.URL
			queryString := url.Query()
			redirectURI := queryString.Get("redirect_uri")

			// if redirectURI == "" {
			// 	return ctx.Results.BadRequest(httpfx.WithErrorMessage("redirect_uri is required"))
			// }

			// Initiate auth flow
			authURL, err := authService.Initiate(
				ctx.Request.Context(),
				authProvider,
				baseURI+"/"+localeParam,
				redirectURI,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Auth initiation failed"),
				)
			}

			logger.DebugContext(ctx.Request.Context(), "Redirecting to auth provider login",
				slog.String("auth_url", authURL),
				slog.String("redirect_uri", redirectURI),
				slog.String("auth_provider", authProvider))

			// FIXME(@eser) Optionally set state in cookie/session
			return ctx.Results.Redirect(authURL)
		}).
		HasSummary("Auth Login").
		HasDescription("Redirects to auth provider login.").
		HasResponse(http.StatusFound)

	routes.
		Route(
			"GET /{locale}/auth/{authProvider}/callback",
			func(ctx *httpfx.Context) httpfx.Result {
				// get auth provider from path
				authProvider := ctx.Request.PathValue("authProvider")

				url := ctx.Request.URL
				queryString := url.Query()
				code := queryString.Get("code")
				state := queryString.Get("state")
				redirectURI := queryString.Get("redirect_uri")

				if code == "" {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("code is required"))
				}

				if state == "" {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("state is required"))
				}

				result, err := authService.AuthHandleCallback(
					ctx.Request.Context(),
					authProvider,
					code,
					state,
					redirectURI,
				)
				if err != nil {
					logger.ErrorContext(ctx.Request.Context(), "Auth callback failed",
						slog.String("error", err.Error()))

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithErrorMessage("Auth callback failed"),
					)
				}

				// Set session cookie for cross-domain SSO
				SetSessionCookie(
					ctx.ResponseWriter,
					result.SessionID,
					result.ExpiresAt,
					authService.Config,
				)

				if result.RedirectURI != "" {
					return ctx.Results.Redirect(result.RedirectURI)
				}

				// Set JWT as cookie or return in response
				return ctx.Results.JSON(map[string]any{
					"token": result.JWT,
					"user":  result.User,
				})
			},
		).
		HasSummary("Auth Callback").
		HasDescription("Handles auth provider callback and returns JWT.").
		HasResponse(http.StatusOK)

	routes.
		Route("POST /{locale}/auth/logout", func(ctx *httpfx.Context) httpfx.Result {
			// Get current session ID from cookie
			sessionID, err := GetSessionIDFromCookie(ctx.Request, authService.Config)
			if err != nil {
				// No session to logout, just clear cookie and return success
				ClearSessionCookie(ctx.ResponseWriter, authService.Config)

				return ctx.Results.JSON(map[string]string{"status": "logged out"})
			}

			// Logout session: invalidate old, create new anonymous with same preferences
			result, err := sessionService.LogoutSession(ctx.Request.Context(), sessionID)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to logout session",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID))
				// Still clear the cookie even if logout fails
				ClearSessionCookie(ctx.ResponseWriter, authService.Config)

				return ctx.Results.JSON(map[string]string{"status": "logged out"})
			}

			// Set new session cookie
			expiresAt := time.Now().Add(24 * time.Hour)
			SetSessionCookie(
				ctx.ResponseWriter,
				result.NewSession.ID,
				expiresAt,
				authService.Config,
			)

			return ctx.Results.JSON(map[string]string{"status": "logged out"})
		}).
		HasSummary("Logout").
		HasDescription("Logs out the user and creates a new anonymous session.").
		HasResponse(http.StatusOK)

	routes.
		Route("POST /{locale}/auth/refresh", func(ctx *httpfx.Context) httpfx.Result {
			// Get current token from Authorization header
			authHeader := ctx.Request.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return ctx.Results.Unauthorized(httpfx.WithErrorMessage("No token provided"))
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			// Use the service to refresh the token
			result, err := authService.RefreshToken(ctx.Request.Context(), tokenStr)
			if err != nil {
				if errors.Is(err, auth.ErrInvalidToken) {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Invalid token"))
				}
				if errors.Is(err, auth.ErrSessionExpired) {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Session expired"))
				}

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to refresh token"),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"token":     result.JWT,
				"expiresAt": result.ExpiresAt.Unix(),
			})
		}).
		HasSummary("Refresh Token").
		HasDescription("Refreshes JWT token before expiration.").
		HasResponse(http.StatusOK)
}
