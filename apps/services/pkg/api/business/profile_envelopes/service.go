package profile_envelopes

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

// Repository defines storage operations for the envelope service.
type Repository interface {
	CreateEnvelope(ctx context.Context, envelope *Envelope) error
	GetEnvelopeByID(ctx context.Context, id string) (*Envelope, error)
	ListEnvelopesByTargetProfileID(
		ctx context.Context,
		profileID string,
		statusFilter string,
		limit int,
	) ([]*Envelope, error)
	UpdateEnvelopeStatus(
		ctx context.Context,
		id string,
		status string,
		now time.Time,
	) error
	UpdateEnvelopeProperties(ctx context.Context, id string, properties any) error
	ListAcceptedInvitations(
		ctx context.Context,
		targetProfileID string,
		invitationKind string,
	) ([]*Envelope, error)
	CountPendingEnvelopes(ctx context.Context, targetProfileID string) (int, error)
}

// Service provides envelope business logic.
type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator func() string
	notifiers   []EnvelopeNotifier
}

// NewService creates a new envelope service.
func NewService(
	logger *logfx.Logger,
	repo Repository,
	idGenerator func() string,
) *Service {
	return &Service{
		logger:      logger,
		repo:        repo,
		idGenerator: idGenerator,
		notifiers:   nil,
	}
}

// RegisterNotifier adds a notifier that will be called when new envelopes are created.
func (s *Service) RegisterNotifier(notifier EnvelopeNotifier) {
	s.notifiers = append(s.notifiers, notifier)
}

// CreateEnvelope creates a new envelope in pending status.
func (s *Service) CreateEnvelope(
	ctx context.Context,
	params *CreateEnvelopeParams,
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

	envelope := &Envelope{
		ID:              s.idGenerator(),
		TargetProfileID: params.TargetProfileID,
		SenderProfileID: params.SenderProfileID,
		SenderUserID:    params.SenderUserID,
		Kind:            params.Kind,
		Status:          StatusPending,
		Title:           params.Title,
		Description:     params.Description,
		Properties:      params.Properties,
		AcceptedAt:      nil,
		RejectedAt:      nil,
		RevokedAt:       nil,
		RedeemedAt:      nil,
		CreatedAt:       now,
		UpdatedAt:       nil,
		DeletedAt:       nil,
	}

	err := s.repo.CreateEnvelope(ctx, envelope)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreate, err)
	}

	s.logger.InfoContext(ctx, "Envelope created",
		slog.String("envelope_id", envelope.ID),
		slog.String("kind", params.Kind),
		slog.String("target_profile_id", params.TargetProfileID))

	// Notify registered notifiers (best-effort, fire-and-forget).
	s.notifyNewEnvelope(ctx, params)

	return envelope, nil
}

// notifyNewEnvelope fires all registered notifiers for a newly created envelope.
// Notifications are best-effort — failures are logged but do not affect the caller.
func (s *Service) notifyNewEnvelope(ctx context.Context, params *CreateEnvelopeParams) {
	if len(s.notifiers) == 0 {
		return
	}

	notification := &EnvelopeNotification{
		TargetProfileID:    params.TargetProfileID,
		EnvelopeTitle:      params.Title,
		SenderProfileTitle: params.SenderProfileTitle,
		Locale:             params.Locale,
	}

	for _, notifier := range s.notifiers {
		notifier.NotifyNewEnvelope(ctx, notification)
	}
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

// GetAcceptedInvitations returns accepted invitation envelopes for a profile,
// optionally filtered by invitation sub-kind in properties.
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

// GetEnvelopeByID returns a single envelope by ID.
func (s *Service) GetEnvelopeByID(ctx context.Context, id string) (*Envelope, error) {
	envelope, err := s.repo.GetEnvelopeByID(ctx, id)
	if err != nil {
		return nil, ErrEnvelopeNotFound
	}

	return envelope, nil
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
