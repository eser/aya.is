package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/queue"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

// Enqueue inserts a new item into the queue.
func (r *Repository) Enqueue(
	ctx context.Context,
	id string,
	itemType queue.ItemType,
	payload map[string]any,
	maxRetries int,
	visibilityTimeoutSecs int,
	visibleAt time.Time,
) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.queries.EnqueueQueueItem(ctx, EnqueueQueueItemParams{
		ID:                    id,
		Type:                  string(itemType),
		Payload:               payloadJSON,
		MaxRetries:            int32(maxRetries),
		VisibilityTimeoutSecs: int32(visibilityTimeoutSecs),
		VisibleAt:             visibleAt,
	})
}

// ClaimNext atomically claims the next available item for processing.
func (r *Repository) ClaimNext(ctx context.Context, workerID string) (*queue.Item, error) {
	row, err := r.queries.ClaimNextQueueItem(ctx, ClaimNextQueueItemParams{
		WorkerID: sql.NullString{String: workerID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return r.rowToQueueItem(row), nil
}

// Complete marks an item as successfully completed.
func (r *Repository) Complete(ctx context.Context, id string, workerID string) error {
	_, err := r.queries.CompleteQueueItem(ctx, CompleteQueueItemParams{
		ID:       id,
		WorkerID: sql.NullString{String: workerID, Valid: true},
	})

	return err
}

// Fail marks an item as failed with error message and backoff.
func (r *Repository) Fail(
	ctx context.Context,
	id string,
	workerID string,
	errorMessage string,
	backoffSeconds int,
) error {
	_, err := r.queries.FailQueueItem(ctx, FailQueueItemParams{
		ID:             id,
		WorkerID:       sql.NullString{String: workerID, Valid: true},
		ErrorMessage:   sql.NullString{String: errorMessage, Valid: true},
		BackoffSeconds: int32(backoffSeconds),
	})

	return err
}

// ListByType returns items of a given type for audit/debugging.
func (r *Repository) ListByType(
	ctx context.Context,
	itemType queue.ItemType,
	limit int,
) ([]*queue.Item, error) {
	rows, err := r.queries.ListQueueItemsByType(ctx, ListQueueItemsByTypeParams{
		Type:       string(itemType),
		LimitCount: int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*queue.Item, len(rows))
	for i, row := range rows {
		result[i] = r.rowToQueueItem(row)
	}

	return result, nil
}

// rowToQueueItem converts a database row to an Item domain object.
func (r *Repository) rowToQueueItem(row *Queue) *queue.Item {
	var payload map[string]any
	if len(row.Payload) > 0 {
		_ = json.Unmarshal(row.Payload, &payload)
	}

	return &queue.Item{
		ID:                    row.ID,
		Type:                  queue.ItemType(row.Type),
		Payload:               payload,
		Status:                queue.ItemStatus(row.Status),
		RetryCount:            int(row.RetryCount),
		MaxRetries:            int(row.MaxRetries),
		VisibleAt:             row.VisibleAt,
		VisibilityTimeoutSecs: int(row.VisibilityTimeoutSecs),
		StartedAt:             vars.ToTimePtr(row.StartedAt),
		CompletedAt:           vars.ToTimePtr(row.CompletedAt),
		FailedAt:              vars.ToTimePtr(row.FailedAt),
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             vars.ToTimePtr(row.UpdatedAt),
		ErrorMessage:          vars.ToStringPtr(row.ErrorMessage),
		WorkerID:              vars.ToStringPtr(row.WorkerID),
	}
}
