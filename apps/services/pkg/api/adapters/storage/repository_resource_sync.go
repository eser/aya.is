package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/resourcesync"
	"github.com/eser/aya.is/services/pkg/lib/vars"
	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
)

// listGitHubResourcesForSyncSQL is the raw SQL query for listing GitHub resources.
// TODO(sqlc): Replace with r.queries.ListGitHubResourcesForSync once the sqlc query is generated.
const listGitHubResourcesForSyncSQL = `
SELECT
  pr.id as resource_id,
  pr.profile_id,
  pr.remote_id as resource_remote_id,
  pr.public_id as resource_public_id,
  pr.properties as resource_properties,
  pl.id as link_id,
  pl.auth_access_token,
  pl.auth_access_token_expires_at,
  pl.auth_refresh_token
FROM "profile_resource" pr
  INNER JOIN "profile_link" pl ON pl.profile_id = pr.profile_id
    AND pl.kind = 'github'
    AND pl.is_managed = true
    AND pl.deleted_at IS NULL
    AND pl.auth_access_token IS NOT NULL
WHERE pr.kind = 'github_repo'
  AND pr.is_managed = true
  AND pr.deleted_at IS NULL
ORDER BY pr.created_at ASC
LIMIT $1
`

// ListGitHubResourcesForSync returns GitHub resources with their profile's managed GitHub access tokens.
func (r *Repository) ListGitHubResourcesForSync(
	ctx context.Context,
	batchSize int,
) ([]*resourcesync.GitHubResourceForSync, error) {
	rows, err := r.dbtx.QueryContext(ctx, listGitHubResourcesForSyncSQL, batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var result []*resourcesync.GitHubResourceForSync

	for rows.Next() {
		var (
			resourceID               string
			profileID                string
			resourceRemoteID         sql.NullString
			resourcePublicID         sql.NullString
			resourceProperties       pqtype.NullRawMessage
			linkID                   string
			authAccessToken          sql.NullString
			authAccessTokenExpiresAt sql.NullTime
			authRefreshToken         sql.NullString
		)

		scanErr := rows.Scan(
			&resourceID,
			&profileID,
			&resourceRemoteID,
			&resourcePublicID,
			&resourceProperties,
			&linkID,
			&authAccessToken,
			&authAccessTokenExpiresAt,
			&authRefreshToken,
		)
		if scanErr != nil {
			return nil, scanErr
		}

		var properties map[string]any
		if resourceProperties.Valid {
			unmarshalErr := json.Unmarshal(resourceProperties.RawMessage, &properties)
			if unmarshalErr != nil {
				properties = nil
			}
		}

		result = append(result, &resourcesync.GitHubResourceForSync{
			ResourceID:               resourceID,
			ProfileID:                profileID,
			ResourceRemoteID:         resourceRemoteID.String,
			ResourcePublicID:         resourcePublicID.String,
			ResourceProperties:       properties,
			LinkID:                   linkID,
			AuthAccessToken:          authAccessToken.String,
			AuthAccessTokenExpiresAt: vars.ToTimePtr(authAccessTokenExpiresAt),
			AuthRefreshToken:         vars.ToStringPtr(authRefreshToken),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// updateProfileResourcePropertiesSQL updates the properties of a profile_resource.
// TODO(sqlc): Replace with r.queries.UpdateProfileResourceProperties once the sqlc query is generated.
const updateProfileResourcePropertiesSQL = `
UPDATE "profile_resource"
SET properties = $1, updated_at = NOW()
WHERE id = $2
  AND deleted_at IS NULL
`

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

	_, err = r.dbtx.ExecContext(
		ctx,
		updateProfileResourcePropertiesSQL,
		pqtype.NullRawMessage{RawMessage: propertiesJSON, Valid: true},
		id,
	)

	return err
}

// updateProfileMembershipPropertiesSQL updates the properties of a profile_membership.
// Uses JSONB || operator to shallow-merge new properties into existing ones,
// preserving keys not present in the update (e.g., "videos" when updating "github").
// Note: profile_membership has no updated_at column.
const updateProfileMembershipPropertiesSQL = `
UPDATE "profile_membership"
SET properties = COALESCE(properties, '{}'::jsonb) || $1::jsonb
WHERE id = $2
  AND deleted_at IS NULL
`

// UpdateProfileMembershipPropertiesForResourceSync updates the properties JSONB of a profile_membership.
func (r *Repository) UpdateProfileMembershipPropertiesForResourceSync(
	ctx context.Context,
	id string,
	properties map[string]any,
) error {
	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	_, err = r.dbtx.ExecContext(
		ctx,
		updateProfileMembershipPropertiesSQL,
		pqtype.NullRawMessage{RawMessage: propertiesJSON, Valid: true},
		id,
	)

	return err
}

// GetMembershipByProfiles looks up a profile_membership between two profiles and returns its ID.
// TODO(sqlc): This requires a new sqlc query "GetMembershipIDBetweenProfiles" that returns
// pm.id instead of pm.kind. The existing GetMembershipBetweenProfiles returns kind.
// Once the sqlc query is generated, replace the raw SQL below with the generated query.
func (r *Repository) GetMembershipByProfiles(
	ctx context.Context,
	profileID string,
	memberProfileID string,
) (string, error) {
	row := r.dbtx.QueryRowContext(
		ctx,
		`SELECT pm.id
		FROM "profile_membership" pm
		WHERE pm.profile_id = $1
		  AND pm.member_profile_id = $2
		  AND pm.deleted_at IS NULL
		  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
		LIMIT 1`,
		profileID,
		memberProfileID,
	)

	var id string

	err := row.Scan(&id)
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
// TODO(sqlc): This requires a new sqlc query "FindProfileLinkProfileByKindAndRemoteID" that
// searches profile_link by kind + remote_id without requiring profile_id.
// Once the sqlc query is generated, replace the raw SQL below with the generated query.
func (r *Repository) GetProfileLinkByRemoteIDForResourceSync(
	ctx context.Context,
	kind string,
	remoteID string,
) (string, error) {
	row := r.dbtx.QueryRowContext(
		ctx,
		`SELECT pl.profile_id
		FROM "profile_link" pl
		WHERE pl.kind = $1
		  AND pl.remote_id = $2
		  AND pl.deleted_at IS NULL
		LIMIT 1`,
		kind,
		remoteID,
	)

	var profileID string

	err := row.Scan(&profileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return profileID, nil
}

// getProfileLinksByRemoteIDsSQL batch-loads profile_links by kind and multiple remote_ids.
// Joins with profile to filter for individual profiles only, since contributor matching
// should find the person, not organizations/products that share the same OAuth token.
const getProfileLinksByRemoteIDsSQL = `
SELECT pl.remote_id, pl.profile_id
FROM "profile_link" pl
  INNER JOIN "profile" p ON p.id = pl.profile_id AND p.kind = 'individual'
WHERE pl.kind = $1
  AND pl.remote_id = ANY($2::TEXT[])
  AND pl.deleted_at IS NULL
`

// GetProfileLinksByRemoteIDsForResourceSync batch-loads profile_links and returns map[remoteID]profileID.
func (r *Repository) GetProfileLinksByRemoteIDsForResourceSync(
	ctx context.Context,
	kind string,
	remoteIDs []string,
) (map[string]string, error) {
	rows, err := r.dbtx.QueryContext(ctx, getProfileLinksByRemoteIDsSQL, kind, pq.Array(remoteIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	result := make(map[string]string, len(remoteIDs))

	for rows.Next() {
		var remoteID, profileID string

		scanErr := rows.Scan(&remoteID, &profileID)
		if scanErr != nil {
			return nil, scanErr
		}

		result[remoteID] = profileID
	}

	return result, rows.Err()
}

// getMembershipsByProfilePairsSQL batch-loads memberships for multiple (profileID, memberProfileID) pairs.
const getMembershipsByProfilePairsSQL = `
SELECT pm.profile_id, pm.member_profile_id, pm.id
FROM "profile_membership" pm
WHERE pm.profile_id = ANY($1::TEXT[])
  AND pm.member_profile_id = ANY($2::TEXT[])
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
`

// GetMembershipsByProfilePairsForResourceSync batch-loads memberships.
// Returns map["profileID:memberProfileID"]membershipID.
func (r *Repository) GetMembershipsByProfilePairsForResourceSync(
	ctx context.Context,
	profileIDs []string,
	memberProfileIDs []string,
) (map[string]string, error) {
	rows, err := r.dbtx.QueryContext(
		ctx,
		getMembershipsByProfilePairsSQL,
		pq.Array(profileIDs),
		pq.Array(memberProfileIDs),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	result := make(map[string]string)

	for rows.Next() {
		var profileID, memberProfileID, membershipID string

		scanErr := rows.Scan(&profileID, &memberProfileID, &membershipID)
		if scanErr != nil {
			return nil, scanErr
		}

		result[profileID+":"+memberProfileID] = membershipID
	}

	return result, rows.Err()
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
