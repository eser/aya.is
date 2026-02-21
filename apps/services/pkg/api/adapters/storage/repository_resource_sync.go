package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/resourcesync"
	"github.com/eser/aya.is/services/pkg/lib/vars"
	"github.com/sqlc-dev/pqtype"
)

// ListGitHubResourcesForSync returns GitHub resources with their profile's managed GitHub access tokens.
func (r *Repository) ListGitHubResourcesForSync(
	ctx context.Context,
	batchSize int,
) ([]*resourcesync.GitHubResourceForSync, error) {
	rows, err := r.queries.ListGitHubResourcesForSync(ctx, ListGitHubResourcesForSyncParams{
		BatchSize: int32(batchSize),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*resourcesync.GitHubResourceForSync, 0, len(rows))

	for _, row := range rows {
		var properties map[string]any
		if row.ResourceProperties.Valid {
			unmarshalErr := json.Unmarshal(row.ResourceProperties.RawMessage, &properties)
			if unmarshalErr != nil {
				properties = nil
			}
		}

		result = append(result, &resourcesync.GitHubResourceForSync{
			ResourceID:               row.ResourceID,
			ProfileID:                row.ProfileID,
			ResourceRemoteID:         row.ResourceRemoteID.String,
			ResourcePublicID:         row.ResourcePublicID.String,
			ResourceProperties:       properties,
			LinkID:                   row.LinkID,
			AuthAccessToken:          row.AuthAccessToken.String,
			AuthAccessTokenExpiresAt: vars.ToTimePtr(row.AuthAccessTokenExpiresAt),
			AuthRefreshToken:         vars.ToStringPtr(row.AuthRefreshToken),
		})
	}

	return result, nil
}

// UpdateProfileResourcePropertiesForResourceSync updates the properties JSONB of a profile_resource.
func (r *Repository) UpdateProfileResourcePropertiesForResourceSync(
	ctx context.Context,
	id string,
	properties map[string]any,
) error {
	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	_, err = r.queries.UpdateProfileResourceProperties(ctx, UpdateProfileResourcePropertiesParams{
		ID:         id,
		Properties: pqtype.NullRawMessage{RawMessage: propertiesJSON, Valid: true},
	})

	return err
}

// UpdateProfileMembershipPropertiesForResourceSync updates the properties JSONB of a profile_membership.
// Uses JSONB || operator to shallow-merge new properties into existing ones,
// preserving keys not present in the update (e.g., "videos" when updating "github").
func (r *Repository) UpdateProfileMembershipPropertiesForResourceSync(
	ctx context.Context,
	id string,
	properties map[string]any,
) error {
	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	_, err = r.queries.MergeProfileMembershipProperties(ctx, MergeProfileMembershipPropertiesParams{
		ID:         id,
		Properties: pqtype.NullRawMessage{RawMessage: propertiesJSON, Valid: true},
	})

	return err
}

// GetMembershipByProfiles looks up a profile_membership between two profiles and returns its ID.
func (r *Repository) GetMembershipByProfiles(
	ctx context.Context,
	profileID string,
	memberProfileID string,
) (string, error) {
	id, err := r.queries.GetMembershipIDBetweenProfiles(ctx, GetMembershipIDBetweenProfilesParams{
		ProfileID:       profileID,
		MemberProfileID: sql.NullString{String: memberProfileID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return id, nil
}

// GetProfileLinkByRemoteIDForResourceSync finds a profile_link by kind and remote_id (globally, not per-profile).
// Returns the profile_id of the link owner so we can match contributors to profiles.
func (r *Repository) GetProfileLinkByRemoteIDForResourceSync(
	ctx context.Context,
	kind string,
	remoteID string,
) (string, error) {
	profileID, err := r.queries.FindProfileLinkProfileByKindAndRemoteID(
		ctx,
		FindProfileLinkProfileByKindAndRemoteIDParams{
			Kind:     kind,
			RemoteID: sql.NullString{String: remoteID, Valid: true},
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return profileID, nil
}

// GetProfileLinksByRemoteIDsForResourceSync batch-loads profile_links and returns map[remoteID]profileID.
// Joins with profile to filter for individual profiles only, since contributor matching
// should find the person, not organizations/products that share the same OAuth token.
func (r *Repository) GetProfileLinksByRemoteIDsForResourceSync(
	ctx context.Context,
	kind string,
	remoteIDs []string,
) (map[string]string, error) {
	rows, err := r.queries.GetProfileLinksByRemoteIDs(ctx, GetProfileLinksByRemoteIDsParams{
		Kind:      kind,
		RemoteIds: remoteIDs,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(remoteIDs))

	for _, row := range rows {
		if row.RemoteID.Valid {
			result[row.RemoteID.String] = row.ProfileID
		}
	}

	return result, nil
}

// GetMembershipsByProfilePairsForResourceSync batch-loads memberships.
// Returns map["profileID:memberProfileID"]membershipID.
func (r *Repository) GetMembershipsByProfilePairsForResourceSync(
	ctx context.Context,
	profileIDs []string,
	memberProfileIDs []string,
) (map[string]string, error) {
	rows, err := r.queries.GetMembershipsByProfilePairs(ctx, GetMembershipsByProfilePairsParams{
		ProfileIds:       profileIDs,
		MemberProfileIds: memberProfileIDs,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)

	for _, row := range rows {
		if row.MemberProfileID.Valid {
			result[row.ProfileID+":"+row.MemberProfileID.String] = row.ID
		}
	}

	return result, nil
}

// resourceSyncAdapter wraps *Repository to satisfy the resourcesync.Repository interface.
// This is needed because method signatures differ (e.g., GetProfileLinkByRemoteID
// has no profileID parameter in the resourcesync interface).
type resourceSyncAdapter struct {
	repo *Repository
}

// NewResourceSyncRepository creates a resourcesync.Repository adapter from a storage Repository.
func NewResourceSyncRepository(repo *Repository) resourcesync.Repository {
	return &resourceSyncAdapter{repo: repo}
}

func (a *resourceSyncAdapter) ListGitHubResourcesForSync(
	ctx context.Context,
	batchSize int,
) ([]*resourcesync.GitHubResourceForSync, error) {
	return a.repo.ListGitHubResourcesForSync(ctx, batchSize)
}

func (a *resourceSyncAdapter) UpdateProfileResourceProperties(
	ctx context.Context,
	id string,
	properties map[string]any,
) error {
	return a.repo.UpdateProfileResourcePropertiesForResourceSync(ctx, id, properties)
}

func (a *resourceSyncAdapter) UpdateProfileMembershipProperties(
	ctx context.Context,
	id string,
	properties map[string]any,
) error {
	return a.repo.UpdateProfileMembershipPropertiesForResourceSync(ctx, id, properties)
}

func (a *resourceSyncAdapter) GetMembershipByProfiles(
	ctx context.Context,
	profileID string,
	memberProfileID string,
) (string, error) {
	return a.repo.GetMembershipByProfiles(ctx, profileID, memberProfileID)
}

func (a *resourceSyncAdapter) GetProfileLinkByRemoteID(
	ctx context.Context,
	kind string,
	remoteID string,
) (string, error) {
	return a.repo.GetProfileLinkByRemoteIDForResourceSync(ctx, kind, remoteID)
}

func (a *resourceSyncAdapter) GetProfileLinksByRemoteIDs(
	ctx context.Context,
	kind string,
	remoteIDs []string,
) (map[string]string, error) {
	return a.repo.GetProfileLinksByRemoteIDsForResourceSync(ctx, kind, remoteIDs)
}

func (a *resourceSyncAdapter) GetMembershipsByProfilePairs(
	ctx context.Context,
	profileIDs []string,
	memberProfileIDs []string,
) (map[string]string, error) {
	return a.repo.GetMembershipsByProfilePairsForResourceSync(ctx, profileIDs, memberProfileIDs)
}
