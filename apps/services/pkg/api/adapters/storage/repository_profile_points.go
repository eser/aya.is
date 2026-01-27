package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

// GetBalance returns the current point balance for a profile.
func (r *Repository) GetBalance(ctx context.Context, profileID string) (int, error) {
	points, err := r.queries.GetProfilePoints(ctx, GetProfilePointsParams{
		ProfileID: profileID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, profile_points.ErrProfileNotFound
		}

		return 0, err
	}

	return int(points), nil
}

// RecordTransaction inserts a new transaction and updates the profile's points atomically.
func (r *Repository) RecordTransaction(
	ctx context.Context,
	id string,
	targetProfileID string,
	originProfileID *string,
	transactionType profile_points.TransactionType,
	triggeringEvent *string,
	description string,
	amount int,
) (*profile_points.Transaction, error) {
	// Start a database transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	queriesTx := r.queries.WithTx(tx)

	// Handle balance updates based on transaction type
	switch transactionType {
	case profile_points.TransactionTypeGain:
		// Add points to target profile
		_, err = queriesTx.AddPointsToProfile(ctx, AddPointsToProfileParams{
			ID:     targetProfileID,
			Amount: int32(amount),
		})
		if err != nil {
			return nil, err
		}

	case profile_points.TransactionTypeSpend:
		// Deduct points from target profile
		rowsAffected, err := queriesTx.DeductPointsFromProfile(ctx, DeductPointsFromProfileParams{
			ID:     targetProfileID,
			Amount: int32(amount),
		})
		if err != nil {
			return nil, err
		}

		if rowsAffected == 0 {
			return nil, profile_points.ErrInsufficientPoints
		}

	case profile_points.TransactionTypeTransfer:
		if originProfileID == nil {
			return nil, errors.New("origin profile ID required for transfer")
		}

		// Deduct from origin
		rowsAffected, err := queriesTx.DeductPointsFromProfile(ctx, DeductPointsFromProfileParams{
			ID:     *originProfileID,
			Amount: int32(amount),
		})
		if err != nil {
			return nil, err
		}

		if rowsAffected == 0 {
			return nil, profile_points.ErrInsufficientPoints
		}

		// Add to target
		_, err = queriesTx.AddPointsToProfile(ctx, AddPointsToProfileParams{
			ID:     targetProfileID,
			Amount: int32(amount),
		})
		if err != nil {
			return nil, err
		}
	}

	// Get the new balance after the update
	newBalance, err := queriesTx.GetProfilePoints(ctx, GetProfilePointsParams{
		ProfileID: targetProfileID,
	})
	if err != nil {
		return nil, err
	}

	// Record the transaction
	row, err := queriesTx.RecordProfilePointTransaction(ctx, RecordProfilePointTransactionParams{
		ID:              id,
		TargetProfileID: targetProfileID,
		OriginProfileID: vars.ToSQLNullString(originProfileID),
		TransactionType: string(transactionType),
		TriggeringEvent: vars.ToSQLNullString(triggeringEvent),
		Description:     description,
		Amount:          int32(amount),
		BalanceAfter:    newBalance,
	})
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.rowToProfilePointTransaction(row), nil
}

// ListTransactionsByProfileID returns transactions for a profile with pagination.
func (r *Repository) ListTransactionsByProfileID(
	ctx context.Context,
	profileID string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profile_points.Transaction], error) {
	limit := cursor.Limit
	if limit <= 0 {
		limit = 20
	}

	rows, err := r.queries.ListProfilePointTransactionsByProfileID(
		ctx,
		ListProfilePointTransactionsByProfileIDParams{
			ProfileID:  profileID,
			LimitCount: int32(limit + 1), // Fetch one extra to determine if there are more
		},
	)
	if err != nil {
		return cursors.Cursored[[]*profile_points.Transaction]{}, err
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	result := make([]*profile_points.Transaction, len(rows))
	for i, row := range rows {
		result[i] = r.rowToProfilePointTransaction(row)
	}

	var nextCursor *string

	if hasMore && len(result) > 0 {
		lastID := result[len(result)-1].ID
		nextCursor = &lastID
	}

	return cursors.WrapResponseWithCursor(result, nextCursor), nil
}

// GetTransactionByID returns a single transaction by ID.
func (r *Repository) GetTransactionByID(
	ctx context.Context,
	id string,
) (*profile_points.Transaction, error) {
	row, err := r.queries.GetProfilePointTransactionByID(ctx, GetProfilePointTransactionByIDParams{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, profile_points.ErrTransactionNotFound
		}

		return nil, err
	}

	return r.rowToProfilePointTransaction(row), nil
}

// rowToProfilePointTransaction converts a database row to a Transaction domain object.
func (r *Repository) rowToProfilePointTransaction(
	row *ProfilePointTransaction,
) *profile_points.Transaction {
	return &profile_points.Transaction{
		ID:              row.ID,
		TargetProfileID: row.TargetProfileID,
		OriginProfileID: vars.ToStringPtr(row.OriginProfileID),
		TransactionType: profile_points.TransactionType(row.TransactionType),
		TriggeringEvent: vars.ToStringPtr(row.TriggeringEvent),
		Description:     row.Description,
		Amount:          int(row.Amount),
		BalanceAfter:    int(row.BalanceAfter),
		CreatedAt:       row.CreatedAt,
	}
}
