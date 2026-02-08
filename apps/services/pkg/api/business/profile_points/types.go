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
	ID              string          `json:"id"`
	TargetProfileID string          `json:"target_profile_id"`
	OriginProfileID *string         `json:"origin_profile_id"`
	TransactionType TransactionType `json:"transaction_type"`
	TriggeringEvent *string         `json:"triggering_event"`
	Description     string          `json:"description"`
	Amount          uint64          `json:"amount"`
	BalanceAfter    uint64          `json:"balance_after"`
	CreatedAt       time.Time       `json:"created_at"`
}

// Balance represents a profile's current point balance.
type Balance struct {
	ProfileID string `json:"profile_id"`
	Points    uint64 `json:"points"`
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
	ActorID         string
	OriginProfileID string
	TargetProfileID string
	Amount          uint64
	Description     string
}

// SpendParams holds parameters for spending points.
type SpendParams struct {
	ActorID         string
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

// ProfileInfo holds basic profile information for display purposes.
type ProfileInfo struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

// PendingAward represents a point award that requires approval.
type PendingAward struct {
	ID              string             `json:"id"`
	TargetProfileID string             `json:"target_profile_id"`
	TargetProfile   *ProfileInfo       `json:"target_profile,omitempty"`
	TriggeringEvent string             `json:"triggering_event"`
	Description     string             `json:"description"`
	Amount          uint64             `json:"amount"`
	Status          PendingAwardStatus `json:"status"`
	ReviewedBy      *string            `json:"reviewed_by"`
	ReviewedAt      *time.Time         `json:"reviewed_at"`
	RejectionReason *string            `json:"rejection_reason"`
	Metadata        map[string]any     `json:"metadata"`
	CreatedAt       time.Time          `json:"created_at"`
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

// Spend costs.
const (
	CostAutoTranslate uint64 = 5
)

// Triggering event identifiers.
const (
	EventStoryPublished    = "STORY_PUBLISHED"
	EventProfileVerified   = "PROFILE_VERIFIED"
	EventFirstContribution = "FIRST_CONTRIBUTION"
	EventAutoTranslate     = "AUTO_TRANSLATE"
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
