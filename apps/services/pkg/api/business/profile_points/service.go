package profile_points

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func DefaultIDGenerator() string {
	return lib.IDsGenerateUnique()
}

// Service provides profile point operations.
type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator IDGenerator
}

// NewService creates a new profile points service.
func NewService(
	logger *logfx.Logger,
	repo Repository,
	idGenerator IDGenerator,
) *Service {
	return &Service{
		logger:      logger,
		repo:        repo,
		idGenerator: idGenerator,
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
	if params.Amount <= 0 {
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

	return tx, nil
}

// TransferPoints transfers points between profiles.
func (s *Service) TransferPoints(ctx context.Context, params TransferParams) (*Transaction, error) {
	if params.Amount <= 0 {
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

	return tx, nil
}

// SpendPoints deducts points from a profile.
func (s *Service) SpendPoints(ctx context.Context, params SpendParams) (*Transaction, error) {
	if params.Amount <= 0 {
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
