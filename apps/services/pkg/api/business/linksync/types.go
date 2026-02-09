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
	RemoteID    string
	PublishedAt time.Time
	Properties  map[string]any
}

// LinkImportForStoryCreation represents an import record ready for story creation.
type LinkImportForStoryCreation struct {
	ID                   string
	ProfileLinkID        string
	RemoteID             string
	Properties           map[string]any
	CreatedAt            time.Time
	ProfileID            string
	ProfileDefaultLocale string
}

// LinkImportWithStory represents an import that has a corresponding managed story (for reconciliation).
type LinkImportWithStory struct {
	ID                   string
	ProfileLinkID        string
	RemoteID             string
	Properties           map[string]any
	CreatedAt            time.Time
	ProfileID            string
	ProfileDefaultLocale string
	StoryID              string
	PublicationID        *string
}

// TokenRefreshResult contains the result of a token refresh.
type TokenRefreshResult struct {
	AccessToken          string
	AccessTokenExpiresAt *time.Time
	RefreshToken         *string
}
