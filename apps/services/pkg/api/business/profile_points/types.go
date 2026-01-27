package profile_points

import (
	"time"
)

// TransactionType represents the type of point transaction.
type TransactionType string

const (
	TransactionTypeGain     TransactionType = "GAIN"
	TransactionTypeTransfer TransactionType = "TRANSFER"
	TransactionTypeSpend    TransactionType = "SPEND"
)

// Transaction represents a single point transaction record.
type Transaction struct {
	ID              string
	TargetProfileID string
	OriginProfileID *string
	TransactionType TransactionType
	TriggeringEvent *string
	Description     string
	Amount          int
	BalanceAfter    int
	CreatedAt       time.Time
}

// Balance represents a profile's current point balance.
type Balance struct {
	ProfileID string
	Points    int
}

// GainParams holds parameters for awarding points.
type GainParams struct {
	TargetProfileID string
	Amount          int
	TriggeringEvent *string
	Description     string
}

// TransferParams holds parameters for transferring points.
type TransferParams struct {
	OriginProfileID string
	TargetProfileID string
	Amount          int
	Description     string
}

// SpendParams holds parameters for spending points.
type SpendParams struct {
	TargetProfileID string
	Amount          int
	TriggeringEvent *string
	Description     string
}

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string
