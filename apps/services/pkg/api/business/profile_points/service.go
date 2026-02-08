package profile_points

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func DefaultIDGenerator() string {
	return lib.IDsGenerateUnique()
}

// Service provides profile point operations.
type Service struct {
	logger       *logfx.Logger
	repo         Repository
	auditService *events.AuditService
	idGenerator  IDGenerator
}

// NewService creates a new profile points service.
func NewService(
	logger *logfx.Logger,
	repo Repository,
	idGenerator IDGenerator,
	auditService *events.AuditService,
) *Service {
	return &Service{
		logger:       logger,
		repo:         repo,
		auditService: auditService,
		idGenerator:  idGenerator,
	}
}

// GetBalance returns the current point balance for a profile (public).
func (s *Service) GetBalance(ctx context.Context, profileID string) (*Balance, error) {
	points, err := s.repo.GetBalance(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w (profileID: %s): %w", ErrFailedToGetBalance, profileID, err)
	}

	return &Balance{
		ProfileID: profileID,
		Points:    points,
	}, nil
}

// GainPoints awards points to a profile (system action).
func (s *Service) GainPoints(ctx context.Context, params GainParams) (*Transaction, error) {
	if params.Amount == 0 {
		return nil, ErrInvalidAmount
	}

	id := s.idGenerator()

	tx, err := s.repo.RecordTransaction(
		ctx,
		id,
		params.TargetProfileID,
		nil,
		TransactionTypeGain,
		params.TriggeringEvent,
		params.Description,
		params.Amount,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToRecordTx, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.PointsGained,
		EntityType: "profile",
		EntityID:   params.TargetProfileID,
		ActorKind:  events.ActorSystem,
		Payload:    map[string]any{"amount": params.Amount, "description": params.Description},
	})

	return tx, nil
}

// TransferPoints transfers points between profiles.
func (s *Service) TransferPoints(ctx context.Context, params TransferParams) (*Transaction, error) {
	if params.Amount == 0 {
		return nil, ErrInvalidAmount
	}

	if params.OriginProfileID == params.TargetProfileID {
		return nil, ErrSelfTransfer
	}

	// Check origin has sufficient balance
	originBalance, err := s.repo.GetBalance(ctx, params.OriginProfileID)
	if err != nil {
		return nil, fmt.Errorf(
			"%w (profileID: %s): %w",
			ErrFailedToGetBalance,
			params.OriginProfileID,
			err,
		)
	}

	if originBalance < params.Amount {
		return nil, ErrInsufficientPoints
	}

	id := s.idGenerator()
	originID := params.OriginProfileID

	tx, err := s.repo.RecordTransaction(
		ctx,
		id,
		params.TargetProfileID,
		&originID,
		TransactionTypeTransfer,
		nil,
		params.Description,
		params.Amount,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToRecordTx, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.PointsTransferred,
		EntityType: "profile",
		EntityID:   params.TargetProfileID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"amount":            params.Amount,
			"origin_profile_id": params.OriginProfileID,
		},
	})

	return tx, nil
}

// SpendPoints deducts points from a profile.
func (s *Service) SpendPoints(ctx context.Context, params SpendParams) (*Transaction, error) {
	if params.Amount == 0 {
		return nil, ErrInvalidAmount
	}

	// Check profile has sufficient balance
	balance, err := s.repo.GetBalance(ctx, params.TargetProfileID)
	if err != nil {
		return nil, fmt.Errorf(
			"%w (profileID: %s): %w",
			ErrFailedToGetBalance,
			params.TargetProfileID,
			err,
		)
	}

	if balance < params.Amount {
		return nil, ErrInsufficientPoints
	}

	id := s.idGenerator()

	tx, err := s.repo.RecordTransaction(
		ctx,
		id,
		params.TargetProfileID,
		nil,
		TransactionTypeSpend,
		params.TriggeringEvent,
		params.Description,
		params.Amount,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToRecordTx, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.PointsSpent,
		EntityType: "profile",
		EntityID:   params.TargetProfileID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"amount": params.Amount, "description": params.Description},
	})

	return tx, nil
}

// ListTransactions returns transactions for a profile with pagination.
func (s *Service) ListTransactions(
	ctx context.Context,
	profileID string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*Transaction], error) {
	result, err := s.repo.ListTransactionsByProfileID(ctx, profileID, cursor)
	if err != nil {
		return cursors.Cursored[[]*Transaction]{}, fmt.Errorf("%w: %w", ErrFailedToListTx, err)
	}

	return result, nil
}

// AwardForEvent creates a pending award for an event (or awards immediately if auto-approve).
func (s *Service) AwardForEvent(
	ctx context.Context,
	event string,
	targetProfileID string,
	metadata map[string]any,
) (*PendingAward, error) {
	category, exists := AwardCategories[event]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrInvalidEvent, event)
	}

	// All events currently require approval
	return s.CreatePendingAward(ctx, CreatePendingAwardParams{
		TargetProfileID: targetProfileID,
		TriggeringEvent: event,
		Description:     category.Description,
		Amount:          category.Amount,
		Metadata:        metadata,
	})
}

// CreatePendingAward creates a new pending award record.
func (s *Service) CreatePendingAward(
	ctx context.Context,
	params CreatePendingAwardParams,
) (*PendingAward, error) {
	if params.Amount == 0 {
		return nil, ErrInvalidAmount
	}

	id := s.idGenerator()

	award, err := s.repo.CreatePendingAward(
		ctx,
		id,
		params.TargetProfileID,
		params.TriggeringEvent,
		params.Description,
		params.Amount,
		params.Metadata,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreatePendingAward, err)
	}

	return award, nil
}

// GetPendingAward returns a pending award by ID.
func (s *Service) GetPendingAward(ctx context.Context, id string) (*PendingAward, error) {
	award, err := s.repo.GetPendingAwardByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPendingAwardNotFound, err)
	}

	return award, nil
}

// ListPendingAwards returns pending awards with optional status filter.
func (s *Service) ListPendingAwards(
	ctx context.Context,
	status *PendingAwardStatus,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*PendingAward], error) {
	result, err := s.repo.ListPendingAwards(ctx, status, cursor)
	if err != nil {
		return cursors.Cursored[[]*PendingAward]{}, fmt.Errorf(
			"%w: %w",
			ErrFailedToListPendingAwards,
			err,
		)
	}

	return result, nil
}

// ApprovePendingAward approves a pending award and awards the points.
func (s *Service) ApprovePendingAward(
	ctx context.Context,
	awardID string,
	reviewerUserID string,
) (*Transaction, error) {
	// Get the pending award first
	award, err := s.repo.GetPendingAwardByID(ctx, awardID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPendingAwardNotFound, err)
	}

	if award.Status != PendingAwardStatusPending {
		return nil, ErrAwardAlreadyProcessed
	}

	transactionID := s.idGenerator()

	tx, err := s.repo.ApprovePendingAward(ctx, awardID, reviewerUserID, transactionID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToApprovePendingAward, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.AwardApproved,
		EntityType: "pending_award",
		EntityID:   awardID,
		ActorID:    &reviewerUserID,
		ActorKind:  events.ActorUser,
	})

	return tx, nil
}

// RejectPendingAward rejects a pending award.
func (s *Service) RejectPendingAward(
	ctx context.Context,
	awardID string,
	reviewerUserID string,
	reason string,
) error {
	// Get the pending award first
	award, err := s.repo.GetPendingAwardByID(ctx, awardID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPendingAwardNotFound, err)
	}

	if award.Status != PendingAwardStatusPending {
		return ErrAwardAlreadyProcessed
	}

	err = s.repo.RejectPendingAward(ctx, awardID, reviewerUserID, reason)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRejectPendingAward, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.AwardRejected,
		EntityType: "pending_award",
		EntityID:   awardID,
		ActorID:    &reviewerUserID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"reason": reason},
	})

	return nil
}

// GetPendingAwardsStats returns statistics about pending awards.
func (s *Service) GetPendingAwardsStats(ctx context.Context) (*PendingAwardsStats, error) {
	stats, err := s.repo.GetPendingAwardsStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetStats, err)
	}

	return stats, nil
}

// BulkApprovePendingAwards approves multiple pending awards in a single batch operation.
func (s *Service) BulkApprovePendingAwards(
	ctx context.Context,
	awardIDs []string,
	reviewerUserID string,
) ([]string, error) {
	if len(awardIDs) == 0 {
		return []string{}, nil
	}

	approvedIDs, err := s.repo.BulkApprovePendingAwards(
		ctx,
		awardIDs,
		reviewerUserID,
		s.idGenerator,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToApprovePendingAward, err)
	}

	return approvedIDs, nil
}

// BulkRejectPendingAwards rejects multiple pending awards in a single batch operation.
func (s *Service) BulkRejectPendingAwards(
	ctx context.Context,
	awardIDs []string,
	reviewerUserID string,
	reason string,
) ([]string, error) {
	if len(awardIDs) == 0 {
		return []string{}, nil
	}

	rejectedIDs, err := s.repo.BulkRejectPendingAwards(
		ctx,
		awardIDs,
		reviewerUserID,
		reason,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToRejectPendingAward, err)
	}

	return rejectedIDs, nil
}
