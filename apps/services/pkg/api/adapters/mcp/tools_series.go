package mcp

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/story_series"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type listSeriesInput struct {
	Locale string `json:"locale,omitempty" jsonschema:"Locale code (default: en)"`
}

type seriesBrief struct {
	Slug             string  `json:"slug"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	SeriesPictureURI *string `json:"series_picture_uri,omitempty"`
}

type listSeriesOutput struct {
	Series []seriesBrief `json:"series"`
}

type getSeriesInput struct {
	Locale string `json:"locale,omitempty" jsonschema:"Locale code (default: en)"`
	Slug   string `json:"slug"             jsonschema:"required,Series slug"`
}

type getSeriesOutput struct {
	Slug             string       `json:"slug"`
	Title            string       `json:"title"`
	Description      string       `json:"description"`
	SeriesPictureURI *string      `json:"series_picture_uri,omitempty"`
	Stories          []storyBrief `json:"stories"`
}

func registerSeriesTools(
	server *mcp.Server,
	storySeriesService *story_series.Service,
	storyService *stories.Service,
) {
	mcp.AddTool(
		server,
		&mcp.Tool{ //nolint:exhaustruct // external SDK type
			Name:        "list_series",
			Description: "Get a list of story series on the aya.is community platform",
		},
		createListSeriesHandler(storySeriesService),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{ //nolint:exhaustruct // external SDK type
			Name:        "get_series",
			Description: "Get detailed information about a story series including its stories",
		},
		createGetSeriesHandler(storySeriesService, storyService),
	)
}

func createListSeriesHandler(
	storySeriesService *story_series.Service,
) func(context.Context, *mcp.CallToolRequest, listSeriesInput) (
	*mcp.CallToolResult,
	listSeriesOutput,
	error,
) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input listSeriesInput,
	) (*mcp.CallToolResult, listSeriesOutput, error) {
		locale := input.Locale
		if locale == "" {
			locale = defaultLocale
		}

		result, err := storySeriesService.List(ctx, locale)
		if err != nil {
			return nil, listSeriesOutput{}, fmt.Errorf("listing series: %w", err)
		}

		output := listSeriesOutput{
			Series: make([]seriesBrief, 0, len(result)),
		}

		for _, series := range result {
			output.Series = append(output.Series, seriesBrief{
				Slug:             series.Slug,
				Title:            series.Title,
				Description:      series.Description,
				SeriesPictureURI: series.SeriesPictureURI,
			})
		}

		return nil, output, nil
	}
}

func createGetSeriesHandler(
	storySeriesService *story_series.Service,
	storyService *stories.Service,
) func(context.Context, *mcp.CallToolRequest, getSeriesInput) (
	*mcp.CallToolResult,
	getSeriesOutput,
	error,
) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input getSeriesInput,
	) (*mcp.CallToolResult, getSeriesOutput, error) {
		locale := input.Locale
		if locale == "" {
			locale = defaultLocale
		}

		series, err := storySeriesService.GetBySlug(ctx, locale, input.Slug)
		if err != nil {
			return nil, getSeriesOutput{}, fmt.Errorf("getting series: %w", err)
		}

		if series == nil {
			return nil, getSeriesOutput{}, story_series.ErrSeriesNotFound
		}

		seriesStories, err := storyService.ListBySeriesID(ctx, locale, series.ID)
		if err != nil {
			return nil, getSeriesOutput{}, fmt.Errorf("listing series stories: %w", err)
		}

		output := getSeriesOutput{
			Slug:             series.Slug,
			Title:            series.Title,
			Description:      series.Description,
			SeriesPictureURI: series.SeriesPictureURI,
			Stories:          make([]storyBrief, 0, len(seriesStories)),
		}

		for _, story := range seriesStories {
			output.Stories = append(output.Stories, storyBrief{
				Slug:            story.Slug,
				Title:           story.Title,
				Summary:         story.Summary,
				Kind:            story.Kind,
				StoryPictureURI: story.StoryPictureURI,
				AuthorName:      nil,
				AuthorSlug:      nil,
			})
		}

		return nil, output, nil
	}
}
