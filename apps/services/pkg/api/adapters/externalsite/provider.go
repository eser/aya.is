package externalsite

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpclient"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/siteimporter"
)

// Provider implements siteimporter.SiteProvider for external static sites
// hosted on GitHub (Jekyll, Hugo, Zola, etc.).
type Provider struct {
	logger     *logfx.Logger
	httpClient *httpclient.Client
}

// NewProvider creates a new external site provider.
func NewProvider(logger *logfx.Logger, httpClient *httpclient.Client) *Provider {
	return &Provider{
		logger:     logger,
		httpClient: httpClient,
	}
}

// Kind returns the provider kind.
func (p *Provider) Kind() string {
	return "external-site"
}

// Check validates a GitHub repository URL and returns connection info.
func (p *Provider) Check(ctx context.Context, rawURL string) (*siteimporter.CheckResult, error) {
	ownerRepo := extractOwnerRepo(rawURL)
	if ownerRepo == "" {
		return nil, fmt.Errorf(
			"%w: could not extract owner/repo from URL",
			siteimporter.ErrInvalidURL,
		)
	}

	// Verify the repo exists and is public via GitHub API
	repoInfo, err := fetchRepoInfo(ctx, p.httpClient, ownerRepo)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", siteimporter.ErrSiteNotFound, err)
	}

	p.logger.DebugContext(ctx, "Validated GitHub repo for external site",
		slog.String("owner_repo", ownerRepo),
		slog.String("name", repoInfo.Name))

	return &siteimporter.CheckResult{
		Valid:    true,
		Username: ownerRepo,
		URI:      "", // URI comes from the frontend dialog's site_url field
		Title:    repoInfo.Name,
	}, nil
}

// FetchAll fetches all markdown files from a GitHub repository.
func (p *Provider) FetchAll(
	ctx context.Context,
	username string, // "owner/repo"
) ([]*siteimporter.ImportItem, error) {
	p.logger.DebugContext(ctx, "Fetching all content from GitHub repo",
		slog.String("repo", username))

	// Get the default branch
	repoInfo, err := fetchRepoInfo(ctx, p.httpClient, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo info: %w", err)
	}

	// Fetch the repo tree recursively
	tree, err := fetchRepoTree(ctx, p.httpClient, username, repoInfo.DefaultBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo tree: %w", err)
	}

	// Filter for markdown files
	mdFiles := filterMarkdownFiles(tree)

	p.logger.DebugContext(ctx, "Found markdown files",
		slog.String("repo", username),
		slog.Int("total_tree_entries", len(tree)),
		slog.Int("markdown_files", len(mdFiles)))

	// Fetch and parse each markdown file
	var items []*siteimporter.ImportItem

	for _, filePath := range mdFiles {
		item, fetchErr := p.fetchAndParseFile(ctx, username, repoInfo.DefaultBranch, filePath)
		if fetchErr != nil {
			p.logger.WarnContext(ctx, "Failed to fetch/parse file, skipping",
				slog.String("file", filePath),
				slog.String("error", fetchErr.Error()))

			continue
		}

		if item != nil {
			items = append(items, item)
		}
	}

	p.logger.DebugContext(ctx, "Completed fetching external site content",
		slog.String("repo", username),
		slog.Int("items", len(items)))

	return items, nil
}

// fetchAndParseFile fetches a single markdown file and converts it to an ImportItem.
func (p *Provider) fetchAndParseFile(
	ctx context.Context,
	ownerRepo string,
	branch string,
	filePath string,
) (*siteimporter.ImportItem, error) {
	content, err := fetchRawFile(ctx, p.httpClient, ownerRepo, branch, filePath)
	if err != nil {
		return nil, err
	}

	fm, body, err := ParseMarkdownFile(content, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Build properties
	props := map[string]any{
		"source_path": filePath,
	}

	if body != "" {
		props["content"] = body
	}

	if fm.Slug != "" {
		props["slug"] = fm.Slug
	}

	if fm.Language != "" {
		props["language"] = fm.Language
	}

	if len(fm.Tags) > 0 {
		props["tags"] = fm.Tags
	}

	// Build description
	description := fm.Description
	if description == "" && len(body) > 0 {
		description = truncateString(body, 200)
	}

	// Determine story kind from directory
	storyKind := storyKindFromPath(filePath)

	// Build the GitHub blob URL for the file
	parts := strings.SplitN(ownerRepo, "/", 2)
	blobURL := fmt.Sprintf("https://github.com/%s/blob/%s/%s", ownerRepo, branch, filePath)

	if len(parts) != 2 {
		blobURL = ""
	}

	item := &siteimporter.ImportItem{
		RemoteID:    filePath,
		Title:       fm.Title,
		Description: description,
		Link:        blobURL,
		StoryKind:   storyKind,
		Properties:  props,
	}

	if fm.Date != nil {
		item.PublishedAt = *fm.Date
	}

	return item, nil
}

// extractOwnerRepo parses a GitHub URL to extract "owner/repo".
// Supports:
//   - https://github.com/owner/repo
//   - github.com/owner/repo
//   - owner/repo
//   - git@github.com:owner/repo.git
func extractOwnerRepo(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	// Handle SSH format: git@github.com:owner/repo.git
	if strings.HasPrefix(rawURL, "git@github.com:") {
		path := strings.TrimPrefix(rawURL, "git@github.com:")
		path = strings.TrimSuffix(path, ".git")

		return normalizeOwnerRepo(path)
	}

	// Remove protocol
	rawURL = strings.TrimPrefix(rawURL, "https://")
	rawURL = strings.TrimPrefix(rawURL, "http://")

	// Remove github.com prefix
	rawURL = strings.TrimPrefix(rawURL, "github.com/")
	rawURL = strings.TrimPrefix(rawURL, "www.github.com/")

	// Remove .git suffix
	rawURL = strings.TrimSuffix(rawURL, ".git")

	// Remove trailing slash
	rawURL = strings.TrimRight(rawURL, "/")

	return normalizeOwnerRepo(rawURL)
}

// normalizeOwnerRepo validates and normalizes an "owner/repo" string.
func normalizeOwnerRepo(path string) string {
	parts := strings.SplitN(path, "/", 3) //nolint:mnd
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return ""
	}

	return parts[0] + "/" + parts[1]
}

// storyKindFromPath maps a file's directory to a story kind.
func storyKindFromPath(filePath string) string {
	lower := strings.ToLower(filePath)

	switch {
	case strings.Contains(lower, "/talks/") || strings.Contains(lower, "/_talks/"):
		return "presentation"
	case strings.Contains(lower, "/projects/") || strings.Contains(lower, "/_projects/"):
		return "content"
	default:
		return "article"
	}
}

// truncateString truncates a string to maxLen characters at a word boundary.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	// Find the last space before maxLen
	truncated := s[:maxLen]
	lastSpace := strings.LastIndex(truncated, " ")

	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}
