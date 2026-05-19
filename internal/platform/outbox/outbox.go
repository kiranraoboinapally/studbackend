package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/eventbus"

	"gorm.io/gorm"
)

// Writer writes outbox events inside an existing DB transaction.
// This guarantees the event is saved atomically with the business state change.
type Writer struct {
	db *gorm.DB
}

func NewWriter(db *gorm.DB) *Writer {
	return &Writer{db: db}
}

// WriteEvent saves an event into shared.outbox_events within the given transaction.
// Call this inside a db.Transaction() callback so it shares the same commit/rollback.
func (w *Writer) WriteEvent(tx *gorm.DB, aggregateType, aggregateID, eventType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("outbox: marshal payload: %w", err)
	}
	evt := domain.OutboxEvent{
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       string(data),
		Published:     false,
		CreatedAt:     time.Now(),
	}
	return tx.Create(&evt).Error
}

// ─── Worker: polls outbox table and dispatches to EventBus ───────────────────

// Worker is a background goroutine that polls the outbox table
// and publishes unpublished events to the in-process event bus.
// This implements the "Transactional Outbox" pattern with at-least-once delivery.
type Worker struct {
	db           *gorm.DB
	bus          *eventbus.Bus
	pollInterval time.Duration
	batchSize    int
}

func NewWorker(db *gorm.DB, bus *eventbus.Bus, pollIntervalSec, batchSize int) *Worker {
	return &Worker{
		db:           db,
		bus:          bus,
		pollInterval: time.Duration(pollIntervalSec) * time.Second,
		batchSize:    batchSize,
	}
}

// Start begins polling in a background goroutine. Cancel the context to stop.
func (w *Worker) Start(ctx context.Context) {
	log.Printf("🔄 Outbox worker started (poll=%v, batch=%d)", w.pollInterval, w.batchSize)
	go w.run(ctx)
}

func (w *Worker) run(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Outbox worker stopped")
			return
		case <-ticker.C:
			w.processBatch(ctx)
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) {
	var events []domain.OutboxEvent
	result := w.db.Where("published = ? AND retry_count < ?", false, 5).
		Order("created_at ASC").
		Limit(w.batchSize).
		Find(&events)

	if result.Error != nil {
		log.Printf("⚠️  Outbox poll error: %v", result.Error)
		return
	}

	for _, evt := range events {
		// Dispatch to event bus
		w.bus.Publish(ctx, eventbus.Event{
			Type:          evt.EventType,
			AggregateType: evt.AggregateType,
			AggregateID:   evt.AggregateID,
			Payload:       evt.Payload, // raw JSON string; handlers will unmarshal
		})

		// Mark as published
		now := time.Now()
		if err := w.db.Model(&evt).Updates(map[string]interface{}{
			"published":    true,
			"published_at": &now,
		}).Error; err != nil {
			log.Printf("⚠️  Outbox: failed to mark event %d as published: %v", evt.ID, err)
			// Increment retry count
			w.db.Model(&evt).Updates(map[string]interface{}{
				"retry_count": gorm.Expr("retry_count + 1"),
				"last_error":  err.Error(),
			})
		}
	}
}
