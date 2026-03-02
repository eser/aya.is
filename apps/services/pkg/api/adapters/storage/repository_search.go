package storage

import (
	"context"
	"database/sql"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

// Search performs a unified search across profiles, stories, and profile pages.
// If profileSlug is provided, search is scoped to that profile only.
func (r *Repository) Search(
	ctx context.Context,
	localeCode string,
	query string,
	profileSlug *string,
	limit int32,
) ([]*profiles.SearchResult, error) {
	results := make([]*profiles.SearchResult, 0)

	// Convert profile slug to sql.NullString
	filterProfileSlug := sql.NullString{} //nolint:exhaustruct
	if profileSlug != nil && *profileSlug != "" {
		filterProfileSlug = sql.NullString{String: *profileSlug, Valid: true}
	}

	searchParams := searchQueryParams{
		query:             query,
		localeCode:        localeCode,
		filterProfileSlug: filterProfileSlug,
		limitCount:        limit,
	}

	profileResults, err := r.searchProfiles(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	results = append(results, profileResults...)

	storyResults, err := r.searchStories(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	results = append(results, storyResults...)

	pageResults, err := r.searchPages(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	results = append(results, pageResults...)

	// Sort by rank (descending)
	sortByRank(results)

	// Limit to requested count
	if int32(len(results)) > limit {
		results = results[:limit]
	}

	return results, nil
}

// searchQueryParams holds common parameters for search sub-queries.
type searchQueryParams struct {
	query             string
	localeCode        string
	filterProfileSlug sql.NullString
	limitCount        int32
}

// searchProfiles searches profiles and returns search results.
func (r *Repository) searchProfiles(
	ctx context.Context,
	params searchQueryParams,
) ([]*profiles.SearchResult, error) {
	profileRows, err := r.queries.SearchProfiles(ctx, SearchProfilesParams{
		Query:             params.query,
		LocaleCode:        params.localeCode,
		FilterProfileSlug: params.filterProfileSlug,
		LimitCount:        params.limitCount,
	})
	if err != nil {
		return nil, err
	}

	results := make([]*profiles.SearchResult, 0, len(profileRows))

	for _, profileRow := range profileRows {
		kind := profileRow.Kind
		results = append(results, &profiles.SearchResult{
			Type:         "profile",
			ID:           profileRow.ID,
			Slug:         profileRow.Slug,
			Title:        profileRow.Title,
			Summary:      &profileRow.Description,
			ImageURI:     vars.ToStringPtr(profileRow.ProfilePictureURI),
			ProfileSlug:  nil,
			ProfileTitle: nil,
			Kind:         &kind,
			Rank:         profileRow.Rank,
		})
	}

	return results, nil
}

// searchStories searches stories and returns search results.
func (r *Repository) searchStories(
	ctx context.Context,
	params searchQueryParams,
) ([]*profiles.SearchResult, error) {
	storyRows, err := r.queries.SearchStories(ctx, SearchStoriesParams{
		Query:             params.query,
		LocaleCode:        params.localeCode,
		FilterProfileSlug: params.filterProfileSlug,
		LimitCount:        params.limitCount,
	})
	if err != nil {
		return nil, err
	}

	results := make([]*profiles.SearchResult, 0, len(storyRows))

	for _, storyRow := range storyRows {
		kind := storyRow.Kind
		results = append(results, &profiles.SearchResult{
			Type:         "story",
			ID:           storyRow.ID,
			Slug:         storyRow.Slug,
			Title:        storyRow.Title,
			Summary:      &storyRow.Summary,
			ImageURI:     vars.ToStringPtr(storyRow.StoryPictureURI),
			ProfileSlug:  vars.ToStringPtr(storyRow.AuthorSlug),
			ProfileTitle: vars.ToStringPtr(storyRow.AuthorTitle),
			Kind:         &kind,
			Rank:         storyRow.Rank,
		})
	}

	return results, nil
}

// searchPages searches profile pages and returns search results.
func (r *Repository) searchPages(
	ctx context.Context,
	params searchQueryParams,
) ([]*profiles.SearchResult, error) {
	pageRows, err := r.queries.SearchProfilePages(ctx, SearchProfilePagesParams{
		Query:             params.query,
		LocaleCode:        params.localeCode,
		FilterProfileSlug: params.filterProfileSlug,
		LimitCount:        params.limitCount,
	})
	if err != nil {
		return nil, err
	}

	results := make([]*profiles.SearchResult, 0, len(pageRows))

	for _, pageRow := range pageRows {
		results = append(results, &profiles.SearchResult{
			Type:         "page",
			ID:           pageRow.ID,
			Slug:         pageRow.Slug,
			Title:        pageRow.Title,
			Summary:      &pageRow.Summary,
			ImageURI:     vars.ToStringPtr(pageRow.CoverPictureURI),
			ProfileSlug:  &pageRow.ProfileSlug,
			ProfileTitle: &pageRow.ProfileTitle,
			Kind:         nil,
			Rank:         pageRow.Rank,
		})
	}

	return results, nil
}

// sortByRank sorts results by rank in descending order.
func sortByRank(results []*profiles.SearchResult) {
	for i := range len(results) - 1 {
		for j := i + 1; j < len(results); j++ {
			if results[j].Rank > results[i].Rank {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}
