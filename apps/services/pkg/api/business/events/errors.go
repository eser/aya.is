package events

import "errors"

// Sentinel errors for audit operations.
var (
	ErrFailedToRecordAudit = errors.New("failed to record audit entry")
	ErrFailedToListAudit   = errors.New("failed to list audit entries")
)

// Sentinel errors for queue operations.
var (
	ErrFailedToEnqueue      = errors.New("failed to enqueue event queue item")
	ErrFailedToClaim        = errors.New("failed to claim event queue item")
	ErrFailedToComplete     = errors.New("failed to complete event queue item")
	ErrFailedToFail         = errors.New("failed to mark event queue item as failed")
	ErrHandlerNotRegistered = errors.New("no handler registered for item type")
	ErrHandlerPanicked      = errors.New("event queue handler panicked")
)
