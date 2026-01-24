package runtime_states

import "time"

// RuntimeState represents a key-value pair in the runtime state store.
type RuntimeState struct {
	Key       string
	Value     string
	UpdatedAt time.Time
}
