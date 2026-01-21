package profilelink_oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
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

const (
	OAuthStateSizeBytes = 32 // 256 bits of entropy for OAuth state
	StateExpiryDuration = 10 * time.Minute
)

var (
	ErrFailedToGenerateState  = errors.New("failed to generate OAuth state")
	ErrFailedToExchangeCode   = errors.New("failed to exchange authorization code")
	ErrFailedToGetChannelInfo = errors.New("failed to get YouTube channel info")
	ErrInvalidState           = errors.New("invalid or expired OAuth state")
	ErrNoChannelFound         = errors.New("no YouTube channel found for this account")
	ErrAccessDenied           = errors.New("access was denied by user")
)

// HTTPClient interface for dependency injection.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ProfileLinkOAuthState contains state information for OAuth flow.
type ProfileLinkOAuthState struct {
	State       string    `json:"state"`
	ProfileSlug string    `json:"profile_slug"`
	Locale      string    `json:"locale"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// ProfileLinkOAuthResult contains the result of a successful OAuth flow.
type ProfileLinkOAuthResult struct {
	Kind                 string
	RemoteID             string // YouTube channel ID
	PublicID             string // Channel handle (@username)
	Title                string // Channel name
	URI                  string // Channel URL
	AccessToken          string
	RefreshToken         string
	AccessTokenExpiresAt *time.Time
	Scope                string
}

// YouTubeOAuthProvider handles YouTube OAuth for profile links.
type YouTubeOAuthProvider struct {
	config     *auth.YouTubeOAuthConfig
	logger     *logfx.Logger
	httpClient HTTPClient
}

// NewYouTubeOAuthProvider creates a new YouTube OAuth provider.
func NewYouTubeOAuthProvider(
	config *auth.YouTubeOAuthConfig,
	logger *logfx.Logger,
	httpClient HTTPClient,
) *YouTubeOAuthProvider {
	return &YouTubeOAuthProvider{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}
}

// InitiateOAuth starts the OAuth flow and returns the authorization URL.
func (y *YouTubeOAuthProvider) InitiateOAuth(
	ctx context.Context,
	redirectURI string,
	profileSlug string,
	locale string,
) (string, string, error) {
	// Generate cryptographically secure random state
	stateBytes := make([]byte, OAuthStateSizeBytes)

	_, stateErr := rand.Read(stateBytes)
	if stateErr != nil {
		return "", "", fmt.Errorf("%w: %w", ErrFailedToGenerateState, stateErr)
	}

	randomState := hex.EncodeToString(stateBytes)

	// Create state object with profile info
	stateObj := ProfileLinkOAuthState{
		State:       randomState,
		ProfileSlug: profileSlug,
		Locale:      locale,
		ExpiresAt:   time.Now().Add(StateExpiryDuration),
	}

	// Encode state as base64 JSON
	stateJSON, jsonErr := json.Marshal(stateObj)
	if jsonErr != nil {
		return "", "", fmt.Errorf("%w: %w", ErrFailedToGenerateState, jsonErr)
	}

	encodedState := base64.URLEncoding.EncodeToString(stateJSON)

	// Build Google OAuth URL
	queryString := url.Values{}
	queryString.Set("client_id", y.config.ClientID)
	queryString.Set("redirect_uri", redirectURI)
	queryString.Set("response_type", "code")
	queryString.Set("scope", y.config.Scope)
	queryString.Set("state", encodedState)
	queryString.Set("access_type", "offline") // Request refresh token
	queryString.Set("prompt", "consent")      // Force consent to get refresh token

	oauthURL := url.URL{
		Scheme:   "https",
		Host:     "accounts.google.com",
		Path:     "/o/oauth2/v2/auth",
		RawQuery: queryString.Encode(),
	}

	y.logger.DebugContext(ctx, "Initiating YouTube OAuth",
		slog.String("profile_slug", profileSlug),
		slog.String("redirect_uri", redirectURI))

	return oauthURL.String(), encodedState, nil
}

// HandleCallback processes the OAuth callback and returns the result.
func (y *YouTubeOAuthProvider) HandleCallback(
	ctx context.Context,
	code string,
	encodedState string,
	redirectURI string,
) (*ProfileLinkOAuthResult, *ProfileLinkOAuthState, error) {
	// Decode and validate state
	stateJSON, decodeErr := base64.URLEncoding.DecodeString(encodedState)
	if decodeErr != nil {
		return nil, nil, fmt.Errorf("%w: invalid state encoding", ErrInvalidState)
	}

	var stateObj ProfileLinkOAuthState
	jsonErr := json.Unmarshal(stateJSON, &stateObj)
	if jsonErr != nil {
		return nil, nil, fmt.Errorf("%w: invalid state format", ErrInvalidState)
	}

	// Check state expiry
	if time.Now().After(stateObj.ExpiresAt) {
		return nil, nil, fmt.Errorf("%w: state expired", ErrInvalidState)
	}

	y.logger.DebugContext(ctx, "Processing YouTube OAuth callback",
		slog.String("profile_slug", stateObj.ProfileSlug))

	// Exchange code for tokens
	tokenResp, tokenErr := y.exchangeCodeForTokens(ctx, code, redirectURI)
	if tokenErr != nil {
		return nil, nil, tokenErr
	}

	// Fetch channel info
	channelInfo, channelErr := y.fetchChannelInfo(ctx, tokenResp.AccessToken)
	if channelErr != nil {
		return nil, nil, channelErr
	}

	// Calculate token expiry
	var expiresAt *time.Time

	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	result := &ProfileLinkOAuthResult{
		Kind:                 "youtube",
		RemoteID:             channelInfo.ID,
		PublicID:             channelInfo.CustomURL,
		Title:                channelInfo.Title,
		URI:                  "https://youtube.com/" + channelInfo.CustomURL,
		AccessToken:          tokenResp.AccessToken,
		RefreshToken:         tokenResp.RefreshToken,
		AccessTokenExpiresAt: expiresAt,
		Scope:                y.config.Scope,
	}

	y.logger.DebugContext(ctx, "YouTube OAuth callback successful",
		slog.String("channel_id", channelInfo.ID),
		slog.String("channel_title", channelInfo.Title))

	return result, &stateObj, nil
}

// tokenResponse represents Google's token endpoint response.
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

func (y *YouTubeOAuthProvider) exchangeCodeForTokens(
	ctx context.Context,
	code string,
	redirectURI string,
) (*tokenResponse, error) {
	values := url.Values{
		"client_id":     {y.config.ClientID},
		"client_secret": {y.config.ClientSecret},
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

	resp, respErr := y.httpClient.Do(req)
	if respErr != nil {
		y.logger.ErrorContext(ctx, "Failed to exchange code for tokens",
			slog.String("error", respErr.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, respErr)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		y.logger.ErrorContext(ctx, "Token exchange failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToExchangeCode, resp.StatusCode)
	}

	var tokenResp tokenResponse
	jsonErr := json.Unmarshal(body, &tokenResp)
	if jsonErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, jsonErr)
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

func (y *YouTubeOAuthProvider) fetchChannelInfo(
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

	resp, respErr := y.httpClient.Do(req)
	if respErr != nil {
		y.logger.ErrorContext(ctx, "Failed to fetch channel info",
			slog.String("error", respErr.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetChannelInfo, respErr)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		y.logger.ErrorContext(ctx, "Channel info request failed",
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

	jsonErr := json.Unmarshal(body, &channelResp)
	if jsonErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetChannelInfo, jsonErr)
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
