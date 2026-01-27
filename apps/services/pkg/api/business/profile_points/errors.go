package profile_points

import "errors"

// Sentinel errors.
var (
	ErrInsufficientPoints  = errors.New("insufficient points for transaction")
	ErrInvalidAmount       = errors.New("amount must be positive")
	ErrSelfTransfer        = errors.New("cannot transfer points to self")
	ErrProfileNotFound     = errors.New("profile not found")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrFailedToGetBalance  = errors.New("failed to get point balance")
	ErrFailedToRecordTx    = errors.New("failed to record transaction")
	ErrFailedToListTx      = errors.New("failed to list transactions")
)
