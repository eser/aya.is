package profile_envelopes

import "time"

// Envelope kinds (broad categories).
const (
	KindInvitation = "invitation"
	KindMessage    = "message"
	KindBadge      = "badge"
	KindPass       = "pass"
)

// Envelope statuses.
const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
	StatusRevoked  = "revoked"
	StatusRedeemed = "redeemed"
)

// Invitation sub-types (stored in properties.invitation_kind).
const (
	InvitationKindTelegramGroup = "telegram_group"
)

// Envelope represents an inbox item sent to a profile.
type Envelope struct {
	ID              string     `json:"id"`
	TargetProfileID string     `json:"target_profile_id"`
	SenderProfileID *string    `json:"sender_profile_id"`
	SenderUserID    *string    `json:"sender_user_id"`
	Kind            string     `json:"kind"`
	Status          string     `json:"status"`
	Title           string     `json:"title"`
	Description     *string    `json:"description"`
	Properties      any        `json:"properties"`
	AcceptedAt      *time.Time `json:"accepted_at"`
	RejectedAt      *time.Time `json:"rejected_at"`
	RevokedAt       *time.Time `json:"revoked_at"`
	RedeemedAt      *time.Time `json:"redeemed_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`

	// Sender profile info (populated from JOINs in list queries).
	SenderProfileSlug       *string `json:"sender_profile_slug"`
	SenderProfileTitle      *string `json:"sender_profile_title"`
	SenderProfilePictureURI *string `json:"sender_profile_picture_uri"`
	SenderProfileKind       *string `json:"sender_profile_kind"`
}

// InvitationProperties holds kind-specific data for invitation envelopes.
type InvitationProperties struct {
	InvitationKind   string  `json:"invitation_kind"`
	TelegramChatID   int64   `json:"telegram_chat_id,omitempty"`
	GroupProfileSlug string  `json:"group_profile_slug,omitempty"`
	GroupName        string  `json:"group_name,omitempty"`
	InviteLink       *string `json:"invite_link,omitempty"`
}

// CreateEnvelopeParams contains parameters for creating an envelope.
type CreateEnvelopeParams struct {
	TargetProfileID string
	SenderProfileID *string
	SenderUserID    *string
	Kind            string
	Title           string
	Description     *string
	Properties      any
}
