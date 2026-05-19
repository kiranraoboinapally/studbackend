package eventbus

import (
	"context"
	"log"
	"sync"
)

// Event represents a domain event that flows through the internal bus.
type Event struct {
	Type          string      // e.g. "student.enrolled", "payment.completed"
	AggregateType string     // e.g. "Student", "Invoice"
	AggregateID   string     // e.g. "42"
	Payload       interface{} // The actual event data (will be JSON-serialized for outbox)
}

// HandlerFunc processes an event. Return error to signal failure (logged, not retried by bus).
type HandlerFunc func(ctx context.Context, event Event) error

// Bus is an in-process, synchronous-fanout event bus.
// Modules subscribe to event types; when Publish is called,
// all matching handlers execute in the caller's goroutine context.
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]HandlerFunc
}

// New creates a new event bus.
func New() *Bus {
	return &Bus{
		handlers: make(map[string][]HandlerFunc),
	}
}

// Subscribe registers a handler for a given event type.
// Multiple handlers can subscribe to the same event type.
func (b *Bus) Subscribe(eventType string, handler HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	log.Printf("📡 EventBus: subscribed to [%s]", eventType)
}

// Publish dispatches an event to all registered handlers synchronously.
// Errors from handlers are logged but do not stop other handlers.
func (b *Bus) Publish(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()

	if len(handlers) == 0 {
		return
	}

	log.Printf("📡 EventBus: publishing [%s] aggregate=%s id=%s to %d handler(s)",
		event.Type, event.AggregateType, event.AggregateID, len(handlers))

	for _, h := range handlers {
		if err := h(ctx, event); err != nil {
			log.Printf("⚠️  EventBus: handler error for [%s]: %v", event.Type, err)
		}
	}
}

// PublishAsync dispatches an event to all handlers in separate goroutines.
func (b *Bus) PublishAsync(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()

	for _, h := range handlers {
		go func(handler HandlerFunc) {
			if err := handler(ctx, event); err != nil {
				log.Printf("⚠️  EventBus: async handler error for [%s]: %v", event.Type, err)
			}
		}(h)
	}
}
