package events

import (
	"context"
	"time"
)

// QueueRepository defines the storage operations for the event queue (port).
type QueueRepository interface {
	// Enqueue inserts a new item into the event queue.
	Enqueue(
		ctx context.Context,
		id string,
		itemType QueueItemType,
		payload map[string]any,
		maxRetries int,
		visibilityTimeoutSecs int,
		visibleAt time.Time,
	) error

	// ClaimNext atomically claims the next available item for processing.
	// Uses CTE with FOR UPDATE SKIP LOCKED for distributed safety.
	// Returns nil, nil if no items are available.
	ClaimNext(ctx context.Context, workerID string) (*QueueItem, error)

	// Complete marks an item as successfully completed.
	// Validates worker_id to prevent stale workers from interfering.
	Complete(ctx context.Context, id string, workerID string) error

	// Fail marks an item as failed with error message and backoff.
	// If retries exhausted, marks as dead. Otherwise reschedules with backoff.
	// Validates worker_id to prevent stale workers from interfering.
	Fail(
		ctx context.Context,
		id string,
		workerID string,
		errorMessage string,
		backoffSeconds int,
	) error

	// ListByType returns items of a given type (for audit/debugging).
	ListByType(ctx context.Context, itemType QueueItemType, limit int) ([]*QueueItem, error)
}
