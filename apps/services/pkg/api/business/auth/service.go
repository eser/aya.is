package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

var (
	ErrProviderNotFound         = errors.New("provider not found")
	ErrFailedToInitiate         = errors.New("failed to initiate")
	ErrFailedToHandleCallback   = errors.New("failed to handle callback")
	ErrFailedToParseBaseURI     = errors.New("failed to parse base URI")
	ErrFailedToParseRedirectURI = errors.New("failed to parse redirect URI")
	ErrUnsafeRedirectURI        = errors.New("redirect URI not in allowed origins")
	ErrInvalidToken             = errors.New("invalid token")
	ErrJWTNotConfigured         = errors.New("JWT not configured")
	ErrInvalidSigningMethod     = errors.New("invalid JWT signing method")
	ErrFailedToGenerateToken    = errors.New("failed to generate token")
	ErrSessionExpired           = errors.New("session expired")
	ErrFailedToGetUser          = errors.New("failed to get user")
	ErrFailedToUpdateSession    = errors.New("failed to update session")
)

type Provider interface {
	// InitiateOAuth builds the OAuth URL with the given state.
	// The caller is responsible for generating the state (e.g., auth.GenerateRandomState for login).
	InitiateOAuth(
		ctx context.Context,
		redirectURI string,
		state string,
	) (authURL string, err error)

	// HandleOAuthCallback exchanges the code for tokens and returns account info.
	// State validation and user/session creation is handled by the service layer.
	HandleOAuthCallback(
		ctx context.Context,
		code string,
		redirectURI string,
	) (OAuthCallbackResult, error)
}

type TokenService interface {
	// ParseToken validates a JWT token and returns the claims
	ParseToken(tokenStr string) (*JWTClaims, error)

	// GenerateToken creates a new JWT token with the given claims
	GenerateToken(claims *JWTClaims) (string, error)
}

type Service struct {
	logger       *logfx.Logger
	tokenService TokenService
	Config       *Config
	providers    map[string]Provider
	userService  *users.Service
}

func NewService(
	logger *logfx.Logger,
	tokenService TokenService,
	config *Config,
	userService *users.Service,
) *Service {
	return &Service{
		logger:       logger,
		tokenService: tokenService,
		Config:       config,
		providers:    map[string]Provider{},
		userService:  userService,
	}
}

func (s *Service) GetProvider(provider string) Provider {
	service, serviceOk := s.providers[provider]
	if !serviceOk {
		return nil
	}

	return service
}

func (s *Service) RegisterProvider(providerName string, provider Provider) {
	s.providers[providerName] = provider
}

// isAllowedOrigin checks whether the given origin is in the CORS allowed origins list.
func (s *Service) isAllowedOrigin(origin string) bool {
	for _, allowed := range s.Config.GetCorsAllowedOrigins() {
		if strings.EqualFold(allowed, origin) {
			return true
		}
	}

	return false
}

func (s *Service) Initiate(
	ctx context.Context,
	providerName string,
	baseURI string,
	finalRedirectURI string,
) (string, error) {
	provider := s.GetProvider(providerName)

	if provider == nil {
		return "", ErrProviderNotFound
	}

	// Generate state for login flow
	state, err := GenerateRandomState()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToInitiate, err)
	}

	callbackURI, err := url.Parse(baseURI)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToParseBaseURI, err)
	}

	callbackURIQueryString := callbackURI.Query()
	callbackURIQueryString.Set("redirect_uri", finalRedirectURI)

	callbackURI.Path += "/auth/" + providerName + "/callback"
	callbackURI.RawQuery = callbackURIQueryString.Encode()

	authURL, err := provider.InitiateOAuth(ctx, callbackURI.String(), state)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToInitiate, err)
	}

	return authURL, nil
}

func (s *Service) AuthHandleCallback(
	ctx context.Context,
	providerName string,
	code string,
	state string,
	redirectURI string,
) (AuthResult, error) {
	provider := s.GetProvider(providerName)

	if provider == nil {
		return AuthResult{}, ErrProviderNotFound
	}

	// Get account info from provider (state validation is service layer responsibility)
	// Pass empty redirect_uri to avoid a mismatch in GitHub's token exchange.
	// The redirect_uri sent to GitHub during authorization includes query params
	// (e.g., ?redirect_uri=<frontend_url>), but the redirectURI here is the frontend URL.
	// Passing a mismatched value causes GitHub to reject the token exchange.
	accountInfo, err := provider.HandleOAuthCallback(ctx, code, "")
	if err != nil {
		return AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToHandleCallback, err)
	}

	// Create/update user
	s.logger.DebugContext(ctx, "Creating/updating user from OAuth",
		slog.String("provider", providerName),
		slog.String("remote_id", accountInfo.RemoteID),
		slog.String("username", accountInfo.Username))

	user, err := s.userService.UpsertGitHubUser(
		ctx,
		accountInfo.RemoteID,
		accountInfo.Email,
		accountInfo.Name,
		accountInfo.Username,
	)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to upsert user",
			slog.String("remote_id", accountInfo.RemoteID),
			slog.String("error", err.Error()))

		return AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToHandleCallback, err)
	}

	// Create session
	now := time.Now()
	expiresAt := now.Add(s.Config.TokenTTL)

	session := users.Session{
		ID:                       string(s.userService.GenerateID()),
		Status:                   users.SessionStatusActive,
		OauthRequestState:        state,
		OauthRequestCodeVerifier: "",
		OauthRedirectURI:         nil,
		LoggedInUserID:           &user.ID,
		LoggedInAt:               &now,
		LastActivityAt:           &now,
		UserAgent:                nil,
		ExpiresAt:                &expiresAt,
		CreatedAt:                now,
		UpdatedAt:                nil,
	}

	s.logger.DebugContext(ctx, "Creating session",
		slog.String("session_id", session.ID),
		slog.String("user_id", user.ID))

	err = s.userService.CreateSession(ctx, &session)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create session",
			slog.String("session_id", session.ID),
			slog.String("error", err.Error()))

		return AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToHandleCallback, err)
	}

	// Generate JWT — only session_id is stored; user is derived from session
	claims := &JWTClaims{
		SessionID: session.ID,
		ExpiresAt: expiresAt.Unix(),
	}

	s.logger.DebugContext(ctx, "Generating JWT token",
		slog.String("session_id", session.ID))

	tokenString, err := s.tokenService.GenerateToken(claims)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate JWT token",
			slog.String("user_id", user.ID),
			slog.String("error", err.Error()))

		return AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToGenerateToken, err)
	}

	// Validate redirect URI against allowed CORS origins to prevent open redirect
	validatedRedirectURI := redirectURI
	if validatedRedirectURI != "" {
		parsedRedirect, parseErr := url.Parse(validatedRedirectURI)
		if parseErr != nil {
			return AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToParseRedirectURI, parseErr)
		}

		redirectOrigin := parsedRedirect.Scheme + "://" + parsedRedirect.Host
		if !s.isAllowedOrigin(redirectOrigin) {
			s.logger.WarnContext(ctx, "Blocked redirect to disallowed origin",
				slog.String("redirect_uri", validatedRedirectURI),
				slog.String("origin", redirectOrigin))

			return AuthResult{}, ErrUnsafeRedirectURI
		}
	}

	authResult := AuthResult{
		User:        user,
		SessionID:   session.ID,
		JWT:         tokenString,
		ExpiresAt:   expiresAt,
		RedirectURI: validatedRedirectURI,
	}

	// Add auth_token to redirect URI
	if authResult.RedirectURI != "" {
		finalRedirectURI, err := url.Parse(authResult.RedirectURI)
		if err != nil {
			return authResult, fmt.Errorf("%w: %w", ErrFailedToParseRedirectURI, err)
		}

		finalRedirectURIQueryString := finalRedirectURI.Query()
		finalRedirectURIQueryString.Set("auth_token", authResult.JWT)

		finalRedirectURI.RawQuery = finalRedirectURIQueryString.Encode()

		authResult.RedirectURI = finalRedirectURI.String()
	}

	s.logger.DebugContext(ctx, "OAuth callback completed successfully",
		slog.String("user_id", user.ID),
		slog.String("session_id", session.ID),
		slog.String("provider", providerName))

	return authResult, nil
}

// GenerateSessionToken creates a new JWT token for a given session.
// Used for cookie-based session check.
func (s *Service) GenerateSessionToken(sessionID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.Config.TokenTTL)

	claims := &JWTClaims{
		SessionID: sessionID,
		ExpiresAt: expiresAt.Unix(),
	}

	tokenString, err := s.tokenService.GenerateToken(claims)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: %w", ErrFailedToGenerateToken, err)
	}

	return tokenString, expiresAt, nil
}

// RefreshToken validates the current JWT token and issues a new one with extended expiration.
func (s *Service) RefreshToken( //nolint:funlen
	ctx context.Context,
	tokenStr string,
) (*AuthResult, error) {
	s.logger.DebugContext(ctx, "Attempting to refresh JWT token")

	// Parse and validate current token using the token service
	claims, err := s.tokenService.ParseToken(tokenStr)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to parse JWT token", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	s.logger.DebugContext(ctx, "JWT token parsed successfully",
		slog.String("session_id", claims.SessionID))

	// Verify session is still active
	session, err := s.userService.GetSessionByID(ctx, claims.SessionID)
	if err != nil || session.Status != users.SessionStatusActive {
		s.logger.WarnContext(ctx, "Session is not active",
			slog.String("session_id", claims.SessionID),
			slog.String("status", string(session.Status)))

		return nil, ErrSessionExpired
	}

	if session.LoggedInUserID == nil {
		return nil, ErrSessionExpired
	}

	// Get user details from session
	user, err := s.userService.GetByID(ctx, *session.LoggedInUserID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user for token refresh",
			slog.String("user_id", *session.LoggedInUserID),
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUser, err)
	}

	// Generate new JWT with extended expiration
	now := time.Now()
	expiresAt := now.Add(s.Config.TokenTTL)

	newClaims := &JWTClaims{
		SessionID: claims.SessionID,
		ExpiresAt: expiresAt.Unix(),
	}

	tokenString, err := s.tokenService.GenerateToken(newClaims)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate new JWT token",
			slog.String("session_id", claims.SessionID),
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGenerateToken, err)
	}

	// Update session's last activity
	updateErr := s.userService.UpdateSessionLoggedInAt(ctx, claims.SessionID, now)
	if updateErr != nil {
		s.logger.WarnContext(ctx, "Failed to update session logged in time",
			slog.String("session_id", claims.SessionID),
			slog.String("error", updateErr.Error()))
		// Don't fail the whole operation for this
	}

	s.logger.DebugContext(ctx, "JWT token refreshed successfully",
		slog.String("session_id", claims.SessionID))

	return &AuthResult{
		User:        user,
		SessionID:   claims.SessionID,
		JWT:         tokenString,
		ExpiresAt:   expiresAt,
		RedirectURI: "",
	}, nil
}
