package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/github"
	"github.com/eser/aya.is/services/pkg/api/adapters/youtube"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// ProfileLinkProviders contains all external service providers.
type ProfileLinkProviders struct {
	YouTube                *youtube.Provider
	GitHub                 *github.Provider
	PendingConnectionStore *profiles.PendingConnectionStore
}

// RegisterHTTPRoutesForProfileLinks registers the OAuth routes for profile links.
func RegisterHTTPRoutesForProfileLinks(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	providers *ProfileLinkProviders,
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
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get variables from path
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")
			providerParam := ctx.Request.PathValue("provider")

			// Validate provider
			if providerParam != "youtube" && providerParam != "github" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage(
						"Unsupported provider. Supported: 'youtube', 'github'.",
					),
				)
			}

			// Get user ID from session
			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

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

			// Get profile to determine its kind
			profile, err := profileService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
			if err != nil || profile == nil {
				return ctx.Results.Error(
					http.StatusNotFound,
					httpfx.WithErrorMessage("Profile not found"),
				)
			}

			// Build the redirect URI for OAuth callback
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

			// Generate state for linking flow (service layer responsibility)
			_, encodedState, err := profiles.CreateProfileLinkState(
				slugParam,
				profile.Kind,
				localeParam,
				redirectOrigin,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to create profile link state",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to initiate OAuth flow"),
				)
			}

			// Initiate OAuth flow based on provider
			var authURL string
			switch providerParam {
			case "youtube":
				authURL, err = providers.YouTube.InitiateOAuth(
					ctx.Request.Context(),
					redirectURI,
					encodedState,
				)
			case "github":
				authURL, err = providers.GitHub.InitiateProfileLinkOAuth(
					ctx.Request.Context(),
					redirectURI,
					encodedState,
				)
			}

			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to initiate OAuth",
					slog.String("error", err.Error()),
					slog.String("provider", providerParam),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to initiate OAuth flow"),
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

	// OAuth callback handler (no locale in path - simpler for OAuth app config)
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
					httpfx.WithErrorMessage("Missing required OAuth state parameter"),
				)
			}

			// Validate provider
			if providerParam != "youtube" && providerParam != "github" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Unsupported provider"),
				)
			}

			// Decode and validate state at service layer
			stateObj, stateErr := profiles.DecodeProfileLinkState(stateParam)
			if stateErr != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to decode state",
					slog.String("error", stateErr.Error()),
					slog.String("provider", providerParam))

				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Invalid OAuth state"),
				)
			}

			// Validate state expiry
			validationErr := profiles.ValidateProfileLinkState(stateObj)
			if validationErr != nil {
				logger.ErrorContext(ctx.Request.Context(), "State validation failed",
					slog.String("error", validationErr.Error()),
					slog.String("provider", providerParam))

				redirectURL := fmt.Sprintf("%s/%s?error=state_expired",
					stateObj.RedirectOrigin, stateObj.Locale)

				return ctx.Results.Redirect(redirectURL)
			}

			// Build the redirect URI (must match what we used in initiate)
			redirectURI := fmt.Sprintf("%s/profiles/_links/callback/%s",
				siteURI, providerParam)

			// Handle the OAuth callback based on provider
			var result auth.OAuthCallbackResult
			var err error

			switch providerParam {
			case "youtube":
				result, err = providers.YouTube.HandleOAuthCallback(
					ctx.Request.Context(),
					code,
					redirectURI,
				)
			case "github":
				result, err = providers.GitHub.HandleOAuthCallback(
					ctx.Request.Context(),
					code,
					redirectURI,
				)
			}

			// Helper to build redirect URL
			buildRedirectURL := func(path string) string {
				return stateObj.RedirectOrigin + path
			}

			// Check for access denied
			if errorParam == "access_denied" {
				logger.InfoContext(ctx.Request.Context(), "User denied OAuth access",
					slog.String("provider", providerParam))

				return ctx.Results.Redirect(
					buildRedirectURL(fmt.Sprintf("/%s?error=access_denied", stateObj.Locale)),
				)
			}

			// Check for missing code
			if code == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Missing authorization code"),
				)
			}

			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "OAuth callback failed",
					slog.String("error", err.Error()),
					slog.String("provider", providerParam))

				return ctx.Results.Redirect(
					buildRedirectURL(fmt.Sprintf("/%s?error=oauth_failed", stateObj.Locale)),
				)
			}

			// For GitHub with organization/product profiles, store pending connection for account selection
			if providerParam == "github" &&
				(stateObj.ProfileKind == "organization" || stateObj.ProfileKind == "product") {
				pendingConn := &profiles.PendingGitHubConnection{
					AccessToken: result.AccessToken,
					Scope:       result.Scope,
					ProfileSlug: stateObj.ProfileSlug,
					ProfileKind: stateObj.ProfileKind,
					Locale:      stateObj.Locale,
				}

				pendingID := providers.PendingConnectionStore.Store(pendingConn)

				logger.DebugContext(
					ctx.Request.Context(),
					"Stored pending GitHub connection for account selection",
					slog.String("pending_id", pendingID),
					slog.String("profile_slug", stateObj.ProfileSlug),
				)

				// Redirect with pending status for frontend to show account selection
				redirectURL := fmt.Sprintf("%s/%s/%s/settings/links?pending=github&pending_id=%s",
					stateObj.RedirectOrigin, stateObj.Locale, stateObj.ProfileSlug, pendingID)

				return ctx.Results.Redirect(redirectURL)
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

			// Determine link kind from provider
			linkKind := providerParam // "youtube" or "github"

			// Check if a link with this remote ID already exists
			existingLink, err := profileService.GetProfileLinkByRemoteID(
				ctx.Request.Context(),
				profileID,
				linkKind,
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
					stateObj.Locale,
					result.Username,
					result.URI,
					result.Name,
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

				logger.DebugContext(ctx.Request.Context(), "Updated OAuth tokens for existing link",
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
					linkKind,
					profileID,
					newOrder,
					stateObj.Locale,
					result.RemoteID,
					result.Username,
					result.URI,
					result.Name,
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

				logger.DebugContext(ctx.Request.Context(), "Created OAuth profile link",
					slog.String("link_id", linkID),
					slog.String("provider", providerParam),
					slog.String("remote_id", result.RemoteID))
			}

			// Redirect to the settings page with success message
			redirectURL := fmt.Sprintf("%s/%s/%s/settings/links?connected=%s",
				stateObj.RedirectOrigin, stateObj.Locale, stateObj.ProfileSlug, providerParam)

			return ctx.Results.Redirect(redirectURL)
		}).
		HasSummary("Profile Link OAuth Callback").
		HasDescription("Handle OAuth callback from providers and create/update profile links.").
		HasResponse(http.StatusTemporaryRedirect)

	// Get available GitHub accounts for selection (user + organizations)
	routes.Route(
		"GET /{locale}/profiles/{slug}/_links/github/accounts",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			pendingID := ctx.Request.URL.Query().Get("pending_id")
			if pendingID == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Missing pending_id parameter"),
				)
			}

			// Get pending connection
			pendingConn := providers.PendingConnectionStore.Get(pendingID)
			if pendingConn == nil {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Pending connection not found or expired"),
				)
			}

			// Fetch user info
			userInfo, err := providers.GitHub.Client().FetchUserInfo(
				ctx.Request.Context(),
				pendingConn.AccessToken,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to fetch GitHub user info",
					slog.String("error", err.Error()))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to fetch GitHub account info"),
				)
			}

			// Fetch organizations
			orgs, err := providers.GitHub.Client().FetchUserOrganizations(
				ctx.Request.Context(),
				pendingConn.AccessToken,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to fetch GitHub organizations",
					slog.String("error", err.Error()))
				// Don't fail, just return empty orgs
				orgs = []*github.OrgInfo{}
			}

			// Build response with user account and organizations
			// Construct html_url from login if not provided
			userHTMLURL := userInfo.HTMLURL
			if userHTMLURL == "" {
				userHTMLURL = "https://github.com/" + userInfo.Login
			}

			accounts := []profiles.GitHubAccount{
				{
					ID:        strconv.FormatInt(userInfo.ID, 10),
					Login:     userInfo.Login,
					Name:      userInfo.Name,
					AvatarURL: userInfo.Avatar,
					HTMLURL:   userHTMLURL,
					Type:      "User",
				},
			}

			for _, org := range orgs {
				name := org.Name
				if name == "" {
					name = org.Login
				}

				// GitHub /user/orgs API doesn't return html_url, construct it from login
				htmlURL := org.HTMLURL
				if htmlURL == "" {
					htmlURL = "https://github.com/" + org.Login
				}

				accounts = append(accounts, profiles.GitHubAccount{
					ID:          strconv.FormatInt(org.ID, 10),
					Login:       org.Login,
					Name:        name,
					AvatarURL:   org.Avatar,
					HTMLURL:     htmlURL,
					Type:        "Organization",
					Description: org.Description,
				})
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]any{
					"accounts":     accounts,
					"profile_kind": pendingConn.ProfileKind,
				},
				"error": nil,
			})
		}).
		HasSummary("Get GitHub Accounts").
		HasDescription("Get available GitHub accounts (user and organizations) for linking.").
		HasResponse(http.StatusOK)

	// Finalize GitHub connection with selected account
	routes.Route(
		"POST /{locale}/profiles/{slug}/_links/github/finalize",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			slugParam := ctx.Request.PathValue("slug")

			// Parse request body
			var reqBody struct {
				PendingID string `json:"pending_id"`
				AccountID string `json:"account_id"`
				Login     string `json:"login"`
				Name      string `json:"name"`
				HTMLURL   string `json:"html_url"`
				Type      string `json:"type"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&reqBody); err != nil {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			if reqBody.PendingID == "" || reqBody.AccountID == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Missing required fields"),
				)
			}

			// Get pending connection
			pendingConn := providers.PendingConnectionStore.Get(reqBody.PendingID)
			if pendingConn == nil {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Pending connection not found or expired"),
				)
			}

			// Verify the pending connection is for this profile
			if pendingConn.ProfileSlug != slugParam {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Pending connection does not match profile"),
				)
			}

			// Get profile ID
			profileID, err := profileService.GetProfileIDBySlug(
				ctx.Request.Context(),
				slugParam,
			)
			if err != nil || profileID == "" {
				return ctx.Results.Error(
					http.StatusNotFound,
					httpfx.WithErrorMessage("Profile not found"),
				)
			}

			// Check if a link with this remote ID already exists
			existingLink, err := profileService.GetProfileLinkByRemoteID(
				ctx.Request.Context(),
				profileID,
				"github",
				reqBody.AccountID,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to check existing link",
					slog.String("error", err.Error()))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to check existing link"),
				)
			}

			if existingLink != nil {
				// Update existing link
				err = profileService.UpdateProfileLinkOAuthTokens(
					ctx.Request.Context(),
					existingLink.ID,
					pendingConn.Locale,
					reqBody.Login,
					reqBody.HTMLURL,
					reqBody.Name,
					pendingConn.Scope,
					pendingConn.AccessToken,
					nil,
					nil,
				)
				if err != nil {
					logger.ErrorContext(ctx.Request.Context(), "Failed to update link",
						slog.String("error", err.Error()))

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithErrorMessage("Failed to update link"),
					)
				}
			} else {
				// Create new link
				linkID := lib.IDsGenerateUnique()
				maxOrder, _ := profileService.GetMaxProfileLinkOrder(ctx.Request.Context(), profileID)

				_, err = profileService.CreateOAuthProfileLink(
					ctx.Request.Context(),
					linkID,
					"github",
					profileID,
					maxOrder+1,
					pendingConn.Locale,
					reqBody.AccountID,
					reqBody.Login,
					reqBody.HTMLURL,
					reqBody.Name,
					"github",
					pendingConn.Scope,
					pendingConn.AccessToken,
					nil,
					nil,
				)
				if err != nil {
					logger.ErrorContext(ctx.Request.Context(), "Failed to create link",
						slog.String("error", err.Error()))

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithErrorMessage("Failed to create link"),
					)
				}
			}

			// Clean up pending connection
			providers.PendingConnectionStore.Delete(reqBody.PendingID)

			logger.DebugContext(ctx.Request.Context(), "Finalized GitHub connection",
				slog.String("profile_slug", slugParam),
				slog.String("account_login", reqBody.Login),
				slog.String("account_type", reqBody.Type))

			return ctx.Results.JSON(map[string]any{
				"data": map[string]string{
					"status": "connected",
				},
				"error": nil,
			})
		}).
		HasSummary("Finalize GitHub Connection").
		HasDescription("Complete the GitHub connection with the selected account.").
		HasResponse(http.StatusOK)
}
