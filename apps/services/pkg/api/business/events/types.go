package events

import (
	"time"
)

// EventStatus represents the lifecycle of an event in the queue.
type EventStatus string

const (
	StatusPending    EventStatus = "pending"
	StatusProcessing EventStatus = "processing"
	StatusCompleted  EventStatus = "completed"
	StatusFailed     EventStatus = "failed"
	StatusDead       EventStatus = "dead"
)

// EventType identifies the kind of event for handler dispatch.
type EventType string

// Common event types.
const (
	EventTypeNewStory     EventType = "NEW_STORY"
	EventTypeStoryUpdated EventType = "STORY_UPDATED"
	EventTypeProfileSync  EventType = "PROFILE_SYNC"
	EventTypeNotification EventType = "NOTIFICATION"
)

// Event represents a queued event in the system.
type Event struct {
	ID                    string
	Type                  EventType
	Payload               map[string]any
	Status                EventStatus
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

// EnqueueParams holds parameters for enqueueing a new event.
type EnqueueParams struct {
	Type                  EventType
	Payload               map[string]any
	MaxRetries            int        // 0 means use default (3)
	VisibilityTimeoutSecs int        // 0 means use default (300 = 5 minutes)
	ScheduledAt           *time.Time // nil means now (immediately eligible)
}

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string
