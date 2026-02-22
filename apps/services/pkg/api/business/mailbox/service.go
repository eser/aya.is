package mailbox

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

// DefaultIDGenerator returns a new unique ID.
func DefaultIDGenerator() string {
	return lib.IDsGenerateUnique()
}

// Service provides mailbox business logic.
type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator func() string
	onCreated   OnEnvelopeCreatedFunc
}

// NewService creates a new mailbox service.
func NewService(
	logger *logfx.Logger,
	repo Repository,
	idGenerator func() string,
) *Service {
	return &Service{
		logger:      logger,
		repo:        repo,
		idGenerator: idGenerator,
		onCreated:   nil,
	}
}

// SetOnCreated sets a callback that fires after each new envelope is persisted.
func (s *Service) SetOnCreated(fn OnEnvelopeCreatedFunc) {
	s.onCreated = fn
}

// GetOrCreateDirectConversation finds or creates a direct conversation between two profiles.
func (s *Service) GetOrCreateDirectConversation(
	ctx context.Context,
	senderProfileID string,
	targetProfileID string,
) (*Conversation, error) {
	// Try to find existing.
	conv, err := s.repo.FindDirectConversation(ctx, senderProfileID, targetProfileID)
	if err == nil && conv != nil {
		return conv, nil
	}

	// Create new direct conversation.
	conv = &Conversation{
		ID:                 s.idGenerator(),
		Kind:               ConversationKindDirect,
		Title:              nil,
		CreatedByProfileID: &senderProfileID,
		CreatedAt:          time.Now(),
	}

	err = s.repo.CreateConversation(ctx, conv)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreate, err)
	}

	// Add sender as participant.
	err = s.repo.AddParticipant(ctx, &Participant{
		ID:             s.idGenerator(),
		ConversationID: conv.ID,
		ProfileID:      senderProfileID,
		JoinedAt:       conv.CreatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreate, err)
	}

	// Add target as participant (skip if same as sender — self-conversation).
	if targetProfileID != senderProfileID {
		err = s.repo.AddParticipant(ctx, &Participant{
			ID:             s.idGenerator(),
			ConversationID: conv.ID,
			ProfileID:      targetProfileID,
			JoinedAt:       conv.CreatedAt,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToCreate, err)
		}
	}

	s.logger.InfoContext(ctx, "Direct conversation created",
		slog.String("conversation_id", conv.ID),
		slog.String("sender_profile_id", senderProfileID),
		slog.String("target_profile_id", targetProfileID))

	return conv, nil
}

// SendMessage sends a message within a direct conversation (creates one if needed).
func (s *Service) SendMessage(
	ctx context.Context,
	params *SendMessageParams,
) (*Envelope, error) {
	validKinds := map[string]bool{
		KindMessage: true,
	}

	if !validKinds[params.Kind] {
		return nil, ErrInvalidEnvelopeKind
	}

	conv, err := s.GetOrCreateDirectConversation(
		ctx,
		params.SenderProfileID,
		params.TargetProfileID,
	)
	if err != nil {
		return nil, err
	}

	// If the conversation already has messages, check that the first envelope has been accepted.
	// The first message acts as a conversation request; replies are blocked until accepted.
	existingEnvelopes, listErr := s.repo.ListEnvelopesByConversation(ctx, conv.ID, 1)
	if listErr == nil && len(existingEnvelopes) > 0 {
		if existingEnvelopes[0].Status != StatusAccepted {
			return nil, ErrConversationPending
		}
	}

	return s.createEnvelopeInConversation(ctx, conv.ID, params)
}

// SendSystemEnvelope creates an envelope (invitation, badge, pass) in a new conversation.
// If a sender profile is present, the conversation is "direct"; otherwise it is "system".
func (s *Service) SendSystemEnvelope(
	ctx context.Context,
	params *SendMessageParams,
) (*Envelope, error) {
	validKinds := map[string]bool{
		KindInvitation: true,
		KindMessage:    true,
		KindBadge:      true,
		KindPass:       true,
	}

	if !validKinds[params.Kind] {
		return nil, ErrInvalidEnvelopeKind
	}

	now := time.Now()

	kind := ConversationKindSystem
	if params.SenderProfileID != "" {
		kind = ConversationKindDirect
	}

	conv := &Conversation{
		ID:                 s.idGenerator(),
		Kind:               kind,
		CreatedByProfileID: &params.SenderProfileID,
		CreatedAt:          now,
	}

	err := s.repo.CreateConversation(ctx, conv)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreate, err)
	}

	// Add sender as participant first (creator appears first in ordering).
	if params.SenderProfileID != "" {
		err = s.repo.AddParticipant(ctx, &Participant{
			ID:             s.idGenerator(),
			ConversationID: conv.ID,
			ProfileID:      params.SenderProfileID,
			JoinedAt:       now,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToCreate, err)
		}
	}

	// Add target as participant.
	err = s.repo.AddParticipant(ctx, &Participant{
		ID:             s.idGenerator(),
		ConversationID: conv.ID,
		ProfileID:      params.TargetProfileID,
		JoinedAt:       now,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreate, err)
	}

	return s.createEnvelopeInConversation(ctx, conv.ID, params)
}

// createEnvelopeInConversation is the internal helper that creates an envelope.
func (s *Service) createEnvelopeInConversation(
	ctx context.Context,
	conversationID string,
	params *SendMessageParams,
) (*Envelope, error) {
	now := time.Now()

	envelope := &Envelope{
		ID:              s.idGenerator(),
		ConversationID:  conversationID,
		TargetProfileID: params.TargetProfileID,
		SenderProfileID: &params.SenderProfileID,
		SenderUserID:    params.SenderUserID,
		Kind:            params.Kind,
		Status:          StatusPending,
		Title:           params.Title,
		Description:     params.Description,
		Properties:      params.Properties,
		ReplyToID:       params.ReplyToID,
		CreatedAt:       now,
	}

	err := s.repo.CreateEnvelope(ctx, envelope)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreate, err)
	}

	// Update conversation timestamp.
	_ = s.repo.UpdateConversationTimestamp(ctx, conversationID)

	s.logger.InfoContext(ctx, "Envelope created",
		slog.String("envelope_id", envelope.ID),
		slog.String("conversation_id", conversationID),
		slog.String("kind", params.Kind),
		slog.String("target_profile_id", params.TargetProfileID))

	if s.onCreated != nil {
		s.onCreated(ctx, envelope, params)
	}

	return envelope, nil
}

// ListConversations returns conversations for a profile.
func (s *Service) ListConversations(
	ctx context.Context,
	profileID string,
	includeArchived bool,
	limit int,
) ([]*Conversation, error) {
	if limit <= 0 {
		limit = 50
	}

	conversations, err := s.repo.ListConversationsForProfile(ctx, profileID, includeArchived, limit)
	if err != nil {
		return nil, fmt.Errorf("list conversations: %w", err)
	}

	return conversations, nil
}

// GetConversation returns a conversation with full details.
func (s *Service) GetConversation(
	ctx context.Context,
	conversationID string,
	requestingProfileID string,
) (*Conversation, []*Envelope, error) {
	conv, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, nil, ErrConversationNotFound
	}

	// Verify requester is a participant.
	_, participantErr := s.repo.GetParticipant(ctx, conversationID, requestingProfileID)
	if participantErr != nil {
		return nil, nil, ErrNotParticipant
	}

	// List participants.
	participants, err := s.repo.ListParticipants(ctx, conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("list participants: %w", err)
	}

	conv.Participants = participants

	// List envelopes.
	const defaultEnvelopeLimit = 100

	envelopeList, err := s.repo.ListEnvelopesByConversation(
		ctx,
		conversationID,
		defaultEnvelopeLimit,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("list envelopes: %w", err)
	}

	// Populate reactions for each envelope.
	for _, env := range envelopeList {
		reactions, reactErr := s.repo.ListReactionsByEnvelope(ctx, env.ID)
		if reactErr == nil {
			env.Reactions = reactions
		}
	}

	return conv, envelopeList, nil
}

// GetConversationParticipants returns participants for a conversation.
func (s *Service) GetConversationParticipants(
	ctx context.Context,
	conversationID string,
) ([]*Participant, error) {
	participants, err := s.repo.ListParticipants(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("list participants: %w", err)
	}

	return participants, nil
}

// MarkConversationRead sets the read cursor to now.
func (s *Service) MarkConversationRead(
	ctx context.Context,
	conversationID string,
	profileID string,
) error {
	return s.repo.UpdateParticipantReadCursor(ctx, conversationID, profileID)
}

// ArchiveConversation archives a conversation for a participant.
func (s *Service) ArchiveConversation(
	ctx context.Context,
	conversationID string,
	profileID string,
) error {
	return s.repo.SetParticipantArchived(ctx, conversationID, profileID, true)
}

// UnarchiveConversation unarchives a conversation for a participant.
func (s *Service) UnarchiveConversation(
	ctx context.Context,
	conversationID string,
	profileID string,
) error {
	return s.repo.SetParticipantArchived(ctx, conversationID, profileID, false)
}

// AcceptEnvelope transitions an envelope from pending to accepted.
func (s *Service) AcceptEnvelope(
	ctx context.Context,
	envelopeID string,
	targetProfileID string,
) error {
	envelope, err := s.repo.GetEnvelopeByID(ctx, envelopeID)
	if err != nil {
		return ErrEnvelopeNotFound
	}

	if envelope.TargetProfileID != targetProfileID {
		return ErrNotTargetProfile
	}

	if envelope.Status != StatusPending {
		return ErrAlreadyProcessed
	}

	err = s.repo.UpdateEnvelopeStatus(ctx, envelopeID, StatusAccepted, time.Now())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdate, err)
	}

	s.logger.InfoContext(ctx, "Envelope accepted",
		slog.String("envelope_id", envelopeID))

	return nil
}

// RejectEnvelope transitions an envelope from pending to rejected.
func (s *Service) RejectEnvelope(
	ctx context.Context,
	envelopeID string,
	targetProfileID string,
) error {
	envelope, err := s.repo.GetEnvelopeByID(ctx, envelopeID)
	if err != nil {
		return ErrEnvelopeNotFound
	}

	if envelope.TargetProfileID != targetProfileID {
		return ErrNotTargetProfile
	}

	if envelope.Status != StatusPending {
		return ErrAlreadyProcessed
	}

	err = s.repo.UpdateEnvelopeStatus(ctx, envelopeID, StatusRejected, time.Now())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdate, err)
	}

	s.logger.InfoContext(ctx, "Envelope rejected",
		slog.String("envelope_id", envelopeID))

	return nil
}

// RevokeEnvelope transitions an envelope from pending to revoked (sender action).
func (s *Service) RevokeEnvelope(ctx context.Context, envelopeID string) error {
	envelope, err := s.repo.GetEnvelopeByID(ctx, envelopeID)
	if err != nil {
		return ErrEnvelopeNotFound
	}

	if envelope.Status != StatusPending {
		return ErrAlreadyProcessed
	}

	err = s.repo.UpdateEnvelopeStatus(ctx, envelopeID, StatusRevoked, time.Now())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdate, err)
	}

	s.logger.InfoContext(ctx, "Envelope revoked",
		slog.String("envelope_id", envelopeID))

	return nil
}

// RedeemEnvelope transitions an envelope from accepted to redeemed and updates properties.
func (s *Service) RedeemEnvelope(
	ctx context.Context,
	envelopeID string,
	updatedProperties any,
) error {
	envelope, err := s.repo.GetEnvelopeByID(ctx, envelopeID)
	if err != nil {
		return ErrEnvelopeNotFound
	}

	if envelope.Status != StatusAccepted {
		return ErrInvalidStatus
	}

	if updatedProperties != nil {
		err = s.repo.UpdateEnvelopeProperties(ctx, envelopeID, updatedProperties)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrFailedToUpdate, err)
		}
	}

	err = s.repo.UpdateEnvelopeStatus(ctx, envelopeID, StatusRedeemed, time.Now())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdate, err)
	}

	s.logger.InfoContext(ctx, "Envelope redeemed",
		slog.String("envelope_id", envelopeID))

	return nil
}

// AddReaction adds an emoji reaction to an envelope.
func (s *Service) AddReaction(
	ctx context.Context,
	envelopeID string,
	profileID string,
	emoji string,
) error {
	if !AllowedReactions[emoji] {
		return ErrInvalidReaction
	}

	envelope, err := s.repo.GetEnvelopeByID(ctx, envelopeID)
	if err != nil {
		return ErrEnvelopeNotFound
	}

	if envelope.Status != StatusAccepted && envelope.Status != StatusRedeemed {
		return ErrInvalidStatus
	}

	reaction := &Reaction{
		ID:         s.idGenerator(),
		EnvelopeID: envelopeID,
		ProfileID:  profileID,
		Emoji:      emoji,
		CreatedAt:  time.Now(),
	}

	return s.repo.AddReaction(ctx, reaction)
}

// RemoveReaction removes an emoji reaction from an envelope.
func (s *Service) RemoveReaction(
	ctx context.Context,
	envelopeID string,
	profileID string,
	emoji string,
) error {
	return s.repo.RemoveReaction(ctx, envelopeID, profileID, emoji)
}

// GetEnvelopeByID returns a single envelope by ID.
func (s *Service) GetEnvelopeByID(ctx context.Context, id string) (*Envelope, error) {
	envelope, err := s.repo.GetEnvelopeByID(ctx, id)
	if err != nil {
		return nil, ErrEnvelopeNotFound
	}

	return envelope, nil
}

// ListEnvelopes returns envelopes for a target profile, optionally filtered by status.
func (s *Service) ListEnvelopes(
	ctx context.Context,
	targetProfileID string,
	statusFilter string,
) ([]*Envelope, error) {
	const defaultLimit = 50

	envelopes, err := s.repo.ListEnvelopesByTargetProfileID(
		ctx,
		targetProfileID,
		statusFilter,
		defaultLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("list envelopes: %w", err)
	}

	return envelopes, nil
}

// GetAcceptedInvitations returns accepted invitation envelopes for a profile.
func (s *Service) GetAcceptedInvitations(
	ctx context.Context,
	targetProfileID string,
	invitationKind string,
) ([]*Envelope, error) {
	envelopes, err := s.repo.ListAcceptedInvitations(ctx, targetProfileID, invitationKind)
	if err != nil {
		return nil, fmt.Errorf("list accepted invitations: %w", err)
	}

	return envelopes, nil
}

// CountPendingEnvelopes returns the count of pending envelopes for a profile.
func (s *Service) CountPendingEnvelopes(
	ctx context.Context,
	targetProfileID string,
) (int, error) {
	count, err := s.repo.CountPendingEnvelopes(ctx, targetProfileID)
	if err != nil {
		return 0, fmt.Errorf("count pending envelopes: %w", err)
	}

	return count, nil
}

// CountUnreadConversations returns the total unread conversation count for a profile.
func (s *Service) CountUnreadConversations(
	ctx context.Context,
	profileID string,
) (int, error) {
	count, err := s.repo.CountUnreadConversations(ctx, profileID)
	if err != nil {
		return 0, fmt.Errorf("count unread conversations: %w", err)
	}

	return count, nil
}
