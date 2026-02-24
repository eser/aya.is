package externalsite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpclient"
)

// escapeOwnerRepo URL-path-escapes the owner and repo segments individually,
// preserving the slash separator. This prevents injection into URL paths.
func escapeOwnerRepo(ownerRepo string) string {
	parts := strings.SplitN(ownerRepo, "/", 2) //nolint:mnd
	if len(parts) < 2 {                        //nolint:mnd
		return neturl.PathEscape(ownerRepo)
	}

	return neturl.PathEscape(parts[0]) + "/" + neturl.PathEscape(parts[1])
}

const (
	githubAPIBase = "https://api.github.com"
	githubRawBase = "https://raw.githubusercontent.com"
)

// repoInfoResponse holds selected fields from GitHub's GET /repos/{owner}/{repo} response.
type repoInfoResponse struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
	Private       bool   `json:"private"`
}

// treeResponse represents GitHub's GET /repos/{owner}/{repo}/git/trees/{sha} response.
type treeResponse struct {
	Tree []treeEntry `json:"tree"`
}

// treeEntry represents a single entry in a GitHub tree.
type treeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"` // "blob" or "tree"
}

// fetchRepoInfo fetches basic repository info from GitHub API.
func fetchRepoInfo(
	ctx context.Context,
	client *httpclient.Client,
	ownerRepo string,
) (*repoInfoResponse, error) {
	url := fmt.Sprintf("%s/repos/%s", githubAPIBase, escapeOwnerRepo(ownerRepo))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repository %q not found", ownerRepo)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var info repoInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if info.Private {
		return nil, fmt.Errorf("repository %q is private", ownerRepo)
	}

	return &info, nil
}

// fetchRepoTree fetches the full recursive tree of a repository branch.
func fetchRepoTree(
	ctx context.Context,
	client *httpclient.Client,
	ownerRepo string,
	branch string,
) ([]treeEntry, error) {
	url := fmt.Sprintf(
		"%s/repos/%s/git/trees/%s?recursive=1",
		githubAPIBase, escapeOwnerRepo(ownerRepo), neturl.PathEscape(branch),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tree: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var tree treeResponse
	if err := json.NewDecoder(resp.Body).Decode(&tree); err != nil {
		return nil, fmt.Errorf("failed to decode tree: %w", err)
	}

	return tree.Tree, nil
}

// markdownExtensions lists recognised markdown file extensions.
var markdownExtensions = []string{".md", ".mdx", ".mdoc", ".markdown"}

// isMarkdownFile checks if a file path has a recognised markdown extension.
func isMarkdownFile(path string) bool {
	lower := strings.ToLower(path)
	for _, ext := range markdownExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}

	return false
}

// stripMarkdownExt removes the markdown extension from a filename, returning the name part.
func stripMarkdownExt(filename string) string {
	lower := strings.ToLower(filename)
	for _, ext := range markdownExtensions {
		if strings.HasSuffix(lower, ext) {
			return filename[:len(filename)-len(ext)]
		}
	}

	return filename
}

// filterMarkdownFiles selects markdown blob files from a tree, skipping section index/README files.
// Recognised extensions: .md, .mdx, .mdoc, .markdown.
// It auto-detects the content root (content/, _posts/, or repo root).
// Hugo page bundles (e.g. content/posts/my-post/index.md) are kept.
func filterMarkdownFiles(tree []treeEntry) []string {
	// Detect content root
	contentRoot := detectContentRoot(tree)

	var files []string

	for _, entry := range tree {
		if entry.Type != "blob" {
			continue
		}

		if !isMarkdownFile(entry.Path) {
			continue
		}

		// Only include files under the content root
		if contentRoot != "" && !strings.HasPrefix(entry.Path, contentRoot) {
			continue
		}

		base := entry.Path[strings.LastIndex(entry.Path, "/")+1:]
		name := strings.ToLower(stripMarkdownExt(base))

		// Always skip _index (Hugo section listings) and README
		if name == "_index" || name == "readme" {
			continue
		}

		// For index files, only keep if it's a Hugo page bundle (nested in a post subdirectory).
		// e.g. content/posts/my-post/index.md → keep (page bundle)
		// e.g. content/index.md → skip (section root)
		if name == "index" {
			pathUnderRoot := entry.Path
			if contentRoot != "" {
				pathUnderRoot = strings.TrimPrefix(entry.Path, contentRoot)
			}

			// Count slashes to determine depth: "posts/my-post/index.md" has 2 slashes → keep
			// "index.md" has 0 slashes → skip
			if strings.Count(pathUnderRoot, "/") < 2 {
				continue
			}
		}

		files = append(files, entry.Path)
	}

	return files
}

// detectContentRoot finds the best content directory in the tree.
// Prefers content/, then _posts/, then empty string (repo root).
func detectContentRoot(tree []treeEntry) string {
	for _, entry := range tree {
		if entry.Type == "tree" && entry.Path == "content" {
			return "content/"
		}
	}

	for _, entry := range tree {
		if entry.Type == "tree" && entry.Path == "_posts" {
			return "_posts/"
		}
	}

	return ""
}

// fetchRawFile fetches a single raw file from GitHub.
func fetchRawFile(
	ctx context.Context,
	client *httpclient.Client,
	ownerRepo string,
	branch string,
	filePath string,
) (string, error) {
	url := fmt.Sprintf(
		"%s/%s/%s/%s",
		githubRawBase,
		escapeOwnerRepo(ownerRepo),
		neturl.PathEscape(branch),
		filePath,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch %s: status %d", filePath, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	return string(body), nil
}
