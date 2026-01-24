package runtime_states

import "errors"

// Sentinel errors.
var (
	ErrStateNotFound       = errors.New("runtime state not found")
	ErrFailedToGet         = errors.New("failed to get runtime state")
	ErrFailedToSet         = errors.New("failed to set runtime state")
	ErrFailedToRemove      = errors.New("failed to remove runtime state")
	ErrInvalidTime         = errors.New("failed to parse time value from runtime state")
	ErrFailedToAcquireLock = errors.New("failed to acquire advisory lock")
	ErrFailedToReleaseLock = errors.New("failed to release advisory lock")
)
