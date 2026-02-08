package events

import "context"

// QueueHandler is a function that processes a specific queue item type.
type QueueHandler func(ctx context.Context, item *QueueItem) error

// HandlerRegistry maps queue item types to their handlers.
type HandlerRegistry struct {
	handlers map[QueueItemType]QueueHandler
}

// NewHandlerRegistry creates a new empty handler registry.
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[QueueItemType]QueueHandler),
	}
}

// Register associates a handler with a queue item type.
func (r *HandlerRegistry) Register(itemType QueueItemType, handler QueueHandler) {
	r.handlers[itemType] = handler
}

// Get returns the handler for a given queue item type, or nil if not registered.
func (r *HandlerRegistry) Get(itemType QueueItemType) QueueHandler {
	handler, exists := r.handlers[itemType]
	if !exists {
		return nil
	}

	return handler
}
