package workers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
)

// Sentinel errors.
var ErrQueueProcessingFailed = errors.New("queue processing failed")

// QueueWorkerConfig holds configuration for the queue worker.
type QueueWorkerConfig struct {
	Enabled      bool          `conf:"enabled"       default:"true"`
	PollInterval time.Duration `conf:"poll_interval" default:"5s"`
	BackoffBase  int           `conf:"backoff_base"  default:"4"`
}

// QueueWorker polls and dispatches items from the queue.
type QueueWorker struct {
	config        *QueueWorkerConfig
	logger        *logfx.Logger
	repo          events.QueueRepository
	registry      *events.HandlerRegistry
	workerID      string
	runtimeStates *runtime_states.Service
}

// NewQueueWorker creates a new queue worker.
func NewQueueWorker(
	config *QueueWorkerConfig,
	logger *logfx.Logger,
	repo events.QueueRepository,
	registry *events.HandlerRegistry,
	workerID string,
	runtimeStates *runtime_states.Service,
) *QueueWorker {
	return &QueueWorker{
		config:        config,
		logger:        logger,
		repo:          repo,
		registry:      registry,
		workerID:      workerID,
		runtimeStates: runtimeStates,
	}
}

// Name returns the worker name.
func (w *QueueWorker) Name() string {
	return "queue"
}

// Interval returns the poll interval.
func (w *QueueWorker) Interval() time.Duration {
	return w.config.PollInterval
}

// Execute runs a single poll cycle: claim an item, dispatch to handler, complete/fail.
func (w *QueueWorker) Execute(ctx context.Context) error {
	// Check if worker is disabled by admin
	disabledKey := "worker." + w.Name() + ".disabled"

	disabled, err := w.runtimeStates.Get(ctx, disabledKey)
	if err == nil && disabled == "true" {
		return nil
	}

	item, err := w.repo.ClaimNext(ctx, w.workerID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrQueueProcessingFailed, err)
	}

	if item == nil {
		return nil
	}

	w.logger.InfoContext(ctx, "Processing queue item",
		slog.String("item_id", item.ID),
		slog.String("type", string(item.Type)),
		slog.Int("retry_count", item.RetryCount))

	handler := w.registry.Get(item.Type)
	if handler == nil {
		errMsg := fmt.Sprintf("no handler registered for item type: %s", item.Type)
		w.logger.ErrorContext(ctx, errMsg, slog.String("item_id", item.ID))

		failErr := w.repo.Fail(ctx, item.ID, w.workerID, errMsg, 0)
		if failErr != nil {
			w.logger.WarnContext(ctx, "failed to mark queue item as failed",
				slog.String("item_id", item.ID),
				slog.String("error", failErr.Error()))
		}

		return nil
	}

	handlerErr := w.executeHandler(ctx, handler, item)
	if handlerErr != nil {
		backoff := events.CalculateBackoff(item.RetryCount, w.config.BackoffBase)

		w.logger.ErrorContext(ctx, "Queue handler failed",
			slog.String("item_id", item.ID),
			slog.String("type", string(item.Type)),
			slog.Int("retry_count", item.RetryCount),
			slog.Int("max_retries", item.MaxRetries),
			slog.Int("backoff_seconds", backoff),
			slog.Any("error", handlerErr))

		failErr := w.repo.Fail(ctx, item.ID, w.workerID, handlerErr.Error(), backoff)
		if failErr != nil {
			w.logger.WarnContext(ctx, "failed to mark queue item as failed after handler error",
				slog.String("item_id", item.ID),
				slog.String("error", failErr.Error()))
		}

		return nil
	}

	err = w.repo.Complete(ctx, item.ID, w.workerID)
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to mark queue item as completed",
			slog.String("item_id", item.ID),
			slog.Any("error", err))

		return nil
	}

	w.logger.DebugContext(ctx, "Queue item completed successfully",
		slog.String("item_id", item.ID),
		slog.String("type", string(item.Type)))

	return nil
}

// executeHandler runs the handler with panic recovery.
func (w *QueueWorker) executeHandler(
	ctx context.Context,
	handler events.QueueHandler,
	item *events.QueueItem,
) (resultErr error) {
	defer func() {
		if rec := recover(); rec != nil {
			resultErr = fmt.Errorf("%w: %v", events.ErrHandlerPanicked, rec)
		}
	}()

	return handler(ctx, item)
}
