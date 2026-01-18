package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
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
	ErrInvalidToken             = errors.New("invalid token")
	ErrJWTNotConfigured         = errors.New("JWT not configured")
	ErrInvalidSigningMethod     = errors.New("invalid JWT signing method")
	ErrFailedToGenerateToken    = errors.New("failed to generate token")
	ErrSessionExpired           = errors.New("session expired")
	ErrFailedToGetUser          = errors.New("failed to get user")
	ErrFailedToUpdateSession    = errors.New("failed to update session")
)

type Provider interface {
	// InitiateOAuth returns the URL to redirect the user to GitHub for login, and the state to track the request.
	InitiateOAuth(
		ctx context.Context,
		redirectURI string,
	) (authURL string, state OAuthState, err error)

	// HandleOAuthCallback exchanges the code for a token, fetches user info, upserts user,
	// creates session, and returns JWT.
	HandleOAuthCallback(
		ctx context.Context,
		code string,
		state string,
		redirectURI string,
	) (AuthResult, error)
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

	callbackURI, err := url.Parse(baseURI)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToParseBaseURI, err)
	}

	callbackURIQueryString := callbackURI.Query()
	callbackURIQueryString.Set("redirect_uri", finalRedirectURI)

	callbackURI.Path += "/auth/" + providerName + "/callback"
	callbackURI.RawQuery = callbackURIQueryString.Encode()

	authURL, _, err := provider.InitiateOAuth(ctx, callbackURI.String())
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

	authResult, err := provider.HandleOAuthCallback(ctx, code, state, redirectURI)
	if err != nil {
		return AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToHandleCallback, err)
	}

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

	return authResult, nil
}

// GenerateSessionToken creates a new JWT token for a given session and user.
// Used for cookie-based session check.
func (s *Service) GenerateSessionToken(sessionID, userID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.Config.TokenTTL)

	claims := &JWTClaims{
		UserID:    userID,
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
		slog.String("user_id", claims.UserID),
		slog.String("session_id", claims.SessionID))

	// Verify session is still active
	session, err := s.userService.GetSessionByID(ctx, claims.SessionID)
	if err != nil || session.Status != users.SessionStatusActive {
		s.logger.WarnContext(ctx, "Session is not active",
			slog.String("session_id", claims.SessionID),
			slog.String("status", string(session.Status)))

		return nil, ErrSessionExpired
	}

	// Get user details
	user, err := s.userService.GetByID(ctx, claims.UserID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user for token refresh",
			slog.String("user_id", claims.UserID),
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUser, err)
	}

	// Generate new JWT with extended expiration
	now := time.Now()
	expiresAt := now.Add(s.Config.TokenTTL)

	newClaims := &JWTClaims{
		UserID:    claims.UserID,
		SessionID: claims.SessionID,
		ExpiresAt: expiresAt.Unix(),
	}

	tokenString, err := s.tokenService.GenerateToken(newClaims)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate new JWT token",
			slog.String("user_id", claims.UserID),
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
		slog.String("user_id", claims.UserID),
		slog.String("session_id", claims.SessionID))

	return &AuthResult{
		User:        user,
		SessionID:   claims.SessionID,
		JWT:         tokenString,
		ExpiresAt:   expiresAt,
		RedirectURI: "",
	}, nil
}
