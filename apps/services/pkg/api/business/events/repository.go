package events

import (
	"context"
	"time"
)

// Repository defines the storage operations for the event queue (port).
type Repository interface {
	// Enqueue inserts a new event into the queue.
	Enqueue(
		ctx context.Context,
		id string,
		eventType EventType,
		payload map[string]any,
		maxRetries int,
		visibilityTimeoutSecs int,
		visibleAt time.Time,
	) error

	// ClaimNext atomically claims the next available event for processing.
	// Uses CTE with FOR UPDATE SKIP LOCKED for distributed safety.
	// Returns nil, nil if no events are available.
	ClaimNext(ctx context.Context, workerID string) (*Event, error)

	// Complete marks an event as successfully completed.
	// Validates worker_id to prevent stale workers from interfering.
	Complete(ctx context.Context, id string, workerID string) error

	// Fail marks an event as failed with error message and backoff.
	// If retries exhausted, marks as dead. Otherwise reschedules with backoff.
	// Validates worker_id to prevent stale workers from interfering.
	Fail(
		ctx context.Context,
		id string,
		workerID string,
		errorMessage string,
		backoffSeconds int,
	) error

	// ListByType returns events of a given type (for audit/debugging).
	ListByType(ctx context.Context, eventType EventType, limit int) ([]*Event, error)
}
