package queue

import "context"

// Handler is a function that processes a specific item type.
type Handler func(ctx context.Context, item *Item) error

// HandlerRegistry maps item types to their handlers.
type HandlerRegistry struct {
	handlers map[ItemType]Handler
}

// NewHandlerRegistry creates a new empty handler registry.
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[ItemType]Handler),
	}
}

// Register associates a handler with an item type.
func (r *HandlerRegistry) Register(itemType ItemType, handler Handler) {
	r.handlers[itemType] = handler
}

// Get returns the handler for a given item type, or nil if not registered.
func (r *HandlerRegistry) Get(itemType ItemType) Handler {
	handler, exists := r.handlers[itemType]
	if !exists {
		return nil
	}

	return handler
}
