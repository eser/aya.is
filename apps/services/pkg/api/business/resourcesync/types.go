package resourcesync

import "time"

// GitHubResourceForSync represents a GitHub repo resource with its associated access token.
type GitHubResourceForSync struct {
	ResourceID               string
	ProfileID                string
	ResourceRemoteID         string // GitHub repo ID
	ResourcePublicID         string // "owner/repo"
	ResourceProperties       map[string]any
	LinkID                   string
	AuthAccessToken          string
	AuthAccessTokenExpiresAt *time.Time
	AuthRefreshToken         *string
}

// GitHubContributorStats holds the GitHub contribution stats for a contributor.
type GitHubContributorStats struct {
	Commits int `json:"commits"`
	PRs     struct {
		Total    int `json:"total"`
		Resolved int `json:"resolved"`
	} `json:"prs"`
	Issues struct {
		Total    int `json:"total"`
		Resolved int `json:"resolved"`
	} `json:"issues"`
	Stars int `json:"stars"`
}
