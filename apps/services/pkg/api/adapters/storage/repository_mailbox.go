package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/mailbox"
	"github.com/eser/aya.is/services/pkg/lib/vars"
	"github.com/sqlc-dev/pqtype"
)

type mailboxAdapter struct {
	repo *Repository
}

// NewMailboxRepository creates a new adapter that implements mailbox.Repository.
func NewMailboxRepository(repo *Repository) mailbox.Repository {
	return &mailboxAdapter{repo: repo}
}

// ── Conversations ──────────────────────────────────────────────────────

func (a *mailboxAdapter) CreateConversation(ctx context.Context, conv *mailbox.Conversation) error {
	return a.repo.queries.CreateMailboxConversation(ctx, CreateMailboxConversationParams{
		ID:                 conv.ID,
		Kind:               conv.Kind,
		Title:              vars.ToSQLNullString(conv.Title),
		CreatedByProfileID: vars.ToSQLNullString(conv.CreatedByProfileID),
	})
}

func (a *mailboxAdapter) GetConversationByID(
	ctx context.Context,
	id string,
) (*mailbox.Conversation, error) {
	row, err := a.repo.queries.GetMailboxConversationByID(
		ctx,
		GetMailboxConversationByIDParams{ID: id},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mailbox.ErrConversationNotFound
		}

		return nil, err
	}

	return conversationFromRow(row), nil
}

func (a *mailboxAdapter) FindDirectConversation(
	ctx context.Context,
	profileA string,
	profileB string,
) (*mailbox.Conversation, error) {
	row, err := a.repo.queries.FindDirectConversation(ctx, FindDirectConversationParams{
		ProfileA: profileA,
		ProfileB: profileB,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mailbox.ErrConversationNotFound
		}

		return nil, err
	}

	return conversationFromRow(row), nil
}

func (a *mailboxAdapter) ListConversationsForProfile(
	ctx context.Context,
	profileID string,
	includeArchived bool,
	limit int,
) ([]*mailbox.Conversation, error) {
	rows, err := a.repo.queries.ListConversationsForProfile(ctx, ListConversationsForProfileParams{
		ProfileID:       profileID,
		IncludeArchived: includeArchived,
		LimitCount:      int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*mailbox.Conversation, 0, len(rows))

	for _, row := range rows {
		conv := &mailbox.Conversation{
			ID:                 row.ID,
			Kind:               row.Kind,
			Title:              vars.ToStringPtr(row.Title),
			CreatedByProfileID: vars.ToStringPtr(row.CreatedByProfileID),
			CreatedAt:          row.CreatedAt,
			UnreadCount:        int(row.UnreadCount),
			IsArchived:         row.IsArchived,
		}

		if row.UpdatedAt.Valid {
			conv.UpdatedAt = &row.UpdatedAt.Time
		}

		if row.LastEnvelopeTitle != "" {
			senderID := vars.ToStringPtr(row.LastEnvelopeSenderProfileID)
			conv.LastEnvelope = &mailbox.EnvelopePreview{
				Title:           row.LastEnvelopeTitle,
				Kind:            row.LastEnvelopeKind,
				CreatedAt:       row.LastEnvelopeAt.Format(time.RFC3339),
				SenderProfileID: senderID,
			}
		}

		result = append(result, conv)
	}

	return result, nil
}

func (a *mailboxAdapter) UpdateConversationTimestamp(ctx context.Context, id string) error {
	return a.repo.queries.UpdateConversationTimestamp(
		ctx,
		UpdateConversationTimestampParams{ID: id},
	)
}

// ── Participants ───────────────────────────────────────────────────────

func (a *mailboxAdapter) AddParticipant(
	ctx context.Context,
	participant *mailbox.Participant,
) error {
	return a.repo.queries.AddMailboxParticipant(ctx, AddMailboxParticipantParams{
		ID:             participant.ID,
		ConversationID: participant.ConversationID,
		ProfileID:      participant.ProfileID,
	})
}

func (a *mailboxAdapter) GetParticipant(
	ctx context.Context,
	conversationID string,
	profileID string,
) (*mailbox.Participant, error) {
	row, err := a.repo.queries.GetMailboxParticipant(ctx, GetMailboxParticipantParams{
		ConversationID: conversationID,
		ProfileID:      profileID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mailbox.ErrNotParticipant
		}

		return nil, err
	}

	return participantFromGetRow(row), nil
}

func (a *mailboxAdapter) ListParticipants(
	ctx context.Context,
	conversationID string,
) ([]*mailbox.Participant, error) {
	rows, err := a.repo.queries.ListMailboxParticipants(ctx, ListMailboxParticipantsParams{
		ConversationID: conversationID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*mailbox.Participant, 0, len(rows))

	for _, row := range rows {
		result = append(result, participantFromListRow(row))
	}

	return result, nil
}

func (a *mailboxAdapter) UpdateParticipantReadCursor(
	ctx context.Context,
	conversationID string,
	profileID string,
) error {
	return a.repo.queries.UpdateParticipantReadCursor(ctx, UpdateParticipantReadCursorParams{
		ConversationID: conversationID,
		ProfileID:      profileID,
	})
}

func (a *mailboxAdapter) SetParticipantArchived(
	ctx context.Context,
	conversationID string,
	profileID string,
	archived bool,
) error {
	return a.repo.queries.SetParticipantArchived(ctx, SetParticipantArchivedParams{
		ConversationID: conversationID,
		ProfileID:      profileID,
		IsArchived:     archived,
	})
}

// ── Envelopes ──────────────────────────────────────────────────────────

func (a *mailboxAdapter) CreateEnvelope(ctx context.Context, envelope *mailbox.Envelope) error {
	var propsJSON pqtype.NullRawMessage

	if envelope.Properties != nil {
		data, err := json.Marshal(envelope.Properties)
		if err != nil {
			return err
		}

		propsJSON = pqtype.NullRawMessage{RawMessage: data, Valid: true}
	}

	return a.repo.queries.CreateMailboxEnvelope(ctx, CreateMailboxEnvelopeParams{
		ID:              envelope.ID,
		ConversationID:  envelope.ConversationID,
		TargetProfileID: envelope.TargetProfileID,
		SenderProfileID: vars.ToSQLNullString(envelope.SenderProfileID),
		SenderUserID:    vars.ToSQLNullString(envelope.SenderUserID),
		Kind:            envelope.Kind,
		Title:           envelope.Title,
		Description:     vars.ToSQLNullString(envelope.Description),
		Properties:      propsJSON,
		ReplyToID:       vars.ToSQLNullString(envelope.ReplyToID),
	})
}

func (a *mailboxAdapter) GetEnvelopeByID(
	ctx context.Context,
	id string,
) (*mailbox.Envelope, error) {
	row, err := a.repo.queries.GetMailboxEnvelopeByID(ctx, GetMailboxEnvelopeByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mailbox.ErrEnvelopeNotFound
		}

		return nil, err
	}

	return envelopeFromRawRow(row), nil
}

func (a *mailboxAdapter) ListEnvelopesByConversation(
	ctx context.Context,
	conversationID string,
	limit int,
) ([]*mailbox.Envelope, error) {
	rows, err := a.repo.queries.ListEnvelopesByConversation(ctx, ListEnvelopesByConversationParams{
		ConversationID: conversationID,
		LimitCount:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*mailbox.Envelope, 0, len(rows))

	for _, row := range rows {
		result = append(result, envelopeFromConversationRow(row))
	}

	return result, nil
}

func (a *mailboxAdapter) ListEnvelopesByTargetProfileID(
	ctx context.Context,
	profileID string,
	statusFilter string,
	limit int,
) ([]*mailbox.Envelope, error) {
	var statusNullStr sql.NullString
	if statusFilter != "" {
		statusNullStr = sql.NullString{String: statusFilter, Valid: true}
	}

	rows, err := a.repo.queries.ListMailboxEnvelopesByTargetProfileID(
		ctx,
		ListMailboxEnvelopesByTargetProfileIDParams{
			TargetProfileID: profileID,
			StatusFilter:    statusNullStr,
			LimitCount:      int32(limit),
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*mailbox.Envelope, 0, len(rows))

	for _, row := range rows {
		result = append(result, envelopeFromListRow(row))
	}

	return result, nil
}

func (a *mailboxAdapter) UpdateEnvelopeStatus(
	ctx context.Context,
	id string,
	status string,
	now time.Time,
) error {
	nullNow := sql.NullTime{Time: now, Valid: true}

	var rowsAffected int64

	var err error

	switch status {
	case mailbox.StatusAccepted:
		rowsAffected, err = a.repo.queries.UpdateMailboxEnvelopeStatusToAccepted(
			ctx, UpdateMailboxEnvelopeStatusToAcceptedParams{Now: nullNow, ID: id},
		)
	case mailbox.StatusRejected:
		rowsAffected, err = a.repo.queries.UpdateMailboxEnvelopeStatusToRejected(
			ctx, UpdateMailboxEnvelopeStatusToRejectedParams{Now: nullNow, ID: id},
		)
	case mailbox.StatusRevoked:
		rowsAffected, err = a.repo.queries.UpdateMailboxEnvelopeStatusToRevoked(
			ctx, UpdateMailboxEnvelopeStatusToRevokedParams{Now: nullNow, ID: id},
		)
	case mailbox.StatusRedeemed:
		rowsAffected, err = a.repo.queries.UpdateMailboxEnvelopeStatusToRedeemed(
			ctx, UpdateMailboxEnvelopeStatusToRedeemedParams{Now: nullNow, ID: id},
		)
	default:
		return mailbox.ErrInvalidStatus
	}

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return mailbox.ErrAlreadyProcessed
	}

	return nil
}

func (a *mailboxAdapter) UpdateEnvelopeProperties(
	ctx context.Context,
	id string,
	properties any,
) error {
	data, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	return a.repo.queries.UpdateMailboxEnvelopeProperties(
		ctx,
		UpdateMailboxEnvelopePropertiesParams{
			ID:         id,
			Properties: pqtype.NullRawMessage{RawMessage: data, Valid: true},
		},
	)
}

func (a *mailboxAdapter) ListAcceptedInvitations(
	ctx context.Context,
	targetProfileID string,
	invitationKind string,
) ([]*mailbox.Envelope, error) {
	var kindFilter sql.NullString
	if invitationKind != "" {
		kindFilter = sql.NullString{String: invitationKind, Valid: true}
	}

	rows, err := a.repo.queries.ListAcceptedMailboxInvitations(
		ctx,
		ListAcceptedMailboxInvitationsParams{
			TargetProfileID: targetProfileID,
			InvitationKind:  kindFilter,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*mailbox.Envelope, 0, len(rows))

	for _, row := range rows {
		result = append(result, envelopeFromInvitationRow(row))
	}

	return result, nil
}

func (a *mailboxAdapter) CountPendingEnvelopes(
	ctx context.Context,
	targetProfileID string,
) (int, error) {
	count, err := a.repo.queries.CountPendingMailboxEnvelopes(
		ctx,
		CountPendingMailboxEnvelopesParams{TargetProfileID: targetProfileID},
	)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// ── Reactions ──────────────────────────────────────────────────────────

func (a *mailboxAdapter) AddReaction(ctx context.Context, reaction *mailbox.Reaction) error {
	return a.repo.queries.AddMailboxReaction(ctx, AddMailboxReactionParams{
		ID:         reaction.ID,
		EnvelopeID: reaction.EnvelopeID,
		ProfileID:  reaction.ProfileID,
		Emoji:      reaction.Emoji,
	})
}

func (a *mailboxAdapter) RemoveReaction(
	ctx context.Context,
	envelopeID string,
	profileID string,
	emoji string,
) error {
	_, err := a.repo.queries.RemoveMailboxReaction(ctx, RemoveMailboxReactionParams{
		EnvelopeID: envelopeID,
		ProfileID:  profileID,
		Emoji:      emoji,
	})

	return err
}

func (a *mailboxAdapter) ListReactionsByEnvelope(
	ctx context.Context,
	envelopeID string,
) ([]*mailbox.Reaction, error) {
	rows, err := a.repo.queries.ListReactionsByEnvelope(ctx, ListReactionsByEnvelopeParams{
		EnvelopeID: envelopeID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*mailbox.Reaction, 0, len(rows))

	for _, row := range rows {
		slug := row.ProfileSlug
		title := row.ProfileTitle

		result = append(result, &mailbox.Reaction{
			ID:           row.ID,
			EnvelopeID:   row.EnvelopeID,
			ProfileID:    row.ProfileID,
			Emoji:        row.Emoji,
			CreatedAt:    row.CreatedAt,
			ProfileSlug:  &slug,
			ProfileTitle: &title,
		})
	}

	return result, nil
}

// ── Counts ─────────────────────────────────────────────────────────────

func (a *mailboxAdapter) CountUnreadConversations(
	ctx context.Context,
	profileID string,
) (int, error) {
	count, err := a.repo.queries.CountUnreadConversations(ctx, CountUnreadConversationsParams{
		ProfileID: profileID,
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// ── Row converters ─────────────────────────────────────────────────────

func conversationFromRow(row *MailboxConversation) *mailbox.Conversation {
	conv := &mailbox.Conversation{
		ID:                 row.ID,
		Kind:               row.Kind,
		Title:              vars.ToStringPtr(row.Title),
		CreatedByProfileID: vars.ToStringPtr(row.CreatedByProfileID),
		CreatedAt:          row.CreatedAt,
	}

	if row.UpdatedAt.Valid {
		conv.UpdatedAt = &row.UpdatedAt.Time
	}

	return conv
}

func participantFromGetRow(row *GetMailboxParticipantRow) *mailbox.Participant {
	slug := row.ProfileSlug
	title := row.ProfileTitle
	kind := row.ProfileKind

	p := &mailbox.Participant{
		ID:             row.ID,
		ConversationID: row.ConversationID,
		ProfileID:      row.ProfileID,
		IsArchived:     row.IsArchived,
		JoinedAt:       row.JoinedAt,
		ProfileSlug:    &slug,
		ProfileTitle:   &title,
		ProfileKind:    &kind,
	}

	if row.LastReadAt.Valid {
		p.LastReadAt = &row.LastReadAt.Time
	}

	if row.LeftAt.Valid {
		p.LeftAt = &row.LeftAt.Time
	}

	p.ProfilePictureURI = vars.ToStringPtr(row.ProfilePictureURI)

	return p
}

func participantFromListRow(row *ListMailboxParticipantsRow) *mailbox.Participant {
	slug := row.ProfileSlug
	title := row.ProfileTitle
	kind := row.ProfileKind

	p := &mailbox.Participant{
		ID:             row.ID,
		ConversationID: row.ConversationID,
		ProfileID:      row.ProfileID,
		IsArchived:     row.IsArchived,
		JoinedAt:       row.JoinedAt,
		ProfileSlug:    &slug,
		ProfileTitle:   &title,
		ProfileKind:    &kind,
	}

	if row.LastReadAt.Valid {
		p.LastReadAt = &row.LastReadAt.Time
	}

	if row.LeftAt.Valid {
		p.LeftAt = &row.LeftAt.Time
	}

	p.ProfilePictureURI = vars.ToStringPtr(row.ProfilePictureURI)

	return p
}

func envelopeFromRawRow(row *MailboxEnvelope) *mailbox.Envelope {
	envelope := &mailbox.Envelope{
		ID:              row.ID,
		ConversationID:  row.ConversationID,
		TargetProfileID: row.TargetProfileID,
		SenderProfileID: vars.ToStringPtr(row.SenderProfileID),
		SenderUserID:    vars.ToStringPtr(row.SenderUserID),
		Kind:            row.Kind,
		Status:          row.Status,
		Title:           row.Title,
		Description:     vars.ToStringPtr(row.Description),
		ReplyToID:       vars.ToStringPtr(row.ReplyToID),
		CreatedAt:       row.CreatedAt,
	}

	if row.Properties.Valid {
		var props any

		_ = json.Unmarshal(row.Properties.RawMessage, &props)

		envelope.Properties = props
	}

	assignEnvelopeTimestamps(
		envelope,
		row.AcceptedAt,
		row.RejectedAt,
		row.RevokedAt,
		row.RedeemedAt,
		row.UpdatedAt,
		row.DeletedAt,
	)

	return envelope
}

func envelopeFromConversationRow(row *ListEnvelopesByConversationRow) *mailbox.Envelope {
	envelope := &mailbox.Envelope{
		ID:              row.ID,
		ConversationID:  row.ConversationID,
		TargetProfileID: row.TargetProfileID,
		SenderProfileID: vars.ToStringPtr(row.SenderProfileID),
		SenderUserID:    vars.ToStringPtr(row.SenderUserID),
		Kind:            row.Kind,
		Status:          row.Status,
		Title:           row.Title,
		Description:     vars.ToStringPtr(row.Description),
		ReplyToID:       vars.ToStringPtr(row.ReplyToID),
		CreatedAt:       row.CreatedAt,
	}

	if row.Properties.Valid {
		var props any

		_ = json.Unmarshal(row.Properties.RawMessage, &props)

		envelope.Properties = props
	}

	assignEnvelopeTimestamps(
		envelope,
		row.AcceptedAt,
		row.RejectedAt,
		row.RevokedAt,
		row.RedeemedAt,
		row.UpdatedAt,
		row.DeletedAt,
	)

	// Sender profile info from JOIN.
	envelope.SenderProfileSlug = vars.ToStringPtr(row.SenderProfileSlug)
	envelope.SenderProfileKind = vars.ToStringPtr(row.SenderProfileKind)
	envelope.SenderProfilePictureURI = vars.ToStringPtr(row.SenderProfilePictureURI)

	if row.SenderProfileTitle != "" {
		envelope.SenderProfileTitle = &row.SenderProfileTitle
	}

	return envelope
}

func envelopeFromListRow(row *ListMailboxEnvelopesByTargetProfileIDRow) *mailbox.Envelope {
	envelope := &mailbox.Envelope{
		ID:              row.ID,
		ConversationID:  row.ConversationID,
		TargetProfileID: row.TargetProfileID,
		SenderProfileID: vars.ToStringPtr(row.SenderProfileID),
		SenderUserID:    vars.ToStringPtr(row.SenderUserID),
		Kind:            row.Kind,
		Status:          row.Status,
		Title:           row.Title,
		Description:     vars.ToStringPtr(row.Description),
		ReplyToID:       vars.ToStringPtr(row.ReplyToID),
		CreatedAt:       row.CreatedAt,
	}

	if row.Properties.Valid {
		var props any

		_ = json.Unmarshal(row.Properties.RawMessage, &props)

		envelope.Properties = props
	}

	assignEnvelopeTimestamps(
		envelope,
		row.AcceptedAt,
		row.RejectedAt,
		row.RevokedAt,
		row.RedeemedAt,
		row.UpdatedAt,
		row.DeletedAt,
	)

	// Sender profile info from JOIN.
	envelope.SenderProfileSlug = vars.ToStringPtr(row.SenderProfileSlug)
	envelope.SenderProfileKind = vars.ToStringPtr(row.SenderProfileKind)
	envelope.SenderProfilePictureURI = vars.ToStringPtr(row.SenderProfilePictureURI)

	if row.SenderProfileTitle != "" {
		envelope.SenderProfileTitle = &row.SenderProfileTitle
	}

	return envelope
}

func envelopeFromInvitationRow(row *ListAcceptedMailboxInvitationsRow) *mailbox.Envelope {
	envelope := &mailbox.Envelope{
		ID:              row.ID,
		ConversationID:  row.ConversationID,
		TargetProfileID: row.TargetProfileID,
		SenderProfileID: vars.ToStringPtr(row.SenderProfileID),
		SenderUserID:    vars.ToStringPtr(row.SenderUserID),
		Kind:            row.Kind,
		Status:          row.Status,
		Title:           row.Title,
		Description:     vars.ToStringPtr(row.Description),
		ReplyToID:       vars.ToStringPtr(row.ReplyToID),
		CreatedAt:       row.CreatedAt,
	}

	if row.Properties.Valid {
		var props any

		_ = json.Unmarshal(row.Properties.RawMessage, &props)

		envelope.Properties = props
	}

	assignEnvelopeTimestamps(
		envelope,
		row.AcceptedAt,
		row.RejectedAt,
		row.RevokedAt,
		row.RedeemedAt,
		row.UpdatedAt,
		row.DeletedAt,
	)

	// Sender profile info from JOIN.
	envelope.SenderProfileSlug = vars.ToStringPtr(row.SenderProfileSlug)
	envelope.SenderProfileKind = vars.ToStringPtr(row.SenderProfileKind)
	envelope.SenderProfilePictureURI = vars.ToStringPtr(row.SenderProfilePictureURI)

	if row.SenderProfileTitle != "" {
		envelope.SenderProfileTitle = &row.SenderProfileTitle
	}

	return envelope
}

func assignEnvelopeTimestamps(
	envelope *mailbox.Envelope,
	acceptedAt, rejectedAt, revokedAt, redeemedAt, updatedAt, deletedAt sql.NullTime,
) {
	if acceptedAt.Valid {
		envelope.AcceptedAt = &acceptedAt.Time
	}

	if rejectedAt.Valid {
		envelope.RejectedAt = &rejectedAt.Time
	}

	if revokedAt.Valid {
		envelope.RevokedAt = &revokedAt.Time
	}

	if redeemedAt.Valid {
		envelope.RedeemedAt = &redeemedAt.Time
	}

	if updatedAt.Valid {
		envelope.UpdatedAt = &updatedAt.Time
	}

	if deletedAt.Valid {
		envelope.DeletedAt = &deletedAt.Time
	}
}
