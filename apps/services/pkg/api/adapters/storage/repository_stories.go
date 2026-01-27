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

func (r *Repository) GetStoryIDBySlugIncludingDeleted(
	ctx context.Context,
	slug string,
) (string, error) {
	row, err := r.queries.GetStoryIDBySlugIncludingDeleted(
		ctx,
		GetStoryIDBySlugIncludingDeletedParams{Slug: slug},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return row, nil
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

// CRUD methods for stories

func (r *Repository) InsertStory(
	ctx context.Context,
	id string,
	authorProfileID string,
	slug string,
	kind string,
	status string,
	isFeatured bool,
	storyPictureURI *string,
	properties map[string]any,
	publishedAt *time.Time,
) (*stories.Story, error) {
	params := InsertStoryParams{
		ID:              id,
		AuthorProfileID: sql.NullString{String: authorProfileID, Valid: true},
		Slug:            slug,
		Kind:            kind,
		Status:          status,
		IsFeatured:      isFeatured,
		StoryPictureURI: vars.ToSQLNullString(storyPictureURI),
		Properties:      vars.ToSQLNullRawMessage(properties),
		PublishedAt:     vars.ToSQLNullTime(publishedAt),
	}

	row, err := r.queries.InsertStory(ctx, params)
	if err != nil {
		return nil, err
	}

	return &stories.Story{
		ID:              row.ID,
		AuthorProfileID: vars.ToStringPtr(row.AuthorProfileID),
		Slug:            row.Slug,
		Kind:            row.Kind,
		Status:          row.Status,
		IsFeatured:      row.IsFeatured,
		StoryPictureURI: vars.ToStringPtr(row.StoryPictureURI),
		PublishedAt:     vars.ToTimePtr(row.PublishedAt),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) InsertStoryTx(
	ctx context.Context,
	storyID string,
	localeCode string,
	title string,
	summary string,
	content string,
) error {
	params := InsertStoryTxParams{
		StoryID:    storyID,
		LocaleCode: localeCode,
		Title:      title,
		Summary:    summary,
		Content:    content,
	}

	return r.queries.InsertStoryTx(ctx, params)
}

func (r *Repository) InsertStoryPublication(
	ctx context.Context,
	id string,
	storyID string,
	profileID string,
	kind string,
	properties map[string]any,
) error {
	params := InsertStoryPublicationParams{
		ID:         id,
		StoryID:    storyID,
		ProfileID:  profileID,
		Kind:       kind,
		Properties: vars.ToSQLNullRawMessage(properties),
	}

	_, err := r.queries.InsertStoryPublication(ctx, params)

	return err
}

func (r *Repository) UpdateStory(
	ctx context.Context,
	id string,
	slug string,
	status string,
	isFeatured bool,
	storyPictureURI *string,
	publishedAt *time.Time,
) error {
	params := UpdateStoryParams{
		ID:              id,
		Slug:            slug,
		Status:          status,
		IsFeatured:      isFeatured,
		StoryPictureURI: vars.ToSQLNullString(storyPictureURI),
		PublishedAt:     vars.ToSQLNullTime(publishedAt),
	}

	_, err := r.queries.UpdateStory(ctx, params)

	return err
}

func (r *Repository) UpdateStoryTx(
	ctx context.Context,
	storyID string,
	localeCode string,
	title string,
	summary string,
	content string,
) error {
	params := UpdateStoryTxParams{
		StoryID:    storyID,
		LocaleCode: localeCode,
		Title:      title,
		Summary:    summary,
		Content:    content,
	}

	_, err := r.queries.UpdateStoryTx(ctx, params)

	return err
}

func (r *Repository) UpsertStoryTx(
	ctx context.Context,
	storyID string,
	localeCode string,
	title string,
	summary string,
	content string,
) error {
	params := UpsertStoryTxParams{
		StoryID:    storyID,
		LocaleCode: localeCode,
		Title:      title,
		Summary:    summary,
		Content:    content,
	}

	return r.queries.UpsertStoryTx(ctx, params)
}

func (r *Repository) RemoveStory(ctx context.Context, id string) error {
	params := RemoveStoryParams{ID: id}
	_, err := r.queries.RemoveStory(ctx, params)

	return err
}

func (r *Repository) GetStoryForEdit(
	ctx context.Context,
	localeCode string,
	id string,
) (*stories.StoryForEdit, error) {
	params := GetStoryForEditParams{
		LocaleCode: localeCode,
		ID:         id,
	}

	row, err := r.queries.GetStoryForEdit(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &stories.StoryForEdit{
		ID:              row.ID,
		AuthorProfileID: vars.ToStringPtr(row.AuthorProfileID),
		Slug:            row.Slug,
		Kind:            row.Kind,
		Status:          row.Status,
		IsFeatured:      row.IsFeatured,
		StoryPictureURI: vars.ToStringPtr(row.StoryPictureURI),
		Title:           row.Title,
		Summary:         row.Summary,
		Content:         row.Content,
		CreatedAt:       row.CreatedAt,
		PublishedAt:     vars.ToTimePtr(row.PublishedAt),
		UpdatedAt:       vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) GetStoryOwnershipForUser(
	ctx context.Context,
	userID string,
	storyID string,
) (*stories.StoryOwnership, error) {
	params := GetStoryOwnershipForUserParams{
		UserID:  userID,
		StoryID: storyID,
	}

	row, err := r.queries.GetStoryOwnershipForUser(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &stories.StoryOwnership{
		ID:              row.ID,
		Slug:            row.Slug,
		AuthorProfileID: vars.ToStringPtr(row.AuthorProfileID),
		UserKind:        row.UserKind.String,
		CanEdit:         row.CanEdit,
	}, nil
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
			PublishedAt:     vars.ToTimePtr(story.PublishedAt),
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
			Points:            uint64(profile.Points),
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
			Points            uint64           `db:"points"              json:"points"`
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
			Points:            publicationProfile.Profile.Points,
		}
	}

	return storyWithChildren, nil
}
