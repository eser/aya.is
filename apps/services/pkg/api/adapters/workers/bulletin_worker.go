package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	bulletinbiz "github.com/eser/aya.is/services/pkg/api/business/bulletin"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
)

const lockIDBulletin int64 = 100012

// BulletinWorker periodically processes bulletin digest windows.
type BulletinWorker struct {
	config        *BulletinConfig
	logger        *logfx.Logger
	service       *bulletinbiz.Service
	runtimeStates *runtime_states.Service
}

// NewBulletinWorker creates a new bulletin worker.
func NewBulletinWorker(
	config *BulletinConfig,
	logger *logfx.Logger,
	service *bulletinbiz.Service,
	runtimeStates *runtime_states.Service,
) *BulletinWorker {
	return &BulletinWorker{
		config:        config,
		logger:        logger,
		service:       service,
		runtimeStates: runtimeStates,
	}
}

// Name returns the worker name.
func (w *BulletinWorker) Name() string {
	return "bulletin"
}

// Interval returns the check interval.
func (w *BulletinWorker) Interval() time.Duration {
	return w.config.CheckInterval
}

// Execute runs a bulletin digest processing cycle.
func (w *BulletinWorker) Execute(ctx context.Context) error {
	// Check if worker is disabled by admin
	disabledKey := "worker." + w.Name() + ".disabled"

	disabled, err := w.runtimeStates.Get(ctx, disabledKey)
	if err == nil && disabled == "true" {
		return workerfx.ErrWorkerSkipped
	}

	// Try advisory lock to prevent concurrent execution
	acquired, lockErr := w.runtimeStates.TryLock(ctx, lockIDBulletin)
	if lockErr != nil {
		w.logger.WarnContext(ctx, "Failed to acquire advisory lock for bulletin",
			slog.Any("error", lockErr))

		return workerfx.ErrWorkerSkipped
	}

	if !acquired {
		w.logger.DebugContext(ctx, "Another instance is running bulletin worker")

		return workerfx.ErrWorkerSkipped
	}

	defer func() {
		releaseErr := w.runtimeStates.ReleaseLock(ctx, lockIDBulletin)
		if releaseErr != nil {
			w.logger.WarnContext(ctx, "Failed to release advisory lock for bulletin",
				slog.String("error", releaseErr.Error()))
		}
	}()

	return w.service.ProcessDigestWindow(ctx)
}
