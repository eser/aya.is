package events

import "context"

// Handler is a function that processes a specific event type.
type Handler func(ctx context.Context, event *Event) error

// HandlerRegistry maps event types to their handlers.
type HandlerRegistry struct {
	handlers map[EventType]Handler
}

// NewHandlerRegistry creates a new empty handler registry.
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[EventType]Handler),
	}
}

// Register associates a handler with an event type.
func (r *HandlerRegistry) Register(eventType EventType, handler Handler) {
	r.handlers[eventType] = handler
}

// Get returns the handler for a given event type, or nil if not registered.
func (r *HandlerRegistry) Get(eventType EventType) Handler {
	handler, exists := r.handlers[eventType]
	if !exists {
		return nil
	}

	return handler
}
