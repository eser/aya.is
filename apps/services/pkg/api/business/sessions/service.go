package sessions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

var (
	ErrSessionNotFound         = errors.New("session not found")
	ErrFailedToGetSession      = errors.New("failed to get session")
	ErrFailedToCreateSession   = errors.New("failed to create session")
	ErrFailedToInvalidate      = errors.New("failed to invalidate session")
	ErrFailedToCopyPreferences = errors.New("failed to copy preferences")
	ErrRateLimitExceeded       = errors.New("rate limit exceeded")
)

// Repository defines the interface for session preference operations.
type Repository interface {
	// Preference operations
	GetPreferences(ctx context.Context, sessionID string) (SessionPreferences, error)
	GetPreference(ctx context.Context, sessionID, key string) (*SessionPreference, error)
	SetPreference(ctx context.Context, sessionID, key, value string) error
	DeletePreference(ctx context.Context, sessionID, key string) error
	CopyPreferences(ctx context.Context, oldSessionID, newSessionID string) error

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
	logger       *logfx.Logger
	config       *Config
	repo         Repository
	userService  *users.Service
	auditService *events.AuditService
	idGen        func() string
}

// NewService creates a new session service.
func NewService(
	logger *logfx.Logger,
	config *Config,
	repo Repository,
	userService *users.Service,
	idGen func() string,
	auditService *events.AuditService,
) *Service {
	return &Service{
		logger:       logger,
		config:       config,
		repo:         repo,
		userService:  userService,
		auditService: auditService,
		idGen:        idGen,
	}
}

// GetSessionByID gets a session by ID (delegates to user service).
func (s *Service) GetSessionByID(ctx context.Context, id string) (*users.Session, error) {
	session, err := s.userService.GetSessionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetSession, err)
	}

	return session, nil
}

// CreateSession creates a new anonymous session.
func (s *Service) CreateSession(ctx context.Context, ipHash string) (*users.Session, error) {
	// Check rate limit
	const rateLimitWindowSeconds = 3600 // 1 hour window

	allowed, err := s.repo.CheckAndIncrementRateLimit(
		ctx,
		ipHash,
		s.config.RateLimit.PerIP,
		rateLimitWindowSeconds,
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
		Status:                   users.SessionStatusActive,
		OauthRequestState:        "", // Not an OAuth session
		OauthRequestCodeVerifier: "",
		OauthRedirectURI:         nil,
		LoggedInUserID:           nil,
		LoggedInAt:               nil,
		LastActivityAt:           nil,
		ExpiresAt:                nil, // Anonymous sessions don't expire by default
		UpdatedAt:                nil,
		UserAgent:                nil,
		OAuthProvider:            nil,
		OAuthAccessToken:         nil,
		OAuthTokenScope:          nil,
		CreatedAt:                now,
	}

	createErr := s.userService.CreateSession(ctx, session)
	if createErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateSession, createErr)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.SessionCreated,
		EntityType: "session",
		EntityID:   session.ID,
		ActorID:    nil,
		ActorKind:  events.ActorSystem,
		SessionID:  nil,
		Payload:    nil,
	})

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
	validationErr := ValidatePreferenceKey(key)
	if validationErr != nil {
		return "", validationErr
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

	deleteErr := s.repo.DeletePreference(ctx, sessionID, key)
	if deleteErr != nil {
		return fmt.Errorf("deleting preference: %w", deleteErr)
	}

	return nil
}

// LogoutResult contains the result of a logout operation.
type LogoutResult struct {
	NewSession *users.Session
}

// LogoutSession invalidates the current session and creates a new anonymous session
// with the same preferences. This ensures the user is logged out across all domains.
func (s *Service) LogoutSession(ctx context.Context, oldSessionID string) (*LogoutResult, error) {
	// Create a new anonymous session (no rate limiting for logout)
	now := time.Now()
	newSession := &users.Session{
		ID:                       s.idGen(),
		Status:                   users.SessionStatusActive,
		OauthRequestState:        "",
		OauthRequestCodeVerifier: "",
		OauthRedirectURI:         nil,
		LoggedInUserID:           nil, // Anonymous
		LoggedInAt:               nil,
		LastActivityAt:           nil,
		ExpiresAt:                nil,
		UpdatedAt:                nil,
		UserAgent:                nil,
		OAuthProvider:            nil,
		OAuthAccessToken:         nil,
		OAuthTokenScope:          nil,
		CreatedAt:                now,
	}

	err := s.userService.CreateSession(ctx, newSession)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateSession, err)
	}

	// Copy preferences from old session to new session
	err = s.repo.CopyPreferences(ctx, oldSessionID, newSession.ID)
	if err != nil {
		// Log but don't fail - the logout is more important than preserving preferences
		s.logger.WarnContext(ctx, "Failed to copy preferences during logout",
			"error", err.Error(),
			"old_session_id", oldSessionID,
			"new_session_id", newSession.ID)
	}

	// Invalidate the old session
	err = s.userService.InvalidateSession(ctx, oldSessionID)
	if err != nil {
		// Log but don't fail - the new session is already created
		s.logger.WarnContext(ctx, "Failed to invalidate old session during logout",
			"error", err.Error(),
			"old_session_id", oldSessionID)
	}

	// Resolve the user who initiated the logout
	oldSession, _ := s.userService.GetSessionByID(ctx, oldSessionID)

	var actorID *string
	if oldSession != nil {
		actorID = oldSession.LoggedInUserID
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.SessionTerminated,
		EntityType: "session",
		EntityID:   oldSessionID,
		ActorID:    actorID,
		ActorKind:  events.ActorUser,
		SessionID:  &oldSessionID,
		Payload:    nil,
	})

	return &LogoutResult{NewSession: newSession}, nil
}
