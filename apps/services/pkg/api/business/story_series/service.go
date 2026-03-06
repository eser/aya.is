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

// GetByID returns a series by its ID with localized text.
func (s *Service) GetByID(ctx context.Context, localeCode string, id string) (*StorySeries, error) {
	series, err := s.repo.GetSeriesByID(ctx, localeCode, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetSeries, id, err)
	}

	return series, nil
}

// GetBySlug returns a series by its slug with localized text.
func (s *Service) GetBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
) (*StorySeries, error) {
	series, err := s.repo.GetSeriesBySlug(ctx, localeCode, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetSeries, slug, err)
	}

	return series, nil
}

// List returns all series with localized text.
func (s *Service) List(ctx context.Context, localeCode string) ([]*StorySeries, error) {
	series, err := s.repo.ListSeries(ctx, localeCode)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListSeries, err)
	}

	return series, nil
}

// Create creates a new series with an initial translation.
func (s *Service) Create(ctx context.Context, params CreateParams) (*StorySeries, error) {
	seriesID := s.idGenerator()

	err := s.repo.InsertSeries(
		ctx,
		seriesID,
		params.Slug,
		params.SeriesPictureURI,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateSeries, err)
	}

	err = s.repo.UpsertSeriesTx(
		ctx,
		seriesID,
		params.LocaleCode,
		params.Title,
		params.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToUpsertTranslation, err)
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

	result, err := s.repo.GetSeriesByID(ctx, params.LocaleCode, seriesID)
	if err != nil {
		return nil, fmt.Errorf("fetching created series: %w", err)
	}

	return result, nil
}

// Update updates series base fields.
func (s *Service) Update(ctx context.Context, seriesID string, params UpdateParams) error {
	rows, err := s.repo.UpdateSeries(
		ctx,
		seriesID,
		params.Slug,
		params.SeriesPictureURI,
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
			"slug": params.Slug,
		},
	})

	return nil
}

// UpsertTranslation creates or updates a series translation.
func (s *Service) UpsertTranslation(
	ctx context.Context,
	seriesID string,
	params TranslationParams,
) error {
	err := s.repo.UpsertSeriesTx(
		ctx,
		seriesID,
		params.LocaleCode,
		params.Title,
		params.Description,
	)
	if err != nil {
		return fmt.Errorf(
			"%w(id: %s, locale: %s): %w",
			ErrFailedToUpsertTranslation,
			seriesID,
			params.LocaleCode,
			err,
		)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StorySeriesUpdated,
		EntityType: "story_series",
		EntityID:   seriesID,
		ActorID:    nil,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"locale": params.LocaleCode,
			"title":  params.Title,
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
