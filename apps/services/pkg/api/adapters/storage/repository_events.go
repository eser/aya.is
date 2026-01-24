package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

// Enqueue inserts a new event into the queue.
func (r *Repository) Enqueue(
	ctx context.Context,
	id string,
	eventType events.EventType,
	payload map[string]any,
	maxRetries int,
	visibilityTimeoutSecs int,
	visibleAt time.Time,
) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.queries.EnqueueEvent(ctx, EnqueueEventParams{
		ID:                    id,
		Type:                  string(eventType),
		Payload:               payloadJSON,
		MaxRetries:            int32(maxRetries),
		VisibilityTimeoutSecs: int32(visibilityTimeoutSecs),
		VisibleAt:             visibleAt,
	})
}

// ClaimNext atomically claims the next available event for processing.
func (r *Repository) ClaimNext(ctx context.Context, workerID string) (*events.Event, error) {
	row, err := r.queries.ClaimNextEvent(ctx, ClaimNextEventParams{
		WorkerID: sql.NullString{String: workerID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return r.rowToEventQueueItem(row), nil
}

// Complete marks an event as successfully completed.
func (r *Repository) Complete(ctx context.Context, id string, workerID string) error {
	_, err := r.queries.CompleteEvent(ctx, CompleteEventParams{
		ID:       id,
		WorkerID: sql.NullString{String: workerID, Valid: true},
	})

	return err
}

// Fail marks an event as failed with error message and backoff.
func (r *Repository) Fail(
	ctx context.Context,
	id string,
	workerID string,
	errorMessage string,
	backoffSeconds int,
) error {
	_, err := r.queries.FailEvent(ctx, FailEventParams{
		ID:             id,
		WorkerID:       sql.NullString{String: workerID, Valid: true},
		ErrorMessage:   sql.NullString{String: errorMessage, Valid: true},
		BackoffSeconds: int32(backoffSeconds),
	})

	return err
}

// ListByType returns events of a given type for audit/debugging.
func (r *Repository) ListByType(
	ctx context.Context,
	eventType events.EventType,
	limit int,
) ([]*events.Event, error) {
	rows, err := r.queries.ListEventsByType(ctx, ListEventsByTypeParams{
		Type:       string(eventType),
		LimitCount: int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*events.Event, len(rows))
	for i, row := range rows {
		result[i] = r.rowToEventQueueItem(row)
	}

	return result, nil
}

// rowToEventQueueItem converts a database row to an Event domain object.
func (r *Repository) rowToEventQueueItem(row *EventQueue) *events.Event {
	var payload map[string]any
	if len(row.Payload) > 0 {
		_ = json.Unmarshal(row.Payload, &payload)
	}

	return &events.Event{
		ID:                    row.ID,
		Type:                  events.EventType(row.Type),
		Payload:               payload,
		Status:                events.EventStatus(row.Status),
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
