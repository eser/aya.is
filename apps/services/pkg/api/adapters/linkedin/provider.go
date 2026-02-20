package linkedin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
)

// Sentinel errors.
var (
	ErrFailedToExchangeCode = errors.New("failed to exchange authorization code")
	ErrFailedToGetUserInfo  = errors.New("failed to get LinkedIn user info")
	ErrFailedToGetOrgs      = errors.New("failed to get LinkedIn organization pages")
)

// HTTPClient interface for dependency injection.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Provider implements OAuth for LinkedIn profile linking.
type Provider struct {
	config     *auth.LinkedInOAuthConfig
	logger     *logfx.Logger
	httpClient HTTPClient
}

// NewProvider creates a new LinkedIn provider.
func NewProvider(
	config *auth.LinkedInOAuthConfig,
	logger *logfx.Logger,
	httpClient HTTPClient,
) *Provider {
	return &Provider{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}
}

// InitiateOAuth builds the OAuth URL with basic scope.
func (p *Provider) InitiateOAuth(
	ctx context.Context,
	redirectURI string,
	state string,
) (string, error) {
	authURL := p.buildAuthURL(redirectURI, state, p.config.Scope)

	p.logger.DebugContext(ctx, "Initiating LinkedIn OAuth",
		slog.String("redirect_uri", redirectURI))

	return authURL, nil
}

// InitiateProfileLinkOAuth builds the OAuth URL with expanded scope for profile linking.
// Uses ProfileLinkScope which includes r_organization_social for organization page access.
func (p *Provider) InitiateProfileLinkOAuth(
	ctx context.Context,
	redirectURI string,
	state string,
) (string, error) {
	authURL := p.buildAuthURL(redirectURI, state, p.config.ProfileLinkScope)

	p.logger.DebugContext(ctx, "Initiating LinkedIn Profile Link OAuth",
		slog.String("redirect_uri", redirectURI))

	return authURL, nil
}

// HandleOAuthCallback exchanges the code for tokens and returns account info.
func (p *Provider) HandleOAuthCallback(
	ctx context.Context,
	code string,
	redirectURI string,
) (auth.OAuthCallbackResult, error) {
	p.logger.DebugContext(ctx, "Processing LinkedIn OAuth callback")

	// Exchange code for tokens
	tokenResp, err := p.exchangeCodeForTokens(ctx, code, redirectURI)
	if err != nil {
		return auth.OAuthCallbackResult{}, err
	}

	// Fetch user info via OpenID Connect
	userInfo, err := p.FetchUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return auth.OAuthCallbackResult{}, err
	}

	// Try /v2/me to get vanityName for constructing profile URL
	vanityName := p.fetchVanityName(ctx, tokenResp.AccessToken)

	profileURI := ""
	username := userInfo.Name

	if vanityName != "" {
		profileURI = "https://www.linkedin.com/in/" + vanityName
		username = vanityName
	}

	// Calculate token expiry
	var expiresAt *time.Time

	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	p.logger.DebugContext(ctx, "LinkedIn OAuth callback successful",
		slog.String("sub", userInfo.Sub),
		slog.String("name", userInfo.Name),
		slog.String("vanity_name", vanityName))

	return auth.OAuthCallbackResult{
		RemoteID:             userInfo.Sub,
		Username:             username,
		Name:                 userInfo.Name,
		Email:                userInfo.Email,
		URI:                  profileURI,
		AccessToken:          tokenResp.AccessToken,
		RefreshToken:         tokenResp.RefreshToken,
		AccessTokenExpiresAt: expiresAt,
		Scope:                tokenResp.Scope,
	}, nil
}

// UserInfo represents LinkedIn user information from the OpenID Connect userinfo endpoint.
type UserInfo struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

// OrgPageInfo represents a LinkedIn organization page.
type OrgPageInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	LogoURL    string `json:"logo_url,omitempty"`
	VanityName string `json:"vanity_name,omitempty"`
	URI        string `json:"uri"`
}

// tokenResponse represents LinkedIn's token endpoint response.
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

func (p *Provider) buildAuthURL(redirectURI, state, scope string) string {
	queryString := url.Values{}
	queryString.Set("response_type", "code")
	queryString.Set("client_id", p.config.ClientID)
	queryString.Set("redirect_uri", redirectURI)
	queryString.Set("state", state)
	queryString.Set("scope", scope)

	oauthURL := url.URL{ //nolint:exhaustruct
		Scheme:   "https",
		Host:     "www.linkedin.com",
		Path:     "/oauth/v2/authorization",
		RawQuery: queryString.Encode(),
	}

	return oauthURL.String()
}

func (p *Provider) exchangeCodeForTokens(
	ctx context.Context,
	code string,
	redirectURI string,
) (*tokenResponse, error) {
	values := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"client_id":     {p.config.ClientID},
		"client_secret": {p.config.ClientSecret},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://www.linkedin.com/oauth/v2/accessToken",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.ErrorContext(ctx, "Failed to exchange code for tokens",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		p.logger.ErrorContext(ctx, "Token exchange failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToExchangeCode, resp.StatusCode)
	}

	var tokenResp tokenResponse

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}

	if tokenResp.AccessToken == "" {
		return nil, ErrFailedToExchangeCode
	}

	return &tokenResp, nil
}

// FetchUserInfo retrieves user profile from the LinkedIn OpenID Connect userinfo endpoint.
func (p *Provider) FetchUserInfo(
	ctx context.Context,
	accessToken string,
) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.linkedin.com/v2/userinfo",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.ErrorContext(ctx, "Failed to fetch LinkedIn user info",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		p.logger.ErrorContext(ctx, "LinkedIn user info request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToGetUserInfo, resp.StatusCode)
	}

	var userInfo UserInfo

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}

	return &userInfo, nil
}

// fetchVanityName calls /v2/me to retrieve the member's vanityName.
// Returns empty string if unavailable (non-fatal).
func (p *Provider) fetchVanityName(ctx context.Context, accessToken string) string {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.linkedin.com/v2/me?projection=(vanityName)",
		nil,
	)
	if err != nil {
		return ""
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.DebugContext(ctx, "Failed to fetch /v2/me for vanityName",
			slog.String("error", err.Error()))

		return ""
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		p.logger.DebugContext(ctx, "LinkedIn /v2/me returned non-200",
			slog.Int("status", resp.StatusCode))

		return ""
	}

	body, _ := io.ReadAll(resp.Body)

	var meResp struct {
		VanityName string `json:"vanityName"`
	}

	if err := json.Unmarshal(body, &meResp); err != nil {
		return ""
	}

	return meResp.VanityName
}

// FetchOrganizationPages retrieves LinkedIn organization pages the user administers.
// Requires r_organization_social scope from the Community Management API product.
// Gracefully returns an empty slice if the scope is not available.
func (p *Provider) FetchOrganizationPages(
	ctx context.Context,
	accessToken string,
) ([]*OrgPageInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.linkedin.com/v2/organizationAcls?q=roleAssignee&role=ADMINISTRATOR&projection=(elements*(organization~(id,localizedName,vanityName,logoV2(original~:playableStreams))))",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetOrgs, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.ErrorContext(ctx, "Failed to fetch LinkedIn organizations",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetOrgs, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	// If scope is not available, LinkedIn returns 403 — treat as empty list
	if resp.StatusCode == http.StatusForbidden {
		p.logger.DebugContext(
			ctx,
			"LinkedIn organization scope not available, returning empty list",
		)

		return []*OrgPageInfo{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		p.logger.ErrorContext(ctx, "LinkedIn organizations request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToGetOrgs, resp.StatusCode)
	}

	var orgResp struct {
		Elements []struct {
			Organization struct {
				ID            int64  `json:"id"`
				LocalizedName string `json:"localizedName"`
				VanityName    string `json:"vanityName"`
			} `json:"organization~"`
		} `json:"elements"`
	}

	if err := json.Unmarshal(body, &orgResp); err != nil {
		p.logger.ErrorContext(ctx, "Failed to parse LinkedIn organizations response",
			slog.String("error", err.Error()))

		return []*OrgPageInfo{}, nil
	}

	orgs := make([]*OrgPageInfo, 0, len(orgResp.Elements))

	for _, elem := range orgResp.Elements {
		org := elem.Organization
		orgID := strconv.FormatInt(org.ID, 10)

		uri := "https://www.linkedin.com/company/" + org.VanityName
		if org.VanityName == "" {
			uri = "https://www.linkedin.com/company/" + orgID
		}

		orgs = append(orgs, &OrgPageInfo{
			ID:         orgID,
			Name:       org.LocalizedName,
			VanityName: org.VanityName,
			URI:        uri,
		})
	}

	return orgs, nil
}
