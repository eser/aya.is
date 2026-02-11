package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/lib/vars"
	"github.com/sqlc-dev/pqtype"
)

// Enqueue inserts a new item into the event queue.
func (r *Repository) Enqueue(
	ctx context.Context,
	id string,
	itemType events.QueueItemType,
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
		Payload:               pqtype.NullRawMessage{RawMessage: payloadJSON, Valid: true},
		MaxRetries:            int32(maxRetries),
		VisibilityTimeoutSecs: int32(visibilityTimeoutSecs),
		VisibleAt:             visibleAt,
	})
}

// ClaimNext atomically claims the next available item for processing.
func (r *Repository) ClaimNext(ctx context.Context, workerID string) (*events.QueueItem, error) {
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
	itemType events.QueueItemType,
	limit int,
) ([]*events.QueueItem, error) {
	rows, err := r.queries.ListQueueItemsByType(ctx, ListQueueItemsByTypeParams{
		Type:       string(itemType),
		LimitCount: int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*events.QueueItem, len(rows))
	for i, row := range rows {
		result[i] = r.rowToQueueItem(row)
	}

	return result, nil
}

// rowToQueueItem converts a database row to a QueueItem domain object.
func (r *Repository) rowToQueueItem(row *EventQueue) *events.QueueItem {
	var payload map[string]any
	if row.Payload.Valid && len(row.Payload.RawMessage) > 0 {
		err := json.Unmarshal(row.Payload.RawMessage, &payload)
		if err != nil {
			slog.Warn(
				"failed to unmarshal queue item payload",
				slog.String("error", err.Error()),
				slog.String("id", row.ID),
			)
		}
	}

	return &events.QueueItem{
		ID:                    row.ID,
		Type:                  events.QueueItemType(row.Type),
		Payload:               payload,
		Status:                events.QueueItemStatus(row.Status),
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
