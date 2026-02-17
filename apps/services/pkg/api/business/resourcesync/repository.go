package resourcesync

import "context"

// Repository defines the data access operations needed for resource syncing.
type Repository interface {
	// ListGitHubResourcesForSync returns GitHub resources with their profile's managed GitHub access tokens.
	ListGitHubResourcesForSync(ctx context.Context, batchSize int) ([]*GitHubResourceForSync, error)

	// UpdateProfileResourceProperties updates the properties JSONB of a profile_resource.
	UpdateProfileResourceProperties(ctx context.Context, id string, properties map[string]any) error

	// UpdateProfileMembershipProperties updates the properties JSONB of a profile_membership.
	UpdateProfileMembershipProperties(
		ctx context.Context,
		id string,
		properties map[string]any,
	) error

	// GetMembershipByProfiles looks up a profile_membership between two profiles.
	GetMembershipByProfiles(
		ctx context.Context,
		profileID string,
		memberProfileID string,
	) (string, error)

	// GetProfileLinkByRemoteID finds a profile_link by kind and remote_id (GitHub user ID).
	GetProfileLinkByRemoteID(ctx context.Context, kind string, remoteID string) (string, error)
}
