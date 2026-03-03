package story_series

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
)

func DefaultIDGenerator() string {
	return lib.IDsGenerateUnique()
}

// Service provides story series operations.
type Service struct {
	logger       *logfx.Logger
	repo         Repository
	idGenerator  IDGenerator
	auditService *events.AuditService
}

// NewService creates a new story series service.
func NewService(
	logger *logfx.Logger,
	repo Repository,
	idGenerator IDGenerator,
	auditService *events.AuditService,
) *Service {
	return &Service{
		logger:       logger,
		repo:         repo,
		idGenerator:  idGenerator,
		auditService: auditService,
	}
}

// GetByID returns a series by its ID.
func (s *Service) GetByID(ctx context.Context, id string) (*StorySeries, error) {
	series, err := s.repo.GetSeriesByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetSeries, id, err)
	}

	return series, nil
}

// GetBySlug returns a series by its slug.
func (s *Service) GetBySlug(ctx context.Context, slug string) (*StorySeries, error) {
	series, err := s.repo.GetSeriesBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetSeries, slug, err)
	}

	return series, nil
}

// List returns all series.
func (s *Service) List(ctx context.Context) ([]*StorySeries, error) {
	series, err := s.repo.ListSeries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListSeries, err)
	}

	return series, nil
}

// Create creates a new series.
func (s *Service) Create(ctx context.Context, params CreateParams) (*StorySeries, error) {
	seriesID := s.idGenerator()

	series, err := s.repo.InsertSeries(
		ctx,
		seriesID,
		params.Slug,
		params.SeriesPictureURI,
		params.Title,
		params.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateSeries, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StorySeriesCreated,
		EntityType: "story_series",
		EntityID:   seriesID,
		ActorID:    nil,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"slug":  params.Slug,
			"title": params.Title,
		},
	})

	return series, nil
}

// Update updates an existing series.
func (s *Service) Update(ctx context.Context, seriesID string, params UpdateParams) error {
	rows, err := s.repo.UpdateSeries(
		ctx,
		seriesID,
		params.Slug,
		params.SeriesPictureURI,
		params.Title,
		params.Description,
	)
	if err != nil {
		return fmt.Errorf("%w(id: %s): %w", ErrFailedToUpdateSeries, seriesID, err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: %s", ErrSeriesNotFound, seriesID)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StorySeriesUpdated,
		EntityType: "story_series",
		EntityID:   seriesID,
		ActorID:    nil,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"slug":  params.Slug,
			"title": params.Title,
		},
	})

	return nil
}

// Delete soft-deletes a series.
func (s *Service) Delete(ctx context.Context, seriesID string) error {
	rows, err := s.repo.RemoveSeries(ctx, seriesID)
	if err != nil {
		return fmt.Errorf("%w(id: %s): %w", ErrFailedToRemoveSeries, seriesID, err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: %s", ErrSeriesNotFound, seriesID)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StorySeriesDeleted,
		EntityType: "story_series",
		EntityID:   seriesID,
		ActorID:    nil,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload:    nil,
	})

	return nil
}
