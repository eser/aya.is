package events

import "errors"

// Sentinel errors.
var (
	ErrFailedToEnqueue      = errors.New("failed to enqueue event")
	ErrFailedToClaim        = errors.New("failed to claim event")
	ErrFailedToComplete     = errors.New("failed to complete event")
	ErrFailedToFail         = errors.New("failed to mark event as failed")
	ErrHandlerNotRegistered = errors.New("no handler registered for event type")
	ErrHandlerPanicked      = errors.New("event handler panicked")
)
