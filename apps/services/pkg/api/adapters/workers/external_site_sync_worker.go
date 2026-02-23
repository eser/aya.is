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

const lockIDExternalSiteFullSync int64 = 100005

// ExternalSiteSyncWorker syncs external site content for managed profile links.
type ExternalSiteSyncWorker struct {
	config              *ExternalSiteSyncConfig
	logger              *logfx.Logger
	syncService         *linksync.Service
	siteImporterService *siteimporter.Service
	storyProcessor      StoryProcessor
	runtimeStates       *runtime_states.Service
	idGenerator         func() string
}

// NewExternalSiteSyncWorker creates a new external site sync worker.
func NewExternalSiteSyncWorker(
	config *ExternalSiteSyncConfig,
	logger *logfx.Logger,
	syncService *linksync.Service,
	siteImporterService *siteimporter.Service,
	storyProcessor StoryProcessor,
	runtimeStates *runtime_states.Service,
	idGenerator func() string,
) *ExternalSiteSyncWorker {
	return &ExternalSiteSyncWorker{
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
func (w *ExternalSiteSyncWorker) Name() string {
	return "external-site-full-sync"
}

// Interval returns the check interval.
func (w *ExternalSiteSyncWorker) Interval() time.Duration {
	return w.config.CheckInterval
}

// Execute checks the distributed schedule and runs a sync cycle if it's time.
func (w *ExternalSiteSyncWorker) Execute(ctx context.Context) error {
	// Check if worker is disabled by admin
	disabledKey := "worker." + w.Name() + ".disabled"

	disabled, err := w.runtimeStates.Get(ctx, disabledKey)
	if err == nil && disabled == "true" {
		return workerfx.ErrWorkerSkipped
	}

	// Check if it's time to run based on persisted schedule
	nextRunKey := "external-site.sync.full_sync_worker.next_run_at"

	nextRunAt, err := w.runtimeStates.GetTime(ctx, nextRunKey)
	if err == nil && time.Now().Before(nextRunAt) {
		return workerfx.ErrWorkerSkipped
	}

	// Try advisory lock to prevent concurrent execution
	acquired, lockErr := w.runtimeStates.TryLock(ctx, lockIDExternalSiteFullSync)
	if lockErr != nil {
		w.logger.WarnContext(ctx, "Failed to acquire advisory lock for external site sync",
			slog.Any("error", lockErr))

		return workerfx.ErrWorkerSkipped
	}

	if !acquired {
		w.logger.DebugContext(ctx, "Another instance is running external site sync")

		return workerfx.ErrWorkerSkipped
	}

	defer func() {
		releaseErr := w.runtimeStates.ReleaseLock(ctx, lockIDExternalSiteFullSync)
		if releaseErr != nil {
			w.logger.WarnContext(ctx, "Failed to release advisory lock for external site sync",
				slog.String("error", releaseErr.Error()))
		}
	}()

	// Claim the next slot before executing
	setErr := w.runtimeStates.SetTime(ctx, nextRunKey, time.Now().Add(w.config.FullSyncInterval))
	if setErr != nil {
		w.logger.WarnContext(ctx, "Failed to set next run time for external site sync",
			slog.String("error", setErr.Error()))
	}

	return w.executeSync(ctx)
}

// executeSync runs the actual sync cycle.
func (w *ExternalSiteSyncWorker) executeSync(ctx context.Context) error {
	w.logger.WarnContext(ctx, "Starting external site sync cycle")

	// Get managed external-site links (no OAuth tokens needed)
	links, err := w.syncService.GetPublicManagedLinks(ctx, "external-site", w.config.BatchSize)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSyncFailed, err)
	}

	if len(links) == 0 {
		w.logger.WarnContext(ctx, "No external site links to sync")

		return nil
	}

	w.logger.WarnContext(ctx, "Processing external site links",
		slog.Int("count", len(links)))

	// Process each link
	for _, link := range links {
		result, syncErr := w.siteImporterService.SyncPublicLink(ctx, link)
		if syncErr != nil {
			w.logger.ErrorContext(ctx, "Failed to sync external site link",
				slog.String("link_id", link.ID),
				slog.String("profile_id", link.ProfileID),
				slog.Any("error", syncErr))

			continue
		}

		if result.Error != nil {
			w.logger.ErrorContext(ctx, "External site sync returned error",
				slog.String("link_id", link.ID),
				slog.Any("error", result.Error))
		} else {
			w.logger.WarnContext(ctx, "Successfully synced external site link",
				slog.String("link_id", link.ID),
				slog.Int("added", result.ItemsAdded),
				slog.Int("deleted", result.ItemsDeleted))
		}
	}

	w.logger.WarnContext(ctx, "Completed external site sync cycle",
		slog.Int("links_processed", len(links)))

	// Process stories: create new ones from imports
	if w.storyProcessor != nil {
		storyErr := w.storyProcessor.ProcessStories(ctx)
		if storyErr != nil {
			w.logger.ErrorContext(ctx, "Failed to process stories after external site sync",
				slog.Any("error", storyErr))
		}
	}

	return nil
}
