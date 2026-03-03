package workers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/youtube"
	"github.com/eser/aya.is/services/pkg/api/business/linksync"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
)

// Advisory lock ID for YouTube live status worker.
const lockIDYouTubeLiveStatus int64 = 100003

// LiveStatusChecker defines the interface for checking live broadcast status and refreshing tokens.
type LiveStatusChecker interface {
	CheckLiveBroadcasts(
		ctx context.Context,
		accessToken string,
	) (*youtube.LiveBroadcastResult, error)
	RefreshAccessToken(
		ctx context.Context,
		refreshToken string,
	) (*linksync.TokenRefreshResult, error)
}

// YouTubeLiveStatusWorker periodically checks managed YouTube links for active live broadcasts
// and updates their online status accordingly.
type YouTubeLiveStatusWorker struct {
	config        *YouTubeLiveStatusConfig
	logger        *logfx.Logger
	syncService   *linksync.Service
	checker       LiveStatusChecker
	runtimeStates *runtime_states.Service
}

// NewYouTubeLiveStatusWorker creates a new YouTube live status worker.
func NewYouTubeLiveStatusWorker(
	config *YouTubeLiveStatusConfig,
	logger *logfx.Logger,
	syncService *linksync.Service,
	checker LiveStatusChecker,
	runtimeStates *runtime_states.Service,
) *YouTubeLiveStatusWorker {
	return &YouTubeLiveStatusWorker{
		config:        config,
		logger:        logger,
		syncService:   syncService,
		checker:       checker,
		runtimeStates: runtimeStates,
	}
}

// Name returns the worker name.
func (w *YouTubeLiveStatusWorker) Name() string {
	return "youtube-live-status"
}

// Interval returns the check interval (how often to poll for schedule readiness).
func (w *YouTubeLiveStatusWorker) Interval() time.Duration {
	return w.config.CheckInterval
}

// Execute checks the distributed schedule and runs a live status check cycle if it's time.
func (w *YouTubeLiveStatusWorker) Execute(ctx context.Context) error {
	// Check if worker is disabled by admin
	disabledKey := "worker." + w.Name() + ".disabled"

	disabled, err := w.runtimeStates.Get(ctx, disabledKey)
	if err == nil && disabled == disabledStateValue {
		return workerfx.ErrWorkerSkipped
	}

	// Check if it's time to run based on persisted schedule
	nextRunKey := w.stateKey("next_run_at")

	nextRunAt, err := w.runtimeStates.GetTime(ctx, nextRunKey)
	if err == nil && time.Now().Before(nextRunAt) {
		return workerfx.ErrWorkerSkipped
	}
	// If ErrStateNotFound or ErrInvalidTime, proceed (first run or corrupted state)

	// Try advisory lock to prevent concurrent execution across instances
	acquired, lockErr := w.runtimeStates.TryLock(ctx, lockIDYouTubeLiveStatus)
	if lockErr != nil {
		w.logger.WarnContext(ctx, "Failed to acquire advisory lock for live status worker",
			slog.Any("error", lockErr))

		return workerfx.ErrWorkerSkipped
	}

	if !acquired {
		w.logger.DebugContext(ctx, "Another instance is running the live status worker")

		return workerfx.ErrWorkerSkipped
	}

	defer func() {
		releaseErr := w.runtimeStates.ReleaseLock(ctx, lockIDYouTubeLiveStatus)
		if releaseErr != nil {
			w.logger.WarnContext(ctx, "failed to release advisory lock for live status worker",
				slog.String("error", releaseErr.Error()))
		}
	}()

	// Claim the next slot before executing
	setErr := w.runtimeStates.SetTime(ctx, nextRunKey, time.Now().Add(w.config.SyncInterval))
	if setErr != nil {
		w.logger.WarnContext(ctx, "failed to set next run time for live status worker",
			slog.String("error", setErr.Error()))
	}

	return w.executeCheck(ctx)
}

func (w *YouTubeLiveStatusWorker) stateKey(suffix string) string {
	return "youtube.live_status_worker." + suffix
}

// executeCheck runs the actual live status check cycle.
func (w *YouTubeLiveStatusWorker) executeCheck(ctx context.Context) error { //nolint:cyclop,funlen
	w.logger.DebugContext(ctx, "Starting YouTube live status check cycle")

	// Get managed YouTube links (need tokens for the API call)
	links, err := w.syncService.GetManagedLinks(ctx, "youtube", w.config.BatchSize)
	if err != nil {
		return fmt.Errorf("getting managed YouTube links: %w", err)
	}

	if len(links) == 0 {
		w.logger.DebugContext(ctx, "No YouTube links to check for live status")

		return nil
	}

	now := time.Now().UTC()
	checkedCount := 0
	liveCount := 0

	for _, link := range links {
		accessToken := link.AuthAccessToken

		// Check if token needs refresh
		if w.needsTokenRefresh(link) {
			refreshedToken, refreshErr := w.refreshTokenIfPossible(ctx, link)
			if refreshErr != nil {
				w.logger.WarnContext(ctx, "Failed to refresh token for live status check",
					slog.String("link_id", link.ID),
					slog.Any("error", refreshErr))

				continue
			}

			if refreshedToken != "" {
				accessToken = refreshedToken
			}
		}

		// Check live broadcast status
		result, checkErr := w.checker.CheckLiveBroadcasts(ctx, accessToken)
		if checkErr != nil {
			w.logger.WarnContext(ctx, "Failed to check live broadcasts",
				slog.String("link_id", link.ID),
				slog.Any("error", checkErr))

			continue
		}

		// Build online properties
		var onlineProperties map[string]any

		if result.IsLive {
			liveCount++

			onlineInfo := map[string]any{
				"broadcast_id":    result.BroadcastID,
				"broadcast_url":   "https://www.youtube.com/watch?v=" + result.BroadcastID,
				"title":           result.Title,
				"thumbnail_url":   result.ThumbnailURL,
				"last_checked_at": now.Format(time.RFC3339),
			}

			if result.StartedAt != nil {
				onlineInfo["started_at"] = result.StartedAt.Format(time.RFC3339)
			}

			onlineProperties = map[string]any{
				"online_information": onlineInfo,
			}
		} else {
			onlineProperties = map[string]any{
				"online_information": map[string]any{
					"last_checked_at": now.Format(time.RFC3339),
				},
			}
		}

		// Update link online status
		updateErr := w.syncService.UpdateLinkOnlineStatus(
			ctx,
			link.ID,
			result.IsLive,
			onlineProperties,
		)
		if updateErr != nil {
			w.logger.WarnContext(ctx, "Failed to update link online status",
				slog.String("link_id", link.ID),
				slog.Any("error", updateErr))

			continue
		}

		checkedCount++
	}

	w.logger.WarnContext(ctx, "Completed YouTube live status check cycle",
		slog.Int("checked", checkedCount),
		slog.Int("live", liveCount),
		slog.Int("total_links", len(links)))

	return nil
}

// refreshTokenIfPossible refreshes the access token if a refresh token is available.
// Returns the new access token (empty string if no refresh was performed) or an error.
func (w *YouTubeLiveStatusWorker) refreshTokenIfPossible(
	ctx context.Context,
	link *linksync.ManagedLink,
) (string, error) {
	if link.AuthRefreshToken == nil {
		// No refresh token — try with existing access token anyway
		w.logger.WarnContext(
			ctx,
			"No refresh token available for live status check, trying existing token",
			slog.String("link_id", link.ID),
		)

		return "", nil
	}

	w.logger.DebugContext(ctx, "Refreshing YouTube access token for live status check",
		slog.String("link_id", link.ID))

	refreshResult, err := w.checker.RefreshAccessToken(ctx, *link.AuthRefreshToken)
	if err != nil {
		return "", fmt.Errorf("refreshing access token: %w", err)
	}

	// Preserve existing refresh token if Google didn't return a new one
	refreshToken := refreshResult.RefreshToken
	if refreshToken == nil {
		refreshToken = link.AuthRefreshToken
	}

	// Update tokens in database
	err = w.syncService.UpdateLinkTokens(
		ctx,
		link.ID,
		refreshResult.AccessToken,
		refreshResult.AccessTokenExpiresAt,
		refreshToken,
	)
	if err != nil {
		return "", fmt.Errorf("updating link tokens: %w", err)
	}

	return refreshResult.AccessToken, nil
}

// needsTokenRefresh checks if the access token needs to be refreshed.
func (w *YouTubeLiveStatusWorker) needsTokenRefresh(link *linksync.ManagedLink) bool {
	if link.AuthAccessTokenExpiresAt == nil {
		// No expiry set, assume token is valid
		return false
	}

	// Refresh if token expires within the buffer time
	return time.Until(*link.AuthAccessTokenExpiresAt) < w.config.TokenRefreshBuffer
}
