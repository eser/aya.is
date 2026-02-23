package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
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
	profileService *profiles.Service,
) {
	routes.
		Route(
			"GET /{locale}/users",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				cursor := cursors.NewCursorFromRequest(ctx.Request)

				records, err := userService.List(ctx.Request.Context(), cursor)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(records)
			},
		).
		HasSummary("List users").
		HasDescription("List users.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/users/{id}",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				idParam := ctx.Request.PathValue("id")

				record, err := userService.GetByID(ctx.Request.Context(), idParam)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

				return ctx.Results.JSON(wrappedResponse)
			},
		).
		HasSummary("Get user by ID").
		HasDescription("Get user by ID.").
		HasResponse(http.StatusOK)

	// --- Auth endpoints ---
	routes.
		Route("GET /{locale}/auth/{authProvider}/login", func(ctx *httpfx.Context) httpfx.Result {
			// get auth provider from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			authProvider := ctx.Request.PathValue("authProvider")

			requestURL := ctx.Request.URL
			queryString := requestURL.Query()
			redirectURI := queryString.Get("redirect_uri")

			if redirectURI == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("redirect_uri is required"))
			}

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

				requestURL := ctx.Request.URL
				queryString := requestURL.Query()
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

				// Auto-sync managed GitHub link for users with an individual profile.
				// This ensures the profile link tokens are refreshed on every login.
				if result.User.IndividualProfileID != nil &&
					result.Session != nil &&
					result.Session.OAuthProvider != nil &&
					*result.Session.OAuthProvider == "github" &&
					result.Session.OAuthAccessToken != nil &&
					result.User.GithubRemoteID != nil {
					githubHandle := ""
					if result.User.GithubHandle != nil {
						githubHandle = *result.User.GithubHandle
					}

					tokenScope := ""
					if result.Session.OAuthTokenScope != nil {
						tokenScope = *result.Session.OAuthTokenScope
					}

					upsertManagedGitHubLink(
						ctx.Request.Context(),
						logger,
						profileService,
						*result.User.IndividualProfileID,
						*result.User.GithubRemoteID,
						githubHandle,
						*result.Session.OAuthAccessToken,
						tokenScope,
					)
				}

				if result.RedirectURI != "" {
					return ctx.Results.Redirect(result.RedirectURI)
				}

				// Set JWT as cookie or return in response
				return ctx.Results.JSON(map[string]any{
					"data": map[string]any{
						"token": result.JWT,
						"user":  result.User,
					},
					"error": nil,
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

				return ctx.Results.JSON(map[string]any{
					"data":  map[string]string{"status": "logged out"},
					"error": nil,
				})
			}

			// Logout session: invalidate old, create new anonymous with same preferences
			result, err := sessionService.LogoutSession(ctx.Request.Context(), sessionID)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to logout session",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID))
				// Still clear the cookie even if logout fails
				ClearSessionCookie(ctx.ResponseWriter, authService.Config)

				return ctx.Results.JSON(map[string]any{
					"data":  map[string]string{"status": "logged out"},
					"error": nil,
				})
			}

			// Set new session cookie
			expiresAt := time.Now().Add(24 * time.Hour)
			SetSessionCookie(
				ctx.ResponseWriter,
				result.NewSession.ID,
				expiresAt,
				authService.Config,
			)

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]string{"status": "logged out"},
				"error": nil,
			})
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
				"data": map[string]any{
					"token":     result.JWT,
					"expiresAt": result.ExpiresAt.Unix(),
				},
				"error": nil,
			})
		}).
		HasSummary("Refresh Token").
		HasDescription("Refreshes JWT token before expiration.").
		HasResponse(http.StatusOK)
}

// upsertManagedGitHubLink creates or updates the managed GitHub profile link
// for a user's individual profile. Called on login (to refresh tokens) and
// on individual profile creation (to auto-create the link).
// Token parameters come from the session's OAuth fields.
func upsertManagedGitHubLink(
	ctx context.Context,
	logger *logfx.Logger,
	profileService *profiles.Service,
	profileID string,
	githubRemoteID string,
	githubHandle string,
	accessToken string,
	tokenScope string,
) {
	uri := "https://github.com/" + githubHandle

	// First check if ANY managed GitHub link already exists for this profile.
	// This prevents duplicates even when remote_id differs (e.g., manual link vs auto-sync).
	managedLink, _ := profileService.GetManagedGitHubLink(ctx, profileID)
	if managedLink != nil {
		err := profileService.UpdateProfileLinkOAuthTokens(
			ctx,
			managedLink.ID,
			"en",
			githubHandle,
			uri,
			"GitHub",
			tokenScope,
			accessToken,
			nil, // accessTokenExpiresAt — GitHub tokens don't expire
			nil, // refreshToken
		)
		if err != nil {
			logger.WarnContext(ctx, "Failed to update managed GitHub link tokens",
				slog.String("profile_id", profileID),
				slog.String("error", err.Error()))
		} else {
			logger.DebugContext(ctx, "Updated managed GitHub link tokens",
				slog.String("profile_id", profileID))
		}

		return
	}

	// Also check by remote_id in case there's a non-managed link with the same remote_id
	existingLink, _ := profileService.GetProfileLinkByRemoteID(
		ctx, profileID, "github", githubRemoteID,
	)

	if existingLink != nil {
		err := profileService.UpdateProfileLinkOAuthTokens(
			ctx,
			existingLink.ID,
			"en",
			githubHandle,
			uri,
			"GitHub",
			tokenScope,
			accessToken,
			nil, // accessTokenExpiresAt — GitHub tokens don't expire
			nil, // refreshToken
		)
		if err != nil {
			logger.WarnContext(ctx, "Failed to update managed GitHub link tokens",
				slog.String("profile_id", profileID),
				slog.String("error", err.Error()))
		} else {
			logger.DebugContext(ctx, "Updated managed GitHub link tokens",
				slog.String("profile_id", profileID))
		}

		return
	}

	// Create new managed GitHub link
	maxOrder, _ := profileService.GetMaxProfileLinkOrder(ctx, profileID)
	linkID := lib.IDsGenerateUnique()

	_, err := profileService.CreateOAuthProfileLink(
		ctx,
		linkID,
		"github",
		profileID,
		maxOrder+1,
		"en",
		githubRemoteID,
		githubHandle,
		uri,
		"GitHub",
		"github",
		tokenScope,
		accessToken,
		nil, // accessTokenExpiresAt
		nil, // refreshToken
	)
	if err != nil {
		logger.WarnContext(ctx, "Failed to create managed GitHub link",
			slog.String("profile_id", profileID),
			slog.String("error", err.Error()))
	} else {
		logger.DebugContext(ctx, "Created managed GitHub link on login",
			slog.String("profile_id", profileID))
	}
}
