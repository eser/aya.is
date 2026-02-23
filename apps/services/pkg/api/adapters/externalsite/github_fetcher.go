package externalsite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpclient"
)

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
	url := fmt.Sprintf("%s/repos/%s", githubAPIBase, ownerRepo)

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
		githubAPIBase, ownerRepo, branch,
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

// filterMarkdownFiles selects .md blob files from a tree, skipping index/README files.
// It auto-detects the content root (content/, _posts/, or repo root).
func filterMarkdownFiles(tree []treeEntry) []string {
	// Detect content root
	contentRoot := detectContentRoot(tree)

	var files []string

	for _, entry := range tree {
		if entry.Type != "blob" {
			continue
		}

		if !strings.HasSuffix(entry.Path, ".md") {
			continue
		}

		// Only include files under the content root
		if contentRoot != "" && !strings.HasPrefix(entry.Path, contentRoot) {
			continue
		}

		// Skip index files and READMEs
		base := entry.Path[strings.LastIndex(entry.Path, "/")+1:]
		lower := strings.ToLower(base)

		if lower == "_index.md" || lower == "readme.md" || lower == "index.md" {
			continue
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
	url := fmt.Sprintf("%s/%s/%s/%s", githubRawBase, ownerRepo, branch, filePath)

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
