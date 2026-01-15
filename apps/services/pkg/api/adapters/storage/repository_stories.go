package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

var ErrFailedToParseStoryWithChildren = errors.New("failed to parse story with children")

func (r *Repository) GetStoryIDBySlug(ctx context.Context, slug string) (string, error) {
	var result string

	err := r.cache.Execute(
		ctx,
		"story_id_by_slug:"+slug,
		&result,
		func(ctx context.Context) (any, error) {
			row, err := r.queries.GetStoryIDBySlug(ctx, GetStoryIDBySlugParams{Slug: slug})
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, nil //nolint:nilnil
				}

				return nil, err
			}

			return row, nil
		},
	)

	return result, err //nolint:wrapcheck
}

func (r *Repository) GetStoryByID(
	ctx context.Context,
	localeCode string,
	id string,
	authorProfileID *string,
) (*stories.StoryWithChildren, error) {
	getStoryByIDParams := GetStoryByIDParams{
		LocaleCode: localeCode,
		ID:         id,
		FilterPublicationProfileID: sql.NullString{
			String: "",
			Valid:  false,
		},
		FilterAuthorProfileID: sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	if authorProfileID != nil {
		getStoryByIDParams.FilterAuthorProfileID = sql.NullString{
			String: *authorProfileID,
			Valid:  true,
		}
	}

	row, err := r.queries.GetStoryByID(ctx, getStoryByIDParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result, err := r.parseStoryWithChildren(
		row.Profile,
		row.ProfileTx,
		row.Story,
		row.StoryTx,
		row.Publications,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) ListStoriesOfPublication(
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*stories.StoryWithChildren], error) {
	var wrappedResponse cursors.Cursored[[]*stories.StoryWithChildren]

	rows, err := r.queries.ListStoriesOfPublication(
		ctx,
		ListStoriesOfPublicationParams{
			LocaleCode: localeCode,
			FilterKind: vars.MapValueToNullString(cursor.Filters, "kind"),
			FilterAuthorProfileID: vars.MapValueToNullString(
				cursor.Filters,
				"author_profile_id",
			),
			FilterPublicationProfileID: vars.MapValueToNullString(
				cursor.Filters,
				"publication_profile_id",
			),
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	result := make([]*stories.StoryWithChildren, len(rows))

	for i, row := range rows {
		storyWithChildren, err := r.parseStoryWithChildren(
			row.Profile,
			row.ProfileTx,
			row.Story,
			row.StoryTx,
			row.Publications,
		)
		if err != nil {
			return wrappedResponse, err
		}

		result[i] = storyWithChildren
	}

	wrappedResponse.Data = result

	if len(result) == cursor.Limit {
		wrappedResponse.CursorPtr = &result[len(result)-1].ID
	}

	return wrappedResponse, nil
}

func (r *Repository) parseStoryWithChildren( //nolint:funlen
	profile Profile,
	profileTx ProfileTx,
	story Story,
	storyTx StoryTx,
	publications json.RawMessage,
) (*stories.StoryWithChildren, error) {
	storyWithChildren := &stories.StoryWithChildren{
		Story: &stories.Story{
			ID:              story.ID,
			AuthorProfileID: vars.ToStringPtr(story.AuthorProfileID),
			Slug:            story.Slug,
			Kind:            story.Kind,
			Status:          story.Status,
			IsFeatured:      story.IsFeatured,
			StoryPictureURI: vars.ToStringPtr(story.StoryPictureURI),
			Title:           storyTx.Title,
			Summary:         storyTx.Summary,
			Content:         storyTx.Content,
			Properties:      vars.ToObject(story.Properties),
			CreatedAt:       story.CreatedAt,
			UpdatedAt:       vars.ToTimePtr(story.UpdatedAt),
			DeletedAt:       vars.ToTimePtr(story.DeletedAt),
		},
		AuthorProfile: &profiles.Profile{
			ID:                profile.ID,
			Slug:              profile.Slug,
			Kind:              profile.Kind,
			CustomDomain:      vars.ToStringPtr(profile.CustomDomain),
			ProfilePictureURI: vars.ToStringPtr(profile.ProfilePictureURI),
			Pronouns:          vars.ToStringPtr(profile.Pronouns),
			Title:             profileTx.Title,
			Description:       profileTx.Description,
			Properties:        vars.ToObject(profile.Properties),
			CreatedAt:         profile.CreatedAt,
			UpdatedAt:         vars.ToTimePtr(profile.UpdatedAt),
			DeletedAt:         vars.ToTimePtr(profile.DeletedAt),
		},
		Publications: nil,
	}

	var publicationProfiles []struct {
		Profile struct {
			CreatedAt         time.Time        `db:"created_at"          json:"created_at"`
			CustomDomain      *string          `db:"custom_domain"       json:"custom_domain"`
			ProfilePictureURI *string          `db:"profile_picture_uri" json:"profile_picture_uri"`
			Pronouns          *string          `db:"pronouns"            json:"pronouns"`
			Properties        *json.RawMessage `db:"properties"          json:"properties"`
			UpdatedAt         *time.Time       `db:"updated_at"          json:"updated_at"`
			DeletedAt         *time.Time       `db:"deleted_at"          json:"deleted_at"`
			ID                string           `db:"id"                  json:"id"`
			Slug              string           `db:"slug"                json:"slug"`
			Kind              string           `db:"kind"                json:"kind"`
		} `json:"profile"`
		ProfileTx struct {
			Properties  *json.RawMessage `db:"properties"  json:"properties"`
			ProfileID   string           `db:"profile_id"  json:"profile_id"`
			LocaleCode  string           `db:"locale_code" json:"locale_code"`
			Title       string           `db:"title"       json:"title"`
			Description string           `db:"description" json:"description"`
		} `json:"profile_tx"`
	}

	err := json.Unmarshal(publications, &publicationProfiles)
	if err != nil {
		r.logger.Error("failed to unmarshal publications", "error", err)

		return nil, fmt.Errorf("%w: %w", ErrFailedToParseStoryWithChildren, err)
	}

	storyWithChildren.Publications = make([]*profiles.Profile, len(publicationProfiles))
	for j, publicationProfile := range publicationProfiles {
		storyWithChildren.Publications[j] = &profiles.Profile{
			ID:                publicationProfile.Profile.ID,
			Slug:              publicationProfile.Profile.Slug,
			Kind:              publicationProfile.Profile.Kind,
			CustomDomain:      publicationProfile.Profile.CustomDomain,
			ProfilePictureURI: publicationProfile.Profile.ProfilePictureURI,
			Pronouns:          publicationProfile.Profile.Pronouns,
			Title:             publicationProfile.ProfileTx.Title,
			Description:       publicationProfile.ProfileTx.Description,
			Properties:        publicationProfile.Profile.Properties,
			CreatedAt:         publicationProfile.Profile.CreatedAt,
			UpdatedAt:         publicationProfile.Profile.UpdatedAt,
			DeletedAt:         publicationProfile.Profile.DeletedAt,
		}
	}

	return storyWithChildren, nil
}
