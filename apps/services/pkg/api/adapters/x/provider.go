package x

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
)

// Sentinel errors.
var (
	ErrFailedToExchangeCode = errors.New("failed to exchange authorization code")
	ErrFailedToGetUserInfo  = errors.New("failed to get X user info")
)

// HTTPClient interface for dependency injection.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Provider implements OAuth 2.0 with PKCE for X (Twitter) profile linking.
type Provider struct {
	config     *auth.XOAuthConfig
	logger     *logfx.Logger
	httpClient HTTPClient
	pkceStore  *profiles.PKCEStore
}

// NewProvider creates a new X provider.
func NewProvider(
	config *auth.XOAuthConfig,
	logger *logfx.Logger,
	httpClient HTTPClient,
	pkceStore *profiles.PKCEStore,
) *Provider {
	return &Provider{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
		pkceStore:  pkceStore,
	}
}

// InitiateProfileLinkOAuth builds the OAuth URL with PKCE and expanded scope for profile linking.
func (p *Provider) InitiateProfileLinkOAuth(
	ctx context.Context,
	redirectURI string,
	state string,
	stateKey string,
) (string, error) {
	codeChallenge, err := p.pkceStore.GeneratePKCE(stateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate PKCE: %w", err)
	}

	authURL := p.buildAuthURL(redirectURI, state, p.config.ProfileLinkScope, codeChallenge)

	p.logger.DebugContext(ctx, "Initiating X Profile Link OAuth with PKCE",
		slog.String("redirect_uri", redirectURI))

	return authURL, nil
}

// HandleOAuthCallback exchanges the code for tokens using PKCE and returns account info.
func (p *Provider) HandleOAuthCallback(
	ctx context.Context,
	code string,
	redirectURI string,
	stateKey string,
) (auth.OAuthCallbackResult, error) {
	p.logger.DebugContext(ctx, "Processing X OAuth callback")

	codeVerifier, err := p.pkceStore.GetAndDelete(stateKey)
	if err != nil {
		return auth.OAuthCallbackResult{}, fmt.Errorf("PKCE verifier not found: %w", err)
	}

	tokenResp, err := p.exchangeCodeForTokens(ctx, code, redirectURI, codeVerifier)
	if err != nil {
		return auth.OAuthCallbackResult{}, err
	}

	userInfo, err := p.fetchUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return auth.OAuthCallbackResult{}, err
	}

	var expiresAt *time.Time

	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	profileURI := "https://x.com/" + userInfo.Username

	p.logger.DebugContext(ctx, "X OAuth callback successful",
		slog.String("id", userInfo.ID),
		slog.String("username", userInfo.Username))

	return auth.OAuthCallbackResult{
		RemoteID:             userInfo.ID,
		Username:             userInfo.Username,
		Name:                 userInfo.Name,
		Email:                "",
		URI:                  profileURI,
		AccessToken:          tokenResp.AccessToken,
		RefreshToken:         tokenResp.RefreshToken,
		AccessTokenExpiresAt: expiresAt,
		Scope:                tokenResp.Scope,
	}, nil
}

// UserInfo represents X user information from /2/users/me.
type UserInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profile_image_url"`
}

// tokenResponse represents X's token endpoint response.
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

func (p *Provider) buildAuthURL(redirectURI, state, scope, codeChallenge string) string {
	queryString := url.Values{}
	queryString.Set("response_type", "code")
	queryString.Set("client_id", p.config.ClientID)
	queryString.Set("redirect_uri", redirectURI)
	queryString.Set("state", state)
	queryString.Set("scope", scope)
	queryString.Set("code_challenge", codeChallenge)
	queryString.Set("code_challenge_method", "S256")

	oauthURL := url.URL{ //nolint:exhaustruct
		Scheme:   "https",
		Host:     "x.com",
		Path:     "/i/oauth2/authorize",
		RawQuery: queryString.Encode(),
	}

	return oauthURL.String()
}

func (p *Provider) exchangeCodeForTokens(
	ctx context.Context,
	code string,
	redirectURI string,
	codeVerifier string,
) (*tokenResponse, error) {
	values := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"code_verifier": {codeVerifier},
		"client_id":     {p.config.ClientID},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.x.com/2/oauth2/token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(p.config.ClientID, p.config.ClientSecret)

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

func (p *Provider) fetchUserInfo(
	ctx context.Context,
	accessToken string,
) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.x.com/2/users/me?user.fields=username,name,profile_image_url",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.ErrorContext(ctx, "Failed to fetch X user info",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		p.logger.ErrorContext(ctx, "X user info request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToGetUserInfo, resp.StatusCode)
	}

	// X API v2 wraps data in a "data" field
	var apiResp struct {
		Data UserInfo `json:"data"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}

	return &apiResp.Data, nil
}
