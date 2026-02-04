package queue

import (
	"time"
)

// ItemStatus represents the lifecycle of an item in the queue.
type ItemStatus string

const (
	StatusPending    ItemStatus = "pending"
	StatusProcessing ItemStatus = "processing"
	StatusCompleted  ItemStatus = "completed"
	StatusFailed     ItemStatus = "failed"
	StatusDead       ItemStatus = "dead"
)

// ItemType identifies the kind of item for handler dispatch.
type ItemType string

// Common item types.
const (
	ItemTypeNewStory     ItemType = "NEW_STORY"
	ItemTypeStoryUpdated ItemType = "STORY_UPDATED"
	ItemTypeProfileSync  ItemType = "PROFILE_SYNC"
	ItemTypeNotification ItemType = "NOTIFICATION"
)

// Item represents a queued item in the system.
type Item struct {
	ID                    string
	Type                  ItemType
	Payload               map[string]any
	Status                ItemStatus
	RetryCount            int
	MaxRetries            int
	VisibleAt             time.Time
	VisibilityTimeoutSecs int
	StartedAt             *time.Time
	CompletedAt           *time.Time
	FailedAt              *time.Time
	CreatedAt             time.Time
	UpdatedAt             *time.Time
	ErrorMessage          *string
	WorkerID              *string
}

// EnqueueParams holds parameters for enqueueing a new item.
type EnqueueParams struct {
	Type                  ItemType
	Payload               map[string]any
	MaxRetries            int        // 0 means use default (3)
	VisibilityTimeoutSecs int        // 0 means use default (300 = 5 minutes)
	ScheduledAt           *time.Time // nil means now (immediately eligible)
}

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string
