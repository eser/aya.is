package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
)

// Sentinel errors.
var (
	ErrFailedToExchangeCode      = errors.New("failed to exchange authorization code")
	ErrFailedToGetUserInfo       = errors.New("failed to get user info")
	ErrFailedToFetchRepos        = errors.New("failed to fetch user repos")
	ErrFailedToFetchRepoInfo     = errors.New("failed to fetch repo info")
	ErrFailedToSearchIssues      = errors.New("failed to search issues")
	ErrFailedToGraphQLBatchQuery = errors.New("failed to execute GraphQL batch query")
	ErrNoUserFound               = errors.New("no user found")
)

// HTTPClient interface for dependency injection.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// UserInfo represents GitHub user information.
type UserInfo struct {
	ID      int64  `json:"id"`
	Login   string `json:"login"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Avatar  string `json:"avatar_url"`
	HTMLURL string `json:"html_url"`
}

// OrgInfo represents GitHub organization information.
type OrgInfo struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Avatar      string `json:"avatar_url"`
	HTMLURL     string `json:"html_url"`
}

// TokenResponse represents GitHub's token endpoint response.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// Client provides GitHub API operations.
type Client struct {
	config     *auth.GitHubAuthProviderConfig
	logger     *logfx.Logger
	httpClient HTTPClient
}

// NewClient creates a new GitHub client.
func NewClient(
	config *auth.GitHubAuthProviderConfig,
	logger *logfx.Logger,
	httpClient HTTPClient,
) *Client {
	return &Client{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}
}

// rateLimitThreshold is the minimum remaining calls before the client
// pauses and waits for the rate limit to reset. Keeps a small buffer
// so a few calls can still succeed while we read the headers.
const rateLimitThreshold = 3

// checkRateLimit inspects GitHub API rate limit headers and sleeps if
// remaining calls are at or below the threshold. Called after every API
// response (successful or not) to preemptively throttle before hitting 429/403.
func (c *Client) checkRateLimit(ctx context.Context, resp *http.Response) {
	remainingStr := resp.Header.Get("X-Ratelimit-Remaining")
	resetStr := resp.Header.Get("X-Ratelimit-Reset")

	if remainingStr == "" || resetStr == "" {
		return
	}

	remaining, err := strconv.Atoi(remainingStr)
	if err != nil {
		return
	}

	if remaining > rateLimitThreshold {
		return
	}

	resetUnix, err := strconv.ParseInt(resetStr, 10, 64)
	if err != nil {
		return
	}

	resetTime := time.Unix(resetUnix, 0)
	wait := time.Until(resetTime) + time.Second // 1s buffer for clock skew

	if wait <= 0 {
		return
	}

	c.logger.WarnContext(ctx, "GitHub API rate limit low, pausing until reset",
		slog.Int("remaining", remaining),
		slog.Duration("wait", wait),
		slog.Time("reset_at", resetTime))

	select {
	case <-time.After(wait):
		c.logger.DebugContext(ctx, "GitHub API rate limit reset, resuming")
	case <-ctx.Done():
	}
}

// Config returns the GitHub configuration.
func (c *Client) Config() *auth.GitHubAuthProviderConfig {
	return c.config
}

// Logger returns the logger.
func (c *Client) Logger() *logfx.Logger {
	return c.logger
}

// BuildAuthURL builds the GitHub OAuth authorization URL.
func (c *Client) BuildAuthURL(redirectURI, state, scope string) string {
	queryString := url.Values{}
	queryString.Set("client_id", c.config.ClientID)
	queryString.Set("redirect_uri", redirectURI)
	queryString.Set("state", state)
	queryString.Set("scope", scope)
	queryString.Set("prompt", "select_account")

	oauthURL := url.URL{ //nolint:exhaustruct
		Scheme:   "https",
		Host:     "github.com",
		Path:     "/login/oauth/authorize",
		RawQuery: queryString.Encode(),
	}

	return oauthURL.String()
}

// ExchangeCodeForToken exchanges an authorization code for an access token.
func (c *Client) ExchangeCodeForToken(
	ctx context.Context,
	code string,
	redirectURI string,
) (*TokenResponse, error) {
	values := url.Values{
		"client_id":     {c.config.ClientID},
		"client_secret": {c.config.ClientSecret},
		"code":          {code},
	}

	if redirectURI != "" {
		values.Set("redirect_uri", redirectURI)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to exchange code for token",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "Token exchange failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToExchangeCode, resp.StatusCode)
	}

	// Try JSON first
	var tokenResp TokenResponse

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		// Fallback to URL-encoded response (older GitHub behavior)
		vals, _ := url.ParseQuery(string(body))
		tokenResp.AccessToken = vals.Get("access_token")
		tokenResp.TokenType = vals.Get("token_type")
		tokenResp.Scope = vals.Get("scope")
	}

	if tokenResp.AccessToken == "" {
		c.logger.ErrorContext(ctx, "No access token received from GitHub")

		return nil, ErrFailedToExchangeCode
	}

	c.logger.DebugContext(ctx, "Successfully obtained GitHub access token")

	return &tokenResp, nil
}

// FetchUserInfo fetches user information from GitHub API.
func (c *Client) FetchUserInfo(
	ctx context.Context,
	accessToken string,
) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-Github-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to fetch user info",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "User info request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToGetUserInfo, resp.StatusCode)
	}

	var userInfo UserInfo

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}

	if userInfo.ID == 0 {
		return nil, ErrNoUserFound
	}

	return &userInfo, nil
}

// FetchUserOrganizations fetches organizations the user belongs to from GitHub API.
func (c *Client) FetchUserOrganizations(
	ctx context.Context,
	accessToken string,
) ([]*OrgInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user/orgs",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-Github-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to fetch user organizations",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "User organizations request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToGetUserInfo, resp.StatusCode)
	}

	var orgs []*OrgInfo

	if err := json.Unmarshal(body, &orgs); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUserInfo, err)
	}

	return orgs, nil
}

// GitHubRepoInfo represents a GitHub repository.
type GitHubRepoInfo struct {
	ID          int64  `json:"id"`
	FullName    string `json:"full_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Private     bool   `json:"private"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
}

// GitHubContributorInfo represents a GitHub contributor.
type GitHubContributorInfo struct {
	ID            int64  `json:"id"`
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}

// GitHubIssueSearchResult represents GitHub issue search results.
type GitHubIssueSearchResult struct {
	TotalCount int `json:"total_count"`
}

// FetchUserRepos fetches repositories accessible to the authenticated user from GitHub API.
func (c *Client) FetchUserRepos(
	ctx context.Context,
	accessToken string,
	affiliation string,
	page int,
	perPage int,
) ([]*GitHubRepoInfo, error) {
	reqURL := fmt.Sprintf(
		"https://api.github.com/user/repos?affiliation=%s&sort=updated&per_page=%d&page=%d",
		url.QueryEscape(affiliation),
		perPage,
		page,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		reqURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchRepos, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-Github-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to fetch user repos",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchRepos, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "User repos request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToFetchRepos, resp.StatusCode)
	}

	var repos []*GitHubRepoInfo

	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchRepos, err)
	}

	return repos, nil
}

// FetchRepoContributors fetches contributors for a specific GitHub repository.
func (c *Client) FetchRepoContributors(
	ctx context.Context,
	accessToken string,
	owner string,
	repo string,
) ([]*GitHubContributorInfo, error) {
	reqURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contributors",
		url.PathEscape(owner),
		url.PathEscape(repo),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		reqURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchRepoInfo, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-Github-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to fetch repo contributors",
			slog.String("error", err.Error()),
			slog.String("owner", owner),
			slog.String("repo", repo))

		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchRepoInfo, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)
	c.checkRateLimit(ctx, resp)

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "Repo contributors request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToFetchRepoInfo, resp.StatusCode)
	}

	var contributors []*GitHubContributorInfo

	if err := json.Unmarshal(body, &contributors); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchRepoInfo, err)
	}

	return contributors, nil
}

// FetchRepoInfo fetches information about a specific GitHub repository.
func (c *Client) FetchRepoInfo(
	ctx context.Context,
	accessToken string,
	owner string,
	repo string,
) (*GitHubRepoInfo, error) {
	reqURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s",
		url.PathEscape(owner),
		url.PathEscape(repo),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		reqURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchRepoInfo, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-Github-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to fetch repo info",
			slog.String("error", err.Error()),
			slog.String("owner", owner),
			slog.String("repo", repo))

		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchRepoInfo, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)
	c.checkRateLimit(ctx, resp)

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "Repo info request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToFetchRepoInfo, resp.StatusCode)
	}

	var repoInfo GitHubRepoInfo

	if err := json.Unmarshal(body, &repoInfo); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchRepoInfo, err)
	}

	return &repoInfo, nil
}

// SearchIssues searches GitHub issues and returns the total count.
func (c *Client) SearchIssues(
	ctx context.Context,
	accessToken string,
	query string,
) (int, error) {
	reqURL := "https://api.github.com/search/issues?q=" + url.QueryEscape(query)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		reqURL,
		nil,
	)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrFailedToSearchIssues, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-Github-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to search issues",
			slog.String("error", err.Error()),
			slog.String("query", query))

		return 0, fmt.Errorf("%w: %w", ErrFailedToSearchIssues, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)
	c.checkRateLimit(ctx, resp)

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "Issue search request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return 0, fmt.Errorf("%w: status %d", ErrFailedToSearchIssues, resp.StatusCode)
	}

	var searchResult GitHubIssueSearchResult

	if err := json.Unmarshal(body, &searchResult); err != nil {
		return 0, fmt.Errorf("%w: %w", ErrFailedToSearchIssues, err)
	}

	return searchResult.TotalCount, nil
}

// graphQLMaxAliasesPerRequest is the maximum number of search aliases to pack
// into a single GraphQL request. Keeps query size reasonable.
const graphQLMaxAliasesPerRequest = 100

// graphQLEndpoint is the GitHub GraphQL API URL.
const graphQLEndpoint = "https://api.github.com/graphql"

// aliasCleanerRegex strips characters that are invalid in GraphQL alias names.
var aliasCleanerRegex = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// sanitizeAlias converts a string to a valid GraphQL alias name.
// GraphQL aliases must match /[_A-Za-z][_0-9A-Za-z]*/.
func sanitizeAlias(s string) string {
	cleaned := aliasCleanerRegex.ReplaceAllString(s, "_")
	if len(cleaned) == 0 || (cleaned[0] >= '0' && cleaned[0] <= '9') {
		cleaned = "a_" + cleaned
	}

	return cleaned
}

// graphQLResponse represents the top-level GraphQL response.
type graphQLResponse struct {
	Data   map[string]graphQLSearchResult `json:"data"`
	Errors []graphQLError                 `json:"errors"`
}

type graphQLSearchResult struct {
	IssueCount int `json:"issueCount"`
}

type graphQLError struct {
	Message string `json:"message"`
}

// SearchIssueCountsBatch executes multiple GitHub search queries in a single
// GraphQL API call using aliases. Returns map[alias]count.
//
// Each entry in queries maps an alias name to a search query string.
// Example: {"eser_prs": "repo:org/repo type:pr author:eser"}
//
// If len(queries) > graphQLMaxAliasesPerRequest, splits into multiple requests.
func (c *Client) SearchIssueCountsBatch(
	ctx context.Context,
	accessToken string,
	queries map[string]string,
) (map[string]int, error) {
	if len(queries) == 0 {
		return map[string]int{}, nil
	}

	// Split into batches
	batches := make([]map[string]string, 0, (len(queries)/graphQLMaxAliasesPerRequest)+1)
	current := make(map[string]string)

	for alias, query := range queries {
		current[alias] = query
		if len(current) >= graphQLMaxAliasesPerRequest {
			batches = append(batches, current)
			current = make(map[string]string)
		}
	}

	if len(current) > 0 {
		batches = append(batches, current)
	}

	// Execute each batch
	result := make(map[string]int, len(queries))

	for i, batch := range batches {
		batchResult, err := c.executeGraphQLSearchBatch(ctx, accessToken, batch)
		if err != nil {
			c.logger.WarnContext(ctx, "GraphQL batch search failed",
				slog.Int("batch", i+1),
				slog.Int("total_batches", len(batches)),
				slog.Any("error", err))

			return nil, err
		}

		for alias, count := range batchResult {
			result[alias] = count
		}
	}

	return result, nil
}

// executeGraphQLSearchBatch sends a single GraphQL request with multiple search aliases.
func (c *Client) executeGraphQLSearchBatch(
	ctx context.Context,
	accessToken string,
	queries map[string]string,
) (map[string]int, error) {
	// Build GraphQL query with aliases
	var queryBuilder strings.Builder

	queryBuilder.WriteString("query {")

	for alias, searchQuery := range queries {
		// Escape double quotes in the search query
		escapedQuery := strings.ReplaceAll(searchQuery, `"`, `\"`)
		queryBuilder.WriteString(fmt.Sprintf(
			` %s: search(query: "%s", type: ISSUE, first: 1) { issueCount }`,
			alias, escapedQuery,
		))
	}

	queryBuilder.WriteString(" }")

	// Build request body
	reqBody, err := json.Marshal(map[string]string{
		"query": queryBuilder.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGraphQLBatchQuery, err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		graphQLEndpoint,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGraphQLBatchQuery, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGraphQLBatchQuery, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)
	c.checkRateLimit(ctx, resp)

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "GraphQL batch search request failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToGraphQLBatchQuery, resp.StatusCode)
	}

	var gqlResp graphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGraphQLBatchQuery, err)
	}

	if len(gqlResp.Errors) > 0 {
		c.logger.WarnContext(ctx, "GraphQL batch search returned errors",
			slog.String("first_error", gqlResp.Errors[0].Message),
			slog.Int("error_count", len(gqlResp.Errors)))

		return nil, fmt.Errorf("%w: %s", ErrFailedToGraphQLBatchQuery, gqlResp.Errors[0].Message)
	}

	// Extract counts from response
	result := make(map[string]int, len(gqlResp.Data))
	for alias, searchResult := range gqlResp.Data {
		result[alias] = searchResult.IssueCount
	}

	c.logger.DebugContext(ctx, "GraphQL batch search completed",
		slog.Int("aliases", len(result)))

	return result, nil
}
