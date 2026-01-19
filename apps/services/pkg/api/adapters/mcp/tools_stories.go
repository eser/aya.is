package mcp

import (
	"context"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var ErrStoryNotFound = errors.New("story not found")

type listStoriesInput struct {
	Locale          string  `json:"locale,omitempty"           jsonschema:"Locale code (default: en)"`
	PublicationSlug *string `json:"publication_slug,omitempty" jsonschema:"Filter by publication profile slug"`
	Limit           int     `json:"limit,omitempty"            jsonschema:"Maximum results (default 20, max 100)"`
	Cursor          *string `json:"cursor,omitempty"           jsonschema:"Pagination cursor for next page"`
}

type storyBrief struct {
	Slug            string  `json:"slug"`
	Title           string  `json:"title"`
	Summary         string  `json:"summary"`
	Kind            string  `json:"kind"`
	Status          string  `json:"status"`
	IsFeatured      bool    `json:"is_featured"`
	StoryPictureURI *string `json:"story_picture_uri,omitempty"`
	AuthorName      *string `json:"author_name,omitempty"`
	AuthorSlug      *string `json:"author_slug,omitempty"`
}

type listStoriesOutput struct {
	Stories    []storyBrief `json:"stories"`
	NextCursor *string      `json:"next_cursor,omitempty"`
}

type getStoryInput struct {
	Locale string `json:"locale,omitempty" jsonschema:"Locale code (default: en)"`
	Slug   string `json:"slug"             jsonschema:"required,Story slug"`
}

type publicationBrief struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type getStoryOutput struct {
	Slug            string             `json:"slug"`
	Title           string             `json:"title"`
	Summary         string             `json:"summary"`
	Content         string             `json:"content"`
	Kind            string             `json:"kind"`
	Status          string             `json:"status"`
	IsFeatured      bool               `json:"is_featured"`
	StoryPictureURI *string            `json:"story_picture_uri,omitempty"`
	AuthorName      *string            `json:"author_name,omitempty"`
	AuthorSlug      *string            `json:"author_slug,omitempty"`
	Publications    []publicationBrief `json:"publications"`
}

func registerStoryTools(server *mcp.Server, storyService *stories.Service) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_stories",
			Description: "Get a list of stories and articles on the AYA platform",
		},
		createListStoriesHandler(storyService),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_story",
			Description: "Get the full contents of a story or article",
		},
		createGetStoryHandler(storyService),
	)
}

func createListStoriesHandler(
	storyService *stories.Service,
) func(context.Context, *mcp.CallToolRequest, listStoriesInput) (
	*mcp.CallToolResult,
	listStoriesOutput,
	error,
) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input listStoriesInput,
	) (*mcp.CallToolResult, listStoriesOutput, error) {
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

		cursor := cursors.NewCursor(limit, input.Cursor)
		cursor.Filters["kind"] = "story"

		var (
			result cursors.Cursored[[]*stories.StoryWithChildren]
			err    error
		)

		if input.PublicationSlug != nil && *input.PublicationSlug != "" {
			result, err = storyService.ListByPublicationProfileSlug(
				ctx,
				locale,
				*input.PublicationSlug,
				cursor,
			)
		} else {
			result, err = storyService.List(ctx, locale, cursor)
		}

		if err != nil {
			return nil, listStoriesOutput{}, err
		}

		output := listStoriesOutput{
			Stories:    make([]storyBrief, 0, len(result.Data)),
			NextCursor: result.CursorPtr,
		}

		for _, story := range result.Data {
			brief := storyBrief{
				Slug:            story.Slug,
				Title:           story.Title,
				Summary:         story.Summary,
				Kind:            story.Kind,
				Status:          story.Status,
				IsFeatured:      story.IsFeatured,
				StoryPictureURI: story.StoryPictureURI,
			}

			if story.AuthorProfile != nil {
				brief.AuthorName = &story.AuthorProfile.Title
				brief.AuthorSlug = &story.AuthorProfile.Slug
			}

			output.Stories = append(output.Stories, brief)
		}

		return nil, output, nil
	}
}

func createGetStoryHandler(
	storyService *stories.Service,
) func(context.Context, *mcp.CallToolRequest, getStoryInput) (
	*mcp.CallToolResult,
	getStoryOutput,
	error,
) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input getStoryInput,
	) (*mcp.CallToolResult, getStoryOutput, error) {
		locale := input.Locale
		if locale == "" {
			locale = defaultLocale
		}

		result, err := storyService.GetBySlug(ctx, locale, input.Slug)
		if err != nil {
			return nil, getStoryOutput{}, err
		}

		if result == nil {
			return nil, getStoryOutput{}, ErrStoryNotFound
		}

		output := getStoryOutput{
			Slug:            result.Slug,
			Title:           result.Title,
			Summary:         result.Summary,
			Content:         result.Content,
			Kind:            result.Kind,
			Status:          result.Status,
			IsFeatured:      result.IsFeatured,
			StoryPictureURI: result.StoryPictureURI,
			Publications:    make([]publicationBrief, 0, len(result.Publications)),
		}

		if result.AuthorProfile != nil {
			output.AuthorName = &result.AuthorProfile.Title
			output.AuthorSlug = &result.AuthorProfile.Slug
		}

		for _, pub := range result.Publications {
			output.Publications = append(output.Publications, publicationBrief{
				Slug: pub.Slug,
				Name: pub.Title,
			})
		}

		return nil, output, nil
	}
}
