package resourcesync

import "time"

// GitHubResourceForSync represents a GitHub repo resource with its associated access token.
type GitHubResourceForSync struct {
	ResourceProperties       map[string]any
	AuthAccessTokenExpiresAt *time.Time
	AuthRefreshToken         *string
	ResourceID               string
	ProfileID                string
	ResourceRemoteID         string // GitHub repo ID
	ResourcePublicID         string // "owner/repo"
	LinkID                   string
	AuthAccessToken          string
}

// GitHubContributorStats holds the GitHub contribution stats for a contributor.
type GitHubContributorStats struct {
	LastSyncedAt time.Time `json:"last_synced_at"`
	PRs          struct {
		Total    int `json:"total"`
		Resolved int `json:"resolved"`
	} `json:"prs"`
	Issues struct {
		Total    int `json:"total"`
		Resolved int `json:"resolved"`
	} `json:"issues"`
	Commits int `json:"commits"`
	Stars   int `json:"stars"`
}
