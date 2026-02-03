package profile_points

import "errors"

// Sentinel errors.
var (
	ErrInsufficientPoints   = errors.New("insufficient points for transaction")
	ErrInvalidAmount        = errors.New("amount must be positive")
	ErrSelfTransfer         = errors.New("cannot transfer points to self")
	ErrProfileNotFound      = errors.New("profile not found")
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrFailedToGetBalance   = errors.New("failed to get point balance")
	ErrFailedToRecordTx     = errors.New("failed to record transaction")
	ErrFailedToListTx       = errors.New("failed to list transactions")
	ErrMissingProfileID     = errors.New("profile_id is required in event payload")
	ErrMissingOriginProfile = errors.New("origin profile ID required for transfer")

	// Pending award errors.
	ErrInvalidEvent                = errors.New("invalid triggering event")
	ErrPendingAwardNotFound        = errors.New("pending award not found")
	ErrAwardAlreadyProcessed       = errors.New("award has already been processed")
	ErrFailedToCreatePendingAward  = errors.New("failed to create pending award")
	ErrFailedToListPendingAwards   = errors.New("failed to list pending awards")
	ErrFailedToApprovePendingAward = errors.New("failed to approve pending award")
	ErrFailedToRejectPendingAward  = errors.New("failed to reject pending award")
	ErrFailedToGetStats            = errors.New("failed to get pending awards stats")
)
