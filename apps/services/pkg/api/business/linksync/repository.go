package linksync

import (
	"context"
	"time"
)

// Repository defines the storage operations for link sync.
type Repository interface {
	// ListManagedLinksForKind returns all managed, non-deleted links of a kind.
	ListManagedLinksForKind(ctx context.Context, kind string, limit int) ([]*ManagedLink, error)

	// GetLatestImportByLinkID returns the most recent import for a link.
	GetLatestImportByLinkID(ctx context.Context, linkID string) (*LinkImport, error)

	// GetLinkImportByRemoteID returns an import by link ID and remote ID.
	GetLinkImportByRemoteID(ctx context.Context, linkID, remoteID string) (*LinkImport, error)

	// CreateLinkImport creates a new link import.
	CreateLinkImport(
		ctx context.Context,
		id, profileLinkID, remoteID string,
		properties map[string]any,
	) error

	// UpdateLinkImport updates an existing link import.
	UpdateLinkImport(ctx context.Context, id string, properties map[string]any) error

	// MarkLinkImportsDeletedExcept soft-deletes imports not in the given remote IDs.
	MarkLinkImportsDeletedExcept(
		ctx context.Context,
		linkID string,
		activeRemoteIDs []string,
	) (int64, error)

	// UpdateLinkTokens updates the OAuth tokens for a link.
	UpdateLinkTokens(
		ctx context.Context,
		linkID string,
		accessToken string,
		accessTokenExpiresAt *time.Time,
		refreshToken *string,
	) error

	// ListImportsForStoryCreation returns imports that don't have corresponding managed stories.
	ListImportsForStoryCreation(
		ctx context.Context,
		kind string,
		limit int,
	) ([]*LinkImportForStoryCreation, error)

	// ListImportsWithExistingStories returns imports that have corresponding managed stories.
	ListImportsWithExistingStories(
		ctx context.Context,
		kind string,
		limit int,
	) ([]*LinkImportWithStory, error)
}
