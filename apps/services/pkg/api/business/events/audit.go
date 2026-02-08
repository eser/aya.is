package events

import (
	"context"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

// ActorKind classifies who triggered the event.
type ActorKind string

const (
	ActorUser   ActorKind = "user"
	ActorSystem ActorKind = "system"
	ActorWorker ActorKind = "worker"
)

// AuditEntry is an immutable record of a business event.
type AuditEntry struct {
	ID         string
	EventType  EventType
	EntityType string
	EntityID   string
	ActorID    *string
	ActorKind  ActorKind
	SessionID  *string
	Payload    map[string]any
	CreatedAt  time.Time
}

// AuditParams holds parameters for recording an audit entry.
type AuditParams struct {
	EventType  EventType
	EntityType string
	EntityID   string
	ActorID    *string
	ActorKind  ActorKind
	SessionID  *string
	Payload    map[string]any
}

// AuditRepository defines storage operations for audit entries (port).
type AuditRepository interface {
	InsertAudit(
		ctx context.Context,
		id string,
		params AuditParams,
	) error

	ListByEntity(
		ctx context.Context,
		entityType string,
		entityID string,
		limit int,
	) ([]*AuditEntry, error)
}

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string

// AuditService records audit entries.
type AuditService struct {
	logger      *logfx.Logger
	repo        AuditRepository
	idGenerator IDGenerator
}

// NewAuditService creates a new audit service.
func NewAuditService(
	logger *logfx.Logger,
	repo AuditRepository,
	idGenerator IDGenerator,
) *AuditService {
	return &AuditService{
		logger:      logger,
		repo:        repo,
		idGenerator: idGenerator,
	}
}

// Record persists an audit entry. Fire-and-forget: errors are logged but not propagated,
// because audit failures must never break business operations.
func (s *AuditService) Record(ctx context.Context, params AuditParams) {
	id := s.idGenerator()

	err := s.repo.InsertAudit(ctx, id, params)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to record audit entry",
			slog.String("event_type", string(params.EventType)),
			slog.String("entity_type", params.EntityType),
			slog.String("entity_id", params.EntityID),
			slog.Any("error", err),
		)
	}
}

// ListByEntity returns audit entries for a given entity.
func (s *AuditService) ListByEntity(
	ctx context.Context,
	entityType string,
	entityID string,
	limit int,
) ([]*AuditEntry, error) {
	return s.repo.ListByEntity(ctx, entityType, entityID, limit)
}
