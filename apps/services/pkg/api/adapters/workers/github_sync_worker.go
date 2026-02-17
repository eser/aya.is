package workers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	"github.com/eser/aya.is/services/pkg/api/business/resourcesync"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
)

const (
	// Advisory lock ID for GitHub resource sync worker.
	lockIDGitHubResourceSync int64 = 100010
)

// GitHubRepoInfoResult holds the repo data fetched from GitHub API.
type GitHubRepoInfoResult struct {
	ID          int64  `json:"id"`
	FullName    string `json:"full_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Private     bool   `json:"private"`
}

// GitHubContributorResult holds a contributor's basic info from GitHub API.
type GitHubContributorResult struct {
	ID            int64  `json:"id"`
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}

// GitHubResourceFetcher defines the interface for fetching GitHub repo data.
// The github.Client will satisfy this interface once FetchRepoInfo, FetchRepoContributors,
// and SearchIssues are added to it.
type GitHubResourceFetcher interface {
	FetchRepoInfo(
		ctx context.Context,
		accessToken string,
		owner string,
		repo string,
	) (*GitHubRepoInfoResult, error)
	FetchRepoContributors(
		ctx context.Context,
		accessToken string,
		owner string,
		repo string,
	) ([]*GitHubContributorResult, error)
	SearchIssues(ctx context.Context, accessToken string, query string) (int, error)
}

// GitHubSyncWorker syncs GitHub contributor stats for registered profile resources.
type GitHubSyncWorker struct {
	config        *GitHubSyncConfig
	logger        *logfx.Logger
	syncService   *resourcesync.Service
	fetcher       GitHubResourceFetcher
	runtimeStates *runtime_states.Service
}

// NewGitHubSyncWorker creates a new GitHub resource sync worker.
func NewGitHubSyncWorker(
	config *GitHubSyncConfig,
	logger *logfx.Logger,
	syncService *resourcesync.Service,
	fetcher GitHubResourceFetcher,
	runtimeStates *runtime_states.Service,
) *GitHubSyncWorker {
	return &GitHubSyncWorker{
		config:        config,
		logger:        logger,
		syncService:   syncService,
		fetcher:       fetcher,
		runtimeStates: runtimeStates,
	}
}

// Name returns the worker name.
func (w *GitHubSyncWorker) Name() string {
	return "github-resource-sync"
}

// Interval returns the check interval (how often to poll for schedule readiness).
func (w *GitHubSyncWorker) Interval() time.Duration {
	return w.config.CheckInterval
}

// Execute checks the distributed schedule and runs a sync cycle if it's time.
func (w *GitHubSyncWorker) Execute(ctx context.Context) error {
	// Check if worker is disabled by admin
	disabledKey := "worker." + w.Name() + ".disabled"

	disabled, err := w.runtimeStates.Get(ctx, disabledKey)
	if err == nil && disabled == "true" {
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
	acquired, lockErr := w.runtimeStates.TryLock(ctx, lockIDGitHubResourceSync)
	if lockErr != nil {
		w.logger.WarnContext(ctx, "Failed to acquire advisory lock",
			slog.String("worker", w.Name()),
			slog.Any("error", lockErr))

		return workerfx.ErrWorkerSkipped
	}

	if !acquired {
		w.logger.DebugContext(ctx, "Another instance is running this worker",
			slog.String("worker", w.Name()))

		return workerfx.ErrWorkerSkipped
	}

	defer func() {
		releaseErr := w.runtimeStates.ReleaseLock(ctx, lockIDGitHubResourceSync)
		if releaseErr != nil {
			w.logger.WarnContext(ctx, "failed to release advisory lock",
				slog.String("worker", w.Name()),
				slog.String("error", releaseErr.Error()))
		}
	}()

	// Claim the next slot before executing
	setErr := w.runtimeStates.SetTime(ctx, nextRunKey, time.Now().Add(w.config.FullSyncInterval))
	if setErr != nil {
		w.logger.WarnContext(ctx, "failed to set next run time",
			slog.String("worker", w.Name()),
			slog.String("error", setErr.Error()))
	}

	return w.executeSync(ctx)
}

func (w *GitHubSyncWorker) stateKey(suffix string) string {
	return "github.resource_sync_worker." + suffix
}

// executeSync runs the actual sync cycle.
func (w *GitHubSyncWorker) executeSync(ctx context.Context) error {
	w.logger.WarnContext(ctx, "Starting GitHub resource sync cycle")

	// Get managed GitHub resources
	resources, err := w.syncService.GetGitHubResourcesForSync(ctx, w.config.BatchSize)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSyncFailed, err)
	}

	if len(resources) == 0 {
		w.logger.WarnContext(ctx, "No GitHub resources to sync")

		return nil
	}

	w.logger.WarnContext(ctx, "Processing GitHub resources",
		slog.Int("count", len(resources)))

	// Process each resource (isolated errors - don't fail the whole batch)
	for _, resource := range resources {
		syncErr := w.syncResource(ctx, resource)
		if syncErr != nil {
			w.logger.ErrorContext(ctx, "Failed to sync GitHub resource",
				slog.String("resource_id", resource.ResourceID),
				slog.String("profile_id", resource.ProfileID),
				slog.String("public_id", resource.ResourcePublicID),
				slog.Any("error", syncErr))
		} else {
			w.logger.WarnContext(ctx, "Successfully synced GitHub resource",
				slog.String("resource_id", resource.ResourceID),
				slog.String("public_id", resource.ResourcePublicID))
		}
	}

	w.logger.WarnContext(ctx, "Completed GitHub resource sync cycle",
		slog.Int("resources_processed", len(resources)))

	return nil
}

// syncResource syncs a single GitHub resource (repo).
func (w *GitHubSyncWorker) syncResource( //nolint:cyclop,funlen
	ctx context.Context,
	resource *resourcesync.GitHubResourceForSync,
) error {
	// Parse owner/repo from ResourcePublicID
	owner, repo, ok := parseOwnerRepo(resource.ResourcePublicID)
	if !ok {
		return fmt.Errorf(
			"invalid resource public_id format: %s",
			resource.ResourcePublicID,
		) //nolint:goerr113
	}

	accessToken := resource.AuthAccessToken

	// Fetch latest repo info (stars, forks, language, description)
	repoInfo, err := w.fetcher.FetchRepoInfo(ctx, accessToken, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to fetch repo info: %w", err)
	}

	// Update resource properties with latest repo data
	resourceProps := mergeResourceProperties(resource.ResourceProperties, repoInfo)

	err = w.syncService.UpdateResourceProperties(ctx, resource.ResourceID, resourceProps)
	if err != nil {
		return fmt.Errorf("failed to update resource properties: %w", err)
	}

	// Fetch contributors and their stats
	contributors, err := w.fetcher.FetchRepoContributors(ctx, accessToken, owner, repo)
	if err != nil {
		w.logger.WarnContext(ctx, "Failed to fetch contributors, skipping contributor sync",
			slog.String("resource_id", resource.ResourceID),
			slog.Any("error", err))

		return nil
	}

	// Process each contributor
	for _, contributor := range contributors {
		w.syncContributor(ctx, resource, accessToken, owner, repo, contributor, repoInfo.Stars)
	}

	return nil
}

// syncContributor processes a single contributor's stats and updates the membership if found.
func (w *GitHubSyncWorker) syncContributor(
	ctx context.Context,
	resource *resourcesync.GitHubResourceForSync,
	accessToken string,
	owner string,
	repo string,
	contributor *GitHubContributorResult,
	stars int,
) {
	// Build contributor stats
	stats := w.buildContributorStats(ctx, accessToken, owner, repo, contributor, stars)

	// Try to match contributor to a profile via profile_link
	contributorRemoteID := strconv.FormatInt(contributor.ID, 10)

	profileID, err := w.syncService.GetProfileLinkByRemoteID(ctx, "github", contributorRemoteID)
	if err != nil {
		// No matching profile found - this is normal for external contributors
		return
	}

	if profileID == "" {
		return
	}

	// Look for a membership between the resource's profile and the contributor's profile
	membershipID, err := w.syncService.GetMembershipByProfiles(ctx, resource.ProfileID, profileID)
	if err != nil {
		// No membership exists - we do NOT auto-create memberships
		return
	}

	if membershipID == "" {
		return
	}

	// Update membership properties with contributor stats
	err = w.updateMembershipWithStats(ctx, membershipID, resource.ResourcePublicID, stats)
	if err != nil {
		w.logger.WarnContext(ctx, "Failed to update membership properties",
			slog.String("membership_id", membershipID),
			slog.String("contributor", contributor.Login),
			slog.Any("error", err))
	}
}

// buildContributorStats fetches and assembles contributor stats from GitHub API.
func (w *GitHubSyncWorker) buildContributorStats(
	ctx context.Context,
	accessToken string,
	owner string,
	repo string,
	contributor *GitHubContributorResult,
	stars int,
) *resourcesync.GitHubContributorStats {
	stats := &resourcesync.GitHubContributorStats{} //nolint:exhaustruct
	stats.Commits = contributor.Contributions
	stats.Stars = stars

	repoQualifier := "repo:" + owner + "/" + repo

	// Fetch total PRs
	totalPRs, err := w.fetcher.SearchIssues(
		ctx, accessToken,
		repoQualifier+" type:pr author:"+contributor.Login,
	)
	if err == nil {
		stats.PRs.Total = totalPRs
	}

	// Fetch resolved (merged) PRs
	resolvedPRs, err := w.fetcher.SearchIssues(
		ctx, accessToken,
		repoQualifier+" type:pr author:"+contributor.Login+" is:merged",
	)
	if err == nil {
		stats.PRs.Resolved = resolvedPRs
	}

	// Fetch total issues
	totalIssues, err := w.fetcher.SearchIssues(
		ctx, accessToken,
		repoQualifier+" type:issue author:"+contributor.Login,
	)
	if err == nil {
		stats.Issues.Total = totalIssues
	}

	// Fetch resolved (closed) issues
	resolvedIssues, err := w.fetcher.SearchIssues(
		ctx, accessToken,
		repoQualifier+" type:issue author:"+contributor.Login+" is:closed",
	)
	if err == nil {
		stats.Issues.Resolved = resolvedIssues
	}

	return stats
}

// updateMembershipWithStats merges GitHub contributor stats into membership properties.
func (w *GitHubSyncWorker) updateMembershipWithStats(
	ctx context.Context,
	membershipID string,
	resourcePublicID string,
	stats *resourcesync.GitHubContributorStats,
) error {
	// Build the properties map with github stats nested under the resource key
	properties := map[string]any{
		"github": map[string]any{
			resourcePublicID: stats,
		},
	}

	return w.syncService.UpdateMembershipProperties(ctx, membershipID, properties)
}

// parseOwnerRepo splits "owner/repo" into its parts.
func parseOwnerRepo(publicID string) (string, string, bool) {
	parts := strings.SplitN(publicID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}

	return parts[0], parts[1], true
}

// mergeResourceProperties merges the latest repo info into existing resource properties.
func mergeResourceProperties(
	existing map[string]any,
	repoInfo *GitHubRepoInfoResult,
) map[string]any {
	if existing == nil {
		existing = make(map[string]any)
	}

	existing["stars"] = repoInfo.Stars
	existing["forks"] = repoInfo.Forks
	existing["language"] = repoInfo.Language
	existing["description"] = repoInfo.Description
	existing["html_url"] = repoInfo.HTMLURL
	existing["name"] = repoInfo.Name
	existing["full_name"] = repoInfo.FullName
	existing["private"] = repoInfo.Private

	return existing
}
