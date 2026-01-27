package profile_points

import (
	"context"

	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

// Repository defines the storage operations for profile points (port).
type Repository interface {
	// GetBalance returns the current point balance for a profile.
	GetBalance(ctx context.Context, profileID string) (int, error)

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
		amount int,
	) (*Transaction, error)

	// ListTransactionsByProfileID returns transactions for a profile with pagination.
	ListTransactionsByProfileID(
		ctx context.Context,
		profileID string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*Transaction], error)

	// GetTransactionByID returns a single transaction by ID.
	GetTransactionByID(ctx context.Context, id string) (*Transaction, error)
}
