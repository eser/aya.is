package mailbox

import (
	"context"
	"time"
)

// Conversation kinds.
const (
	ConversationKindDirect = "direct"
	ConversationKindSystem = "system"
)

// Envelope kinds.
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

// AllowedReactions is the curated set of emoji reactions.
var AllowedReactions = map[string]bool{
	"👍":  true,
	"❤️": true,
	"😂":  true,
	"😮":  true,
	"😢":  true,
	"🔥":  true,
	"🎉":  true,
}

// Conversation represents a messaging thread between participants.
type Conversation struct {
	ID                 string     `json:"id"`
	Kind               string     `json:"kind"`
	Title              *string    `json:"title"`
	CreatedByProfileID *string    `json:"created_by_profile_id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`

	// Populated via queries.
	Participants []*Participant   `json:"participants,omitempty"`
	LastEnvelope *EnvelopePreview `json:"last_envelope,omitempty"`
	UnreadCount  int              `json:"unread_count"`
	IsArchived   bool             `json:"is_archived"`
}

// EnvelopePreview is a summary of the latest envelope in a conversation.
type EnvelopePreview struct {
	Message         *string `json:"message"`
	Kind            string  `json:"kind"`
	CreatedAt       string  `json:"created_at"`
	SenderProfileID *string `json:"sender_profile_id"`
}

// Participant represents a profile's membership in a conversation.
type Participant struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversation_id"`
	ProfileID      string     `json:"profile_id"`
	LastReadAt     *time.Time `json:"last_read_at"`
	IsArchived     bool       `json:"is_archived"`
	JoinedAt       time.Time  `json:"joined_at"`
	LeftAt         *time.Time `json:"left_at"`

	// Populated via JOINs.
	ProfileSlug       *string `json:"profile_slug,omitempty"`
	ProfileTitle      *string `json:"profile_title,omitempty"`
	ProfilePictureURI *string `json:"profile_picture_uri,omitempty"`
	ProfileKind       *string `json:"profile_kind,omitempty"`
}

// Envelope represents an inbox item (message, invitation, badge, pass).
type Envelope struct {
	ID              string     `json:"id"`
	ConversationID  string     `json:"conversation_id"`
	TargetProfileID string     `json:"target_profile_id"`
	SenderProfileID *string    `json:"sender_profile_id"`
	SenderUserID    *string    `json:"sender_user_id"`
	Kind            string     `json:"kind"`
	Status          string     `json:"status"`
	Message         *string    `json:"message"`
	Properties      any        `json:"properties"`
	ReplyToID       *string    `json:"reply_to_id"`
	AcceptedAt      *time.Time `json:"accepted_at"`
	RejectedAt      *time.Time `json:"rejected_at"`
	RevokedAt       *time.Time `json:"revoked_at"`
	RedeemedAt      *time.Time `json:"redeemed_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`

	// Populated via JOINs.
	SenderProfileSlug       *string     `json:"sender_profile_slug,omitempty"`
	SenderProfileTitle      *string     `json:"sender_profile_title,omitempty"`
	SenderProfilePictureURI *string     `json:"sender_profile_picture_uri,omitempty"`
	SenderProfileKind       *string     `json:"sender_profile_kind,omitempty"`
	Reactions               []*Reaction `json:"reactions,omitempty"`
}

// Reaction represents an emoji reaction on an envelope.
type Reaction struct {
	ID         string    `json:"id"`
	EnvelopeID string    `json:"envelope_id"`
	ProfileID  string    `json:"profile_id"`
	Emoji      string    `json:"emoji"`
	CreatedAt  time.Time `json:"created_at"`

	ProfileSlug  *string `json:"profile_slug,omitempty"`
	ProfileTitle *string `json:"profile_title,omitempty"`
}

// InvitationProperties holds kind-specific data for invitation envelopes.
type InvitationProperties struct {
	InvitationKind   string  `json:"invitation_kind"`
	TelegramChatID   int64   `json:"telegram_chat_id,omitempty"`
	GroupProfileSlug string  `json:"group_profile_slug,omitempty"`
	GroupName        string  `json:"group_name,omitempty"`
	InviteLink       *string `json:"invite_link,omitempty"`
}

// SendMessageParams contains parameters for sending a message or system envelope.
type SendMessageParams struct {
	SenderProfileID   string
	TargetProfileID   string
	SenderUserID      *string
	Kind              string
	ConversationTitle string
	Message           *string
	Properties        any
	ReplyToID         *string

	// Notification context (optional, used for notifying the recipient).
	SenderProfileTitle string
	Locale             string
}

// OnEnvelopeCreatedFunc is a callback invoked after a new envelope is persisted.
// Implementations must be best-effort — failures should be logged, never propagated.
type OnEnvelopeCreatedFunc func(ctx context.Context, envelope *Envelope, params *SendMessageParams)

// OnEnvelopeAcceptedFunc is a callback invoked after an envelope is accepted.
// Implementations must be best-effort — failures should be logged, never propagated.
type OnEnvelopeAcceptedFunc func(ctx context.Context, envelope *Envelope)
