package telegram

import "time"

// TelegramLinkToken represents a pending link token for account connection.
type TelegramLinkToken struct {
	ID              string
	Token           string
	ProfileID       string
	ProfileSlug     string
	CreatedByUserID string
	CreatedAt       time.Time
	ExpiresAt       time.Time
	ConsumedAt      *time.Time
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
