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
	"github.com/sqlc-dev/pqtype"
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

func (r *Repository) GetStoryIDBySlugForViewer(
	ctx context.Context,
	slug string,
	viewerUserID *string,
) (string, error) {
	params := GetStoryIDBySlugForViewerParams{
		Slug: slug,
		ViewerUserID: sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	if viewerUserID != nil {
		params.ViewerUserID = sql.NullString{
			String: *viewerUserID,
			Valid:  true,
		}
	}

	row, err := r.queries.GetStoryIDBySlugForViewer(ctx, params)
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

	if t, ok := row.PublishedAt.(time.Time); ok {
		result.PublishedAt = &t
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

		if t, ok := row.PublishedAt.(time.Time); ok {
			storyWithChildren.PublishedAt = &t
		}

		result[i] = storyWithChildren
	}

	wrappedResponse.Data = result

	if len(result) == cursor.Limit {
		wrappedResponse.CursorPtr = &result[len(result)-1].ID
	}

	return wrappedResponse, nil
}

func (r *Repository) ListStoriesByAuthorProfileID(
	ctx context.Context,
	localeCode string,
	authorProfileID string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*stories.StoryWithChildren], error) {
	var wrappedResponse cursors.Cursored[[]*stories.StoryWithChildren]

	rows, err := r.queries.ListStoriesByAuthorProfileID(
		ctx,
		ListStoriesByAuthorProfileIDParams{
			LocaleCode:      localeCode,
			AuthorProfileID: authorProfileID,
			FilterKind:      vars.MapValueToNullString(cursor.Filters, "kind"),
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	result := make([]*stories.StoryWithChildren, len(rows))

	for i, row := range rows {
		storyWithChildren, err := r.parseStoryWithChildrenOptionalPublications(
			row.Profile,
			row.ProfileTx,
			row.Story,
			row.StoryTx,
			row.Publications,
		)
		if err != nil {
			return wrappedResponse, err
		}

		if t, ok := row.PublishedAt.(time.Time); ok {
			storyWithChildren.PublishedAt = &t
		}

		result[i] = storyWithChildren
	}

	wrappedResponse.Data = result

	if len(result) == cursor.Limit {
		wrappedResponse.CursorPtr = &result[len(result)-1].ID
	}

	return wrappedResponse, nil
}

func (r *Repository) ListActivityStories(
	ctx context.Context,
	localeCode string,
	filterAuthorProfileID *string,
) ([]*stories.StoryWithChildren, error) {
	params := ListActivityStoriesParams{
		LocaleCode: localeCode,
		FilterAuthorProfileID: sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	if filterAuthorProfileID != nil {
		params.FilterAuthorProfileID = sql.NullString{
			String: *filterAuthorProfileID,
			Valid:  true,
		}
	}

	rows, err := r.queries.ListActivityStories(ctx, params)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		if t, ok := row.PublishedAt.(time.Time); ok {
			storyWithChildren.PublishedAt = &t
		}

		result[i] = storyWithChildren
	}

	return result, nil
}

// Story CRUD methods

func (r *Repository) InsertStory(
	ctx context.Context,
	id string,
	authorProfileID string,
	slug string,
	kind string,
	storyPictureURI *string,
	properties map[string]any,
	isManaged bool,
	remoteID *string,
) (*stories.Story, error) {
	params := InsertStoryParams{
		ID:              id,
		AuthorProfileID: sql.NullString{String: authorProfileID, Valid: true},
		Slug:            slug,
		Kind:            kind,
		StoryPictureURI: vars.ToSQLNullString(storyPictureURI),
		Properties:      vars.ToSQLNullRawMessage(properties),
		IsManaged:       isManaged,
		RemoteID:        vars.ToSQLNullString(remoteID),
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
		IsManaged:       row.IsManaged,
		StoryPictureURI: vars.ToStringPtr(row.StoryPictureURI),
		SeriesID:        vars.ToStringPtr(row.SeriesID),
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
	isFeatured bool,
	publishedAt *time.Time,
	properties map[string]any,
) error {
	params := InsertStoryPublicationParams{
		ID:          id,
		StoryID:     storyID,
		ProfileID:   profileID,
		Kind:        kind,
		IsFeatured:  isFeatured,
		PublishedAt: vars.ToSQLNullTime(publishedAt),
		Properties:  vars.ToSQLNullRawMessage(properties),
	}

	_, err := r.queries.InsertStoryPublication(ctx, params)

	return err
}

func (r *Repository) UpdateStory(
	ctx context.Context,
	id string,
	slug string,
	storyPictureURI *string,
) error {
	params := UpdateStoryParams{
		ID:              id,
		Slug:            slug,
		StoryPictureURI: vars.ToSQLNullString(storyPictureURI),
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

func (r *Repository) DeleteStoryTx(
	ctx context.Context,
	storyID string,
	localeCode string,
) error {
	params := DeleteStoryTxParams{
		StoryID:    storyID,
		LocaleCode: localeCode,
	}
	_, err := r.queries.DeleteStoryTx(ctx, params)

	return err
}

func (r *Repository) ListStoryTxLocales(
	ctx context.Context,
	storyID string,
) ([]string, error) {
	params := ListStoryTxLocalesParams{
		StoryID: storyID,
	}

	return r.queries.ListStoryTxLocales(ctx, params)
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
		ID:                row.ID,
		AuthorProfileID:   vars.ToStringPtr(row.AuthorProfileID),
		AuthorProfileSlug: vars.ToStringPtr(row.AuthorProfileSlug),
		Slug:              row.Slug,
		Kind:              row.Kind,
		LocaleCode:        row.LocaleCode,
		StoryPictureURI:   vars.ToStringPtr(row.StoryPictureURI),
		Title:             row.Title,
		Summary:           row.Summary,
		Content:           row.Content,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         vars.ToTimePtr(row.UpdatedAt),
		IsManaged:         row.IsManaged,
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

// Publication management methods

func (r *Repository) ListStoryPublications(
	ctx context.Context,
	localeCode string,
	storyID string,
) ([]*stories.StoryPublication, error) {
	params := ListStoryPublicationsParams{
		LocaleCode: localeCode,
		StoryID:    storyID,
	}

	rows, err := r.queries.ListStoryPublications(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*stories.StoryPublication, len(rows))
	for i, row := range rows {
		result[i] = &stories.StoryPublication{
			ID:                row.ID,
			StoryID:           row.StoryID,
			ProfileID:         row.ProfileID,
			ProfileSlug:       row.ProfileSlug,
			ProfileTitle:      row.ProfileTitle,
			ProfilePictureURI: vars.ToStringPtr(row.ProfilePictureURI),
			ProfileKind:       row.ProfileKind,
			Kind:              row.Kind,
			IsFeatured:        row.IsFeatured,
			PublishedAt:       vars.ToTimePtr(row.PublishedAt),
			CreatedAt:         row.CreatedAt,
		}
	}

	return result, nil
}

func (r *Repository) GetStoryPublicationProfileID(
	ctx context.Context,
	publicationID string,
) (string, error) {
	params := GetStoryPublicationProfileIDParams{ID: publicationID}

	profileID, err := r.queries.GetStoryPublicationProfileID(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return profileID, nil
}

func (r *Repository) UpdateStoryPublication(
	ctx context.Context,
	id string,
	isFeatured bool,
) error {
	params := UpdateStoryPublicationParams{
		ID:         id,
		IsFeatured: isFeatured,
	}

	_, err := r.queries.UpdateStoryPublication(ctx, params)

	return err
}

func (r *Repository) UpdateStoryPublicationDate(
	ctx context.Context,
	id string,
	publishedAt time.Time,
) error {
	params := UpdateStoryPublicationDateParams{
		ID:          id,
		PublishedAt: sql.NullTime{Time: publishedAt, Valid: true},
	}

	_, err := r.queries.UpdateStoryPublicationDate(ctx, params)

	return err
}

func (r *Repository) RemoveStoryPublication(ctx context.Context, id string) error {
	params := RemoveStoryPublicationParams{ID: id}
	_, err := r.queries.RemoveStoryPublication(ctx, params)

	return err
}

func (r *Repository) CountStoryPublications(ctx context.Context, storyID string) (int64, error) {
	params := CountStoryPublicationsParams{StoryID: storyID}

	return r.queries.CountStoryPublications(ctx, params)
}

func (r *Repository) GetStoryFirstPublishedAt(
	ctx context.Context,
	storyID string,
) (*time.Time, error) {
	params := GetStoryFirstPublishedAtParams{StoryID: storyID}

	result, err := r.queries.GetStoryFirstPublishedAt(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	if result == nil {
		return nil, nil //nolint:nilnil
	}

	// The result is interface{} from sqlc; try to cast to time.Time
	if t, ok := result.(time.Time); ok {
		return &t, nil
	}

	return nil, nil //nolint:nilnil
}

func (r *Repository) GetUserMembershipForProfile(
	ctx context.Context,
	userID string,
	profileID string,
) (string, error) {
	params := GetUserMembershipForProfileParams{
		UserID:    userID,
		ProfileID: profileID,
	}

	kind, err := r.queries.GetUserMembershipForProfile(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return kind, nil
}

func (r *Repository) InvalidateStorySlugCache(ctx context.Context, slug string) error {
	return r.cache.Invalidate(ctx, "story_id_by_slug:"+slug)
}

func (r *Repository) parseStoryWithChildrenOptionalPublications(
	profile Profile,
	profileTx ProfileTx,
	story Story,
	storyTx StoryTx,
	publications pqtype.NullRawMessage,
) (*stories.StoryWithChildren, error) {
	if !publications.Valid || publications.RawMessage == nil ||
		string(publications.RawMessage) == "null" {
		return &stories.StoryWithChildren{
			Story: &stories.Story{
				ID:              story.ID,
				AuthorProfileID: vars.ToStringPtr(story.AuthorProfileID),
				Slug:            story.Slug,
				Kind:            story.Kind,
				StoryPictureURI: vars.ToStringPtr(story.StoryPictureURI),
				SeriesID:        vars.ToStringPtr(story.SeriesID),
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
			Publications: []*profiles.Profile{},
		}, nil
	}

	return r.parseStoryWithChildren(profile, profileTx, story, storyTx, publications)
}

func (r *Repository) parseStoryWithChildren( //nolint:funlen
	profile Profile,
	profileTx ProfileTx,
	story Story,
	storyTx StoryTx,
	publications pqtype.NullRawMessage,
) (*stories.StoryWithChildren, error) {
	storyWithChildren := &stories.StoryWithChildren{
		Story: &stories.Story{
			ID:              story.ID,
			AuthorProfileID: vars.ToStringPtr(story.AuthorProfileID),
			Slug:            story.Slug,
			Kind:            story.Kind,
			StoryPictureURI: vars.ToStringPtr(story.StoryPictureURI),
			SeriesID:        vars.ToStringPtr(story.SeriesID),
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

	if !publications.Valid || len(publications.RawMessage) == 0 {
		storyWithChildren.Publications = []*profiles.Profile{}

		return storyWithChildren, nil
	}

	err := json.Unmarshal(publications.RawMessage, &publicationProfiles)
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
