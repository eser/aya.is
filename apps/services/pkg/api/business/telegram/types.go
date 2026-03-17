package telegram

import "time"

// ExternalCode represents a short-lived code tied to an external system (e.g. Telegram, Discord).
// All platform-specific data lives in the Properties map (JSONB in the database).
type ExternalCode struct {
	CreatedAt      time.Time
	ExpiresAt      time.Time
	Properties     map[string]any
	ConsumedAt     *time.Time
	ID             string
	Code           string
	ExternalSystem string
}

// LinkResult contains the result of a successful account linking.
type LinkResult struct {
	ProfileID        string
	ProfileSlug      string
	TelegramUsername string
	TelegramUserID   int64
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
	ProfileSlug    string
	ProfileTitle   string
	LinkTitle      string
	LinkURI        string
	LinkPublicID   string
	LinkVisibility string
}

// GroupRegisterCodeData contains the resolved data from a group registration code.
type GroupRegisterCodeData struct {
	ChatTitle string
	ChatID    int64
}

// EligibleTelegramGroup represents a Telegram group resource the user can join via team membership.
type EligibleTelegramGroup struct {
	ResourceID   string
	GroupTitle   string
	ProfileSlug  string
	ProfileTitle string
	ChatID       int64
}

// CreateProfileLinkParams contains parameters for creating a telegram profile link.
type CreateProfileLinkParams struct {
	ID               string
	ProfileID        string
	RemoteID         string // Telegram user ID as string
	PublicID         string // Telegram username
	URI              string // https://t.me/<username>
	Visibility       string // Link visibility (public, followers, etc.)
	AddedByProfileID string
	Order            int
}
