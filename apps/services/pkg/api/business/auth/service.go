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

// FirstAuthInfoSetter is an optional interface for providers that accept user info
// from the first authorization (e.g., Apple sends name/email only on first auth).
type FirstAuthInfoSetter interface {
	SetFirstAuthUserInfo(name, email string)
}

type TokenService interface {
	// ParseToken validates a JWT token and returns the claims
	ParseToken(tokenStr string) (*JWTClaims, error)

	// GenerateToken creates a new JWT token with the given claims
	GenerateToken(claims *JWTClaims) (string, error)
}

// CustomDomainChecker checks whether a domain is a registered custom domain.
// Used for redirect URI origin validation during OAuth callback.
type CustomDomainChecker func(ctx context.Context, domain string) bool

type Service struct {
	logger              *logfx.Logger
	tokenService        TokenService
	Config              *Config
	providers           map[string]Provider
	userService         *users.Service
	customDomainChecker CustomDomainChecker
}

func NewService(
	logger *logfx.Logger,
	tokenService TokenService,
	config *Config,
	userService *users.Service,
	customDomainChecker CustomDomainChecker,
) *Service {
	return &Service{
		logger:              logger,
		tokenService:        tokenService,
		Config:              config,
		providers:           map[string]Provider{},
		userService:         userService,
		customDomainChecker: customDomainChecker,
	}
}

// TokenService returns the token service for JWT operations.
func (s *Service) TokenService() TokenService {
	return s.tokenService
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

// IsAllowedOrigin checks whether the given origin is in the CORS allowed origins list.
func (s *Service) IsAllowedOrigin(origin string) bool {
	for _, allowed := range s.Config.GetCorsAllowedOrigins() {
		if strings.EqualFold(allowed, origin) {
			return true
		}
	}

	return false
}

// isAllowedRedirectOrigin checks whether the origin is allowed for redirect URIs.
// It first checks the configured CORS allowed origins, then falls back to
// checking registered custom domains in the database.
func (s *Service) isAllowedRedirectOrigin(ctx context.Context, origin string) bool {
	// Check config-defined origins first (fast path)
	if s.IsAllowedOrigin(origin) {
		return true
	}

	// Fall back to custom domain check
	if s.customDomainChecker != nil {
		parsed, err := url.Parse(origin)
		if err != nil {
			return false
		}

		domain := strings.TrimPrefix(parsed.Hostname(), "www.")
		if domain != "" {
			return s.customDomainChecker(ctx, domain)
		}
	}

	return false
}

// SetProviderFirstAuthInfo sets first-authorization user info on a provider, if it supports it.
// Used for providers like Apple that send user name/email only on first authorization.
func (s *Service) SetProviderFirstAuthInfo(providerName, name, email string) {
	provider := s.GetProvider(providerName)
	if provider == nil {
		return
	}

	if setter, ok := provider.(FirstAuthInfoSetter); ok {
		setter.SetFirstAuthUserInfo(name, email)
	}
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

	// Create/update user (provider-aware)
	user, err := s.upsertUserFromOAuth(ctx, providerName, accountInfo)
	if err != nil {
		return AuthResult{}, err
	}

	// Create session and generate JWT
	session, tokenString, expiresAt, err := s.createSessionAndToken(
		ctx,
		user,
		providerName,
		state,
		&accountInfo,
	)
	if err != nil {
		return AuthResult{}, err
	}

	// Validate redirect URI against allowed origins and custom domains to prevent open redirect
	validateErr := s.validateRedirectURI(ctx, redirectURI)
	if validateErr != nil {
		return AuthResult{}, validateErr
	}

	authResult := AuthResult{
		User:        user,
		Session:     session,
		SessionID:   session.ID,
		JWT:         tokenString,
		ExpiresAt:   expiresAt,
		RedirectURI: redirectURI,
	}

	// Append auth_token query parameter to redirect URI
	authResult.RedirectURI, err = appendAuthToken(authResult.RedirectURI, authResult.JWT)
	if err != nil {
		return authResult, err
	}

	s.logger.DebugContext(ctx, "OAuth callback completed successfully",
		slog.String("user_id", user.ID),
		slog.String("session_id", session.ID),
		slog.String("provider", providerName))

	return authResult, nil
}

// upsertUserFromOAuth creates or updates the user based on the OAuth provider and account info.
func (s *Service) upsertUserFromOAuth(
	ctx context.Context,
	providerName string,
	accountInfo OAuthCallbackResult,
) (*users.User, error) {
	s.logger.DebugContext(ctx, "Creating/updating user from OAuth",
		slog.String("provider", providerName),
		slog.String("remote_id", accountInfo.RemoteID),
		slog.String("username", accountInfo.Username))

	var (
		user *users.User
		err  error
	)

	switch providerName {
	case "github":
		user, err = s.userService.UpsertGitHubUser(
			ctx,
			accountInfo.RemoteID,
			accountInfo.Email,
			accountInfo.Name,
			accountInfo.Username,
			accountInfo.ProfilePictureURI,
		)
	case "apple":
		user, err = s.userService.UpsertAppleUser(
			ctx,
			accountInfo.RemoteID,
			accountInfo.Email,
			accountInfo.Name,
			accountInfo.ProfilePictureURI,
		)
	default:
		return nil, ErrProviderNotFound
	}

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to upsert user",
			slog.String("remote_id", accountInfo.RemoteID),
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToHandleCallback, err)
	}

	return user, nil
}

// createSessionAndToken creates a new session for the user and generates a JWT token.
func (s *Service) createSessionAndToken(
	ctx context.Context,
	user *users.User,
	providerName string,
	state string,
	accountInfo *OAuthCallbackResult,
) (*users.Session, string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.Config.TokenTTL)

	oauthProvider := providerName
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
		OAuthProvider:            &oauthProvider,
		OAuthAccessToken:         &accountInfo.AccessToken,
		OAuthTokenScope:          &accountInfo.Scope,
	}

	s.logger.DebugContext(ctx, "Creating session",
		slog.String("session_id", session.ID),
		slog.String("user_id", user.ID))

	err := s.userService.CreateSession(ctx, &session)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create session",
			slog.String("session_id", session.ID),
			slog.String("error", err.Error()))

		return nil, "", time.Time{}, fmt.Errorf("%w: %w", ErrFailedToHandleCallback, err)
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

		return nil, "", time.Time{}, fmt.Errorf("%w: %w", ErrFailedToGenerateToken, err)
	}

	return &session, tokenString, expiresAt, nil
}

// validateRedirectURI checks that the redirect URI origin is in the allowed list.
func (s *Service) validateRedirectURI(ctx context.Context, redirectURI string) error {
	if redirectURI == "" {
		return nil
	}

	parsedRedirect, parseErr := url.Parse(redirectURI)
	if parseErr != nil {
		return fmt.Errorf("%w: %w", ErrFailedToParseRedirectURI, parseErr)
	}

	redirectOrigin := parsedRedirect.Scheme + "://" + parsedRedirect.Host

	if !s.isAllowedRedirectOrigin(ctx, redirectOrigin) {
		s.logger.WarnContext(ctx, "Blocked redirect to disallowed origin",
			slog.String("redirect_uri", redirectURI),
			slog.String("origin", redirectOrigin))

		return ErrUnsafeRedirectURI
	}

	return nil
}

// appendAuthToken adds the auth_token query parameter to the redirect URI.
func appendAuthToken(redirectURI string, token string) (string, error) {
	if redirectURI == "" {
		return "", nil
	}

	finalRedirectURI, err := url.Parse(redirectURI)
	if err != nil {
		return redirectURI, fmt.Errorf("%w: %w", ErrFailedToParseRedirectURI, err)
	}

	finalRedirectURIQueryString := finalRedirectURI.Query()
	finalRedirectURIQueryString.Set("auth_token", token)

	finalRedirectURI.RawQuery = finalRedirectURIQueryString.Encode()

	return finalRedirectURI.String(), nil
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
		Session:     session,
		SessionID:   claims.SessionID,
		JWT:         tokenString,
		ExpiresAt:   expiresAt,
		RedirectURI: "",
	}, nil
}
