package linksync

import (
	"time"
)

// LinkImport represents an imported item from a profile link.
type LinkImport struct {
	ID            string
	ProfileLinkID string
	RemoteID      string
	Properties    map[string]any
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
}

// ManagedLink represents a profile link with OAuth tokens for syncing.
type ManagedLink struct {
	ID                       string
	ProfileID                string
	Kind                     string
	RemoteID                 string
	AuthAccessToken          string
	AuthAccessTokenExpiresAt *time.Time
	AuthRefreshToken         *string
}

// SyncResult represents the result of syncing a single link.
type SyncResult struct {
	LinkID       string
	ItemsAdded   int
	ItemsUpdated int
	ItemsDeleted int
	Error        error
}

// RemoteStoryItem represents a story item fetched from a remote provider.
type RemoteStoryItem struct {
	RemoteID     string
	Title        string
	Description  string
	PublishedAt  time.Time
	ThumbnailURL string
	Duration     string
	ViewCount    int64
	LikeCount    int64
	Properties   map[string]any
}

// TokenRefreshResult contains the result of a token refresh.
type TokenRefreshResult struct {
	AccessToken          string
	AccessTokenExpiresAt *time.Time
	RefreshToken         *string
}
