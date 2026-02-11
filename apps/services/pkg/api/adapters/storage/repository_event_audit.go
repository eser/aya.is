package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/sqlc-dev/pqtype"
)

// InsertAudit persists an audit entry to the database.
func (r *Repository) InsertAudit(
	ctx context.Context,
	id string,
	params events.AuditParams,
) error {
	var payloadJSON pqtype.NullRawMessage

	if params.Payload != nil {
		data, err := json.Marshal(params.Payload)
		if err != nil {
			return err
		}

		payloadJSON = pqtype.NullRawMessage{RawMessage: data, Valid: true}
	}

	return r.queries.InsertEventAudit(ctx, InsertEventAuditParams{
		ID:         id,
		EventType:  string(params.EventType),
		EntityType: params.EntityType,
		EntityID:   sql.NullString{String: params.EntityID, Valid: params.EntityID != ""},
		ActorID:    toNullString(params.ActorID),
		ActorKind:  string(params.ActorKind),
		SessionID:  toNullString(params.SessionID),
		Payload:    payloadJSON,
	})
}

// ListByEntity returns audit entries for a given entity type and ID.
func (r *Repository) ListByEntity(
	ctx context.Context,
	entityType string,
	entityID string,
	limit int,
) ([]*events.AuditEntry, error) {
	rows, err := r.queries.ListEventAuditByEntity(ctx, ListEventAuditByEntityParams{
		EntityType: entityType,
		EntityID:   sql.NullString{String: entityID, Valid: entityID != ""},
		LimitCount: int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*events.AuditEntry, len(rows))
	for i, row := range rows {
		result[i] = r.rowToAuditEntry(row)
	}

	return result, nil
}

// rowToAuditEntry converts a database row to an AuditEntry domain object.
func (r *Repository) rowToAuditEntry(row *EventAudit) *events.AuditEntry {
	var payload map[string]any
	if row.Payload.Valid && len(row.Payload.RawMessage) > 0 {
		err := json.Unmarshal(row.Payload.RawMessage, &payload)
		if err != nil {
			slog.Warn(
				"failed to unmarshal audit entry payload",
				slog.String("error", err.Error()),
				slog.String("id", row.ID),
			)
		}
	}

	var actorID *string
	if row.ActorID.Valid {
		actorID = &row.ActorID.String
	}

	var sessionID *string
	if row.SessionID.Valid {
		sessionID = &row.SessionID.String
	}

	return &events.AuditEntry{
		ID:         row.ID,
		EventType:  events.EventType(row.EventType),
		EntityType: row.EntityType,
		EntityID:   row.EntityID.String,
		ActorID:    actorID,
		ActorKind:  events.ActorKind(row.ActorKind),
		SessionID:  sessionID,
		Payload:    payload,
		CreatedAt:  row.CreatedAt,
	}
}

// toNullString converts a *string to sql.NullString.
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{} //nolint:exhaustruct
	}

	return sql.NullString{String: *s, Valid: true}
}
