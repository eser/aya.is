package events

import (
	"time"
)

// QueueItemStatus represents the lifecycle of an item in the event queue.
type QueueItemStatus string

const (
	QueueStatusPending    QueueItemStatus = "pending"
	QueueStatusProcessing QueueItemStatus = "processing"
	QueueStatusCompleted  QueueItemStatus = "completed"
	QueueStatusFailed     QueueItemStatus = "failed"
	QueueStatusDead       QueueItemStatus = "dead"
)

// QueueItemType identifies the kind of item for handler dispatch.
type QueueItemType string

// Common queue item types.
const (
	QueueItemTypeNewStory     QueueItemType = "NEW_STORY"
	QueueItemTypeStoryUpdated QueueItemType = "STORY_UPDATED"
	QueueItemTypeProfileSync  QueueItemType = "PROFILE_SYNC"
	QueueItemTypeNotification QueueItemType = "NOTIFICATION"
)

// QueueItem represents an item in the event queue.
type QueueItem struct {
	ID                    string
	Type                  QueueItemType
	Payload               map[string]any
	Status                QueueItemStatus
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

// QueueEnqueueParams holds parameters for enqueueing a new item.
type QueueEnqueueParams struct {
	Type                  QueueItemType
	Payload               map[string]any
	MaxRetries            int        // 0 means use default (3)
	VisibilityTimeoutSecs int        // 0 means use default (300 = 5 minutes)
	ScheduledAt           *time.Time // nil means now (immediately eligible)
}
