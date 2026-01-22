package http

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/profilelink_oauth"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// RegisterHTTPRoutesForProfileLinkOAuth registers the OAuth routes for profile links.
func RegisterHTTPRoutesForProfileLinkOAuth(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	youtubeProvider *profilelink_oauth.YouTubeOAuthProvider,
	siteURI string,
) {
	// Initiate OAuth flow for connecting a provider to a profile link
	// Returns JSON with auth_url for frontend to redirect to
	routes.Route(
		"POST /{locale}/profiles/{slug}/_links/connect/{provider}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			// Get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")
			providerParam := ctx.Request.PathValue("provider")

			// Validate provider
			if providerParam != "youtube" {
				return ctx.Results.BadRequest(
					httpfx.WithPlainText(
						"Unsupported provider. Only 'youtube' is currently supported.",
					),
				)
			}

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			// Check if user can edit this profile
			canEdit, err := profileService.CanUserEditProfile(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Permission check failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", *session.LoggedInUserID),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to check permissions"),
				)
			}

			if !canEdit {
				return ctx.Results.Error(
					http.StatusForbidden,
					httpfx.WithPlainText("You do not have permission to edit this profile"),
				)
			}

			// Build the redirect URI for OAuth callback (no locale - simpler for Google Console config)
			redirectURI := fmt.Sprintf("%s/profiles/_links/callback/%s",
				siteURI, providerParam)

			// Get the origin from Referer header for redirect after callback
			referer := ctx.Request.Header.Get("Referer")
			redirectOrigin := ""
			if referer != "" {
				if parsedURL, parseErr := url.Parse(referer); parseErr == nil {
					redirectOrigin = fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
				}
			}

			// Initiate OAuth flow
			authURL, _, err := youtubeProvider.InitiateOAuth(
				ctx.Request.Context(),
				redirectURI,
				slugParam,
				localeParam,
				redirectOrigin,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to initiate OAuth",
					slog.String("error", err.Error()),
					slog.String("provider", providerParam),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to initiate OAuth flow"),
				)
			}

			logger.DebugContext(ctx.Request.Context(), "Generated OAuth URL",
				slog.String("provider", providerParam),
				slog.String("slug", slugParam),
				slog.String("auth_url", authURL))

			// Return the auth URL for frontend to redirect
			return ctx.Results.JSON(map[string]any{
				"data": map[string]string{
					"auth_url": authURL,
				},
				"error": nil,
			})
		}).
		HasSummary("Initiate Profile Link OAuth").
		HasDescription("Start the OAuth flow to connect a social media account to a profile. Returns auth_url for frontend redirect.").
		HasResponse(http.StatusOK)

	// OAuth callback handler (no locale in path - simpler for Google Console config)
	routes.Route(
		"GET /profiles/_links/callback/{provider}",
		func(ctx *httpfx.Context) httpfx.Result {
			// Get variables from path and query
			providerParam := ctx.Request.PathValue("provider")

			// Get OAuth callback parameters
			code := ctx.Request.URL.Query().Get("code")
			stateParam := ctx.Request.URL.Query().Get("state")
			errorParam := ctx.Request.URL.Query().Get("error")

			// Validate required parameters (state needed for redirect origin)
			if stateParam == "" {
				return ctx.Results.BadRequest(
					httpfx.WithPlainText("Missing required OAuth state parameter"),
				)
			}

			// Validate provider
			if providerParam != "youtube" {
				return ctx.Results.BadRequest(
					httpfx.WithPlainText("Unsupported provider"),
				)
			}

			// Build the redirect URI (must match what we used in initiate)
			redirectURI := fmt.Sprintf("%s/profiles/_links/callback/%s",
				siteURI, providerParam)

			// Handle the OAuth callback (this also parses the state)
			result, stateObj, err := youtubeProvider.HandleCallback(
				ctx.Request.Context(),
				code,
				stateParam,
				redirectURI,
			)

			// Helper to build redirect URL with fallbacks
			buildRedirectURL := func(path string) string {
				origin := ""
				if stateObj != nil && stateObj.RedirectOrigin != "" {
					origin = stateObj.RedirectOrigin
				}

				return origin + path
			}

			// Check for access denied
			if errorParam == "access_denied" {
				logger.InfoContext(ctx.Request.Context(), "User denied OAuth access",
					slog.String("provider", providerParam))

				locale := "tr"
				if stateObj != nil && stateObj.Locale != "" {
					locale = stateObj.Locale
				}

				return ctx.Results.Redirect(
					buildRedirectURL(fmt.Sprintf("/%s?error=access_denied", locale)),
				)
			}

			// Check for missing code
			if code == "" {
				return ctx.Results.BadRequest(
					httpfx.WithPlainText("Missing authorization code"),
				)
			}

			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "OAuth callback failed",
					slog.String("error", err.Error()),
					slog.String("provider", providerParam))

				locale := "tr"
				if stateObj != nil && stateObj.Locale != "" {
					locale = stateObj.Locale
				}

				return ctx.Results.Redirect(
					buildRedirectURL(fmt.Sprintf("/%s?error=oauth_failed", locale)),
				)
			}

			// Get profile ID from slug
			profileID, err := profileService.GetProfileIDBySlug(
				ctx.Request.Context(),
				stateObj.ProfileSlug,
			)
			if err != nil || profileID == "" {
				logger.ErrorContext(ctx.Request.Context(), "Profile not found",
					slog.String("slug", stateObj.ProfileSlug))

				redirectURL := fmt.Sprintf(
					"%s/%s?error=profile_not_found",
					stateObj.RedirectOrigin,
					stateObj.Locale,
				)

				return ctx.Results.Redirect(redirectURL)
			}

			// Check if a link with this remote ID already exists
			existingLink, err := profileService.GetProfileLinkByRemoteID(
				ctx.Request.Context(),
				profileID,
				result.Kind,
				result.RemoteID,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to check existing link",
					slog.String("error", err.Error()),
					slog.String("profile_id", profileID),
					slog.String("remote_id", result.RemoteID))

				redirectURL := fmt.Sprintf("%s/%s/%s/settings/links?error=oauth_failed",
					stateObj.RedirectOrigin, stateObj.Locale, stateObj.ProfileSlug)

				return ctx.Results.Redirect(redirectURL)
			}

			var expiresAt *sql.NullTime
			if result.AccessTokenExpiresAt != nil {
				expiresAt = &sql.NullTime{Time: *result.AccessTokenExpiresAt, Valid: true}
			}

			if existingLink != nil {
				// Update existing link with new tokens
				err = profileService.UpdateProfileLinkOAuthTokens(
					ctx.Request.Context(),
					existingLink.ID,
					result.PublicID,
					result.URI,
					result.Title,
					result.Scope,
					result.AccessToken,
					expiresAt,
					&result.RefreshToken,
				)
				if err != nil {
					logger.ErrorContext(ctx.Request.Context(), "Failed to update OAuth tokens",
						slog.String("error", err.Error()),
						slog.String("link_id", existingLink.ID))

					redirectURL := fmt.Sprintf("%s/%s/%s/settings/links?error=update_failed",
						stateObj.RedirectOrigin, stateObj.Locale, stateObj.ProfileSlug)

					return ctx.Results.Redirect(redirectURL)
				}

				logger.InfoContext(ctx.Request.Context(), "Updated OAuth tokens for existing link",
					slog.String("link_id", existingLink.ID),
					slog.String("provider", providerParam))
			} else {
				// Create new OAuth profile link
				linkID := lib.IDsGenerateUnique()

				// Get the next order value
				maxOrder, _ := profileService.GetMaxProfileLinkOrder(ctx.Request.Context(), profileID)
				newOrder := maxOrder + 1

				_, err = profileService.CreateOAuthProfileLink(
					ctx.Request.Context(),
					linkID,
					result.Kind,
					profileID,
					newOrder,
					result.RemoteID,
					result.PublicID,
					result.URI,
					result.Title,
					providerParam,
					result.Scope,
					result.AccessToken,
					expiresAt,
					&result.RefreshToken,
				)
				if err != nil {
					logger.ErrorContext(ctx.Request.Context(), "Failed to create OAuth profile link",
						slog.String("error", err.Error()),
						slog.String("profile_id", profileID))

					redirectURL := fmt.Sprintf("%s/%s/%s/settings/links?error=create_failed",
						stateObj.RedirectOrigin, stateObj.Locale, stateObj.ProfileSlug)

					return ctx.Results.Redirect(redirectURL)
				}

				logger.InfoContext(ctx.Request.Context(), "Created OAuth profile link",
					slog.String("link_id", linkID),
					slog.String("provider", providerParam),
					slog.String("channel_id", result.RemoteID))
			}

			// Redirect to the settings page with success message
			redirectURL := fmt.Sprintf("%s/%s/%s/settings/links?connected=%s",
				stateObj.RedirectOrigin, stateObj.Locale, stateObj.ProfileSlug, providerParam)

			return ctx.Results.Redirect(redirectURL)
		}).
		HasSummary("Profile Link OAuth Callback").
		HasDescription("Handle OAuth callback from providers and create/update profile links.").
		HasResponse(http.StatusTemporaryRedirect)
}
