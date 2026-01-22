package linksync

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string

// Service provides link sync operations.
type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator IDGenerator
}

// NewService creates a new link sync service.
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

// GetManagedLinks returns managed links for a given kind.
func (s *Service) GetManagedLinks(
	ctx context.Context,
	kind string,
	limit int,
) ([]*ManagedLink, error) {
	links, err := s.repo.ListManagedLinksForKind(ctx, kind, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetLinks, err)
	}

	return links, nil
}

// GetLastSyncTime returns the created_at time of the most recent import for a link.
// Returns nil if no imports exist (indicating a full sync is needed).
func (s *Service) GetLastSyncTime(ctx context.Context, linkID string) (*time.Time, error) {
	latestImport, err := s.repo.GetLatestImportByLinkID(ctx, linkID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetLatestImport, err)
	}

	return &latestImport.CreatedAt, nil
}

// UpsertImport creates or updates a link import.
func (s *Service) UpsertImport(
	ctx context.Context,
	profileLinkID string,
	remoteID string,
	properties map[string]any,
) error {
	// Check if import exists
	existing, err := s.repo.GetLinkImportByRemoteID(ctx, profileLinkID, remoteID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %w", ErrFailedToUpsertImport, err)
	}

	if existing != nil {
		// Update existing
		err = s.repo.UpdateLinkImport(ctx, existing.ID, properties)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrFailedToUpsertImport, err)
		}

		return nil
	}

	// Create new
	id := s.idGenerator()

	err = s.repo.CreateLinkImport(ctx, id, profileLinkID, remoteID, properties)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpsertImport, err)
	}

	return nil
}

// MarkDeletedImports marks imports as deleted that are not in the active list.
func (s *Service) MarkDeletedImports(
	ctx context.Context,
	linkID string,
	activeRemoteIDs []string,
) (int64, error) {
	count, err := s.repo.MarkLinkImportsDeletedExcept(ctx, linkID, activeRemoteIDs)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrFailedToMarkDeleted, err)
	}

	return count, nil
}

// UpdateLinkTokens updates the OAuth tokens for a link.
func (s *Service) UpdateLinkTokens(
	ctx context.Context,
	linkID string,
	accessToken string,
	accessTokenExpiresAt *time.Time,
	refreshToken *string,
) error {
	err := s.repo.UpdateLinkTokens(ctx, linkID, accessToken, accessTokenExpiresAt, refreshToken)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateTokens, err)
	}

	return nil
}
