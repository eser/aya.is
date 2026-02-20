package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	envelopes "github.com/eser/aya.is/services/pkg/api/business/profile_envelopes"
	"github.com/eser/aya.is/services/pkg/lib/vars"
	"github.com/sqlc-dev/pqtype"
)

type envelopeAdapter struct {
	repo *Repository
}

// NewEnvelopeRepository creates a new adapter that implements envelopes.Repository.
func NewEnvelopeRepository(repo *Repository) envelopes.Repository {
	return &envelopeAdapter{repo: repo}
}

func (a *envelopeAdapter) CreateEnvelope(ctx context.Context, envelope *envelopes.Envelope) error {
	var propsJSON pqtype.NullRawMessage

	if envelope.Properties != nil {
		data, err := json.Marshal(envelope.Properties)
		if err != nil {
			return err
		}

		propsJSON = pqtype.NullRawMessage{RawMessage: data, Valid: true}
	}

	return a.repo.queries.CreateProfileEnvelope(ctx, CreateProfileEnvelopeParams{
		ID:              envelope.ID,
		TargetProfileID: envelope.TargetProfileID,
		SenderProfileID: vars.ToSQLNullString(envelope.SenderProfileID),
		SenderUserID:    vars.ToSQLNullString(envelope.SenderUserID),
		Kind:            envelope.Kind,
		Title:           envelope.Title,
		Description:     vars.ToSQLNullString(envelope.Description),
		Properties:      propsJSON,
	})
}

func (a *envelopeAdapter) GetEnvelopeByID(
	ctx context.Context,
	id string,
) (*envelopes.Envelope, error) {
	row, err := a.repo.queries.GetProfileEnvelopeByID(ctx, GetProfileEnvelopeByIDParams{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, envelopes.ErrEnvelopeNotFound
		}

		return nil, err
	}

	return rowToEnvelope(row), nil
}

func (a *envelopeAdapter) ListEnvelopesByTargetProfileID(
	ctx context.Context,
	profileID string,
	statusFilter string,
	limit int,
) ([]*envelopes.Envelope, error) {
	var statusNullStr sql.NullString
	if statusFilter != "" {
		statusNullStr = sql.NullString{String: statusFilter, Valid: true}
	}

	rows, err := a.repo.queries.ListProfileEnvelopesByTargetProfileID(
		ctx,
		ListProfileEnvelopesByTargetProfileIDParams{
			TargetProfileID: profileID,
			StatusFilter:    statusNullStr,
			LimitCount:      int32(limit),
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*envelopes.Envelope, 0, len(rows))
	for _, row := range rows {
		result = append(result, listRowToEnvelope(row))
	}

	return result, nil
}

func (a *envelopeAdapter) UpdateEnvelopeStatus(
	ctx context.Context,
	id string,
	status string,
	now time.Time,
) error {
	nullNow := sql.NullTime{Time: now, Valid: true}

	var rowsAffected int64

	var err error

	switch status {
	case envelopes.StatusAccepted:
		rowsAffected, err = a.repo.queries.UpdateProfileEnvelopeStatusToAccepted(
			ctx, UpdateProfileEnvelopeStatusToAcceptedParams{Now: nullNow, ID: id},
		)
	case envelopes.StatusRejected:
		rowsAffected, err = a.repo.queries.UpdateProfileEnvelopeStatusToRejected(
			ctx, UpdateProfileEnvelopeStatusToRejectedParams{Now: nullNow, ID: id},
		)
	case envelopes.StatusRevoked:
		rowsAffected, err = a.repo.queries.UpdateProfileEnvelopeStatusToRevoked(
			ctx, UpdateProfileEnvelopeStatusToRevokedParams{Now: nullNow, ID: id},
		)
	case envelopes.StatusRedeemed:
		rowsAffected, err = a.repo.queries.UpdateProfileEnvelopeStatusToRedeemed(
			ctx, UpdateProfileEnvelopeStatusToRedeemedParams{Now: nullNow, ID: id},
		)
	default:
		return envelopes.ErrInvalidStatus
	}

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return envelopes.ErrAlreadyProcessed
	}

	return nil
}

func (a *envelopeAdapter) UpdateEnvelopeProperties(
	ctx context.Context,
	id string,
	properties any,
) error {
	data, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	return a.repo.queries.UpdateProfileEnvelopeProperties(
		ctx,
		UpdateProfileEnvelopePropertiesParams{
			ID:         id,
			Properties: pqtype.NullRawMessage{RawMessage: data, Valid: true},
		},
	)
}

func (a *envelopeAdapter) ListAcceptedInvitations(
	ctx context.Context,
	targetProfileID string,
	invitationKind string,
) ([]*envelopes.Envelope, error) {
	var kindFilter sql.NullString
	if invitationKind != "" {
		kindFilter = sql.NullString{String: invitationKind, Valid: true}
	}

	rows, err := a.repo.queries.ListAcceptedInvitations(ctx, ListAcceptedInvitationsParams{
		TargetProfileID: targetProfileID,
		InvitationKind:  kindFilter,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*envelopes.Envelope, 0, len(rows))
	for _, row := range rows {
		result = append(result, invitationRowToEnvelope(row))
	}

	return result, nil
}

func (a *envelopeAdapter) CountPendingEnvelopes(
	ctx context.Context,
	targetProfileID string,
) (int, error) {
	count, err := a.repo.queries.CountPendingProfileEnvelopes(
		ctx,
		CountPendingProfileEnvelopesParams{TargetProfileID: targetProfileID},
	)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// rowToEnvelope converts a sqlc-generated ProfileEnvelope to the business type.
func rowToEnvelope(row *ProfileEnvelope) *envelopes.Envelope {
	envelope := &envelopes.Envelope{
		ID:              row.ID,
		TargetProfileID: row.TargetProfileID,
		SenderProfileID: vars.ToStringPtr(row.SenderProfileID),
		SenderUserID:    vars.ToStringPtr(row.SenderUserID),
		Kind:            row.Kind,
		Status:          row.Status,
		Title:           row.Title,
		Description:     vars.ToStringPtr(row.Description),
		CreatedAt:       row.CreatedAt,
		AcceptedAt:      nil,
		RejectedAt:      nil,
		RevokedAt:       nil,
		RedeemedAt:      nil,
		UpdatedAt:       nil,
		DeletedAt:       nil,
	}

	if row.Properties.Valid {
		var props any

		_ = json.Unmarshal(row.Properties.RawMessage, &props)

		envelope.Properties = props
	}

	if row.AcceptedAt.Valid {
		envelope.AcceptedAt = &row.AcceptedAt.Time
	}

	if row.RejectedAt.Valid {
		envelope.RejectedAt = &row.RejectedAt.Time
	}

	if row.RevokedAt.Valid {
		envelope.RevokedAt = &row.RevokedAt.Time
	}

	if row.RedeemedAt.Valid {
		envelope.RedeemedAt = &row.RedeemedAt.Time
	}

	if row.UpdatedAt.Valid {
		envelope.UpdatedAt = &row.UpdatedAt.Time
	}

	if row.DeletedAt.Valid {
		envelope.DeletedAt = &row.DeletedAt.Time
	}

	return envelope
}

// listRowToEnvelope converts a ListProfileEnvelopesByTargetProfileIDRow to the business type.
func listRowToEnvelope(row *ListProfileEnvelopesByTargetProfileIDRow) *envelopes.Envelope {
	envelope := &envelopes.Envelope{
		ID:              row.ID,
		TargetProfileID: row.TargetProfileID,
		SenderProfileID: vars.ToStringPtr(row.SenderProfileID),
		SenderUserID:    vars.ToStringPtr(row.SenderUserID),
		Kind:            row.Kind,
		Status:          row.Status,
		Title:           row.Title,
		Description:     vars.ToStringPtr(row.Description),
		CreatedAt:       row.CreatedAt,
	}

	if row.Properties.Valid {
		var props any

		_ = json.Unmarshal(row.Properties.RawMessage, &props)

		envelope.Properties = props
	}

	if row.AcceptedAt.Valid {
		envelope.AcceptedAt = &row.AcceptedAt.Time
	}

	if row.RejectedAt.Valid {
		envelope.RejectedAt = &row.RejectedAt.Time
	}

	if row.RevokedAt.Valid {
		envelope.RevokedAt = &row.RevokedAt.Time
	}

	if row.RedeemedAt.Valid {
		envelope.RedeemedAt = &row.RedeemedAt.Time
	}

	if row.UpdatedAt.Valid {
		envelope.UpdatedAt = &row.UpdatedAt.Time
	}

	if row.DeletedAt.Valid {
		envelope.DeletedAt = &row.DeletedAt.Time
	}

	// Sender profile info from JOIN
	envelope.SenderProfileSlug = vars.ToStringPtr(row.SenderProfileSlug)
	envelope.SenderProfileKind = vars.ToStringPtr(row.SenderProfileKind)
	envelope.SenderProfilePictureURI = vars.ToStringPtr(row.SenderProfilePictureURI)

	if row.SenderProfileTitle != "" {
		envelope.SenderProfileTitle = &row.SenderProfileTitle
	}

	return envelope
}

// invitationRowToEnvelope converts a ListAcceptedInvitationsRow to the business type.
func invitationRowToEnvelope(row *ListAcceptedInvitationsRow) *envelopes.Envelope {
	envelope := &envelopes.Envelope{
		ID:              row.ID,
		TargetProfileID: row.TargetProfileID,
		SenderProfileID: vars.ToStringPtr(row.SenderProfileID),
		SenderUserID:    vars.ToStringPtr(row.SenderUserID),
		Kind:            row.Kind,
		Status:          row.Status,
		Title:           row.Title,
		Description:     vars.ToStringPtr(row.Description),
		CreatedAt:       row.CreatedAt,
	}

	if row.Properties.Valid {
		var props any

		_ = json.Unmarshal(row.Properties.RawMessage, &props)

		envelope.Properties = props
	}

	if row.AcceptedAt.Valid {
		envelope.AcceptedAt = &row.AcceptedAt.Time
	}

	if row.RejectedAt.Valid {
		envelope.RejectedAt = &row.RejectedAt.Time
	}

	if row.RevokedAt.Valid {
		envelope.RevokedAt = &row.RevokedAt.Time
	}

	if row.RedeemedAt.Valid {
		envelope.RedeemedAt = &row.RedeemedAt.Time
	}

	if row.UpdatedAt.Valid {
		envelope.UpdatedAt = &row.UpdatedAt.Time
	}

	if row.DeletedAt.Valid {
		envelope.DeletedAt = &row.DeletedAt.Time
	}

	// Sender profile info from JOIN
	envelope.SenderProfileSlug = vars.ToStringPtr(row.SenderProfileSlug)
	envelope.SenderProfileKind = vars.ToStringPtr(row.SenderProfileKind)
	envelope.SenderProfilePictureURI = vars.ToStringPtr(row.SenderProfilePictureURI)

	if row.SenderProfileTitle != "" {
		envelope.SenderProfileTitle = &row.SenderProfileTitle
	}

	return envelope
}
