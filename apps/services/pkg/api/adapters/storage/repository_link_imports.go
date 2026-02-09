package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/linksync"
	"github.com/eser/aya.is/services/pkg/lib/vars"
	"github.com/sqlc-dev/pqtype"
)

// ListManagedLinksForKind returns all managed, non-deleted links of a kind.
func (r *Repository) ListManagedLinksForKind(
	ctx context.Context,
	kind string,
	limit int,
) ([]*linksync.ManagedLink, error) {
	rows, err := r.queries.ListManagedLinksForKind(ctx, ListManagedLinksForKindParams{
		Kind:       kind,
		LimitCount: int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*linksync.ManagedLink, len(rows))
	for i, row := range rows {
		result[i] = &linksync.ManagedLink{
			ID:                       row.ID,
			ProfileID:                row.ProfileID,
			Kind:                     row.Kind,
			RemoteID:                 row.RemoteID.String,
			AuthAccessToken:          row.AuthAccessToken.String,
			AuthAccessTokenExpiresAt: vars.ToTimePtr(row.AuthAccessTokenExpiresAt),
			AuthRefreshToken:         vars.ToStringPtr(row.AuthRefreshToken),
		}
	}

	return result, nil
}

// GetLatestImportByLinkID returns the most recent import for a link.
func (r *Repository) GetLatestImportByLinkID(
	ctx context.Context,
	linkID string,
) (*linksync.LinkImport, error) {
	row, err := r.queries.GetLatestImportByLinkID(ctx, GetLatestImportByLinkIDParams{
		ProfileLinkID: linkID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, err
	}

	return r.rowToLinkImport(row), nil
}

// GetLinkImportByRemoteID returns an import by link ID and remote ID.
func (r *Repository) GetLinkImportByRemoteID(
	ctx context.Context,
	linkID string,
	remoteID string,
) (*linksync.LinkImport, error) {
	row, err := r.queries.GetLinkImportByRemoteID(ctx, GetLinkImportByRemoteIDParams{
		ProfileLinkID: linkID,
		RemoteID:      sql.NullString{String: remoteID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, err
	}

	return r.rowToLinkImport(row), nil
}

// CreateLinkImport creates a new link import.
func (r *Repository) CreateLinkImport(
	ctx context.Context,
	id string,
	profileLinkID string,
	remoteID string,
	properties map[string]any,
) error {
	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	return r.queries.CreateLinkImport(ctx, CreateLinkImportParams{
		ID:            id,
		ProfileLinkID: profileLinkID,
		RemoteID:      sql.NullString{String: remoteID, Valid: true},
		Properties:    pqtype.NullRawMessage{RawMessage: propertiesJSON, Valid: true},
	})
}

// UpdateLinkImport updates an existing link import.
func (r *Repository) UpdateLinkImport(
	ctx context.Context,
	id string,
	properties map[string]any,
) error {
	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	_, err = r.queries.UpdateLinkImport(ctx, UpdateLinkImportParams{
		ID:         id,
		Properties: pqtype.NullRawMessage{RawMessage: propertiesJSON, Valid: true},
	})

	return err
}

// MarkLinkImportsDeletedExcept soft-deletes imports not in the given remote IDs.
func (r *Repository) MarkLinkImportsDeletedExcept(
	ctx context.Context,
	linkID string,
	activeRemoteIDs []string,
) (int64, error) {
	return r.queries.MarkLinkImportsDeletedExcept(ctx, MarkLinkImportsDeletedExceptParams{
		ProfileLinkID:   linkID,
		ActiveRemoteIds: activeRemoteIDs,
	})
}

// UpdateLinkTokens updates the OAuth tokens for a link.
func (r *Repository) UpdateLinkTokens(
	ctx context.Context,
	linkID string,
	accessToken string,
	accessTokenExpiresAt *time.Time,
	refreshToken *string,
) error {
	var expiresAtSQL sql.NullTime
	if accessTokenExpiresAt != nil {
		expiresAtSQL = sql.NullTime{Time: *accessTokenExpiresAt, Valid: true}
	}

	var refreshTokenSQL sql.NullString
	if refreshToken != nil {
		refreshTokenSQL = sql.NullString{String: *refreshToken, Valid: true}
	}

	_, err := r.queries.UpdateProfileLinkTokens(ctx, UpdateProfileLinkTokensParams{
		ID:                       linkID,
		AuthAccessToken:          sql.NullString{String: accessToken, Valid: true},
		AuthAccessTokenExpiresAt: expiresAtSQL,
		AuthRefreshToken:         refreshTokenSQL,
	})

	return err
}

// ListImportsForStoryCreation returns imports that don't have corresponding managed stories.
func (r *Repository) ListImportsForStoryCreation(
	ctx context.Context,
	kind string,
	limit int,
) ([]*linksync.LinkImportForStoryCreation, error) {
	rows, err := r.queries.ListLinkImportsForStoryCreation(
		ctx,
		ListLinkImportsForStoryCreationParams{
			Kind:       kind,
			LimitCount: int32(limit),
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*linksync.LinkImportForStoryCreation, len(rows))

	for i, row := range rows {
		var properties map[string]any
		if row.Properties.Valid {
			_ = json.Unmarshal(row.Properties.RawMessage, &properties)
		}

		result[i] = &linksync.LinkImportForStoryCreation{
			ID:                   row.ID,
			ProfileLinkID:        row.ProfileLinkID,
			RemoteID:             row.RemoteID.String,
			Properties:           properties,
			CreatedAt:            row.CreatedAt,
			ProfileID:            row.ProfileID,
			ProfileDefaultLocale: row.ProfileDefaultLocale,
		}
	}

	return result, nil
}

// ListImportsWithExistingStories returns imports that have matching managed stories.
func (r *Repository) ListImportsWithExistingStories(
	ctx context.Context,
	kind string,
	limit int,
) ([]*linksync.LinkImportWithStory, error) {
	rows, err := r.queries.ListImportsWithExistingStories(
		ctx,
		ListImportsWithExistingStoriesParams{
			Kind:       kind,
			LimitCount: int32(limit),
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*linksync.LinkImportWithStory, len(rows))

	for i, row := range rows {
		var properties map[string]any
		if row.Properties.Valid {
			_ = json.Unmarshal(row.Properties.RawMessage, &properties)
		}

		item := &linksync.LinkImportWithStory{
			ID:                   row.ID,
			ProfileLinkID:        row.ProfileLinkID,
			RemoteID:             row.RemoteID.String,
			Properties:           properties,
			CreatedAt:            row.CreatedAt,
			ProfileID:            row.ProfileID,
			ProfileDefaultLocale: row.ProfileDefaultLocale,
			StoryID:              row.StoryID,
		}

		if row.PublicationID.Valid {
			item.PublicationID = &row.PublicationID.String
		}

		result[i] = item
	}

	return result, nil
}

// rowToLinkImport converts a database row to a LinkImport domain object.
func (r *Repository) rowToLinkImport(row *ProfileLinkImport) *linksync.LinkImport {
	var properties map[string]any
	if row.Properties.Valid {
		_ = json.Unmarshal(row.Properties.RawMessage, &properties)
	}

	return &linksync.LinkImport{
		ID:            row.ID,
		ProfileLinkID: row.ProfileLinkID,
		RemoteID:      row.RemoteID.String,
		Properties:    properties,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:     vars.ToTimePtr(row.DeletedAt),
	}
}
