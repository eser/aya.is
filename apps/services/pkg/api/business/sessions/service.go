package sessions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

var (
	ErrSessionNotFound       = errors.New("session not found")
	ErrFailedToGetSession    = errors.New("failed to get session")
	ErrFailedToCreateSession = errors.New("failed to create session")
	ErrRateLimitExceeded     = errors.New("rate limit exceeded")
)

// Repository defines the interface for session preference operations.
type Repository interface {
	// Preference operations
	GetPreferences(ctx context.Context, sessionID string) (SessionPreferences, error)
	GetPreference(ctx context.Context, sessionID, key string) (*SessionPreference, error)
	SetPreference(ctx context.Context, sessionID, key, value string) error
	DeletePreference(ctx context.Context, sessionID, key string) error

	// Rate limiting operations
	CheckAndIncrementRateLimit(
		ctx context.Context,
		ipHash string,
		limit int,
		windowSeconds int,
	) (bool, error)
}

// Service handles session-related business logic.
type Service struct {
	logger      *logfx.Logger
	config      *Config
	repo        Repository
	userService *users.Service
	idGen       func() string
}

// NewService creates a new session service.
func NewService(
	logger *logfx.Logger,
	config *Config,
	repo Repository,
	userService *users.Service,
	idGen func() string,
) *Service {
	return &Service{
		logger:      logger,
		config:      config,
		repo:        repo,
		userService: userService,
		idGen:       idGen,
	}
}

// GetSessionByID gets a session by ID (delegates to user service).
func (s *Service) GetSessionByID(ctx context.Context, id string) (*users.Session, error) {
	return s.userService.GetSessionByID(ctx, id)
}

// CreateSession creates a new anonymous session.
func (s *Service) CreateSession(ctx context.Context, ipHash string) (*users.Session, error) {
	// Check rate limit
	allowed, err := s.repo.CheckAndIncrementRateLimit(
		ctx,
		ipHash,
		s.config.RateLimit.PerIP,
		3600, // 1 hour window
	)
	if err != nil {
		s.logger.WarnContext(ctx, "Rate limit check failed", "error", err.Error())
		// Continue anyway - don't block on rate limit errors
	} else if !allowed {
		return nil, ErrRateLimitExceeded
	}

	now := time.Now()
	session := &users.Session{
		ID:                       s.idGen(),
		Status:                   "active",
		OauthRequestState:        "", // Not an OAuth session
		OauthRequestCodeVerifier: "",
		OauthRedirectURI:         nil,
		LoggedInUserID:           nil,
		LoggedInAt:               nil,
		ExpiresAt:                nil, // Anonymous sessions don't expire by default
		CreatedAt:                now,
		UpdatedAt:                nil,
	}

	if err := s.userService.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateSession, err)
	}

	return session, nil
}

// GetPreferences gets all preferences for a session.
func (s *Service) GetPreferences(
	ctx context.Context,
	sessionID string,
) (SessionPreferences, error) {
	prefs, err := s.repo.GetPreferences(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetPreference, err)
	}

	return prefs, nil
}

// GetPreference gets a single preference for a session.
func (s *Service) GetPreference(ctx context.Context, sessionID, key string) (string, error) {
	if err := ValidatePreferenceKey(key); err != nil {
		return "", err
	}

	pref, err := s.repo.GetPreference(ctx, sessionID, key)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGetPreference, err)
	}

	if pref == nil {
		return "", ErrPreferenceNotFound
	}

	return pref.Value, nil
}

// SetPreference sets a preference for a session.
func (s *Service) SetPreference(ctx context.Context, sessionID, key, value string) error {
	err := ValidatePreference(key, value)
	if err != nil {
		return err
	}

	err = s.repo.SetPreference(ctx, sessionID, key, value)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSetPreference, err)
	}

	return nil
}

// SetPreferences sets multiple preferences for a session.
func (s *Service) SetPreferences(
	ctx context.Context,
	sessionID string,
	prefs SessionPreferences,
) error {
	for key, value := range prefs {
		err := s.SetPreference(ctx, sessionID, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeletePreference deletes a preference for a session.
func (s *Service) DeletePreference(ctx context.Context, sessionID, key string) error {
	err := ValidatePreferenceKey(key)
	if err != nil {
		return err
	}

	return s.repo.DeletePreference(ctx, sessionID, key)
}
