package externalsite

import (
	"context"
	"fmt"
	"log/slog"
	neturl "net/url"
	"strings"
	"time"

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

	frontMatter, body, err := ParseMarkdownFile(content, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	props := buildImportProperties(frontMatter, body, ownerRepo, branch, filePath)

	// Build description — strip HTML and markdown for clean plain text
	description := SanitizeDescription(frontMatter.Description)

	const maxDescriptionLen = 200

	if description == "" && len(body) > 0 {
		description = truncateString(SanitizeDescription(body), maxDescriptionLen)
	}

	// Determine story kind from directory
	storyKind := storyKindFromPath(filePath)

	// Build the GitHub blob URL for the file
	parts := strings.SplitN(ownerRepo, "/", ownerRepoParts)
	blobURL := fmt.Sprintf(
		"https://github.com/%s/blob/%s/%s",
		escapeOwnerRepo(ownerRepo),
		neturl.PathEscape(branch),
		filePath,
	)

	if len(parts) != ownerRepoParts {
		blobURL = ""
	}

	var publishedAt time.Time
	if frontMatter.Date != nil {
		publishedAt = *frontMatter.Date
	}

	item := &siteimporter.ImportItem{
		RemoteID:     filePath,
		Title:        frontMatter.Title,
		Description:  description,
		Link:         blobURL,
		ThumbnailURL: "",
		StoryKind:    storyKind,
		Properties:   props,
		PublishedAt:  publishedAt,
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
	if path, found := strings.CutPrefix(rawURL, "git@github.com:"); found {
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
	const maxSplitParts = 3

	parts := strings.SplitN(path, "/", maxSplitParts)
	if len(parts) < ownerRepoParts || parts[0] == "" || parts[1] == "" {
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

// buildImportProperties constructs the properties map for an import item from frontmatter and body.
func buildImportProperties(
	frontMatter *ParsedFrontmatter,
	body string,
	ownerRepo string,
	branch string,
	filePath string,
) map[string]any {
	props := map[string]any{
		"source_path": filePath,
	}

	if body != "" {
		sanitized := SanitizeContent(body)
		props["content"] = ResolveRelativeImages(sanitized, ownerRepo, branch, filePath)
	}

	if frontMatter.Slug != "" {
		props["slug"] = frontMatter.Slug
	}

	if frontMatter.Language != "" {
		props["language"] = frontMatter.Language
	}

	if len(frontMatter.Tags) > 0 {
		props["tags"] = frontMatter.Tags
	}

	return props
}

// truncateString truncates a string to maxLen characters at a word boundary.
func truncateString(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	// Find the last space before maxLen
	truncated := text[:maxLen]
	lastSpace := strings.LastIndex(truncated, " ")

	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}
