package youtube

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
)

// Sentinel errors.
var (
	ErrFailedToExchangeCode   = errors.New("failed to exchange authorization code")
	ErrFailedToGetChannelInfo = errors.New("failed to get YouTube channel info")
	ErrNoChannelFound         = errors.New("no YouTube channel found")
)

// HTTPClient interface for dependency injection.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Provider implements auth.Provider for YouTube.
// Note: YouTube is typically used for profile linking, not user authentication.
// The service layer handles the difference based on the state parameter.
type Provider struct {
	config     *auth.YouTubeOAuthConfig
	logger     *logfx.Logger
	httpClient HTTPClient
}

// NewProvider creates a new YouTube provider.
func NewProvider(
	config *auth.YouTubeOAuthConfig,
	logger *logfx.Logger,
	httpClient HTTPClient,
) *Provider {
	return &Provider{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}
}

// InitiateOAuth builds the OAuth URL with the given state.
// Implements auth.Provider interface.
func (p *Provider) InitiateOAuth(
	ctx context.Context,
	redirectURI string,
	state string,
) (string, error) {
	// Build Google OAuth URL
	queryString := url.Values{}
	queryString.Set("client_id", p.config.ClientID)
	queryString.Set("redirect_uri", redirectURI)
	queryString.Set("response_type", "code")
	queryString.Set("scope", p.config.Scope)
	queryString.Set("state", state)
	queryString.Set("access_type", "offline") // Request refresh token
	queryString.Set("prompt", "consent")      // Force consent to get refresh token

	oauthURL := url.URL{ //nolint:exhaustruct
		Scheme:   "https",
		Host:     "accounts.google.com",
		Path:     "/o/oauth2/v2/auth",
		RawQuery: queryString.Encode(),
	}

	p.logger.DebugContext(ctx, "Initiating YouTube OAuth",
		slog.String("redirect_uri", redirectURI))

	return oauthURL.String(), nil
}

// HandleOAuthCallback exchanges the code for tokens and returns account info.
// Implements auth.Provider interface.
// State validation is handled by the service layer.
func (p *Provider) HandleOAuthCallback(
	ctx context.Context,
	code string,
	redirectURI string,
) (auth.OAuthCallbackResult, error) {
	p.logger.DebugContext(ctx, "Processing YouTube OAuth callback")

	// Exchange code for tokens
	tokenResp, err := p.exchangeCodeForTokens(ctx, code, redirectURI)
	if err != nil {
		return auth.OAuthCallbackResult{}, err
	}

	// Fetch channel info
	channelInfo, err := p.fetchChannelInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return auth.OAuthCallbackResult{}, err
	}

	// Calculate token expiry
	var expiresAt *time.Time

	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	p.logger.DebugContext(ctx, "YouTube OAuth callback successful",
		slog.String("channel_id", channelInfo.ID),
		slog.String("channel_title", channelInfo.Title))

	return auth.OAuthCallbackResult{
		RemoteID:             channelInfo.ID,
		Username:             channelInfo.CustomURL,
		Name:                 channelInfo.Title,
		Email:                "", // YouTube channels don't have email
		URI:                  "https://youtube.com/" + channelInfo.CustomURL,
		AccessToken:          tokenResp.AccessToken,
		RefreshToken:         tokenResp.RefreshToken,
		AccessTokenExpiresAt: expiresAt,
		Scope:                p.config.Scope,
	}, nil
}

// tokenResponse represents Google's token endpoint response.
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

func (p *Provider) exchangeCodeForTokens(
	ctx context.Context,
	code string,
	redirectURI string,
) (*tokenResponse, error) {
	values := url.Values{
		"client_id":     {p.config.ClientID},
		"client_secret": {p.config.ClientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	}

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://oauth2.googleapis.com/token",
		strings.NewReader(values.Encode()),
	)
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

// channelInfo represents YouTube channel information.
type channelInfo struct {
	ID        string
	Title     string
	CustomURL string
}

func (p *Provider) fetchChannelInfo(
	ctx context.Context,
	accessToken string,
) (*channelInfo, error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://www.googleapis.com/youtube/v3/channels?part=snippet&mine=true",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.ErrorContext(ctx, "Failed to fetch channel info",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetChannelInfo, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		p.logger.ErrorContext(ctx, "Channel info request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToGetChannelInfo, resp.StatusCode)
	}

	var channelResp struct {
		Items []struct {
			ID      string `json:"id"`
			Snippet struct {
				Title     string `json:"title"`
				CustomURL string `json:"customUrl"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &channelResp); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetChannelInfo, err)
	}

	if len(channelResp.Items) == 0 {
		return nil, ErrNoChannelFound
	}

	channel := channelResp.Items[0]

	return &channelInfo{
		ID:        channel.ID,
		Title:     channel.Snippet.Title,
		CustomURL: channel.Snippet.CustomURL,
	}, nil
}
