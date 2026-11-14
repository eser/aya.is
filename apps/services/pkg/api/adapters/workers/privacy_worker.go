package workers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/privacy"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
)

// ErrInvalidPrivacyPayload is returned for a queue item that does not carry an export id.
var ErrInvalidPrivacyPayload = errors.New("privacy export item has no export_id")

// PrivacyCleanupConfig holds configuration for the expired export cleanup worker.
type PrivacyCleanupConfig struct {
	Enabled       bool          `conf:"enabled"        default:"true"`
	CheckInterval time.Duration `conf:"check_interval" default:"1h"`
}

// PrivacyExportHandler handles queue items that build a personal data archive.
type PrivacyExportHandler struct {
	logger  *logfx.Logger
	service *privacy.Service
}

// NewPrivacyExportHandler creates the queue handler.
func NewPrivacyExportHandler(logger *logfx.Logger, service *privacy.Service) *PrivacyExportHandler {
	return &PrivacyExportHandler{logger: logger, service: service}
}

// HandleExport builds one archive. A failure that cannot be fixed by retrying is recorded on the export
// itself and swallowed; anything else is returned so the queue retries with backoff.
func (h *PrivacyExportHandler) HandleExport(ctx context.Context, item *events.QueueItem) error {
	exportID, _ := item.Payload["export_id"].(string)
	if exportID == "" {
		h.logger.ErrorContext(ctx, "privacy export item without export_id",
			slog.String("item_id", item.ID))

		return nil
	}

	err := h.service.ProcessExport(ctx, exportID, privacy.IsFinalAttempt(item))
	if err != nil {
		return fmt.Errorf("processing export %s: %w", exportID, err)
	}

	return nil
}

// RegisterHandlers registers the handler for its queue item type.
func (h *PrivacyExportHandler) RegisterHandlers(registry *events.HandlerRegistry) {
	registry.Register(events.QueueItemTypePrivacyExport, h.HandleExport)
}

// PrivacyCleanupWorker removes personal data archives once their retention period is over.
type PrivacyCleanupWorker struct {
	config        *PrivacyCleanupConfig
	logger        *logfx.Logger
	service       *privacy.Service
	runtimeStates *runtime_states.Service
}

// NewPrivacyCleanupWorker creates the cleanup worker.
func NewPrivacyCleanupWorker(
	config *PrivacyCleanupConfig,
	logger *logfx.Logger,
	service *privacy.Service,
	runtimeStates *runtime_states.Service,
) *PrivacyCleanupWorker {
	return &PrivacyCleanupWorker{
		config:        config,
		logger:        logger,
		service:       service,
		runtimeStates: runtimeStates,
	}
}

// Name returns the worker name.
func (w *PrivacyCleanupWorker) Name() string {
	return "privacy-cleanup"
}

// Interval returns the check interval.
func (w *PrivacyCleanupWorker) Interval() time.Duration {
	return w.config.CheckInterval
}

// Execute removes one batch of expired archives.
func (w *PrivacyCleanupWorker) Execute(ctx context.Context) error {
	disabledKey := "worker." + w.Name() + ".disabled"

	disabled, err := w.runtimeStates.Get(ctx, disabledKey)
	if err == nil && disabled == disabledStateValue {
		return workerfx.ErrWorkerSkipped
	}

	handled, err := w.service.CleanupExpiredExports(ctx)
	if err != nil {
		return fmt.Errorf("cleaning up expired exports: %w", err)
	}

	if handled > 0 {
		w.logger.InfoContext(
			ctx,
			"removed expired personal data archives",
			slog.Int("count", handled),
		)
	}

	return nil
}
