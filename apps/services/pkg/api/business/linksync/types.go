package linksync

import (
	"time"
)

// LinkImport represents an imported item from a profile link.
type LinkImport struct {
	CreatedAt     time.Time
	Properties    map[string]any
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
	ID            string
	ProfileLinkID string
	RemoteID      string
}

// ManagedLink represents a profile link with OAuth tokens for syncing.
type ManagedLink struct {
	AuthAccessTokenExpiresAt *time.Time
	AuthRefreshToken         *string
	ID                       string
	ProfileID                string
	Kind                     string
	RemoteID                 string
	AuthAccessToken          string
	IsOnline                 bool
}

// SyncResult represents the result of syncing a single link.
type SyncResult struct {
	Error        error
	LinkID       string
	ItemsAdded   int
	ItemsUpdated int
	ItemsDeleted int
}

// RemoteStoryItem represents a story item fetched from a remote provider.
type RemoteStoryItem struct {
	PublishedAt time.Time
	Properties  map[string]any
	RemoteID    string
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
	CreatedAt            time.Time
	Properties           map[string]any
	PublicationID        *string
	ID                   string
	ProfileLinkID        string
	RemoteID             string
	ProfileID            string
	ProfileDefaultLocale string
	StoryID              string
}

// PublicManagedLink represents a managed link without OAuth tokens (e.g. SpeakerDeck).
type PublicManagedLink struct {
	ID            string
	ProfileID     string
	Kind          string
	RemoteID      string
	URI           string
	ContentFolder string
}

// TokenRefreshResult contains the result of a token refresh.
type TokenRefreshResult struct {
	AccessTokenExpiresAt *time.Time
	RefreshToken         *string
	AccessToken          string
}
