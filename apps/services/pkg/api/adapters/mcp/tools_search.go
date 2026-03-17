package mcp

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchInput struct {
	Profile *string `json:"profile,omitempty" jsonschema:"Optional profile slug to scope search within"`
	Locale  string  `json:"locale,omitempty"  jsonschema:"Locale code (default: en)"`
	Query   string  `json:"query"             jsonschema:"required,Search query text"`
	Limit   int     `json:"limit,omitempty"   jsonschema:"Maximum results (default 20, max 100)"`
}

type searchResultBrief struct {
	Summary      *string `json:"summary,omitempty"`
	Kind         *string `json:"kind,omitempty"`
	ImageURI     *string `json:"image_uri,omitempty"`
	ProfileSlug  *string `json:"profile_slug,omitempty"`
	ProfileTitle *string `json:"profile_title,omitempty"`
	Type         string  `json:"type"`
	Slug         string  `json:"slug"`
	Title        string  `json:"title"`
}

type searchOutput struct {
	Results []searchResultBrief `json:"results"`
}

func registerSearchTools(server *mcp.Server, profileService *profiles.Service) {
	mcp.AddTool(
		server,
		&mcp.Tool{ //nolint:exhaustruct // external SDK type
			Name:        "search",
			Description: "Search across profiles, stories, and pages on the aya.is community platform",
		},
		createSearchHandler(profileService),
	)
}

func createSearchHandler(
	profileService *profiles.Service,
) func(context.Context, *mcp.CallToolRequest, searchInput) (
	*mcp.CallToolResult,
	searchOutput,
	error,
) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input searchInput,
	) (*mcp.CallToolResult, searchOutput, error) {
		locale := input.Locale
		if locale == "" {
			locale = defaultLocale
		}

		limit := input.Limit
		if limit <= 0 {
			limit = defaultLimit
		}

		if limit > maxLimit {
			limit = maxLimit
		}

		results, err := profileService.Search(ctx, locale, input.Query, input.Profile, int32(limit))
		if err != nil {
			return nil, searchOutput{}, fmt.Errorf("searching: %w", err)
		}

		output := searchOutput{
			Results: make([]searchResultBrief, 0, len(results)),
		}

		for _, result := range results {
			output.Results = append(output.Results, searchResultBrief{
				Type:         result.Type,
				Slug:         result.Slug,
				Title:        result.Title,
				Summary:      result.Summary,
				Kind:         result.Kind,
				ImageURI:     result.ImageURI,
				ProfileSlug:  result.ProfileSlug,
				ProfileTitle: result.ProfileTitle,
			})
		}

		return nil, output, nil
	}
}
