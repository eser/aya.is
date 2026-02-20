package telegram

import "time"

// ExternalCode represents a short-lived code tied to an external system (e.g. Telegram, Discord).
// All platform-specific data lives in the Properties map (JSONB in the database).
type ExternalCode struct {
	ID             string
	Code           string
	ExternalSystem string
	Properties     map[string]any
	CreatedAt      time.Time
	ExpiresAt      time.Time
	ConsumedAt     *time.Time
}

// LinkResult contains the result of a successful account linking.
type LinkResult struct {
	ProfileID        string
	ProfileSlug      string
	TelegramUserID   int64
	TelegramUsername string
}

// ProfileLinkInfo represents a linked Telegram account.
type ProfileLinkInfo struct {
	ID        string
	ProfileID string
	RemoteID  string
	PublicID  string
}

// RawGroupTelegramLink is an intermediate type returned from the repository,
// before visibility filtering is applied in the service layer.
type RawGroupTelegramLink struct {
	ProfileSlug    string
	ProfileTitle   string
	MembershipKind string
	LinkTitle      string
	LinkURI        string
	LinkPublicID   string
	LinkVisibility string
}

// GroupTelegramLink represents a Telegram link from a group profile visible to the user.
type GroupTelegramLink struct {
	ProfileSlug  string
	ProfileTitle string
	LinkTitle    string
	LinkURI      string
	LinkPublicID string
}

// CreateProfileLinkParams contains parameters for creating a telegram profile link.
type CreateProfileLinkParams struct {
	ID               string
	ProfileID        string
	RemoteID         string // Telegram user ID as string
	PublicID         string // Telegram username
	URI              string // https://t.me/<username>
	Order            int
	AddedByProfileID string
}
