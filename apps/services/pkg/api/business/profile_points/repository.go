package profile_points

import (
	"context"

	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

// Repository defines the storage operations for profile points (port).
type Repository interface {
	// GetBalance returns the current point balance for a profile.
	GetBalance(ctx context.Context, profileID string) (uint64, error)

	// RecordTransaction inserts a new transaction and updates the profile's points.
	// This should be done atomically (within a database transaction).
	RecordTransaction(
		ctx context.Context,
		id string,
		targetProfileID string,
		originProfileID *string,
		transactionType TransactionType,
		triggeringEvent *string,
		description string,
		amount uint64,
	) (*Transaction, error)

	// ListTransactionsByProfileID returns transactions for a profile with pagination.
	ListTransactionsByProfileID(
		ctx context.Context,
		profileID string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*Transaction], error)

	// GetTransactionByID returns a single transaction by ID.
	GetTransactionByID(ctx context.Context, id string) (*Transaction, error)

	// CreatePendingAward creates a new pending award record.
	CreatePendingAward(
		ctx context.Context,
		id string,
		targetProfileID string,
		triggeringEvent string,
		description string,
		amount uint64,
		metadata map[string]any,
	) (*PendingAward, error)

	// GetPendingAwardByID returns a pending award by ID.
	GetPendingAwardByID(ctx context.Context, id string) (*PendingAward, error)

	// ListPendingAwards returns pending awards with optional status filter.
	ListPendingAwards(
		ctx context.Context,
		status *PendingAwardStatus,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*PendingAward], error)

	// ApprovePendingAward marks a pending award as approved and awards the points.
	ApprovePendingAward(
		ctx context.Context,
		awardID string,
		reviewerUserID string,
		transactionID string,
	) (*Transaction, error)

	// RejectPendingAward marks a pending award as rejected.
	RejectPendingAward(
		ctx context.Context,
		awardID string,
		reviewerUserID string,
		reason string,
	) error

	// GetPendingAwardsStats returns statistics about pending awards.
	GetPendingAwardsStats(ctx context.Context) (*PendingAwardsStats, error)

	// BulkApprovePendingAwards approves multiple pending awards in a single transaction.
	// Returns the list of successfully approved award IDs.
	BulkApprovePendingAwards(
		ctx context.Context,
		awardIDs []string,
		reviewerUserID string,
		idGenerator IDGenerator,
	) ([]string, error)

	// BulkRejectPendingAwards rejects multiple pending awards in a single transaction.
	// Returns the list of successfully rejected award IDs.
	BulkRejectPendingAwards(
		ctx context.Context,
		awardIDs []string,
		reviewerUserID string,
		reason string,
	) ([]string, error)
}

// PendingAwardsStats holds statistics about pending awards.
type PendingAwardsStats struct {
	TotalPending  uint64            `json:"total_pending"`
	TotalApproved uint64            `json:"total_approved"`
	TotalRejected uint64            `json:"total_rejected"`
	PointsAwarded uint64            `json:"points_awarded"`
	ByEventType   map[string]uint64 `json:"by_event_type"`
}
