package runtime_states

import (
	"context"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

const timeFormat = time.RFC3339Nano

// Service provides runtime state operations.
type Service struct {
	logger *logfx.Logger
	repo   Repository
}

// NewService creates a new runtime state service.
func NewService(
	logger *logfx.Logger,
	repo Repository,
) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

// Get retrieves the value for a given key.
// Returns ErrStateNotFound if the key does not exist.
func (s *Service) Get(ctx context.Context, key string) (string, error) {
	state, err := s.repo.GetState(ctx, key)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGet, err)
	}

	if state == nil {
		return "", ErrStateNotFound
	}

	return state.Value, nil
}

// GetTime retrieves a time value for a given key.
// Returns ErrStateNotFound if the key does not exist.
func (s *Service) GetTime(ctx context.Context, key string) (time.Time, error) {
	value, err := s.Get(ctx, key)
	if err != nil {
		return time.Time{}, err
	}

	t, parseErr := time.Parse(timeFormat, value)
	if parseErr != nil {
		return time.Time{}, fmt.Errorf("%w: %w", ErrInvalidTime, parseErr)
	}

	return t, nil
}

// Set upserts a value for a given key.
func (s *Service) Set(ctx context.Context, key string, value string) error {
	err := s.repo.SetState(ctx, key, value)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSet, err)
	}

	return nil
}

// SetTime upserts a time value for a given key.
func (s *Service) SetTime(ctx context.Context, key string, t time.Time) error {
	return s.Set(ctx, key, t.Format(timeFormat))
}

// Remove removes a key from the runtime state store.
func (s *Service) Remove(ctx context.Context, key string) error {
	err := s.repo.RemoveState(ctx, key)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRemove, err)
	}

	return nil
}

// ListByPrefix returns all runtime state entries matching a key prefix.
func (s *Service) ListByPrefix(ctx context.Context, prefix string) ([]*RuntimeState, error) {
	states, err := s.repo.ListStatesByPrefix(ctx, prefix)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGet, err)
	}

	return states, nil
}

// TryLock attempts to acquire a session-scoped advisory lock (non-blocking).
// Returns true if the lock was acquired, false if held by another session.
func (s *Service) TryLock(ctx context.Context, lockID int64) (bool, error) {
	acquired, err := s.repo.TryAdvisoryLock(ctx, lockID)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrFailedToAcquireLock, err)
	}

	return acquired, nil
}

// ReleaseLock releases a previously acquired advisory lock.
func (s *Service) ReleaseLock(ctx context.Context, lockID int64) error {
	err := s.repo.ReleaseAdvisoryLock(ctx, lockID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToReleaseLock, err)
	}

	return nil
}
