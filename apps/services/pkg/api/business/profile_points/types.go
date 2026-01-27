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
	Amount          uint64
	BalanceAfter    uint64
	CreatedAt       time.Time
}

// Balance represents a profile's current point balance.
type Balance struct {
	ProfileID string
	Points    uint64
}

// GainParams holds parameters for awarding points.
type GainParams struct {
	TargetProfileID string
	Amount          uint64
	TriggeringEvent *string
	Description     string
}

// TransferParams holds parameters for transferring points.
type TransferParams struct {
	OriginProfileID string
	TargetProfileID string
	Amount          uint64
	Description     string
}

// SpendParams holds parameters for spending points.
type SpendParams struct {
	TargetProfileID string
	Amount          uint64
	TriggeringEvent *string
	Description     string
}

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string

// PendingAwardStatus represents the status of a pending point award.
type PendingAwardStatus string

const (
	PendingAwardStatusPending  PendingAwardStatus = "pending"
	PendingAwardStatusApproved PendingAwardStatus = "approved"
	PendingAwardStatusRejected PendingAwardStatus = "rejected"
)

// PendingAward represents a point award that requires approval.
type PendingAward struct {
	ID              string
	TargetProfileID string
	TriggeringEvent string
	Description     string
	Amount          uint64
	Status          PendingAwardStatus
	ReviewedBy      *string
	ReviewedAt      *time.Time
	RejectionReason *string
	Metadata        map[string]any
	CreatedAt       time.Time
}

// CreatePendingAwardParams holds parameters for creating a pending award.
type CreatePendingAwardParams struct {
	TargetProfileID string
	TriggeringEvent string
	Description     string
	Amount          uint64
	Metadata        map[string]any
}

// AwardCategory defines the configuration for a type of point award.
type AwardCategory struct {
	Event       string
	Amount      uint64
	AutoApprove bool
	Description string
}

// Triggering event identifiers.
const (
	EventStoryPublished    = "STORY_PUBLISHED"
	EventProfileVerified   = "PROFILE_VERIFIED"
	EventFirstContribution = "FIRST_CONTRIBUTION"
)

// AwardCategories defines the point amounts and rules for each event type.
var AwardCategories = map[string]AwardCategory{
	EventStoryPublished: {
		Event:       EventStoryPublished,
		Amount:      10,
		AutoApprove: false,
		Description: "Published a new story",
	},
	EventProfileVerified: {
		Event:       EventProfileVerified,
		Amount:      50,
		AutoApprove: false,
		Description: "Verified profile",
	},
	EventFirstContribution: {
		Event:       EventFirstContribution,
		Amount:      25,
		AutoApprove: false,
		Description: "First contribution to a project",
	},
}
