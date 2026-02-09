package workers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/linksync"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
)

// Sentinel errors.
var ErrSyncFailed = errors.New("sync failed")

// SyncMode defines the type of sync operation.
type SyncMode string

const (
	SyncModeFull        SyncMode = "full"
	SyncModeIncremental SyncMode = "incremental"

	// Advisory lock IDs for YouTube sync workers.
	lockIDYouTubeFullSync        int64 = 100001
	lockIDYouTubeIncrementalSync int64 = 100002
)

// RemoteStoryFetcher defines the interface for fetching stories from remote providers.
type RemoteStoryFetcher interface {
	FetchRemoteStories(
		ctx context.Context,
		accessToken string,
		remoteSourceID string,
		publishedAfter *time.Time,
		maxResults int,
	) ([]*linksync.RemoteStoryItem, error)

	RefreshAccessToken(
		ctx context.Context,
		refreshToken string,
	) (*linksync.TokenRefreshResult, error)
}

// StoryProcessor handles story creation and reconciliation from link imports.
type StoryProcessor interface {
	ProcessStories(ctx context.Context) error
}

// YouTubeSyncWorker syncs YouTube videos for managed profile links.
type YouTubeSyncWorker struct {
	config         *YouTubeSyncConfig
	logger         *logfx.Logger
	syncService    *linksync.Service
	fetcher        RemoteStoryFetcher
	idGenerator    func() string
	runtimeStates  *runtime_states.Service
	storyProcessor StoryProcessor
	mode           SyncMode
	syncInterval   time.Duration
	lockID         int64
}

// NewYouTubeFullSyncWorker creates a new YouTube full sync worker.
func NewYouTubeFullSyncWorker(
	config *YouTubeSyncConfig,
	logger *logfx.Logger,
	syncService *linksync.Service,
	fetcher RemoteStoryFetcher,
	idGenerator func() string,
	runtimeStates *runtime_states.Service,
	storyProcessor StoryProcessor,
) *YouTubeSyncWorker {
	return &YouTubeSyncWorker{
		config:         config,
		logger:         logger,
		syncService:    syncService,
		fetcher:        fetcher,
		idGenerator:    idGenerator,
		runtimeStates:  runtimeStates,
		storyProcessor: storyProcessor,
		mode:           SyncModeFull,
		syncInterval:   config.FullSyncInterval,
		lockID:         lockIDYouTubeFullSync,
	}
}

// NewYouTubeIncrementalSyncWorker creates a new YouTube incremental sync worker.
func NewYouTubeIncrementalSyncWorker(
	config *YouTubeSyncConfig,
	logger *logfx.Logger,
	syncService *linksync.Service,
	fetcher RemoteStoryFetcher,
	idGenerator func() string,
	runtimeStates *runtime_states.Service,
	storyProcessor StoryProcessor,
) *YouTubeSyncWorker {
	return &YouTubeSyncWorker{
		config:         config,
		logger:         logger,
		syncService:    syncService,
		fetcher:        fetcher,
		idGenerator:    idGenerator,
		runtimeStates:  runtimeStates,
		storyProcessor: storyProcessor,
		mode:           SyncModeIncremental,
		syncInterval:   config.IncrementalSyncInterval,
		lockID:         lockIDYouTubeIncrementalSync,
	}
}

// Name returns the worker name.
func (w *YouTubeSyncWorker) Name() string {
	return "youtube-" + string(w.mode) + "-sync"
}

// Interval returns the check interval (how often to poll for schedule readiness).
func (w *YouTubeSyncWorker) Interval() time.Duration {
	return w.config.CheckInterval
}

// Execute checks the distributed schedule and runs a sync cycle if it's time.
func (w *YouTubeSyncWorker) Execute(ctx context.Context) error {
	// Check if worker is disabled by admin
	disabledKey := "worker." + w.Name() + ".disabled"

	disabled, err := w.runtimeStates.Get(ctx, disabledKey)
	if err == nil && disabled == "true" {
		return nil
	}

	// Check if it's time to run based on persisted schedule
	nextRunKey := w.stateKey("next_run_at")

	nextRunAt, err := w.runtimeStates.GetTime(ctx, nextRunKey)
	if err == nil && time.Now().Before(nextRunAt) {
		return nil
	}
	// If ErrStateNotFound or ErrInvalidTime, proceed (first run or corrupted state)

	// Try advisory lock to prevent concurrent execution across instances
	acquired, lockErr := w.runtimeStates.TryLock(ctx, w.lockID)
	if lockErr != nil {
		w.logger.WarnContext(ctx, "Failed to acquire advisory lock",
			slog.String("mode", string(w.mode)),
			slog.Any("error", lockErr))

		return nil
	}

	if !acquired {
		w.logger.DebugContext(ctx, "Another instance is running this worker",
			slog.String("mode", string(w.mode)))

		return nil
	}

	defer func() {
		_ = w.runtimeStates.ReleaseLock(ctx, w.lockID)
	}()

	// Claim the next slot before executing
	_ = w.runtimeStates.SetTime(ctx, nextRunKey, time.Now().Add(w.syncInterval))

	return w.executeSync(ctx)
}

func (w *YouTubeSyncWorker) stateKey(suffix string) string {
	return "youtube.sync." + string(w.mode) + "_sync_worker." + suffix
}

// executeSync runs the actual sync cycle.
func (w *YouTubeSyncWorker) executeSync(ctx context.Context) error {
	w.logger.WarnContext(ctx, "Starting YouTube sync cycle",
		slog.String("mode", string(w.mode)))

	// Get managed YouTube links
	links, err := w.syncService.GetManagedLinks(ctx, "youtube", w.config.BatchSize)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSyncFailed, err)
	}

	if len(links) == 0 {
		w.logger.WarnContext(ctx, "No YouTube links to sync")

		return nil
	}

	w.logger.WarnContext(ctx, "Processing YouTube links",
		slog.String("mode", string(w.mode)),
		slog.Int("count", len(links)))

	// Process each link (isolated errors - don't fail the whole batch)
	for _, link := range links {
		result := w.syncLink(ctx, link)

		if result.Error != nil {
			w.logger.ErrorContext(ctx, "Failed to sync YouTube link",
				slog.String("link_id", link.ID),
				slog.String("profile_id", link.ProfileID),
				slog.String("mode", string(w.mode)),
				slog.Any("error", result.Error))
		} else {
			w.logger.WarnContext(ctx, "Successfully synced YouTube link",
				slog.String("link_id", link.ID),
				slog.String("mode", string(w.mode)),
				slog.Int("added", result.ItemsAdded),
				slog.Int("updated", result.ItemsUpdated),
				slog.Int("deleted", result.ItemsDeleted))
		}
	}

	w.logger.WarnContext(ctx, "Completed YouTube sync cycle",
		slog.String("mode", string(w.mode)),
		slog.Int("links_processed", len(links)))

	// Process stories: create new ones from imports and reconcile existing ones
	if w.storyProcessor != nil {
		err := w.storyProcessor.ProcessStories(ctx)
		if err != nil {
			w.logger.ErrorContext(ctx, "Failed to process stories after sync",
				slog.String("mode", string(w.mode)),
				slog.Any("error", err))
		}
	}

	return nil
}

// syncLink syncs a single YouTube link.
func (w *YouTubeSyncWorker) syncLink( //nolint:cyclop,funlen
	ctx context.Context,
	link *linksync.ManagedLink,
) *linksync.SyncResult {
	result := &linksync.SyncResult{LinkID: link.ID} //nolint:exhaustruct

	accessToken := link.AuthAccessToken

	// Check if token needs refresh
	if w.needsTokenRefresh(link) {
		if link.AuthRefreshToken == nil {
			// No refresh token - try with existing access token anyway
			// The API will return 401 if it's truly expired
			w.logger.WarnContext(ctx, "No refresh token available, trying existing access token",
				slog.String("link_id", link.ID))
		} else {
			w.logger.DebugContext(ctx, "Refreshing YouTube access token",
				slog.String("link_id", link.ID))

			refreshResult, err := w.fetcher.RefreshAccessToken(ctx, *link.AuthRefreshToken)
			if err != nil {
				result.Error = err

				return result
			}

			// Update tokens in database
			err = w.syncService.UpdateLinkTokens(
				ctx,
				link.ID,
				refreshResult.AccessToken,
				refreshResult.AccessTokenExpiresAt,
				refreshResult.RefreshToken,
			)
			if err != nil {
				result.Error = err

				return result
			}

			accessToken = refreshResult.AccessToken
		}
	}

	// Determine publishedAfter based on sync mode
	var publishedAfter *time.Time

	if w.mode == SyncModeIncremental {
		// Get last sync time for incremental sync
		var err error

		publishedAfter, err = w.syncService.GetLastSyncTime(ctx, link.ID)
		if err != nil {
			result.Error = err

			return result
		}
	}
	// For full sync, publishedAfter remains nil (fetch all)

	// Determine max results based on sync mode
	maxResults := w.config.StoriesPerLink
	if w.mode == SyncModeFull {
		maxResults = w.config.FullSyncMaxStories
	}

	// Fetch stories from YouTube
	stories, err := w.fetcher.FetchRemoteStories(
		ctx,
		accessToken,
		link.RemoteID,
		publishedAfter,
		maxResults,
	)
	if err != nil {
		result.Error = err

		return result
	}

	// Track active remote IDs for deletion marking (only used in full sync)
	activeRemoteIDs := make([]string, 0, len(stories))

	// Upsert each story
	for _, story := range stories {
		err = w.syncService.UpsertImport(ctx, link.ID, story.RemoteID, story.Properties)
		if err != nil {
			// Log but continue with other stories
			w.logger.WarnContext(ctx, "Failed to upsert story",
				slog.String("link_id", link.ID),
				slog.String("remote_id", story.RemoteID),
				slog.Any("error", err))

			continue
		}

		activeRemoteIDs = append(activeRemoteIDs, story.RemoteID)
		result.ItemsAdded++ // Simplified - not distinguishing add vs update
	}

	// Mark deleted items only in full sync mode
	if w.mode == SyncModeFull && len(activeRemoteIDs) > 0 {
		deletedCount, err := w.syncService.MarkDeletedImports(ctx, link.ID, activeRemoteIDs)
		if err != nil {
			w.logger.WarnContext(ctx, "Failed to mark deleted imports",
				slog.String("link_id", link.ID),
				slog.Any("error", err))
		} else {
			result.ItemsDeleted = int(deletedCount)
		}
	}

	return result
}

// needsTokenRefresh checks if the access token needs to be refreshed.
func (w *YouTubeSyncWorker) needsTokenRefresh(link *linksync.ManagedLink) bool {
	if link.AuthAccessTokenExpiresAt == nil {
		// No expiry set, assume token is valid
		return false
	}

	// Refresh if token expires within the buffer time
	return time.Until(*link.AuthAccessTokenExpiresAt) < w.config.TokenRefreshBuffer
}
