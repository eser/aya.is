package runtime_states

import "time"

// RuntimeState represents a key-value pair in the runtime state store.
type RuntimeState struct {
	UpdatedAt time.Time
	Key       string
	Value     string
}
