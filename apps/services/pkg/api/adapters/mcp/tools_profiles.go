package mcp

import (
	"context"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultLocale = "en"
	defaultLimit  = 20
	maxLimit      = 100
)

type listProfilesInput struct {
	Locale string  `json:"locale,omitempty" jsonschema:"Locale code (default: en)"`
	Kind   string  `json:"kind,omitempty"   jsonschema:"Filter by kind: individual, organization, product"`
	Limit  int     `json:"limit,omitempty"  jsonschema:"Maximum results (default 20, max 100)"`
	Cursor *string `json:"cursor,omitempty" jsonschema:"Pagination cursor for next page"`
}

type profileBrief struct {
	Slug              string  `json:"slug"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Kind              string  `json:"kind"`
	CustomDomain      *string `json:"custom_domain,omitempty"`
	ProfilePictureURI *string `json:"profile_picture_uri,omitempty"`
}

type pageBrief struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type linkBrief struct {
	Kind  string `json:"kind"`
	Title string `json:"title"`
	URI   string `json:"uri"`
}

type listProfilesOutput struct {
	Profiles   []profileBrief `json:"profiles"`
	NextCursor *string        `json:"next_cursor,omitempty"`
}

type getProfileInput struct {
	Locale string `json:"locale,omitempty" jsonschema:"Locale code (default: en)"`
	Slug   string `json:"slug"             jsonschema:"required,Profile slug"`
}

type getProfileOutput struct {
	Slug              string      `json:"slug"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	Kind              string      `json:"kind"`
	CustomDomain      *string     `json:"custom_domain,omitempty"`
	ProfilePictureURI *string     `json:"profile_picture_uri,omitempty"`
	Pronouns          *string     `json:"pronouns,omitempty"`
	Pages             []pageBrief `json:"pages"`
	Links             []linkBrief `json:"links"`
}

type getProfilePageInput struct {
	Locale      string `json:"locale,omitempty" jsonschema:"Locale code (default: en)"`
	ProfileSlug string `json:"profile_slug"     jsonschema:"required,Profile slug"`
	PageSlug    string `json:"page_slug"        jsonschema:"required,Page slug"`
}

type getProfilePageOutput struct {
	Slug            string  `json:"slug"`
	Title           string  `json:"title"`
	Summary         string  `json:"summary"`
	Content         string  `json:"content"`
	CoverPictureURI *string `json:"cover_picture_uri,omitempty"`
}

func registerProfileTools(server *mcp.Server, profileService *profiles.Service) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_profiles",
			Description: "Get a list of profiles (people, organizations, products) on the AYA platform",
		},
		createListProfilesHandler(profileService),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_profile",
			Description: "Get detailed information about a specific profile including their pages and links",
		},
		createGetProfileHandler(profileService),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_profile_page",
			Description: "Get the full contents of a profile page",
		},
		createGetProfilePageHandler(profileService),
	)
}

func createListProfilesHandler(
	profileService *profiles.Service,
) func(context.Context, *mcp.CallToolRequest, listProfilesInput) (
	*mcp.CallToolResult,
	listProfilesOutput,
	error,
) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input listProfilesInput,
	) (*mcp.CallToolResult, listProfilesOutput, error) {
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
		if input.Kind != "" {
			cursor.Filters["kind"] = input.Kind
		}

		result, err := profileService.List(ctx, locale, cursor)
		if err != nil {
			return nil, listProfilesOutput{}, err
		}

		output := listProfilesOutput{
			Profiles:   make([]profileBrief, 0, len(result.Data)),
			NextCursor: result.CursorPtr,
		}

		for _, profile := range result.Data {
			output.Profiles = append(output.Profiles, profileBrief{
				Slug:              profile.Slug,
				Name:              profile.Title,
				Description:       profile.Description,
				Kind:              profile.Kind,
				CustomDomain:      profile.CustomDomain,
				ProfilePictureURI: profile.ProfilePictureURI,
			})
		}

		return nil, output, nil
	}
}

func createGetProfileHandler(
	profileService *profiles.Service,
) func(context.Context, *mcp.CallToolRequest, getProfileInput) (
	*mcp.CallToolResult,
	getProfileOutput,
	error,
) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input getProfileInput,
	) (*mcp.CallToolResult, getProfileOutput, error) {
		locale := input.Locale
		if locale == "" {
			locale = defaultLocale
		}

		result, err := profileService.GetBySlugEx(ctx, locale, input.Slug)
		if err != nil {
			return nil, getProfileOutput{}, err
		}

		if result == nil {
			return nil, getProfileOutput{}, profiles.ErrProfileNotFound
		}

		output := getProfileOutput{
			Slug:              result.Slug,
			Name:              result.Title,
			Description:       result.Description,
			Kind:              result.Kind,
			CustomDomain:      result.CustomDomain,
			ProfilePictureURI: result.ProfilePictureURI,
			Pronouns:          result.Pronouns,
			Pages:             make([]pageBrief, 0, len(result.Pages)),
			Links:             make([]linkBrief, 0, len(result.Links)),
		}

		for _, page := range result.Pages {
			output.Pages = append(output.Pages, pageBrief{
				Slug:    page.Slug,
				Title:   page.Title,
				Summary: page.Summary,
			})
		}

		for _, link := range result.Links {
			output.Links = append(output.Links, linkBrief{
				Kind:  link.Kind,
				Title: link.Title,
				URI:   link.URI,
			})
		}

		return nil, output, nil
	}
}

func createGetProfilePageHandler(
	profileService *profiles.Service,
) func(context.Context, *mcp.CallToolRequest, getProfilePageInput) (
	*mcp.CallToolResult,
	getProfilePageOutput,
	error,
) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input getProfilePageInput,
	) (*mcp.CallToolResult, getProfilePageOutput, error) {
		locale := input.Locale
		if locale == "" {
			locale = defaultLocale
		}

		result, err := profileService.GetPageBySlug(ctx, locale, input.ProfileSlug, input.PageSlug)
		if err != nil {
			return nil, getProfilePageOutput{}, err
		}

		if result == nil {
			return nil, getProfilePageOutput{}, profiles.ErrProfileNotFound
		}

		output := getProfilePageOutput{
			Slug:            result.Slug,
			Title:           result.Title,
			Summary:         result.Summary,
			Content:         result.Content,
			CoverPictureURI: result.CoverPictureURI,
		}

		return nil, output, nil
	}
}
