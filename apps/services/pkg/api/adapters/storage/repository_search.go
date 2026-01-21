package storage

import (
	"context"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

// Search performs a unified search across profiles, stories, and profile pages.
func (r *Repository) Search(
	ctx context.Context,
	localeCode string,
	query string,
	limit int32,
) ([]*profiles.SearchResult, error) {
	results := make([]*profiles.SearchResult, 0)

	// Search profiles
	profileRows, err := r.queries.SearchProfiles(ctx, SearchProfilesParams{
		Query:      query,
		LocaleCode: localeCode,
		LimitCount: limit,
	})
	if err != nil {
		return nil, err
	}

	for _, p := range profileRows {
		results = append(results, &profiles.SearchResult{
			Type:     "profile",
			ID:       p.ID,
			Slug:     p.Slug,
			Title:    p.Title,
			Summary:  &p.Description,
			ImageURI: vars.ToStringPtr(p.ProfilePictureURI),
			Rank:     p.Rank,
		})
	}

	// Search stories
	storyRows, err := r.queries.SearchStories(ctx, SearchStoriesParams{
		Query:      query,
		LocaleCode: localeCode,
		LimitCount: limit,
	})
	if err != nil {
		return nil, err
	}

	for _, s := range storyRows {
		results = append(results, &profiles.SearchResult{
			Type:        "story",
			ID:          s.ID,
			Slug:        s.Slug,
			Title:       s.Title,
			Summary:     &s.Summary,
			ImageURI:    vars.ToStringPtr(s.StoryPictureURI),
			ProfileSlug: vars.ToStringPtr(s.AuthorSlug),
			Rank:        s.Rank,
		})
	}

	// Search profile pages
	pageRows, err := r.queries.SearchProfilePages(ctx, SearchProfilePagesParams{
		Query:      query,
		LocaleCode: localeCode,
		LimitCount: limit,
	})
	if err != nil {
		return nil, err
	}

	for _, pp := range pageRows {
		results = append(results, &profiles.SearchResult{
			Type:        "page",
			ID:          pp.ID,
			Slug:        pp.Slug,
			Title:       pp.Title,
			Summary:     &pp.Summary,
			ImageURI:    vars.ToStringPtr(pp.CoverPictureURI),
			ProfileSlug: &pp.ProfileSlug,
			Rank:        pp.Rank,
		})
	}

	// Sort by rank (descending)
	// Since each query is already sorted, we merge-sort by rank
	sortByRank(results)

	// Limit to requested count
	if int32(len(results)) > limit {
		results = results[:limit]
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
