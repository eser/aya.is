package workerfx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

// Sentinel errors.
var (
	ErrWorkerPanicked = errors.New("worker panicked during execution")
)

// Runner manages the execution loop for a worker.
type Runner struct {
	worker Worker
	logger *logfx.Logger
	status WorkerStatus
	mu     sync.RWMutex
}

// NewRunner creates a new worker runner.
func NewRunner(worker Worker, logger *logfx.Logger) *Runner {
	return &Runner{ //nolint:exhaustruct
		worker: worker,
		logger: logger,
		status: WorkerStatus{ //nolint:exhaustruct
			Name: worker.Name(),
		},
	}
}

// Run starts the worker execution loop.
// Blocks until context is canceled.
func (r *Runner) Run(ctx context.Context) error {
	r.logger.InfoContext(ctx, "Starting worker",
		slog.String("worker", r.worker.Name()),
		slog.Duration("interval", r.worker.Interval()))

	// Run immediately on start
	r.executeWorker(ctx)

	// If interval is 0, run continuously without delay
	if r.worker.Interval() == 0 {
		for {
			select {
			case <-ctx.Done():
				r.logger.InfoContext(ctx, "Worker stopped",
					slog.String("worker", r.worker.Name()))

				return nil
			default:
				r.executeWorker(ctx)
			}
		}
	}

	// Run with interval
	ticker := time.NewTicker(r.worker.Interval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.logger.InfoContext(ctx, "Worker stopped",
				slog.String("worker", r.worker.Name()))

			return nil
		case <-ticker.C:
			r.executeWorker(ctx)
		}
	}
}

// Status returns the current worker status.
func (r *Runner) Status() WorkerStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.status
}

// executeWorker runs a single worker cycle with panic recovery.
func (r *Runner) executeWorker(ctx context.Context) {
	r.mu.Lock()
	r.status.IsRunning = true
	r.mu.Unlock()

	start := time.Now()

	defer func() {
		duration := time.Since(start)

		r.mu.Lock()
		r.status.IsRunning = false
		r.status.LastRun = start
		r.status.LastDuration = duration
		r.status.RunCount++
		r.mu.Unlock()

		if rec := recover(); rec != nil {
			err := fmt.Errorf("%w: %v", ErrWorkerPanicked, rec)

			r.mu.Lock()
			r.status.LastError = err
			r.status.ErrorCount++
			r.mu.Unlock()

			r.logger.ErrorContext(ctx, "Worker panicked",
				slog.String("worker", r.worker.Name()),
				slog.Duration("duration", duration),
				slog.Any("panic", rec))
		}
	}()

	err := r.worker.Execute(ctx)
	duration := time.Since(start)

	r.mu.Lock()

	r.status.LastError = err
	if err != nil {
		r.status.ErrorCount++
	}

	r.mu.Unlock()

	if err != nil {
		r.logger.ErrorContext(ctx, "Worker execution failed",
			slog.String("worker", r.worker.Name()),
			slog.Duration("duration", duration),
			slog.Any("error", err))
	} else {
		r.logger.DebugContext(ctx, "Worker execution completed",
			slog.String("worker", r.worker.Name()),
			slog.Duration("duration", duration))
	}
}
