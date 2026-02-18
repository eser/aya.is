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
	// SearchIssueCountsBatch executes multiple search queries in a single GraphQL request.
	// Returns map[alias]count. Falls back to per-query REST if GraphQL is unavailable.
	SearchIssueCountsBatch(
		ctx context.Context,
		accessToken string,
		queries map[string]string,
	) (map[string]int, error)
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

// membershipStatsAccumulator aggregates GitHub stats across multiple repos for a single membership.
type membershipStatsAccumulator struct {
	commits        int
	prsTotal       int
	prsResolved    int
	issuesTotal    int
	issuesResolved int
	stars          int
}

// contributorRepoStats holds per-repo data collected during the collect phase.
type contributorRepoStats struct {
	owner       string
	repo        string
	stars       int
	commits     int
	accessToken string
}

// contributorInfo tracks a contributor's appearances across all resources.
type contributorInfo struct {
	login string
	// Key: resource profileID → repos this contributor appeared in under that profile.
	reposByProfile map[string][]contributorRepoStats
}

// executeSync runs the actual sync cycle using a batch-match-first approach:
//  1. Collect: fetch repo info + contributors for all resources
//  2. Match: batch DB queries to find profiles and memberships
//  3. Stats: fetch search stats ONLY for matched contributors
//  4. Flush: write aggregated stats per membership
func (w *GitHubSyncWorker) executeSync(ctx context.Context) error { //nolint:cyclop,funlen
	w.logger.WarnContext(ctx, "Starting GitHub resource sync cycle")

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

	// ── Phase 1: Collect ──
	// Fetch repo info + contributors for all resources, build contributor map.
	contributorMap := make(map[string]*contributorInfo) // key: GitHub remote ID

	for _, resource := range resources {
		w.collectResourceData(ctx, resource, contributorMap)
	}

	w.logger.WarnContext(ctx, "Collected contributors across all resources",
		slog.Int("unique_contributors", len(contributorMap)))

	if len(contributorMap) == 0 {
		w.logger.WarnContext(ctx, "No contributors found, skipping match phase")

		return nil
	}

	// ── Phase 2: Batch Match ──
	// Batch-load profile links and memberships to find which contributors are tracked.
	membershipStats, matchedCount := w.batchMatchContributors(ctx, contributorMap)

	w.logger.WarnContext(ctx, "Matched contributors to memberships",
		slog.Int("matched_pairs", matchedCount),
		slog.Int("memberships", len(membershipStats)))

	// ── Phase 3: Flush ──
	for membershipID, acc := range membershipStats {
		flushErr := w.flushMembershipStats(ctx, membershipID, acc)
		if flushErr != nil {
			w.logger.WarnContext(ctx, "Failed to flush membership stats",
				slog.String("membership_id", membershipID),
				slog.Any("error", flushErr))
		}
	}

	w.logger.WarnContext(ctx, "Completed GitHub resource sync cycle",
		slog.Int("resources_processed", len(resources)),
		slog.Int("memberships_updated", len(membershipStats)))

	return nil
}

// collectResourceData fetches repo info and contributors for a single resource,
// updates resource properties, and collects contributor appearances.
func (w *GitHubSyncWorker) collectResourceData(
	ctx context.Context,
	resource *resourcesync.GitHubResourceForSync,
	contributorMap map[string]*contributorInfo,
) {
	owner, repo, ok := parseOwnerRepo(resource.ResourcePublicID)
	if !ok {
		w.logger.ErrorContext(ctx, "Invalid resource public_id format",
			slog.String("resource_id", resource.ResourceID),
			slog.String("public_id", resource.ResourcePublicID))

		return
	}

	accessToken := resource.AuthAccessToken

	// Fetch latest repo info
	repoInfo, err := w.fetcher.FetchRepoInfo(ctx, accessToken, owner, repo)
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to fetch repo info",
			slog.String("resource_id", resource.ResourceID),
			slog.String("public_id", resource.ResourcePublicID),
			slog.Any("error", err))

		return
	}

	// Update resource properties
	resourceProps := mergeResourceProperties(resource.ResourceProperties, repoInfo)

	updateErr := w.syncService.UpdateResourceProperties(ctx, resource.ResourceID, resourceProps)
	if updateErr != nil {
		w.logger.WarnContext(ctx, "Failed to update resource properties",
			slog.String("resource_id", resource.ResourceID),
			slog.Any("error", updateErr))
	}

	// Fetch contributors
	contributors, err := w.fetcher.FetchRepoContributors(ctx, accessToken, owner, repo)
	if err != nil {
		w.logger.WarnContext(ctx, "Failed to fetch contributors, skipping",
			slog.String("resource_id", resource.ResourceID),
			slog.Any("error", err))

		return
	}

	// Collect each contributor's appearance under this resource's profile
	for _, c := range contributors {
		remoteID := strconv.FormatInt(c.ID, 10)

		info, exists := contributorMap[remoteID]
		if !exists {
			info = &contributorInfo{
				login:          c.Login,
				reposByProfile: make(map[string][]contributorRepoStats),
			}
			contributorMap[remoteID] = info
		}

		info.reposByProfile[resource.ProfileID] = append(
			info.reposByProfile[resource.ProfileID],
			contributorRepoStats{
				owner:       owner,
				repo:        repo,
				stars:       repoInfo.Stars,
				commits:     c.Contributions,
				accessToken: accessToken,
			},
		)
	}
}

// batchMatchContributors batch-loads profile links and memberships, then fetches
// search stats ONLY for matched contributors. Returns accumulated membership stats.
func (w *GitHubSyncWorker) batchMatchContributors( //nolint:cyclop,funlen
	ctx context.Context,
	contributorMap map[string]*contributorInfo,
) (map[string]*membershipStatsAccumulator, int) {
	membershipStats := make(map[string]*membershipStatsAccumulator)

	// Batch load profile links: remoteID → profileID
	allRemoteIDs := make([]string, 0, len(contributorMap))
	for remoteID := range contributorMap {
		allRemoteIDs = append(allRemoteIDs, remoteID)
	}

	profileLinkMap, err := w.syncService.GetProfileLinksByRemoteIDs(ctx, "github", allRemoteIDs)
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to batch load profile links",
			slog.Any("error", err))

		return membershipStats, 0
	}

	w.logger.WarnContext(ctx, "Batch loaded profile links",
		slog.Int("total_contributors", len(allRemoteIDs)),
		slog.Int("matched_profiles", len(profileLinkMap)))

	if len(profileLinkMap) == 0 {
		return membershipStats, 0
	}

	// Collect unique profile ID pairs for membership batch
	resourceProfileIDSet := make(map[string]bool)
	contributorProfileIDSet := make(map[string]bool)

	for remoteID, info := range contributorMap {
		contribProfileID, found := profileLinkMap[remoteID]
		if !found {
			continue
		}

		for profileID := range info.reposByProfile {
			resourceProfileIDSet[profileID] = true
			contributorProfileIDSet[contribProfileID] = true
		}
	}

	resourceProfileIDs := mapSetToSlice(resourceProfileIDSet)
	contributorProfileIDs := mapSetToSlice(contributorProfileIDSet)

	membershipMap, err := w.syncService.GetMembershipsByProfilePairs(
		ctx, resourceProfileIDs, contributorProfileIDs,
	)
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to batch load memberships",
			slog.Any("error", err))

		return membershipStats, 0
	}

	w.logger.WarnContext(ctx, "Batch loaded memberships",
		slog.Int("matched_memberships", len(membershipMap)))

	// Build search queries and accumulate commits/stars for matched contributors.
	// We do this in two passes:
	//  1. Collect all (commits, stars) and build GraphQL search queries
	//  2. Execute batch search and distribute results

	matchedCount := 0

	// searchQueries maps GraphQL alias → search query string.
	searchQueries := make(map[string]string)
	// aliasToMembership maps GraphQL alias → membershipID for distributing results.
	aliasToMembership := make(map[string]string)
	// Track which access token to use (pick any valid one from matched contributors).
	var batchAccessToken string

	for remoteID, info := range contributorMap {
		contribProfileID, found := profileLinkMap[remoteID]
		if !found {
			continue
		}

		for profileID, repos := range info.reposByProfile {
			membershipID, found := membershipMap[profileID+":"+contribProfileID]
			if !found {
				continue
			}

			matchedCount++

			acc, exists := membershipStats[membershipID]
			if !exists {
				acc = &membershipStatsAccumulator{}
				membershipStats[membershipID] = acc
			}

			for _, r := range repos {
				// Accumulate commits and stars directly (already fetched from REST).
				acc.commits += r.commits
				acc.stars += r.stars

				if batchAccessToken == "" {
					batchAccessToken = r.accessToken
				}

				// Build 4 search queries per (contributor, repo) pair.
				repoQ := "repo:" + r.owner + "/" + r.repo
				prefix := sanitizeAlias(remoteID + "_" + r.owner + "_" + r.repo)

				searchQueries[prefix+"_pr"] = repoQ + " type:pr author:" + info.login
				aliasToMembership[prefix+"_pr"] = membershipID

				searchQueries[prefix+"_prm"] = repoQ + " type:pr author:" + info.login + " is:merged"
				aliasToMembership[prefix+"_prm"] = membershipID

				searchQueries[prefix+"_iss"] = repoQ + " type:issue author:" + info.login
				aliasToMembership[prefix+"_iss"] = membershipID

				searchQueries[prefix+"_issc"] = repoQ + " type:issue author:" + info.login + " is:closed"
				aliasToMembership[prefix+"_issc"] = membershipID
			}
		}
	}

	// Execute batch search via GraphQL (or fall back to REST on failure).
	if len(searchQueries) > 0 && batchAccessToken != "" {
		w.executeBatchSearch(
			ctx,
			batchAccessToken,
			searchQueries,
			aliasToMembership,
			membershipStats,
		)
	}

	return membershipStats, matchedCount
}

// sanitizeAlias converts a string to a valid GraphQL alias name.
// This is a local re-export of the pattern used by the GitHub client.
func sanitizeAlias(s string) string {
	var b strings.Builder

	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			b.WriteRune(c)
		} else {
			b.WriteByte('_')
		}
	}

	cleaned := b.String()
	if len(cleaned) == 0 || (cleaned[0] >= '0' && cleaned[0] <= '9') {
		return "a_" + cleaned
	}

	return cleaned
}

// executeBatchSearch runs all search queries via GraphQL batch. On failure,
// falls back to individual REST calls.
func (w *GitHubSyncWorker) executeBatchSearch(
	ctx context.Context,
	accessToken string,
	searchQueries map[string]string,
	aliasToMembership map[string]string,
	membershipStats map[string]*membershipStatsAccumulator,
) {
	w.logger.WarnContext(ctx, "Executing GraphQL batch search",
		slog.Int("query_count", len(searchQueries)))

	results, err := w.fetcher.SearchIssueCountsBatch(ctx, accessToken, searchQueries)
	if err != nil {
		w.logger.WarnContext(ctx, "GraphQL batch search failed, falling back to REST",
			slog.Any("error", err))

		w.fallbackToRESTSearch(ctx, accessToken, searchQueries, aliasToMembership, membershipStats)

		return
	}

	w.logger.WarnContext(ctx, "GraphQL batch search completed",
		slog.Int("results", len(results)))

	// Distribute results back to membership accumulators.
	for alias, count := range results {
		membershipID, found := aliasToMembership[alias]
		if !found {
			continue
		}

		acc := membershipStats[membershipID]
		if acc == nil {
			continue
		}

		// Parse stat type from alias suffix.
		switch {
		case strings.HasSuffix(alias, "_prm"):
			acc.prsResolved += count
		case strings.HasSuffix(alias, "_pr"):
			acc.prsTotal += count
		case strings.HasSuffix(alias, "_issc"):
			acc.issuesResolved += count
		case strings.HasSuffix(alias, "_iss"):
			acc.issuesTotal += count
		}
	}
}

// fallbackToRESTSearch executes search queries one at a time via REST API
// when GraphQL batch search fails.
func (w *GitHubSyncWorker) fallbackToRESTSearch(
	ctx context.Context,
	accessToken string,
	searchQueries map[string]string,
	aliasToMembership map[string]string,
	membershipStats map[string]*membershipStatsAccumulator,
) {
	for alias, query := range searchQueries {
		membershipID, found := aliasToMembership[alias]
		if !found {
			continue
		}

		acc := membershipStats[membershipID]
		if acc == nil {
			continue
		}

		count, err := w.fetcher.SearchIssues(ctx, accessToken, query)
		if err != nil {
			w.logger.WarnContext(ctx, "REST search fallback failed",
				slog.String("alias", alias),
				slog.Any("error", err))

			continue
		}

		switch {
		case strings.HasSuffix(alias, "_prm"):
			acc.prsResolved += count
		case strings.HasSuffix(alias, "_pr"):
			acc.prsTotal += count
		case strings.HasSuffix(alias, "_issc"):
			acc.issuesResolved += count
		case strings.HasSuffix(alias, "_iss"):
			acc.issuesTotal += count
		}
	}
}

// flushMembershipStats writes aggregated GitHub stats to a membership's properties.
// The flat format matches what the frontend expects: {"github": {"commits": N, "prs": {...}, ...}}.
// Uses JSONB merge so non-github keys (e.g., "videos") are preserved.
func (w *GitHubSyncWorker) flushMembershipStats(
	ctx context.Context,
	membershipID string,
	acc *membershipStatsAccumulator,
) error {
	properties := map[string]any{
		"github": &resourcesync.GitHubContributorStats{
			Commits: acc.commits,
			PRs: struct {
				Total    int `json:"total"`
				Resolved int `json:"resolved"`
			}{Total: acc.prsTotal, Resolved: acc.prsResolved},
			Issues: struct {
				Total    int `json:"total"`
				Resolved int `json:"resolved"`
			}{Total: acc.issuesTotal, Resolved: acc.issuesResolved},
			Stars: acc.stars,
		},
	}

	return w.syncService.UpdateMembershipProperties(ctx, membershipID, properties)
}

// mapSetToSlice converts a map[string]bool set to a string slice.
func mapSetToSlice(set map[string]bool) []string {
	result := make([]string, 0, len(set))
	for k := range set {
		result = append(result, k)
	}

	return result
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
