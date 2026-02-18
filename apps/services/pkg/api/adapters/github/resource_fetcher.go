package github

import (
	"context"

	"github.com/eser/aya.is/services/pkg/api/adapters/workers"
)

// ResourceFetcherAdapter adapts github.Client to implement workers.GitHubResourceFetcher.
type ResourceFetcherAdapter struct {
	client *Client
}

// NewResourceFetcherAdapter creates a new adapter wrapping the given Client.
func NewResourceFetcherAdapter(client *Client) workers.GitHubResourceFetcher {
	if client == nil {
		return nil
	}

	return &ResourceFetcherAdapter{client: client}
}

// FetchRepoInfo fetches information about a specific repository.
func (a *ResourceFetcherAdapter) FetchRepoInfo(
	ctx context.Context,
	accessToken string,
	owner string,
	repo string,
) (*workers.GitHubRepoInfoResult, error) {
	info, err := a.client.FetchRepoInfo(ctx, accessToken, owner, repo)
	if err != nil {
		return nil, err
	}

	return &workers.GitHubRepoInfoResult{
		ID:          info.ID,
		FullName:    info.FullName,
		Name:        info.Name,
		Description: info.Description,
		HTMLURL:     info.HTMLURL,
		Language:    info.Language,
		Stars:       info.Stars,
		Forks:       info.Forks,
		Private:     info.Private,
	}, nil
}

// FetchRepoContributors fetches contributors for a repository.
func (a *ResourceFetcherAdapter) FetchRepoContributors(
	ctx context.Context,
	accessToken string,
	owner string,
	repo string,
) ([]*workers.GitHubContributorResult, error) {
	contributors, err := a.client.FetchRepoContributors(ctx, accessToken, owner, repo)
	if err != nil {
		return nil, err
	}

	results := make([]*workers.GitHubContributorResult, len(contributors))
	for i, c := range contributors {
		results[i] = &workers.GitHubContributorResult{
			ID:            c.ID,
			Login:         c.Login,
			Contributions: c.Contributions,
		}
	}

	return results, nil
}

// SearchIssues searches for issues/PRs using GitHub search API and returns total count.
func (a *ResourceFetcherAdapter) SearchIssues(
	ctx context.Context,
	accessToken string,
	query string,
) (int, error) {
	return a.client.SearchIssues(ctx, accessToken, query)
}

// SearchIssueCountsBatch executes multiple search queries in a single GraphQL request.
func (a *ResourceFetcherAdapter) SearchIssueCountsBatch(
	ctx context.Context,
	accessToken string,
	queries map[string]string,
) (map[string]int, error) {
	return a.client.SearchIssueCountsBatch(ctx, accessToken, queries)
}
