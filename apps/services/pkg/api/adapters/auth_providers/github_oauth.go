package auth_providers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

const (
	ExpirePeriod        = 24 * time.Hour
	OAuthStateSizeBytes = 32 // 256 bits of entropy for OAuth state
)

var (
	ErrFailedToGetAccessToken = errors.New("failed to get access token")
	ErrFailedToCreateUser     = errors.New("failed to create user")
	ErrFailedToCreateSession  = errors.New("failed to create session")
	ErrFailedToGetUserInfo    = errors.New("failed to get user info")
	ErrFailedToGenerateOAuth  = errors.New("failed to generate OAuth state")
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GitHubAuthProvider struct {
	config       *auth.GitHubAuthProviderConfig
	logger       *logfx.Logger
	httpClient   HTTPClient
	tokenService auth.TokenService
	userService  *users.Service
}

func NewGitHubAuthProvider(
	config *auth.GitHubAuthProviderConfig,
	logger *logfx.Logger,
	httpClient HTTPClient,
	tokenService auth.TokenService,
	userService *users.Service,
) *GitHubAuthProvider {
	return &GitHubAuthProvider{
		config:       config,
		logger:       logger,
		httpClient:   httpClient,
		tokenService: tokenService,
		userService:  userService,
	}
}

func (g *GitHubAuthProvider) InitiateOAuth(
	ctx context.Context,
	redirectURI string,
) (string, auth.OAuthState, error) {
	// Generate cryptographically secure random state
	stateBytes := make([]byte, OAuthStateSizeBytes)

	_, stateErr := rand.Read(stateBytes)
	if stateErr != nil {
		return "", auth.OAuthState{}, fmt.Errorf("%w: %w", ErrFailedToGenerateOAuth, stateErr)
	}

	state := hex.EncodeToString(stateBytes)

	queryString := url.Values{}
	queryString.Set("client_id", g.config.ClientID)
	queryString.Set("redirect_uri", redirectURI)
	queryString.Set("state", state)
	queryString.Set("scope", g.config.Scope)

	oauthAuthorizeURL := url.URL{ //nolint:exhaustruct
		Scheme:   "https",
		Host:     "github.com",
		Path:     "/login/oauth/authorize",
		RawQuery: queryString.Encode(),
	}

	oauthState := auth.OAuthState{
		State:       state,
		RedirectURI: redirectURI,
	}

	return oauthAuthorizeURL.String(), oauthState, nil
}

func (g *GitHubAuthProvider) HandleOAuthCallback( //nolint:funlen
	ctx context.Context,
	code string,
	state string,
	redirectURI string,
) (_ auth.AuthResult, err error) {
	g.logger.DebugContext(ctx, "Starting GitHub OAuth callback",
		slog.String("state", state))

	// 1. Exchange code for access token
	values := url.Values{
		"client_id":     {g.config.ClientID},
		"client_secret": {g.config.ClientSecret},
		"code":          {code},
	}
	tokenReq, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		strings.NewReader(values.Encode()),
	)
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResp, tokenRespErr := g.httpClient.Do(tokenReq)
	if tokenRespErr != nil {
		g.logger.ErrorContext(ctx, "Failed to exchange code for access token",
			slog.String("error", tokenRespErr.Error()))

		return auth.AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToGetAccessToken, tokenRespErr)
	}

	defer func() {
		closeErr := tokenResp.Body.Close()
		if closeErr != nil {
			g.logger.WarnContext(ctx, "Failed to close response body",
				slog.String("error", closeErr.Error()))
		}
	}()

	body, _ := io.ReadAll(tokenResp.Body)
	vals, _ := url.ParseQuery(string(body))

	accessToken := vals.Get("access_token")
	if accessToken == "" {
		g.logger.ErrorContext(ctx, "No access token received from GitHub")

		return auth.AuthResult{}, ErrFailedToGetAccessToken
	}

	g.logger.DebugContext(ctx, "Successfully obtained GitHub access token")

	// 2. Fetch user info from GitHub
	userReq, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	userReq.Header.Set("Authorization", "Bearer "+accessToken)

	userResp, userRespErr := g.httpClient.Do(userReq)
	if userRespErr != nil {
		return auth.AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, userRespErr)
	}

	defer func() {
		err = userResp.Body.Close()
	}()

	var ghUser struct {
		Login  string `json:"login"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar_url"`
		ID     int64  `json:"id"`
	}

	ghUserErr := json.NewDecoder(userResp.Body).Decode(&ghUser)
	if ghUserErr != nil {
		return auth.AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, ghUserErr)
	}

	// 3. Upsert user in DB
	ghRemoteID := strconv.FormatInt(ghUser.ID, 10)
	now := time.Now()
	expiresAt := now.Add(ExpirePeriod)

	g.logger.DebugContext(ctx, "Creating/updating user in database",
		slog.String("github_remote_id", ghRemoteID),
		slog.String("github_handle", ghUser.Login))

	user, upsertErr := g.userService.UpsertGitHubUser(
		ctx,
		ghRemoteID,
		ghUser.Email,
		ghUser.Name,
		ghUser.Login,
	)
	if upsertErr != nil {
		g.logger.ErrorContext(ctx, "Failed to upsert user",
			slog.String("github_remote_id", ghRemoteID),
			slog.String("error", upsertErr.Error()))

		return auth.AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToCreateUser, upsertErr)
	}

	// 4. Create session in DB
	session := users.Session{
		ID:                       string(g.userService.GenerateID()),
		Status:                   users.SessionStatusActive,
		OauthRequestState:        state,
		OauthRequestCodeVerifier: "", // PKCE not used here
		OauthRedirectURI:         nil,
		LoggedInUserID:           &user.ID,
		LoggedInAt:               &now,
		LastActivityAt:           &now,
		UserAgent:                nil, // Will be set on first authenticated request
		ExpiresAt:                &expiresAt,
		CreatedAt:                now,
		UpdatedAt:                nil,
	}

	g.logger.DebugContext(ctx, "Creating session",
		slog.String("session_id", session.ID),
		slog.String("user_id", user.ID))

	createSessionErr := g.userService.CreateSession(ctx, &session)
	if createSessionErr != nil {
		g.logger.ErrorContext(ctx, "Failed to create session",
			slog.String("session_id", session.ID),
			slog.String("error", createSessionErr.Error()))

		return auth.AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToCreateSession, createSessionErr)
	}

	// 5. Issue JWT
	claims := &auth.JWTClaims{
		UserID:    user.ID,
		SessionID: session.ID,
		ExpiresAt: expiresAt.Unix(),
	}

	g.logger.DebugContext(ctx, "Generating JWT token",
		slog.String("user_id", user.ID),
		slog.String("session_id", session.ID))

	tokenString, tokenStringErr := g.tokenService.GenerateToken(claims)
	if tokenStringErr != nil {
		g.logger.ErrorContext(ctx, "Failed to generate JWT token",
			slog.String("user_id", user.ID),
			slog.String("error", tokenStringErr.Error()))

		return auth.AuthResult{}, fmt.Errorf(
			"%w: %w",
			auth.ErrFailedToGenerateToken,
			tokenStringErr,
		)
	}

	g.logger.DebugContext(ctx, "GitHub OAuth callback completed successfully",
		slog.String("user_id", user.ID),
		slog.String("session_id", session.ID),
		slog.String("github_handle", ghUser.Login))

	return auth.AuthResult{
		User:        user,
		SessionID:   session.ID,
		JWT:         tokenString,
		ExpiresAt:   expiresAt,
		RedirectURI: redirectURI,
	}, nil
}
