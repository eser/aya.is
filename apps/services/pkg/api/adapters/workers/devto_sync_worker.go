package workers

import (
  "context"
  "fmt"
  "log/slog"
  "time"

  "github.com/eser/aya.is/services/pkg/ajan/logfx"
  "github.com/eser/aya.is/services/pkg/ajan/workerfx"
  "github.com/eser/aya.is/services/pkg/api/business/linksync"
  "github.com/eser/aya.is/services/pkg/api/business/runtime_states"
  "github.com/eser/aya.is/services/pkg/api/business/siteimporter"
)

const lockIDDevtoFullSync int64 = 100007

// DevtoSyncWorker syncs Dev.to articles for managed profile links.
type DevtoSyncWorker struct {
  config              *DevtoSyncConfig
  logger              *logfx.Logger
  syncService         *linksync.Service
  siteImporterService *siteimporter.Service
  storyProcessor      StoryProcessor
  runtimeStates       *runtime_states.Service
  idGenerator         func() string
}

// NewDevtoSyncWorker creates a new Dev.to sync worker.
func NewDevtoSyncWorker(
  config *DevtoSyncConfig,
  logger *logfx.Logger,
  syncService *linksync.Service,
  siteImporterService *siteimporter.Service,
  storyProcessor StoryProcessor,
  runtimeStates *runtime_states.Service,
  idGenerator func() string,
) *DevtoSyncWorker {
  return &DevtoSyncWorker{
    config:              config,
    logger:              logger,
    syncService:         syncService,
    siteImporterService: siteImporterService,
    storyProcessor:      storyProcessor,
    runtimeStates:       runtimeStates,
    idGenerator:         idGenerator,
  }
}

// Name returns the worker name.
func (w *DevtoSyncWorker) Name() string {
  return "devto-full-sync"
}

// Interval returns the check interval.
func (w *DevtoSyncWorker) Interval() time.Duration {
  return w.config.CheckInterval
}

// Execute checks the distributed schedule and runs a sync cycle if it's time.
func (w *DevtoSyncWorker) Execute(ctx context.Context) error {
  // Check if worker is disabled by admin
  disabledKey := "worker." + w.Name() + ".disabled"

  disabled, err := w.runtimeStates.Get(ctx, disabledKey)
  if err == nil && disabled == disabledStateValue {
    return workerfx.ErrWorkerSkipped
  }

  // Check if it's time to run based on persisted schedule
  nextRunKey := "devto.sync.full_sync_worker.next_run_at"

  nextRunAt, err := w.runtimeStates.GetTime(ctx, nextRunKey)
  if err == nil && time.Now().Before(nextRunAt) {
    return workerfx.ErrWorkerSkipped
  }

  // Try advisory lock to prevent concurrent execution
  acquired, lockErr := w.runtimeStates.TryLock(ctx, lockIDDevtoFullSync)
  if lockErr != nil {
    w.logger.WarnContext(ctx, "Failed to acquire advisory lock for Dev.to sync",
      slog.Any("error", lockErr))

    return workerfx.ErrWorkerSkipped
  }

  if !acquired {
    w.logger.DebugContext(ctx, "Another instance is running Dev.to sync")

    return workerfx.ErrWorkerSkipped
  }

  defer func() {
    releaseErr := w.runtimeStates.ReleaseLock(ctx, lockIDDevtoFullSync)
    if releaseErr != nil {
      w.logger.WarnContext(ctx, "Failed to release advisory lock for Dev.to sync",
        slog.String("error", releaseErr.Error()))
    }
  }()

  // Claim the next slot before executing
  setErr := w.runtimeStates.SetTime(ctx, nextRunKey, time.Now().Add(w.config.FullSyncInterval))
  if setErr != nil {
    w.logger.WarnContext(ctx, "Failed to set next run time for Dev.to sync",
      slog.String("error", setErr.Error()))
  }

  return w.executeSync(ctx)
}

// executeSync runs the actual sync cycle.
func (w *DevtoSyncWorker) executeSync(ctx context.Context) error {
  w.logger.DebugContext(ctx, "Starting Dev.to sync cycle")

  // Get managed Dev.to links (no OAuth tokens needed)
  links, err := w.syncService.GetPublicManagedLinks(ctx, "devto", w.config.BatchSize)
  if err != nil {
    return fmt.Errorf("%w: %w", ErrSyncFailed, err)
  }

  if len(links) == 0 {
    w.logger.DebugContext(ctx, "No Dev.to links to sync")

    return nil
  }

  w.logger.DebugContext(ctx, "Processing Dev.to links",
    slog.Int("count", len(links)))

  // Process each link
  for _, link := range links {
    result, syncErr := w.siteImporterService.SyncPublicLink(ctx, link)
    if syncErr != nil {
      w.logger.ErrorContext(ctx, "Failed to sync Dev.to link",
        slog.String("link_id", link.ID),
        slog.String("profile_id", link.ProfileID),
        slog.Any("error", syncErr))

      continue
    }

    if result.Error != nil {
      w.logger.ErrorContext(ctx, "Dev.to sync returned error",
        slog.String("link_id", link.ID),
        slog.Any("error", result.Error))
    } else {
      w.logger.DebugContext(ctx, "Successfully synced Dev.to link",
        slog.String("link_id", link.ID),
        slog.Int("added", result.ItemsAdded),
        slog.Int("deleted", result.ItemsDeleted))
    }
  }

  w.logger.DebugContext(ctx, "Completed Dev.to sync cycle",
    slog.Int("links_processed", len(links)))

  // Process stories: create new ones from imports
  if w.storyProcessor != nil {
    storyErr := w.storyProcessor.ProcessStories(ctx)
    if storyErr != nil {
      w.logger.ErrorContext(ctx, "Failed to process stories after Dev.to sync",
        slog.Any("error", storyErr))
    }
  }

  return nil
}
