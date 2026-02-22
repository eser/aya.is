package mailbox

import (
	"context"
	"time"
)

// Repository defines storage operations for the mailbox service.
type Repository interface {
	// Conversations
	CreateConversation(ctx context.Context, conv *Conversation) error
	GetConversationByID(ctx context.Context, id string) (*Conversation, error)
	FindDirectConversation(
		ctx context.Context,
		profileA string,
		profileB string,
	) (*Conversation, error)
	ListConversationsForProfile(
		ctx context.Context,
		profileID string,
		includeArchived bool,
		limit int,
	) ([]*Conversation, error)
	UpdateConversationTimestamp(ctx context.Context, id string) error
	RemoveConversation(ctx context.Context, conversationID string) error

	// Participants
	AddParticipant(ctx context.Context, participant *Participant) error
	GetParticipant(
		ctx context.Context,
		conversationID string,
		profileID string,
	) (*Participant, error)
	ListParticipants(ctx context.Context, conversationID string) ([]*Participant, error)
	UpdateParticipantReadCursor(ctx context.Context, conversationID string, profileID string) error
	SetParticipantArchived(
		ctx context.Context,
		conversationID string,
		profileID string,
		archived bool,
	) error

	// Envelopes
	CreateEnvelope(ctx context.Context, envelope *Envelope) error
	GetEnvelopeByID(ctx context.Context, id string) (*Envelope, error)
	ListEnvelopesByConversation(
		ctx context.Context,
		conversationID string,
		limit int,
	) ([]*Envelope, error)
	ListEnvelopesByTargetProfileID(
		ctx context.Context,
		profileID string,
		statusFilter string,
		limit int,
	) ([]*Envelope, error)
	UpdateEnvelopeStatus(ctx context.Context, id string, status string, now time.Time) error
	UpdateEnvelopeProperties(ctx context.Context, id string, properties any) error
	ListAcceptedInvitations(
		ctx context.Context,
		targetProfileID string,
		invitationKind string,
	) ([]*Envelope, error)
	CountPendingEnvelopes(ctx context.Context, targetProfileID string) (int, error)

	// Reactions
	AddReaction(ctx context.Context, reaction *Reaction) error
	RemoveReaction(ctx context.Context, envelopeID string, profileID string, emoji string) error
	ListReactionsByEnvelope(ctx context.Context, envelopeID string) ([]*Reaction, error)

	// Counts
	CountUnreadConversations(ctx context.Context, profileID string) (int, error)
}
