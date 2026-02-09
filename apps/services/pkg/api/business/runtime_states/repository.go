package runtime_states

import "context"

// Repository defines the storage operations for runtime state (port).
type Repository interface {
	// GetState retrieves a runtime state entry by key.
	// Returns nil, nil if the key does not exist.
	GetState(ctx context.Context, key string) (*RuntimeState, error)

	// SetState upserts a runtime state entry.
	SetState(ctx context.Context, key string, value string) error

	// RemoveState removes a runtime state entry.
	RemoveState(ctx context.Context, key string) error

	// TryAdvisoryLock attempts to acquire a session-scoped advisory lock (non-blocking).
	// Returns true if the lock was acquired, false if held by another session.
	TryAdvisoryLock(ctx context.Context, lockID int64) (bool, error)

	// ReleaseAdvisoryLock releases a previously acquired advisory lock.
	ReleaseAdvisoryLock(ctx context.Context, lockID int64) error

	// ListStatesByPrefix returns all runtime state entries matching a key prefix.
	ListStatesByPrefix(ctx context.Context, prefix string) ([]*RuntimeState, error)
}
