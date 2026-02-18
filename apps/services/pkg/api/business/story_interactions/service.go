package story_interactions

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

func DefaultIDGenerator() string {
	return lib.IDsGenerateUnique()
}

// Service provides story interaction operations.
type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator IDGenerator
}

// NewService creates a new story interactions service.
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

// SetInteraction upserts a generic (non-RSVP) interaction.
func (s *Service) SetInteraction(
	ctx context.Context,
	storyID string,
	profileID string,
	kind string,
) (*StoryInteraction, error) {
	id := s.idGenerator()

	interaction, err := s.repo.UpsertInteraction(ctx, id, storyID, profileID, kind)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToSetInteraction, err)
	}

	return interaction, nil
}

// SetRSVP sets an RSVP interaction, enforcing mutual exclusivity among RSVP kinds.
// It removes any existing RSVP interactions before setting the new one.
func (s *Service) SetRSVP(
	ctx context.Context,
	storyID string,
	profileID string,
	kind InteractionKind,
) (*StoryInteraction, error) {
	if !IsRSVPKind(kind) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidInteractionKind, kind)
	}

	// Remove existing RSVP interactions for mutual exclusivity
	_, err := s.repo.RemoveInteractionsByKinds(ctx, storyID, profileID, RSVPKindsCSV())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToRemoveInteraction, err)
	}

	id := s.idGenerator()

	interaction, err := s.repo.UpsertInteraction(ctx, id, storyID, profileID, string(kind))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToSetInteraction, err)
	}

	return interaction, nil
}

// RemoveInteraction removes a specific interaction.
func (s *Service) RemoveInteraction(
	ctx context.Context,
	storyID string,
	profileID string,
	kind string,
) error {
	_, err := s.repo.RemoveInteraction(ctx, storyID, profileID, kind)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRemoveInteraction, err)
	}

	return nil
}

// GetInteraction returns a specific interaction for a profile on a story.
func (s *Service) GetInteraction(
	ctx context.Context,
	storyID string,
	profileID string,
	kind string,
) (*StoryInteraction, error) {
	interaction, err := s.repo.GetInteraction(ctx, storyID, profileID, kind)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetInteraction, err)
	}

	return interaction, nil
}

// ListForProfile returns all interactions a profile has on a story.
func (s *Service) ListForProfile(
	ctx context.Context,
	storyID string,
	profileID string,
) ([]*StoryInteraction, error) {
	interactions, err := s.repo.ListInteractionsForProfile(ctx, storyID, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListInteractions, err)
	}

	return interactions, nil
}

// ListInteractions lists interactions on a story with profile info, optionally filtered by kind.
func (s *Service) ListInteractions(
	ctx context.Context,
	localeCode string,
	storyID string,
	filterKind *string,
) ([]*InteractionWithProfile, error) {
	interactions, err := s.repo.ListInteractions(ctx, localeCode, storyID, filterKind)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListInteractions, err)
	}

	return interactions, nil
}

// CountInteractions returns interaction counts grouped by kind for a story.
func (s *Service) CountInteractions(
	ctx context.Context,
	storyID string,
) ([]*InteractionCount, error) {
	counts, err := s.repo.CountInteractionsByKind(ctx, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCountInteractions, err)
	}

	return counts, nil
}
