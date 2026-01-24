package events

import (
	"context"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

const (
	DefaultMaxRetries            = 3
	DefaultVisibilityTimeoutSecs = 300 // 5 minutes
)

// Service provides event queue operations.
type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator IDGenerator
}

// NewService creates a new event queue service.
func NewService(
	logger *logfx.Logger,
	repo Repository,
	idGenerator IDGenerator,
) *Service {
	return &Service{
		logger:      logger,
		repo:        repo,
		idGenerator: idGenerator,
	}
}

// Enqueue adds a new event to the queue for later processing.
func (s *Service) Enqueue(ctx context.Context, params EnqueueParams) (string, error) {
	id := s.idGenerator()

	maxRetries := params.MaxRetries
	if maxRetries == 0 {
		maxRetries = DefaultMaxRetries
	}

	visibilityTimeoutSecs := params.VisibilityTimeoutSecs
	if visibilityTimeoutSecs == 0 {
		visibilityTimeoutSecs = DefaultVisibilityTimeoutSecs
	}

	visibleAt := time.Now()
	if params.ScheduledAt != nil {
		visibleAt = *params.ScheduledAt
	}

	err := s.repo.Enqueue(
		ctx,
		id,
		params.Type,
		params.Payload,
		maxRetries,
		visibilityTimeoutSecs,
		visibleAt,
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToEnqueue, err)
	}

	return id, nil
}

// CalculateBackoff returns the delay in seconds for exponential backoff.
// Formula: baseSeconds * 2^retryCount (e.g., base=4: 8s, 16s, 32s, 64s...).
func CalculateBackoff(retryCount int, baseSeconds int) int {
	if baseSeconds == 0 {
		baseSeconds = 4
	}

	delay := baseSeconds
	for range retryCount {
		delay *= 2
	}

	return delay
}
