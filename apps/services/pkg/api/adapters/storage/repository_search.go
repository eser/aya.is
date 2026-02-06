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
	filterProfileSlug := sql.NullString{}
	if profileSlug != nil && *profileSlug != "" {
		filterProfileSlug = sql.NullString{String: *profileSlug, Valid: true}
	}

	// Search profiles
	profileRows, err := r.queries.SearchProfiles(ctx, SearchProfilesParams{
		Query:             query,
		LocaleCode:        localeCode,
		FilterProfileSlug: filterProfileSlug,
		LimitCount:        limit,
	})
	if err != nil {
		return nil, err
	}

	for _, p := range profileRows {
		kind := p.Kind
		results = append(results, &profiles.SearchResult{
			Type:     "profile",
			ID:       p.ID,
			Slug:     p.Slug,
			Title:    p.Title,
			Summary:  &p.Description,
			ImageURI: vars.ToStringPtr(p.ProfilePictureURI),
			Kind:     &kind,
			Rank:     p.Rank,
		})
	}

	// Search stories
	storyRows, err := r.queries.SearchStories(ctx, SearchStoriesParams{
		Query:             query,
		LocaleCode:        localeCode,
		FilterProfileSlug: filterProfileSlug,
		LimitCount:        limit,
	})
	if err != nil {
		return nil, err
	}

	for _, s := range storyRows {
		kind := s.Kind
		results = append(results, &profiles.SearchResult{
			Type:         "story",
			ID:           s.ID,
			Slug:         s.Slug,
			Title:        s.Title,
			Summary:      &s.Summary,
			ImageURI:     vars.ToStringPtr(s.StoryPictureURI),
			ProfileSlug:  vars.ToStringPtr(s.AuthorSlug),
			ProfileTitle: vars.ToStringPtr(s.AuthorTitle),
			Kind:         &kind,
			Rank:         s.Rank,
		})
	}

	// Search profile pages
	pageRows, err := r.queries.SearchProfilePages(ctx, SearchProfilePagesParams{
		Query:             query,
		LocaleCode:        localeCode,
		FilterProfileSlug: filterProfileSlug,
		LimitCount:        limit,
	})
	if err != nil {
		return nil, err
	}

	for _, pp := range pageRows {
		results = append(results, &profiles.SearchResult{
			Type:         "page",
			ID:           pp.ID,
			Slug:         pp.Slug,
			Title:        pp.Title,
			Summary:      &pp.Summary,
			ImageURI:     vars.ToStringPtr(pp.CoverPictureURI),
			ProfileSlug:  &pp.ProfileSlug,
			ProfileTitle: &pp.ProfileTitle,
			Rank:         pp.Rank,
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
