package workers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
)

// Sentinel errors.
var ErrEventProcessingFailed = errors.New("event processing failed")

// EventQueueWorkerConfig holds configuration for the event queue worker.
type EventQueueWorkerConfig struct {
	Enabled      bool          `conf:"enabled"       default:"true"`
	PollInterval time.Duration `conf:"poll_interval" default:"5s"`
	BackoffBase  int           `conf:"backoff_base"  default:"4"`
}

// EventQueueWorker polls and dispatches events from the queue.
type EventQueueWorker struct {
	config   *EventQueueWorkerConfig
	logger   *logfx.Logger
	repo     events.Repository
	registry *events.HandlerRegistry
	workerID string
}

// NewEventQueueWorker creates a new event queue worker.
func NewEventQueueWorker(
	config *EventQueueWorkerConfig,
	logger *logfx.Logger,
	repo events.Repository,
	registry *events.HandlerRegistry,
	workerID string,
) *EventQueueWorker {
	return &EventQueueWorker{
		config:   config,
		logger:   logger,
		repo:     repo,
		registry: registry,
		workerID: workerID,
	}
}

// Name returns the worker name.
func (w *EventQueueWorker) Name() string {
	return "event-queue"
}

// Interval returns the poll interval.
func (w *EventQueueWorker) Interval() time.Duration {
	return w.config.PollInterval
}

// Execute runs a single poll cycle: claim an event, dispatch to handler, complete/fail.
func (w *EventQueueWorker) Execute(ctx context.Context) error {
	event, err := w.repo.ClaimNext(ctx, w.workerID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrEventProcessingFailed, err)
	}

	if event == nil {
		return nil
	}

	w.logger.InfoContext(ctx, "Processing event",
		slog.String("event_id", event.ID),
		slog.String("type", string(event.Type)),
		slog.Int("retry_count", event.RetryCount))

	handler := w.registry.Get(event.Type)
	if handler == nil {
		errMsg := fmt.Sprintf("no handler registered for event type: %s", event.Type)
		w.logger.ErrorContext(ctx, errMsg, slog.String("event_id", event.ID))

		_ = w.repo.Fail(ctx, event.ID, w.workerID, errMsg, 0)

		return nil
	}

	handlerErr := w.executeHandler(ctx, handler, event)
	if handlerErr != nil {
		backoff := events.CalculateBackoff(event.RetryCount, w.config.BackoffBase)

		w.logger.ErrorContext(ctx, "Event handler failed",
			slog.String("event_id", event.ID),
			slog.String("type", string(event.Type)),
			slog.Int("retry_count", event.RetryCount),
			slog.Int("max_retries", event.MaxRetries),
			slog.Int("backoff_seconds", backoff),
			slog.Any("error", handlerErr))

		_ = w.repo.Fail(ctx, event.ID, w.workerID, handlerErr.Error(), backoff)

		return nil
	}

	err = w.repo.Complete(ctx, event.ID, w.workerID)
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to mark event as completed",
			slog.String("event_id", event.ID),
			slog.Any("error", err))

		return nil
	}

	w.logger.DebugContext(ctx, "Event completed successfully",
		slog.String("event_id", event.ID),
		slog.String("type", string(event.Type)))

	return nil
}

// executeHandler runs the handler with panic recovery.
func (w *EventQueueWorker) executeHandler(
	ctx context.Context,
	handler events.Handler,
	event *events.Event,
) (resultErr error) {
	defer func() {
		if rec := recover(); rec != nil {
			resultErr = fmt.Errorf("%w: %v", events.ErrHandlerPanicked, rec)
		}
	}()

	return handler(ctx, event)
}
