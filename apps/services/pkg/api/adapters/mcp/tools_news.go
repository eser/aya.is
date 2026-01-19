package mcp

import (
	"context"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var ErrNewsNotFound = errors.New("news item not found")

const newsKind = "news"

type listNewsInput struct {
	Locale          string  `json:"locale,omitempty"           jsonschema:"description=Locale code (default: en)"`
	PublicationSlug *string `json:"publication_slug,omitempty" jsonschema:"description=Filter by publication profile slug"`
	Limit           int     `json:"limit,omitempty"            jsonschema:"description=Maximum results (default 20 max 100)"`
	Cursor          *string `json:"cursor,omitempty"           jsonschema:"description=Pagination cursor for next page"`
}

type newsBrief struct {
	Slug            string  `json:"slug"`
	Title           string  `json:"title"`
	Summary         string  `json:"summary"`
	IsFeatured      bool    `json:"is_featured"`
	StoryPictureURI *string `json:"story_picture_uri,omitempty"`
	AuthorName      *string `json:"author_name,omitempty"`
	AuthorSlug      *string `json:"author_slug,omitempty"`
}

type listNewsOutput struct {
	News       []newsBrief `json:"news"`
	NextCursor *string     `json:"next_cursor,omitempty"`
}

type getNewsInput struct {
	Locale string `json:"locale,omitempty" jsonschema:"description=Locale code (default: en)"`
	Slug   string `json:"slug"             jsonschema:"required,description=News item slug"`
}

type getNewsOutput struct {
	Slug            string             `json:"slug"`
	Title           string             `json:"title"`
	Summary         string             `json:"summary"`
	Content         string             `json:"content"`
	IsFeatured      bool               `json:"is_featured"`
	StoryPictureURI *string            `json:"story_picture_uri,omitempty"`
	AuthorName      *string            `json:"author_name,omitempty"`
	AuthorSlug      *string            `json:"author_slug,omitempty"`
	Publications    []publicationBrief `json:"publications"`
}

func registerNewsTools(server *mcp.Server, storyService *stories.Service) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_news",
			Description: "Get a list of news items on the AYA platform",
		},
		createListNewsHandler(storyService),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_news",
			Description: "Get the full contents of a news item",
		},
		createGetNewsHandler(storyService),
	)
}

func createListNewsHandler(
	storyService *stories.Service,
) func(context.Context, *mcp.CallToolRequest, listNewsInput) (
	*mcp.CallToolResult,
	listNewsOutput,
	error,
) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input listNewsInput,
	) (*mcp.CallToolResult, listNewsOutput, error) {
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
		cursor.Filters["kind"] = newsKind

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
			return nil, listNewsOutput{}, err
		}

		output := listNewsOutput{
			News:       make([]newsBrief, 0, len(result.Data)),
			NextCursor: result.CursorPtr,
		}

		for _, story := range result.Data {
			brief := newsBrief{
				Slug:            story.Slug,
				Title:           story.Title,
				Summary:         story.Summary,
				IsFeatured:      story.IsFeatured,
				StoryPictureURI: story.StoryPictureURI,
			}

			if story.AuthorProfile != nil {
				brief.AuthorName = &story.AuthorProfile.Title
				brief.AuthorSlug = &story.AuthorProfile.Slug
			}

			output.News = append(output.News, brief)
		}

		return nil, output, nil
	}
}

func createGetNewsHandler(
	storyService *stories.Service,
) func(context.Context, *mcp.CallToolRequest, getNewsInput) (
	*mcp.CallToolResult,
	getNewsOutput,
	error,
) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input getNewsInput,
	) (*mcp.CallToolResult, getNewsOutput, error) {
		locale := input.Locale
		if locale == "" {
			locale = defaultLocale
		}

		result, err := storyService.GetBySlug(ctx, locale, input.Slug)
		if err != nil {
			return nil, getNewsOutput{}, err
		}

		if result == nil {
			return nil, getNewsOutput{}, ErrNewsNotFound
		}

		if result.Kind != newsKind {
			return nil, getNewsOutput{}, ErrNewsNotFound
		}

		output := getNewsOutput{
			Slug:            result.Slug,
			Title:           result.Title,
			Summary:         result.Summary,
			Content:         result.Content,
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
