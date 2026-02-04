package queue

import "errors"

// Sentinel errors.
var (
	ErrFailedToEnqueue      = errors.New("failed to enqueue queue item")
	ErrFailedToClaim        = errors.New("failed to claim queue item")
	ErrFailedToComplete     = errors.New("failed to complete queue item")
	ErrFailedToFail         = errors.New("failed to mark queue item as failed")
	ErrHandlerNotRegistered = errors.New("no handler registered for item type")
	ErrHandlerPanicked      = errors.New("queue handler panicked")
)
