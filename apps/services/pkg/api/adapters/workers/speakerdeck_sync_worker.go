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

const lockIDSpeakerDeckFullSync int64 = 100003

// SpeakerDeckSyncWorker syncs SpeakerDeck presentations for managed profile links.
type SpeakerDeckSyncWorker struct {
	config              *SpeakerDeckSyncConfig
	logger              *logfx.Logger
	syncService         *linksync.Service
	siteImporterService *siteimporter.Service
	storyProcessor      StoryProcessor
	runtimeStates       *runtime_states.Service
	idGenerator         func() string
}

// NewSpeakerDeckSyncWorker creates a new SpeakerDeck sync worker.
func NewSpeakerDeckSyncWorker(
	config *SpeakerDeckSyncConfig,
	logger *logfx.Logger,
	syncService *linksync.Service,
	siteImporterService *siteimporter.Service,
	storyProcessor StoryProcessor,
	runtimeStates *runtime_states.Service,
	idGenerator func() string,
) *SpeakerDeckSyncWorker {
	return &SpeakerDeckSyncWorker{
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
func (w *SpeakerDeckSyncWorker) Name() string {
	return "speakerdeck-full-sync"
}

// Interval returns the check interval.
func (w *SpeakerDeckSyncWorker) Interval() time.Duration {
	return w.config.CheckInterval
}

// Execute checks the distributed schedule and runs a sync cycle if it's time.
func (w *SpeakerDeckSyncWorker) Execute(ctx context.Context) error {
	// Check if worker is disabled by admin
	disabledKey := "worker." + w.Name() + ".disabled"

	disabled, err := w.runtimeStates.Get(ctx, disabledKey)
	if err == nil && disabled == "true" {
		return workerfx.ErrWorkerSkipped
	}

	// Check if it's time to run based on persisted schedule
	nextRunKey := "speakerdeck.sync.full_sync_worker.next_run_at"

	nextRunAt, err := w.runtimeStates.GetTime(ctx, nextRunKey)
	if err == nil && time.Now().Before(nextRunAt) {
		return workerfx.ErrWorkerSkipped
	}

	// Try advisory lock to prevent concurrent execution
	acquired, lockErr := w.runtimeStates.TryLock(ctx, lockIDSpeakerDeckFullSync)
	if lockErr != nil {
		w.logger.WarnContext(ctx, "Failed to acquire advisory lock for SpeakerDeck sync",
			slog.Any("error", lockErr))

		return workerfx.ErrWorkerSkipped
	}

	if !acquired {
		w.logger.DebugContext(ctx, "Another instance is running SpeakerDeck sync")

		return workerfx.ErrWorkerSkipped
	}

	defer func() {
		releaseErr := w.runtimeStates.ReleaseLock(ctx, lockIDSpeakerDeckFullSync)
		if releaseErr != nil {
			w.logger.WarnContext(ctx, "Failed to release advisory lock for SpeakerDeck sync",
				slog.String("error", releaseErr.Error()))
		}
	}()

	// Claim the next slot before executing
	setErr := w.runtimeStates.SetTime(ctx, nextRunKey, time.Now().Add(w.config.FullSyncInterval))
	if setErr != nil {
		w.logger.WarnContext(ctx, "Failed to set next run time for SpeakerDeck sync",
			slog.String("error", setErr.Error()))
	}

	return w.executeSync(ctx)
}

// executeSync runs the actual sync cycle.
func (w *SpeakerDeckSyncWorker) executeSync(ctx context.Context) error {
	w.logger.WarnContext(ctx, "Starting SpeakerDeck sync cycle")

	// Get managed SpeakerDeck links (no OAuth tokens needed)
	links, err := w.syncService.GetPublicManagedLinks(ctx, "speakerdeck", w.config.BatchSize)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSyncFailed, err)
	}

	if len(links) == 0 {
		w.logger.WarnContext(ctx, "No SpeakerDeck links to sync")

		return nil
	}

	w.logger.WarnContext(ctx, "Processing SpeakerDeck links",
		slog.Int("count", len(links)))

	// Process each link
	for _, link := range links {
		result, syncErr := w.siteImporterService.SyncPublicLink(ctx, link)
		if syncErr != nil {
			w.logger.ErrorContext(ctx, "Failed to sync SpeakerDeck link",
				slog.String("link_id", link.ID),
				slog.String("profile_id", link.ProfileID),
				slog.Any("error", syncErr))

			continue
		}

		if result.Error != nil {
			w.logger.ErrorContext(ctx, "SpeakerDeck sync returned error",
				slog.String("link_id", link.ID),
				slog.Any("error", result.Error))
		} else {
			w.logger.WarnContext(ctx, "Successfully synced SpeakerDeck link",
				slog.String("link_id", link.ID),
				slog.Int("added", result.ItemsAdded),
				slog.Int("deleted", result.ItemsDeleted))
		}
	}

	w.logger.WarnContext(ctx, "Completed SpeakerDeck sync cycle",
		slog.Int("links_processed", len(links)))

	// Process stories: create new ones from imports
	if w.storyProcessor != nil {
		err := w.storyProcessor.ProcessStories(ctx)
		if err != nil {
			w.logger.ErrorContext(ctx, "Failed to process stories after SpeakerDeck sync",
				slog.Any("error", err))
		}
	}

	return nil
}
