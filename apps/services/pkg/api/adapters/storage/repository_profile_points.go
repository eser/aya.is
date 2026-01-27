package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

// GetBalance returns the current point balance for a profile.
func (r *Repository) GetBalance(ctx context.Context, profileID string) (uint64, error) {
	points, err := r.queries.GetProfilePoints(ctx, GetProfilePointsParams{
		ProfileID: profileID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, profile_points.ErrProfileNotFound
		}

		return 0, err
	}

	return uint64(points), nil
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
	amount uint64,
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
		Amount:          uint64(row.Amount),
		BalanceAfter:    uint64(row.BalanceAfter),
		CreatedAt:       row.CreatedAt,
	}
}

// CreatePendingAward creates a new pending award record.
func (r *Repository) CreatePendingAward(
	ctx context.Context,
	id string,
	targetProfileID string,
	triggeringEvent string,
	description string,
	amount uint64,
	metadata map[string]any,
) (*profile_points.PendingAward, error) {
	row, err := r.queries.CreatePendingAward(ctx, CreatePendingAwardParams{
		ID:              id,
		TargetProfileID: targetProfileID,
		TriggeringEvent: triggeringEvent,
		Description:     description,
		Amount:          int32(amount),
		Metadata:        vars.ToSQLNullRawMessage(metadata),
	})
	if err != nil {
		return nil, err
	}

	return r.rowToPendingAward(row), nil
}

// GetPendingAwardByID returns a pending award by ID.
func (r *Repository) GetPendingAwardByID(
	ctx context.Context,
	id string,
) (*profile_points.PendingAward, error) {
	row, err := r.queries.GetPendingAwardByID(ctx, GetPendingAwardByIDParams{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, profile_points.ErrPendingAwardNotFound
		}

		return nil, err
	}

	return r.rowToPendingAward(row), nil
}

// ListPendingAwards returns pending awards with optional status filter.
func (r *Repository) ListPendingAwards(
	ctx context.Context,
	status *profile_points.PendingAwardStatus,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profile_points.PendingAward], error) {
	limit := cursor.Limit
	if limit <= 0 {
		limit = 20
	}

	var (
		rows []*ProfilePointPendingAward
		err  error
	)

	if status != nil {
		rows, err = r.queries.ListPendingAwardsByStatus(ctx, ListPendingAwardsByStatusParams{
			Status:     string(*status),
			LimitCount: int32(limit + 1),
		})
	} else {
		rows, err = r.queries.ListPendingAwards(ctx, ListPendingAwardsParams{
			Status:     sql.NullString{Valid: false},
			LimitCount: int32(limit + 1),
		})
	}

	if err != nil {
		return cursors.Cursored[[]*profile_points.PendingAward]{}, err
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	result := make([]*profile_points.PendingAward, len(rows))
	for i, row := range rows {
		result[i] = r.rowToPendingAward(row)
	}

	var nextCursor *string

	if hasMore && len(result) > 0 {
		lastID := result[len(result)-1].ID
		nextCursor = &lastID
	}

	return cursors.WrapResponseWithCursor(result, nextCursor), nil
}

// ApprovePendingAward marks a pending award as approved and awards the points.
func (r *Repository) ApprovePendingAward(
	ctx context.Context,
	awardID string,
	reviewerUserID string,
	transactionID string,
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

	// Get the pending award
	award, err := queriesTx.GetPendingAwardByID(ctx, GetPendingAwardByIDParams{
		ID: awardID,
	})
	if err != nil {
		return nil, err
	}

	// Update the pending award status
	err = queriesTx.ApprovePendingAward(ctx, ApprovePendingAwardParams{
		ID:         awardID,
		ReviewedBy: sql.NullString{String: reviewerUserID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// Add points to the target profile
	_, err = queriesTx.AddPointsToProfile(ctx, AddPointsToProfileParams{
		ID:     award.TargetProfileID,
		Amount: award.Amount,
	})
	if err != nil {
		return nil, err
	}

	// Get the new balance
	newBalance, err := queriesTx.GetProfilePoints(ctx, GetProfilePointsParams{
		ProfileID: award.TargetProfileID,
	})
	if err != nil {
		return nil, err
	}

	// Record the transaction
	triggeringEvent := award.TriggeringEvent

	txRow, err := queriesTx.RecordProfilePointTransaction(ctx, RecordProfilePointTransactionParams{
		ID:              transactionID,
		TargetProfileID: award.TargetProfileID,
		OriginProfileID: sql.NullString{Valid: false},
		TransactionType: string(profile_points.TransactionTypeGain),
		TriggeringEvent: sql.NullString{String: triggeringEvent, Valid: true},
		Description:     award.Description,
		Amount:          award.Amount,
		BalanceAfter:    newBalance,
	})
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.rowToProfilePointTransaction(txRow), nil
}

// RejectPendingAward marks a pending award as rejected.
func (r *Repository) RejectPendingAward(
	ctx context.Context,
	awardID string,
	reviewerUserID string,
	reason string,
) error {
	return r.queries.RejectPendingAward(ctx, RejectPendingAwardParams{
		ID:              awardID,
		ReviewedBy:      sql.NullString{String: reviewerUserID, Valid: true},
		RejectionReason: sql.NullString{String: reason, Valid: reason != ""},
	})
}

// GetPendingAwardsStats returns statistics about pending awards.
func (r *Repository) GetPendingAwardsStats(
	ctx context.Context,
) (*profile_points.PendingAwardsStats, error) {
	row, err := r.queries.GetPendingAwardsStats(ctx)
	if err != nil {
		return nil, err
	}

	// Get breakdown by event type
	eventRows, err := r.queries.GetPendingAwardsStatsByEventType(ctx)
	if err != nil {
		return nil, err
	}

	byEventType := make(map[string]uint64)
	for _, er := range eventRows {
		byEventType[er.TriggeringEvent] = uint64(er.Count)
	}

	// Handle PointsAwarded which may be int64 or nil
	var pointsAwarded uint64

	if row.PointsAwarded != nil {
		switch v := row.PointsAwarded.(type) {
		case int64:
			pointsAwarded = uint64(v)
		case float64:
			pointsAwarded = uint64(v)
		}
	}

	return &profile_points.PendingAwardsStats{
		TotalPending:  uint64(row.TotalPending),
		TotalApproved: uint64(row.TotalApproved),
		TotalRejected: uint64(row.TotalRejected),
		PointsAwarded: pointsAwarded,
		ByEventType:   byEventType,
	}, nil
}

// toMetadataMap safely converts any to map[string]any for metadata fields.
func toMetadataMap(v any) map[string]any {
	if v == nil {
		return nil
	}

	if m, ok := v.(map[string]any); ok {
		return m
	}

	return nil
}

// rowToPendingAward converts a database row to a PendingAward domain object.
func (r *Repository) rowToPendingAward(
	row *ProfilePointPendingAward,
) *profile_points.PendingAward {
	var reviewedAt *time.Time
	if row.ReviewedAt.Valid {
		reviewedAt = &row.ReviewedAt.Time
	}

	return &profile_points.PendingAward{
		ID:              row.ID,
		TargetProfileID: row.TargetProfileID,
		TriggeringEvent: row.TriggeringEvent,
		Description:     row.Description,
		Amount:          uint64(row.Amount),
		Status:          profile_points.PendingAwardStatus(row.Status),
		ReviewedBy:      vars.ToStringPtr(row.ReviewedBy),
		ReviewedAt:      reviewedAt,
		RejectionReason: vars.ToStringPtr(row.RejectionReason),
		Metadata:        toMetadataMap(vars.ToObject(row.Metadata)),
		CreatedAt:       row.CreatedAt,
	}
}
