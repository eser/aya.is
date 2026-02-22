package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/eser/aya.is/services/pkg/lib/vars"
	"github.com/sqlc-dev/pqtype"
)

func (r *Repository) GetProfileIDBySlug(ctx context.Context, slug string) (string, error) {
	var result string

	err := r.cache.Execute(
		ctx,
		"profile_id_by_slug:"+slug,
		&result,
		func(ctx context.Context) (any, error) {
			row, err := r.queries.GetProfileIDBySlug(ctx, GetProfileIDBySlugParams{Slug: slug})
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

// GetFeatureRelationsVisibility returns the relations module visibility for a profile.
func (r *Repository) GetFeatureRelationsVisibility(
	ctx context.Context,
	profileID string,
) (string, error) {
	visibility, err := r.queries.GetProfileFeatureRelationsVisibility(
		ctx,
		GetProfileFeatureRelationsVisibilityParams{
			ID: profileID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "public", nil
		}

		return "public", err
	}

	return visibility, nil
}

// GetFeatureLinksVisibility returns the links module visibility for a profile.
func (r *Repository) GetFeatureLinksVisibility(
	ctx context.Context,
	profileID string,
) (string, error) {
	visibility, err := r.queries.GetProfileFeatureLinksVisibility(
		ctx,
		GetProfileFeatureLinksVisibilityParams{
			ID: profileID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "public", nil
		}

		return "public", err
	}

	return visibility, nil
}

func (r *Repository) CheckProfileSlugExists(ctx context.Context, slug string) (bool, error) {
	var result bool

	err := r.cache.Execute(
		ctx,
		"profile_slug_exists:"+slug,
		&result,
		func(ctx context.Context) (any, error) {
			exists, err := r.queries.CheckProfileSlugExists(
				ctx,
				CheckProfileSlugExistsParams{Slug: slug},
			)
			if err != nil {
				return nil, err
			}

			return exists, nil
		},
	)

	return result, err //nolint:wrapcheck
}

func (r *Repository) CheckProfileSlugExistsIncludingDeleted(
	ctx context.Context,
	slug string,
) (bool, error) {
	exists, err := r.queries.CheckProfileSlugExistsIncludingDeleted(
		ctx,
		CheckProfileSlugExistsIncludingDeletedParams{Slug: slug},
	)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) CheckPageSlugExistsIncludingDeleted(
	ctx context.Context,
	profileID string,
	pageSlug string,
) (bool, error) {
	exists, err := r.queries.CheckPageSlugExistsIncludingDeleted(
		ctx,
		CheckPageSlugExistsIncludingDeletedParams{
			ProfileID: profileID,
			PageSlug:  pageSlug,
		},
	)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) GetCustomDomainByDomain(
	ctx context.Context,
	domain string,
) (*profiles.ProfileCustomDomain, error) {
	var result *profiles.ProfileCustomDomain

	err := r.cache.Execute(
		ctx,
		"custom_domain_by_domain:"+domain,
		&result,
		func(ctx context.Context) (any, error) {
			row, err := r.queries.GetCustomDomainByDomain(
				ctx,
				GetCustomDomainByDomainParams{Domain: domain},
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, nil //nolint:nilnil
				}

				return nil, err
			}

			return &profiles.ProfileCustomDomain{
				ID:            row.ID,
				ProfileID:     row.ProfileID,
				Domain:        row.Domain,
				DefaultLocale: vars.ToStringPtr(row.DefaultLocale),
				CreatedAt:     row.CreatedAt,
				UpdatedAt:     vars.ToTimePtr(row.UpdatedAt),
			}, nil
		},
	)

	return result, err //nolint:wrapcheck
}

func (r *Repository) ListCustomDomainsByProfileID(
	ctx context.Context,
	profileID string,
) ([]*profiles.ProfileCustomDomain, error) {
	rows, err := r.queries.ListCustomDomainsByProfileID(
		ctx,
		ListCustomDomainsByProfileIDParams{ProfileID: profileID},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.ProfileCustomDomain, 0, len(rows))
	for _, row := range rows {
		result = append(result, &profiles.ProfileCustomDomain{
			ID:            row.ID,
			ProfileID:     row.ProfileID,
			Domain:        row.Domain,
			DefaultLocale: vars.ToStringPtr(row.DefaultLocale),
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     vars.ToTimePtr(row.UpdatedAt),
		})
	}

	return result, nil
}

func (r *Repository) CreateCustomDomain(
	ctx context.Context,
	id string,
	profileID string,
	domain string,
	defaultLocale *string,
) error {
	return r.queries.CreateCustomDomain(ctx, CreateCustomDomainParams{
		ID:            id,
		ProfileID:     profileID,
		Domain:        domain,
		DefaultLocale: vars.ToSQLNullString(defaultLocale),
	})
}

func (r *Repository) UpdateCustomDomain(
	ctx context.Context,
	id string,
	domain string,
	defaultLocale *string,
) error {
	_, err := r.queries.UpdateCustomDomain(ctx, UpdateCustomDomainParams{
		ID:            id,
		Domain:        domain,
		DefaultLocale: vars.ToSQLNullString(defaultLocale),
	})

	return err
}

func (r *Repository) DeleteCustomDomain(
	ctx context.Context,
	id string,
) error {
	_, err := r.queries.DeleteCustomDomain(ctx, DeleteCustomDomainParams{ID: id})

	return err
}

func (r *Repository) GetProfileIdentifierByID(
	ctx context.Context,
	id string,
) (*profiles.ProfileBrief, error) {
	row, err := r.queries.GetProfileIdentifierByID(ctx, GetProfileIdentifierByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &profiles.ProfileBrief{
		ID:                row.ID,
		Slug:              row.Slug,
		Kind:              row.Kind,
		ProfilePictureURI: vars.ToStringPtr(row.ProfilePictureURI),
	}

	return result, nil
}

func (r *Repository) GetProfileByID(
	ctx context.Context,
	localeCode string,
	id string,
) (*profiles.Profile, error) {
	row, err := r.queries.GetProfileByID(
		ctx,
		GetProfileByIDParams{
			LocaleCode: localeCode,
			ID:         id,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &profiles.Profile{
		ID:   row.Profile.ID,
		Slug: row.Profile.Slug,
		Kind: row.Profile.Kind,

		ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
		Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
		Title:             row.ProfileTx.Title,
		Description:       row.ProfileTx.Description,
		DefaultLocale:     row.Profile.DefaultLocale,
		Properties:        vars.ToObject(row.Profile.Properties),
		CreatedAt:         row.Profile.CreatedAt,
		UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
		DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
		Points:            uint64(row.Profile.Points),
		FeatureRelations:  row.Profile.FeatureRelations,
		FeatureLinks:      row.Profile.FeatureLinks,
		FeatureQA:         row.Profile.FeatureQa,
	}

	return result, nil
}

func (r *Repository) ListProfiles(
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profiles.Profile], error) {
	var wrappedResponse cursors.Cursored[[]*profiles.Profile]

	rows, err := r.queries.ListProfiles(
		ctx,
		ListProfilesParams{
			LocaleCode: localeCode,
			FilterKind: vars.MapValueToNullString(cursor.Filters, "kind"),
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	result := make([]*profiles.Profile, len(rows))
	for i, row := range rows {
		result[i] = &profiles.Profile{
			ID:   row.Profile.ID,
			Slug: row.Profile.Slug,
			Kind: row.Profile.Kind,

			ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
			Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
			Title:             row.ProfileTx.Title,
			Description:       row.ProfileTx.Description,
			DefaultLocale:     row.Profile.DefaultLocale,
			Properties:        vars.ToObject(row.Profile.Properties),
			CreatedAt:         row.Profile.CreatedAt,
			UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
			DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
			FeatureRelations:  row.Profile.FeatureRelations,
			FeatureLinks:      row.Profile.FeatureLinks,
			FeatureQA:         row.Profile.FeatureQa,
		}
	}

	wrappedResponse.Data = result

	if len(result) == cursor.Limit && len(result) > 0 {
		wrappedResponse.CursorPtr = &result[len(result)-1].ID
	}

	return wrappedResponse, nil
}

func (r *Repository) ListProfilePagesByProfileID(
	ctx context.Context,
	localeCode string,
	profileID string,
) ([]*profiles.ProfilePageBrief, error) {
	rows, err := r.queries.ListProfilePagesByProfileID(
		ctx,
		ListProfilePagesByProfileIDParams{
			LocaleCode: localeCode,
			ProfileID:  profileID,
		},
	)
	if err != nil {
		return nil, err
	}

	profilePages := make([]*profiles.ProfilePageBrief, len(rows))
	for i, row := range rows {
		profilePages[i] = &profiles.ProfilePageBrief{
			ID:              row.ID,
			Slug:            row.Slug,
			CoverPictureURI: vars.ToStringPtr(row.CoverPictureURI),
			Title:           row.Title,
			Summary:         row.Summary,
			Visibility:      profiles.PageVisibility(row.Visibility),
		}
	}

	return profilePages, nil
}

func (r *Repository) GetProfilePageByProfileIDAndSlug(
	ctx context.Context,
	localeCode string,
	profileID string,
	pageSlug string,
) (*profiles.ProfilePage, error) {
	row, err := r.queries.GetProfilePageByProfileIDAndSlug(
		ctx,
		GetProfilePageByProfileIDAndSlugParams{
			LocaleCode: localeCode,
			ProfileID:  profileID,
			PageSlug:   pageSlug,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &profiles.ProfilePage{
		ID:               row.ID,
		Slug:             row.Slug,
		LocaleCode:       row.LocaleCode,
		CoverPictureURI:  vars.ToStringPtr(row.CoverPictureURI),
		Title:            row.Title,
		Summary:          row.Summary,
		Content:          row.Content,
		SortOrder:        row.Order,
		Visibility:       profiles.PageVisibility(row.Visibility),
		PublishedAt:      vars.ToTimePtr(row.PublishedAt),
		AddedByProfileID: vars.ToStringPtr(row.AddedByProfileID),
	}

	return result, nil
}

func (r *Repository) ListProfilePagesByProfileIDForViewer(
	ctx context.Context,
	localeCode string,
	profileID string,
	viewerUserID *string,
) ([]*profiles.ProfilePageBrief, error) {
	params := ListProfilePagesByProfileIDForViewerParams{
		LocaleCode: localeCode,
		ProfileID:  profileID,
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

	rows, err := r.queries.ListProfilePagesByProfileIDForViewer(ctx, params)
	if err != nil {
		return nil, err
	}

	profilePages := make([]*profiles.ProfilePageBrief, len(rows))
	for i, row := range rows {
		profilePages[i] = &profiles.ProfilePageBrief{
			ID:              row.ID,
			Slug:            row.Slug,
			CoverPictureURI: vars.ToStringPtr(row.CoverPictureURI),
			Title:           row.Title,
			Summary:         row.Summary,
			Visibility:      profiles.PageVisibility(row.Visibility),
		}
	}

	return profilePages, nil
}

func (r *Repository) GetProfilePageByProfileIDAndSlugForViewer(
	ctx context.Context,
	localeCode string,
	profileID string,
	pageSlug string,
	viewerUserID *string,
) (*profiles.ProfilePage, error) {
	params := GetProfilePageByProfileIDAndSlugForViewerParams{
		LocaleCode: localeCode,
		ProfileID:  profileID,
		PageSlug:   pageSlug,
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

	row, err := r.queries.GetProfilePageByProfileIDAndSlugForViewer(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &profiles.ProfilePage{
		ID:               row.ID,
		Slug:             row.Slug,
		LocaleCode:       row.LocaleCode,
		CoverPictureURI:  vars.ToStringPtr(row.CoverPictureURI),
		Title:            row.Title,
		Summary:          row.Summary,
		Content:          row.Content,
		SortOrder:        row.Order,
		Visibility:       profiles.PageVisibility(row.Visibility),
		PublishedAt:      vars.ToTimePtr(row.PublishedAt),
		AddedByProfileID: vars.ToStringPtr(row.AddedByProfileID),
	}

	return result, nil
}

func (r *Repository) ListProfileLinksByProfileID(
	ctx context.Context,
	localeCode string,
	profileID string,
) ([]*profiles.ProfileLinkBrief, error) {
	rows, err := r.queries.ListProfileLinksByProfileID(
		ctx,
		ListProfileLinksByProfileIDParams{
			LocaleCode: localeCode,
			ProfileID:  profileID,
		},
	)
	if err != nil {
		return nil, err
	}

	profileLinks := make([]*profiles.ProfileLinkBrief, len(rows))
	for i, row := range rows {
		profileLinks[i] = &profiles.ProfileLinkBrief{
			ID:          row.ID,
			Kind:        row.Kind,
			Order:       int(row.Order),
			IsManaged:   row.IsManaged,
			IsVerified:  row.IsVerified,
			IsFeatured:  row.IsFeatured,
			Visibility:  profiles.LinkVisibility(row.Visibility),
			PublicID:    row.PublicID.String,
			URI:         row.URI.String,
			Title:       row.Title,
			Icon:        row.Icon,
			Group:       row.Group.String,
			Description: row.Description.String,
		}
	}

	return profileLinks, nil
}

func (r *Repository) ListProfileContributions( //nolint:funlen
	ctx context.Context,
	localeCode string,
	profileID string,
	kinds []string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profiles.ProfileMembership], error) {
	var wrappedResponse cursors.Cursored[[]*profiles.ProfileMembership]

	rows, err := r.queries.ListProfileMemberships(
		ctx,
		ListProfileMembershipsParams{
			LocaleCode:      localeCode,
			FilterProfileID: sql.NullString{String: "", Valid: false},
			FilterProfileKind: sql.NullString{
				String: strings.Join(kinds, ","),
				Valid:  true,
			},
			FilterMemberProfileID:       sql.NullString{String: profileID, Valid: true},
			FilterMemberProfileKind:     sql.NullString{String: "", Valid: false},
			FilterMembershipKindExclude: sql.NullString{String: "follower", Valid: true},
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	profileMemberships := make([]*profiles.ProfileMembership, len(rows))
	for i, row := range rows { //nolint:dupl
		profileMemberships[i] = &profiles.ProfileMembership{
			ID:         row.ProfileMembership.ID,
			Kind:       row.ProfileMembership.Kind,
			StartedAt:  vars.ToTimePtr(row.ProfileMembership.StartedAt),
			FinishedAt: vars.ToTimePtr(row.ProfileMembership.FinishedAt),
			Properties: vars.ToObject(row.ProfileMembership.Properties),
			Profile: &profiles.Profile{
				ID:   row.Profile.ID,
				Slug: row.Profile.Slug,
				Kind: row.Profile.Kind,

				ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
				Title:             row.ProfileTx.Title,
				Description:       row.ProfileTx.Description,
				DefaultLocale:     row.Profile.DefaultLocale,
				Properties:        vars.ToObject(row.Profile.Properties),
				CreatedAt:         row.Profile.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
			},
			MemberProfile: &profiles.Profile{
				ID:   row.Profile_2.ID,
				Slug: row.Profile_2.Slug,
				Kind: row.Profile_2.Kind,

				ProfilePictureURI: vars.ToStringPtr(row.Profile_2.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile_2.Pronouns),
				Title:             row.ProfileTx_2.Title,
				Description:       row.ProfileTx_2.Description,
				DefaultLocale:     row.Profile_2.DefaultLocale,
				Properties:        vars.ToObject(row.Profile_2.Properties),
				CreatedAt:         row.Profile_2.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile_2.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile_2.DeletedAt),
				Points:            uint64(row.Profile_2.Points),
			},
		}
	}

	wrappedResponse.Data = profileMemberships

	if len(profileMemberships) == cursor.Limit {
		wrappedResponse.CursorPtr = &profileMemberships[len(profileMemberships)-1].ID
	}

	return wrappedResponse, nil
}

//nolint:funlen,dupl
func (r *Repository) ListProfileMembers(
	ctx context.Context,
	localeCode string,
	profileID string,
	kinds []string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profiles.ProfileMembership], error) {
	var wrappedResponse cursors.Cursored[[]*profiles.ProfileMembership]

	rows, err := r.queries.ListProfileMemberships(
		ctx,
		ListProfileMembershipsParams{
			LocaleCode:            localeCode,
			FilterProfileID:       sql.NullString{String: profileID, Valid: true},
			FilterProfileKind:     sql.NullString{String: "", Valid: false},
			FilterMemberProfileID: sql.NullString{String: "", Valid: false},
			FilterMemberProfileKind: sql.NullString{
				String: strings.Join(kinds, ","),
				Valid:  true,
			},
			FilterMembershipKindExclude: sql.NullString{String: "follower", Valid: true},
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	profileMemberships := make([]*profiles.ProfileMembership, len(rows))
	for i, row := range rows {
		profileMemberships[i] = &profiles.ProfileMembership{
			ID:         row.ProfileMembership.ID,
			Kind:       row.ProfileMembership.Kind,
			StartedAt:  vars.ToTimePtr(row.ProfileMembership.StartedAt),
			FinishedAt: vars.ToTimePtr(row.ProfileMembership.FinishedAt),
			Properties: vars.ToObject(row.ProfileMembership.Properties),
			Profile: &profiles.Profile{
				ID:   row.Profile.ID,
				Slug: row.Profile.Slug,
				Kind: row.Profile.Kind,

				ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
				Title:             row.ProfileTx.Title,
				Description:       row.ProfileTx.Description,
				DefaultLocale:     row.Profile.DefaultLocale,
				Properties:        vars.ToObject(row.Profile.Properties),
				CreatedAt:         row.Profile.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
			},
			MemberProfile: &profiles.Profile{
				ID:   row.Profile_2.ID,
				Slug: row.Profile_2.Slug,
				Kind: row.Profile_2.Kind,

				ProfilePictureURI: vars.ToStringPtr(row.Profile_2.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile_2.Pronouns),
				Title:             row.ProfileTx_2.Title,
				Description:       row.ProfileTx_2.Description,
				DefaultLocale:     row.Profile_2.DefaultLocale,
				Properties:        vars.ToObject(row.Profile_2.Properties),
				CreatedAt:         row.Profile_2.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile_2.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile_2.DeletedAt),
				Points:            uint64(row.Profile_2.Points),
			},
		}
	}

	wrappedResponse.Data = profileMemberships

	if len(profileMemberships) == cursor.Limit {
		wrappedResponse.CursorPtr = &profileMemberships[len(profileMemberships)-1].ID
	}

	return wrappedResponse, nil
}

func (r *Repository) GetProfileMembershipsByMemberProfileID(
	ctx context.Context,
	localeCode string,
	memberProfileID string,
) ([]*profiles.ProfileMembership, error) {
	rows, err := r.queries.GetProfileMembershipsByMemberProfileID(
		ctx,
		GetProfileMembershipsByMemberProfileIDParams{
			LocaleCode:      localeCode,
			MemberProfileID: sql.NullString{String: memberProfileID, Valid: true},
		},
	)
	if err != nil {
		return nil, err
	}

	memberships := make([]*profiles.ProfileMembership, len(rows))
	for i, row := range rows {
		memberships[i] = &profiles.ProfileMembership{
			ID:         row.MembershipID,
			Kind:       row.MembershipKind,
			StartedAt:  vars.ToTimePtr(row.StartedAt),
			FinishedAt: vars.ToTimePtr(row.FinishedAt),
			Properties: vars.ToObject(row.MembershipProperties),
			Profile: &profiles.Profile{
				ID:   row.Profile.ID,
				Slug: row.Profile.Slug,
				Kind: row.Profile.Kind,

				ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
				Title:             row.ProfileTx.Title,
				Description:       row.ProfileTx.Description,
				DefaultLocale:     row.Profile.DefaultLocale,
				Properties:        vars.ToObject(row.Profile.Properties),
				CreatedAt:         row.Profile.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
			},
			// MemberProfile is not needed for this use case since we're filtering by member profile ID
			MemberProfile: nil,
		}
	}

	return memberships, nil
}

func (r *Repository) CreateProfile(
	ctx context.Context,
	id string,
	slug string,
	kind string,
	defaultLocale string,
	profilePictureURI *string,
	pronouns *string,
	properties map[string]any,
) error {
	params := CreateProfileParams{
		ID:                id,
		Slug:              slug,
		Kind:              kind,
		DefaultLocale:     defaultLocale,
		ProfilePictureURI: vars.ToSQLNullString(profilePictureURI),
		Pronouns:          vars.ToSQLNullString(pronouns),
		Properties:        vars.ToSQLNullRawMessage(properties),
	}

	return r.queries.CreateProfile(ctx, params)
}

func (r *Repository) CreateProfileTx(
	ctx context.Context,
	profileID string,
	localeCode string,
	title string,
	description string,
	properties map[string]any,
) error {
	params := CreateProfileTxParams{
		ProfileID:   profileID,
		LocaleCode:  localeCode,
		Title:       title,
		Description: description,
		Properties:  vars.ToSQLNullRawMessage(properties),
	}

	return r.queries.CreateProfileTx(ctx, params)
}

func (r *Repository) UpdateProfile(
	ctx context.Context,
	id string,
	profilePictureURI *string,
	pronouns *string,
	properties map[string]any,
	featureRelations *string,
	featureLinks *string,
	featureQA *string,
) error {
	params := UpdateProfileParams{
		ID:                id,
		ProfilePictureURI: vars.ToSQLNullString(profilePictureURI),
		Pronouns:          vars.ToSQLNullString(pronouns),
		Properties:        vars.ToSQLNullRawMessage(properties),
		FeatureRelations:  vars.ToSQLNullString(featureRelations),
		FeatureLinks:      vars.ToSQLNullString(featureLinks),
		FeatureQa:         vars.ToSQLNullString(featureQA),
	}

	_, err := r.queries.UpdateProfile(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateProfileTx(
	ctx context.Context,
	profileID string,
	localeCode string,
	title string,
	description string,
	properties map[string]any,
) error {
	params := UpdateProfileTxParams{
		ProfileID:   profileID,
		LocaleCode:  localeCode,
		Title:       title,
		Description: description,
		Properties:  vars.ToSQLNullRawMessage(properties),
	}

	_, err := r.queries.UpdateProfileTx(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpsertProfileTx(
	ctx context.Context,
	profileID string,
	localeCode string,
	title string,
	description string,
	properties map[string]any,
) error {
	params := UpsertProfileTxParams{
		ProfileID:   profileID,
		LocaleCode:  localeCode,
		Title:       title,
		Description: description,
		Properties:  vars.ToSQLNullRawMessage(properties),
	}

	err := r.queries.UpsertProfileTx(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateProfileMembership(
	ctx context.Context,
	id string,
	profileID string,
	memberProfileID *string,
	kind string,
	properties map[string]any,
) error {
	params := CreateProfileMembershipParams{
		ID:              id,
		ProfileID:       profileID,
		MemberProfileID: vars.ToSQLNullString(memberProfileID),
		Kind:            kind,
		Properties:      vars.ToSQLNullRawMessage(properties),
	}

	err := r.queries.CreateProfileMembership(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListProfileMembershipsForSettings(
	ctx context.Context,
	localeCode string,
	profileID string,
) ([]*profiles.ProfileMembershipWithMember, error) {
	rows, err := r.queries.ListProfileMembershipsForSettings(
		ctx,
		ListProfileMembershipsForSettingsParams{
			LocaleCode: localeCode,
			ProfileID:  profileID,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.ProfileMembershipWithMember, 0, len(rows))

	for _, row := range rows {
		membership := &profiles.ProfileMembershipWithMember{
			ID:              row.ID,
			ProfileID:       row.ProfileID,
			MemberProfileID: vars.ToStringPtr(row.MemberProfileID),
			Kind:            row.Kind,
			Properties:      vars.ToObject(row.Properties),
			StartedAt:       vars.ToTimePtr(row.StartedAt),
			FinishedAt:      vars.ToTimePtr(row.FinishedAt),
			MemberProfile: &profiles.ProfileBrief{
				ID:                row.Profile.ID,
				Slug:              row.Profile.Slug,
				Kind:              row.Profile.Kind,
				ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
				Title:             row.ProfileTx.Title,
				Description:       row.ProfileTx.Description,
			},
		}
		result = append(result, membership)
	}

	return result, nil
}

func (r *Repository) GetProfileMembershipByID(
	ctx context.Context,
	id string,
) (*profiles.ProfileMembership, error) {
	row, err := r.queries.GetProfileMembershipByID(ctx, GetProfileMembershipByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &profiles.ProfileMembership{
		ID:              row.ID,
		ProfileID:       row.ProfileID,
		MemberProfileID: vars.ToStringPtr(row.MemberProfileID),
		Kind:            row.Kind,
		Properties:      vars.ToObject(row.Properties),
		StartedAt:       vars.ToTimePtr(row.StartedAt),
		FinishedAt:      vars.ToTimePtr(row.FinishedAt),
	}, nil
}

func (r *Repository) GetProfileMembershipByProfileAndMember(
	ctx context.Context,
	profileID string,
	memberProfileID string,
) (*profiles.ProfileMembership, error) {
	row, err := r.queries.GetProfileMembershipByProfileAndMember(
		ctx,
		GetProfileMembershipByProfileAndMemberParams{
			ProfileID:       profileID,
			MemberProfileID: sql.NullString{String: memberProfileID, Valid: true},
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &profiles.ProfileMembership{
		ID:              row.ID,
		ProfileID:       row.ProfileID,
		MemberProfileID: vars.ToStringPtr(row.MemberProfileID),
		Kind:            row.Kind,
		Properties:      vars.ToObject(row.Properties),
		StartedAt:       vars.ToTimePtr(row.StartedAt),
		FinishedAt:      vars.ToTimePtr(row.FinishedAt),
	}, nil
}

func (r *Repository) UpdateProfileMembership(
	ctx context.Context,
	id string,
	kind string,
) error {
	_, err := r.queries.UpdateProfileMembership(ctx, UpdateProfileMembershipParams{
		ID:   id,
		Kind: kind,
	})

	return err
}

func (r *Repository) DeleteProfileMembership(
	ctx context.Context,
	id string,
) error {
	_, err := r.queries.DeleteProfileMembership(ctx, DeleteProfileMembershipParams{ID: id})

	return err
}

func (r *Repository) CountProfileOwners(
	ctx context.Context,
	profileID string,
) (int64, error) {
	return r.queries.CountProfileOwners(ctx, CountProfileOwnersParams{ProfileID: profileID})
}

func (r *Repository) SearchUsersForMembership(
	ctx context.Context,
	localeCode string,
	profileID string,
	query string,
) ([]*profiles.UserSearchResult, error) {
	rows, err := r.queries.SearchUsersForMembership(ctx, SearchUsersForMembershipParams{
		LocaleCode: localeCode,
		ProfileID:  profileID,
		Query:      sql.NullString{String: query, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.UserSearchResult, 0, len(rows))

	for _, row := range rows {
		name := row.Name
		user := &profiles.UserSearchResult{
			UserID:              row.UserID,
			Email:               row.Email.String,
			Name:                &name,
			IndividualProfileID: vars.ToStringPtr(row.IndividualProfileID),
			Profile: &profiles.ProfileBrief{
				ID:                row.Profile.ID,
				Slug:              row.Profile.Slug,
				Kind:              row.Profile.Kind,
				ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
				Title:             row.ProfileTx.Title,
				Description:       row.ProfileTx.Description,
			},
		}
		result = append(result, user)
	}

	return result, nil
}

func (r *Repository) GetProfileOwnershipForUser(
	ctx context.Context,
	userID string,
	profileSlug string,
) (*profiles.ProfileOwnership, error) {
	row, err := r.queries.GetProfileOwnershipForUser(
		ctx,
		GetProfileOwnershipForUserParams{
			UserID:      userID,
			ProfileSlug: profileSlug,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &profiles.ProfileOwnership{
		ProfileID:   row.ID,
		ProfileSlug: row.Slug,
		ProfileKind: row.ProfileKind,
		UserKind:    row.UserKind,
		CanEdit:     row.CanEdit,
	}

	return result, nil
}

func (r *Repository) GetUserBriefInfo(
	ctx context.Context,
	userID string,
) (*profiles.UserBriefInfo, error) {
	var result profiles.UserBriefInfo

	err := r.cache.Execute(
		ctx,
		"user_brief_info:"+userID,
		&result,
		func(ctx context.Context) (any, error) {
			row, err := r.queries.GetUserBriefInfoByID(
				ctx,
				GetUserBriefInfoByIDParams{UserID: userID},
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, nil //nolint:nilnil
				}

				return nil, err
			}

			var individualProfileID *string
			if row.IndividualProfileID.Valid {
				individualProfileID = &row.IndividualProfileID.String
			}

			return profiles.UserBriefInfo{
				Kind:                row.Kind,
				IndividualProfileID: individualProfileID,
			}, nil
		},
	)

	return &result, err //nolint:wrapcheck
}

func (r *Repository) GetUserProfilePermissions(
	ctx context.Context,
	userID string,
) ([]*profiles.ProfilePermission, error) {
	rows, err := r.queries.GetUserProfilePermissions(
		ctx,
		GetUserProfilePermissionsParams{UserID: userID},
	)
	if err != nil {
		return nil, err
	}

	permissions := make([]*profiles.ProfilePermission, len(rows))

	for i, row := range rows {
		userKind := ""
		if row.UserKind.Valid {
			userKind = row.UserKind.String
		}

		permissions[i] = &profiles.ProfilePermission{
			ProfileID:      row.ID,
			ProfileSlug:    row.Slug,
			ProfileKind:    row.ProfileKind,
			MembershipKind: row.MembershipKind,
			UserKind:       userKind,
		}
	}

	return permissions, nil
}

func (r *Repository) GetProfileTxByID(
	ctx context.Context,
	profileID string,
) ([]*profiles.ProfileTx, error) {
	rows, err := r.queries.GetProfileTxByID(ctx, GetProfileTxByIDParams{ID: profileID})
	if err != nil {
		return nil, err
	}

	translations := make([]*profiles.ProfileTx, len(rows))
	for i, row := range rows {
		translations[i] = &profiles.ProfileTx{
			ProfileID:   row.ProfileTx.ProfileID,
			LocaleCode:  row.ProfileTx.LocaleCode,
			Title:       row.ProfileTx.Title,
			Description: row.ProfileTx.Description,
			Properties:  vars.ToObject(row.ProfileTx.Properties),
		}
	}

	return translations, nil
}

func (r *Repository) GetProfileLink(
	ctx context.Context,
	localeCode string,
	id string,
) (*profiles.ProfileLink, error) {
	row, err := r.queries.GetProfileLink(ctx, GetProfileLinkParams{
		LocaleCode: localeCode,
		ID:         id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	// Convert Icon string to *string (empty string -> nil)
	var iconPtr *string
	if row.Icon != "" {
		iconPtr = &row.Icon
	}

	result := &profiles.ProfileLink{
		ID:               row.ID,
		Kind:             row.Kind,
		ProfileID:        row.ProfileID,
		Order:            int(row.Order),
		IsManaged:        row.IsManaged,
		IsVerified:       row.IsVerified,
		IsFeatured:       row.IsFeatured,
		Visibility:       profiles.LinkVisibility(row.Visibility),
		RemoteID:         vars.ToStringPtr(row.RemoteID),
		PublicID:         vars.ToStringPtr(row.PublicID),
		URI:              vars.ToStringPtr(row.URI),
		Title:            row.Title,
		Icon:             iconPtr,
		Group:            vars.ToStringPtr(row.Group),
		Description:      vars.ToStringPtr(row.Description),
		AddedByProfileID: vars.ToStringPtr(row.AddedByProfileID),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:        vars.ToTimePtr(row.DeletedAt),
	}

	return result, nil
}

func (r *Repository) CreateProfileLink(
	ctx context.Context,
	id string,
	kind string,
	profileID string,
	order int,
	uri *string,
	isFeatured bool,
	visibility profiles.LinkVisibility,
	addedByProfileID *string,
) (*profiles.ProfileLink, error) {
	row, err := r.queries.CreateProfileLink(ctx, CreateProfileLinkParams{
		ID:                        id,
		Kind:                      kind,
		ProfileID:                 profileID,
		LinkOrder:                 int32(order),
		IsManaged:                 false, // For manually added links
		IsVerified:                false, // Will be verified later if needed
		IsFeatured:                isFeatured,
		Visibility:                string(visibility),
		RemoteID:                  sql.NullString{Valid: false},
		PublicID:                  sql.NullString{Valid: false},
		URI:                       vars.ToSQLNullString(uri),
		AuthProvider:              sql.NullString{Valid: false},
		AuthAccessTokenScope:      sql.NullString{Valid: false},
		AuthAccessToken:           sql.NullString{Valid: false},
		AuthAccessTokenExpiresAt:  sql.NullTime{Valid: false},
		AuthRefreshToken:          sql.NullString{Valid: false},
		AuthRefreshTokenExpiresAt: sql.NullTime{Valid: false},
		AddedByProfileID:          vars.ToSQLNullString(addedByProfileID),
	})
	if err != nil {
		return nil, err
	}

	result := &profiles.ProfileLink{
		ID:               row.ID,
		Kind:             row.Kind,
		ProfileID:        row.ProfileID,
		Order:            int(row.Order),
		IsManaged:        row.IsManaged,
		IsVerified:       row.IsVerified,
		IsFeatured:       row.IsFeatured,
		Visibility:       profiles.LinkVisibility(row.Visibility),
		RemoteID:         vars.ToStringPtr(row.RemoteID),
		PublicID:         vars.ToStringPtr(row.PublicID),
		URI:              vars.ToStringPtr(row.URI),
		AddedByProfileID: vars.ToStringPtr(row.AddedByProfileID),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:        vars.ToTimePtr(row.DeletedAt),
	}

	return result, nil
}

func (r *Repository) UpdateProfileLink(
	ctx context.Context,
	id string,
	kind string,
	order int,
	uri *string,
	isFeatured bool,
	visibility profiles.LinkVisibility,
) error {
	params := UpdateProfileLinkParams{
		ID:         id,
		Kind:       kind,
		LinkOrder:  int32(order),
		URI:        vars.ToSQLNullString(uri),
		IsFeatured: isFeatured,
		Visibility: string(visibility),
	}

	_, err := r.queries.UpdateProfileLink(ctx, params)

	return err
}

func (r *Repository) DeleteProfileLink(
	ctx context.Context,
	id string,
) error {
	_, err := r.queries.DeleteProfileLink(ctx, DeleteProfileLinkParams{ID: id})

	return err
}

func (r *Repository) GetProfilePage(
	ctx context.Context,
	id string,
) (*profiles.ProfilePage, error) {
	row, err := r.queries.GetProfilePage(ctx, GetProfilePageParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &profiles.ProfilePage{
		ID:               row.ID,
		Slug:             row.Slug,
		CoverPictureURI:  vars.ToStringPtr(row.CoverPictureURI),
		Visibility:       profiles.PageVisibility(row.Visibility),
		PublishedAt:      vars.ToTimePtr(row.PublishedAt),
		AddedByProfileID: vars.ToStringPtr(row.AddedByProfileID),
		// Note: Title, Summary, Content need to be fetched from profile_page_tx table
	}

	return result, nil
}

func (r *Repository) CreateProfilePage(
	ctx context.Context,
	id string,
	slug string,
	profileID string,
	order int,
	coverPictureURI *string,
	publishedAt *string,
	addedByProfileID *string,
	visibility string,
) (*profiles.ProfilePage, error) {
	var publishedAtTime sql.NullTime
	if publishedAt != nil {
		// Convert string to time if needed
		publishedAtTime = sql.NullTime{Valid: false}
	}

	row, err := r.queries.CreateProfilePage(ctx, CreateProfilePageParams{
		ID:               id,
		Slug:             slug,
		ProfileID:        profileID,
		PageOrder:        int32(order),
		CoverPictureURI:  vars.ToSQLNullString(coverPictureURI),
		PublishedAt:      publishedAtTime,
		AddedByProfileID: vars.ToSQLNullString(addedByProfileID),
		Visibility:       visibility,
	})
	if err != nil {
		return nil, err
	}

	result := &profiles.ProfilePage{
		ID:               row.ID,
		Slug:             row.Slug,
		CoverPictureURI:  vars.ToStringPtr(row.CoverPictureURI),
		Visibility:       profiles.PageVisibility(row.Visibility),
		PublishedAt:      vars.ToTimePtr(row.PublishedAt),
		AddedByProfileID: vars.ToStringPtr(row.AddedByProfileID),
		// Note: Title, Summary, Content need to be fetched from profile_page_tx table
	}

	return result, nil
}

func (r *Repository) CreateProfilePageTx(
	ctx context.Context,
	profilePageID string,
	localeCode string,
	title string,
	summary string,
	content string,
) error {
	params := CreateProfilePageTxParams{
		ProfilePageID: profilePageID,
		LocaleCode:    localeCode,
		Title:         title,
		Summary:       summary,
		Content:       content,
	}

	return r.queries.CreateProfilePageTx(ctx, params)
}

func (r *Repository) UpdateProfilePage(
	ctx context.Context,
	id string,
	slug string,
	order int,
	coverPictureURI *string,
	publishedAt *string,
	visibility string,
) error {
	var publishedAtTime sql.NullTime
	if publishedAt != nil {
		// Convert string to time if needed
		publishedAtTime = sql.NullTime{Valid: false}
	}

	params := UpdateProfilePageParams{
		ID:              id,
		Slug:            slug,
		PageOrder:       int32(order),
		CoverPictureURI: vars.ToSQLNullString(coverPictureURI),
		PublishedAt:     publishedAtTime,
		Visibility:      visibility,
	}

	_, err := r.queries.UpdateProfilePage(ctx, params)

	return err
}

func (r *Repository) UpdateProfilePageTx(
	ctx context.Context,
	profilePageID string,
	localeCode string,
	title string,
	summary string,
	content string,
) error {
	params := UpdateProfilePageTxParams{
		ProfilePageID: profilePageID,
		LocaleCode:    localeCode,
		Title:         title,
		Summary:       summary,
		Content:       content,
	}

	_, err := r.queries.UpdateProfilePageTx(ctx, params)

	return err
}

func (r *Repository) UpsertProfilePageTx(
	ctx context.Context,
	profilePageID string,
	localeCode string,
	title string,
	summary string,
	content string,
) error {
	params := UpsertProfilePageTxParams{
		ProfilePageID: profilePageID,
		LocaleCode:    localeCode,
		Title:         title,
		Summary:       summary,
		Content:       content,
	}

	return r.queries.UpsertProfilePageTx(ctx, params)
}

func (r *Repository) DeleteProfilePageTx(
	ctx context.Context,
	profilePageID string,
	localeCode string,
) error {
	params := DeleteProfilePageTxParams{
		ProfilePageID: profilePageID,
		LocaleCode:    localeCode,
	}
	_, err := r.queries.DeleteProfilePageTx(ctx, params)

	return err
}

func (r *Repository) ListProfilePageTxLocales(
	ctx context.Context,
	profilePageID string,
) ([]string, error) {
	params := ListProfilePageTxLocalesParams{
		ProfilePageID: profilePageID,
	}

	return r.queries.ListProfilePageTxLocales(ctx, params)
}

func (r *Repository) DeleteProfilePage(
	ctx context.Context,
	id string,
) error {
	_, err := r.queries.DeleteProfilePage(ctx, DeleteProfilePageParams{ID: id})

	return err
}

// OAuth Profile Link methods

func (r *Repository) GetProfileLinkByRemoteID(
	ctx context.Context,
	profileID string,
	kind string,
	remoteID string,
) (*profiles.ProfileLink, error) {
	row, err := r.queries.GetProfileLinkByRemoteID(ctx, GetProfileLinkByRemoteIDParams{
		ProfileID: profileID,
		Kind:      kind,
		RemoteID:  sql.NullString{String: remoteID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &profiles.ProfileLink{
		ID:         row.ID,
		Kind:       row.Kind,
		ProfileID:  row.ProfileID,
		Order:      int(row.Order),
		IsManaged:  row.IsManaged,
		IsVerified: row.IsVerified,
		IsFeatured: row.IsFeatured,
		Visibility: profiles.LinkVisibility(row.Visibility),
		RemoteID:   vars.ToStringPtr(row.RemoteID),
		PublicID:   vars.ToStringPtr(row.PublicID),
		URI:        vars.ToStringPtr(row.URI),
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:  vars.ToTimePtr(row.DeletedAt),
	}

	return result, nil
}

// ClearNonManagedProfileLinkRemoteID nulls out remote_id on non-managed links
// with matching (profile_id, kind, remote_id) to avoid unique constraint violations
// when creating a new managed link.
func (r *Repository) ClearNonManagedProfileLinkRemoteID(
	ctx context.Context,
	profileID string,
	kind string,
	remoteID string,
) error {
	_, err := r.queries.ClearNonManagedProfileLinkRemoteID(
		ctx,
		ClearNonManagedProfileLinkRemoteIDParams{
			ProfileID: profileID,
			Kind:      kind,
			RemoteID:  sql.NullString{String: remoteID, Valid: true},
		},
	)

	return err
}

// isProfileLinkRemoteIDInUseSQL checks if a remote_id is already used by another profile's active managed link.
const isProfileLinkRemoteIDInUseSQL = `
SELECT EXISTS(
  SELECT 1
  FROM "profile_link"
  WHERE kind = $1
    AND remote_id = $2
    AND profile_id != $3
    AND is_managed = TRUE
    AND deleted_at IS NULL
)
`

// IsManagedProfileLinkRemoteIDInUse returns true if the given remote_id is already used
// by another profile's active link of the same kind.
func (r *Repository) IsManagedProfileLinkRemoteIDInUse(
	ctx context.Context,
	kind string,
	remoteID string,
	excludeProfileID string,
) (bool, error) {
	var exists bool

	err := r.dbtx.QueryRowContext(
		ctx,
		isProfileLinkRemoteIDInUseSQL,
		kind,
		remoteID,
		excludeProfileID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) CreateOAuthProfileLink(
	ctx context.Context,
	id string,
	kind string,
	profileID string,
	order int,
	remoteID string,
	publicID string,
	uri string,
	authProvider string,
	authScope string,
	accessToken string,
	accessTokenExpiresAt *sql.NullTime,
	refreshToken *string,
) (*profiles.ProfileLink, error) {
	var refreshTokenSQL sql.NullString
	if refreshToken != nil {
		refreshTokenSQL = sql.NullString{String: *refreshToken, Valid: true}
	}

	var expiresAtSQL sql.NullTime
	if accessTokenExpiresAt != nil {
		expiresAtSQL = *accessTokenExpiresAt
	}

	row, err := r.queries.CreateProfileLink(ctx, CreateProfileLinkParams{
		ID:                       id,
		Kind:                     kind,
		ProfileID:                profileID,
		LinkOrder:                int32(order),
		IsManaged:                true, // OAuth links are managed
		IsVerified:               true, // OAuth links are verified
		IsFeatured:               true, // Featured by default
		Visibility:               string(profiles.LinkVisibilityPublic),
		RemoteID:                 sql.NullString{String: remoteID, Valid: true},
		PublicID:                 sql.NullString{String: publicID, Valid: true},
		URI:                      sql.NullString{String: uri, Valid: true},
		AuthProvider:             sql.NullString{String: authProvider, Valid: true},
		AuthAccessTokenScope:     sql.NullString{String: authScope, Valid: true},
		AuthAccessToken:          sql.NullString{String: accessToken, Valid: true},
		AuthAccessTokenExpiresAt: expiresAtSQL,
		AuthRefreshToken:         refreshTokenSQL,
		AuthRefreshTokenExpiresAt: sql.NullTime{
			Valid: false,
		}, // Google doesn't provide refresh token expiry
		AddedByProfileID: sql.NullString{Valid: false}, // OAuth links don't track added_by
	})
	if err != nil {
		return nil, err
	}

	result := &profiles.ProfileLink{
		ID:               row.ID,
		Kind:             row.Kind,
		ProfileID:        row.ProfileID,
		Order:            int(row.Order),
		IsManaged:        row.IsManaged,
		IsVerified:       row.IsVerified,
		IsFeatured:       row.IsFeatured,
		Visibility:       profiles.LinkVisibility(row.Visibility),
		RemoteID:         vars.ToStringPtr(row.RemoteID),
		PublicID:         vars.ToStringPtr(row.PublicID),
		URI:              vars.ToStringPtr(row.URI),
		AddedByProfileID: vars.ToStringPtr(row.AddedByProfileID),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:        vars.ToTimePtr(row.DeletedAt),
	}

	return result, nil
}

func (r *Repository) UpdateProfileLinkOAuthTokens(
	ctx context.Context,
	id string,
	publicID string,
	uri string,
	authScope string,
	accessToken string,
	accessTokenExpiresAt *sql.NullTime,
	refreshToken *string,
) error {
	var refreshTokenSQL sql.NullString
	if refreshToken != nil {
		refreshTokenSQL = sql.NullString{String: *refreshToken, Valid: true}
	}

	var expiresAtSQL sql.NullTime
	if accessTokenExpiresAt != nil {
		expiresAtSQL = *accessTokenExpiresAt
	}

	_, err := r.queries.UpdateProfileLinkOAuthTokens(ctx, UpdateProfileLinkOAuthTokensParams{
		ID:                       id,
		PublicID:                 sql.NullString{String: publicID, Valid: true},
		URI:                      sql.NullString{String: uri, Valid: true},
		AuthAccessTokenScope:     sql.NullString{String: authScope, Valid: true},
		AuthAccessToken:          sql.NullString{String: accessToken, Valid: true},
		AuthAccessTokenExpiresAt: expiresAtSQL,
		AuthRefreshToken:         refreshTokenSQL,
	})

	return err
}

func (r *Repository) GetMaxProfileLinkOrder(
	ctx context.Context,
	profileID string,
) (int, error) {
	result, err := r.queries.GetMaxProfileLinkOrder(ctx, GetMaxProfileLinkOrderParams{
		ProfileID: profileID,
	})
	if err != nil {
		return 0, err
	}

	// The result is interface{} containing int64
	if maxOrder, ok := result.(int64); ok {
		return int(maxOrder), nil
	}

	return 0, nil
}

func (r *Repository) ListAllProfilesForAdmin(
	ctx context.Context,
	localeCode string,
	filterKind string,
	limit int,
	offset int,
) ([]*profiles.Profile, error) {
	var filterKindSQL sql.NullString
	if filterKind != "" {
		filterKindSQL = sql.NullString{String: filterKind, Valid: true}
	}

	rows, err := r.queries.ListAllProfilesForAdmin(ctx, ListAllProfilesForAdminParams{
		LocaleCode:  localeCode,
		FilterKind:  filterKindSQL,
		LimitCount:  int32(limit),
		OffsetCount: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.Profile, 0, len(rows))

	for _, row := range rows {
		// HasTranslation is returned as interface{} from the query, need to convert to bool
		hasTranslation := false
		if ht, ok := row.HasTranslation.(bool); ok {
			hasTranslation = ht
		}

		profile := &profiles.Profile{
			ID:                row.ID,
			Slug:              row.Slug,
			Kind:              row.Kind,
			ProfilePictureURI: vars.ToStringPtr(row.ProfilePictureURI),
			Pronouns:          vars.ToStringPtr(row.Pronouns),
			Title:             row.Title,
			Description:       row.Description,
			Properties:        vars.ToObject(row.Properties),
			Points:            uint64(row.Points),
			CreatedAt:         row.CreatedAt,
			UpdatedAt:         vars.ToTimePtr(row.UpdatedAt),
			HasTranslation:    hasTranslation,
		}
		result = append(result, profile)
	}

	return result, nil
}

func (r *Repository) CountAllProfilesForAdmin(
	ctx context.Context,
	filterKind string,
) (int64, error) {
	var filterKindSQL sql.NullString
	if filterKind != "" {
		filterKindSQL = sql.NullString{String: filterKind, Valid: true}
	}

	count, err := r.queries.CountAllProfilesForAdmin(ctx, CountAllProfilesForAdminParams{
		FilterKind: filterKindSQL,
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetAdminProfileBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
) (*profiles.Profile, error) {
	row, err := r.queries.GetAdminProfileBySlug(ctx, GetAdminProfileBySlugParams{
		LocaleCode: localeCode,
		Slug:       slug,
	})
	if err != nil {
		return nil, err
	}

	hasTranslation := false
	if ht, ok := row.HasTranslation.(bool); ok {
		hasTranslation = ht
	}

	return &profiles.Profile{
		ID:                row.ID,
		Slug:              row.Slug,
		Kind:              row.Kind,
		ProfilePictureURI: vars.ToStringPtr(row.ProfilePictureURI),
		Pronouns:          vars.ToStringPtr(row.Pronouns),
		Title:             row.Title,
		Description:       row.Description,
		Properties:        vars.ToObject(row.Properties),
		Points:            uint64(row.Points),
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         vars.ToTimePtr(row.UpdatedAt),
		HasTranslation:    hasTranslation,
	}, nil
}

func (r *Repository) GetMembershipBetweenProfiles(
	ctx context.Context,
	profileID string,
	memberProfileID string,
) (profiles.MembershipKind, error) {
	var result string

	err := r.cache.Execute(
		ctx,
		"membership_kind:"+profileID+":"+memberProfileID,
		&result,
		func(ctx context.Context) (any, error) {
			kind, err := r.queries.GetMembershipBetweenProfiles(
				ctx,
				GetMembershipBetweenProfilesParams{
					ProfileID:       profileID,
					MemberProfileID: sql.NullString{String: memberProfileID, Valid: true},
				},
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return "", nil
				}

				return nil, err
			}

			return kind, nil
		},
	)

	return profiles.MembershipKind(result), err
}

func (r *Repository) ListFeaturedProfileLinksByProfileID(
	ctx context.Context,
	localeCode string,
	profileID string,
) ([]*profiles.ProfileLinkBrief, error) {
	rows, err := r.queries.ListFeaturedProfileLinksByProfileID(
		ctx,
		ListFeaturedProfileLinksByProfileIDParams{
			LocaleCode: localeCode,
			ProfileID:  profileID,
		},
	)
	if err != nil {
		return nil, err
	}

	profileLinks := make([]*profiles.ProfileLinkBrief, len(rows))
	for i, row := range rows {
		profileLinks[i] = &profiles.ProfileLinkBrief{
			ID:          row.ID,
			Kind:        row.Kind,
			IsManaged:   row.IsManaged,
			IsVerified:  row.IsVerified,
			IsFeatured:  row.IsFeatured,
			Visibility:  profiles.LinkVisibility(row.Visibility),
			PublicID:    row.PublicID.String,
			URI:         row.URI.String,
			Title:       row.Title,
			Icon:        row.Icon,
			Group:       row.Group,
			Description: row.Description,
		}
	}

	return profileLinks, nil
}

func (r *Repository) ListAllProfileLinksByProfileID(
	ctx context.Context,
	localeCode string,
	profileID string,
) ([]*profiles.ProfileLinkBrief, error) {
	rows, err := r.queries.ListAllProfileLinksByProfileID(
		ctx,
		ListAllProfileLinksByProfileIDParams{
			LocaleCode: localeCode,
			ProfileID:  profileID,
		},
	)
	if err != nil {
		return nil, err
	}

	profileLinks := make([]*profiles.ProfileLinkBrief, len(rows))
	for i, row := range rows {
		profileLinks[i] = &profiles.ProfileLinkBrief{
			ID:          row.ID,
			Kind:        row.Kind,
			IsManaged:   row.IsManaged,
			IsVerified:  row.IsVerified,
			IsFeatured:  row.IsFeatured,
			Visibility:  profiles.LinkVisibility(row.Visibility),
			PublicID:    row.PublicID.String,
			URI:         row.URI.String,
			Title:       row.Title,
			Icon:        row.Icon,
			Group:       row.Group,
			Description: row.Description,
		}
	}

	return profileLinks, nil
}

// ListProfileLinksByProfileIDForEditing returns all profile links for editing (settings page).
// This is an alias for ListAllProfileLinksByProfileID since we need all links for editing.
func (r *Repository) ListProfileLinksByProfileIDForEditing(
	ctx context.Context,
	localeCode string,
	profileID string,
) ([]*profiles.ProfileLinkBrief, error) {
	return r.ListAllProfileLinksByProfileID(ctx, localeCode, profileID)
}

func (r *Repository) UpsertProfileLinkTx(
	ctx context.Context,
	profileLinkID string,
	localeCode string,
	title string,
	icon *string,
	group *string,
	description *string,
) error {
	params := UpsertProfileLinkTxParams{
		ProfileLinkID: profileLinkID,
		LocaleCode:    localeCode,
		Title:         title,
		Icon:          vars.ToSQLNullString(icon),
		LinkGroup:     vars.ToSQLNullString(group),
		Description:   vars.ToSQLNullString(description),
	}

	return r.queries.UpsertProfileLinkTx(ctx, params)
}

// Profile Resource methods

// storageProfileResourceToBusinessResource converts a storage profile resource to a business profile resource.
func storageProfileResourceToBusinessResource(
	row *ListProfileResourcesByProfileIDRow,
) *profiles.ProfileResource {
	resource := &profiles.ProfileResource{
		ID:               row.ID,
		ProfileID:        row.ProfileID,
		Kind:             row.Kind,
		IsManaged:        row.IsManaged,
		RemoteID:         vars.ToStringPtr(row.RemoteID),
		PublicID:         vars.ToStringPtr(row.PublicID),
		URL:              vars.ToStringPtr(row.URL),
		Title:            row.Title,
		Description:      vars.ToStringPtr(row.Description),
		Properties:       vars.ToObject(row.Properties),
		AddedByProfileID: row.AddedByProfileID,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:        vars.ToTimePtr(row.DeletedAt),
	}

	// Populate the AddedByProfile brief if join data is present
	if row.AddedBySlug.Valid {
		resource.AddedByProfile = &profiles.ProfileBrief{
			ID:                row.AddedByProfileID,
			Slug:              row.AddedBySlug.String,
			Kind:              row.AddedByKind.String,
			Title:             row.AddedByTitle,
			Description:       row.AddedByDescription,
			ProfilePictureURI: vars.ToStringPtr(row.AddedByProfilePictureURI),
		}
	}

	return resource
}

func (r *Repository) ListProfileResourcesByProfileID(
	ctx context.Context,
	profileID string,
) ([]*profiles.ProfileResource, error) {
	rows, err := r.queries.ListProfileResourcesByProfileID(
		ctx,
		ListProfileResourcesByProfileIDParams{
			ProfileID: profileID,
		},
	)
	if err != nil {
		return nil, err
	}

	resources := make([]*profiles.ProfileResource, 0, len(rows))

	for _, row := range rows {
		resources = append(resources, storageProfileResourceToBusinessResource(row))
	}

	return resources, nil
}

func (r *Repository) GetProfileResourceByID(
	ctx context.Context,
	id string,
) (*profiles.ProfileResource, error) {
	row, err := r.queries.GetProfileResourceByID(ctx, GetProfileResourceByIDParams{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &profiles.ProfileResource{
		ID:               row.ID,
		ProfileID:        row.ProfileID,
		Kind:             row.Kind,
		IsManaged:        row.IsManaged,
		RemoteID:         vars.ToStringPtr(row.RemoteID),
		PublicID:         vars.ToStringPtr(row.PublicID),
		URL:              vars.ToStringPtr(row.URL),
		Title:            row.Title,
		Description:      vars.ToStringPtr(row.Description),
		Properties:       vars.ToObject(row.Properties),
		AddedByProfileID: row.AddedByProfileID,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:        vars.ToTimePtr(row.DeletedAt),
	}, nil
}

func (r *Repository) GetProfileResourceByRemoteID(
	ctx context.Context,
	profileID string,
	kind string,
	remoteID string,
) (*profiles.ProfileResource, error) {
	row, err := r.queries.GetProfileResourceByRemoteID(ctx, GetProfileResourceByRemoteIDParams{
		ProfileID: profileID,
		Kind:      kind,
		RemoteID:  sql.NullString{String: remoteID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &profiles.ProfileResource{
		ID:               row.ID,
		ProfileID:        row.ProfileID,
		Kind:             row.Kind,
		IsManaged:        row.IsManaged,
		RemoteID:         vars.ToStringPtr(row.RemoteID),
		PublicID:         vars.ToStringPtr(row.PublicID),
		URL:              vars.ToStringPtr(row.URL),
		Title:            row.Title,
		Description:      vars.ToStringPtr(row.Description),
		Properties:       vars.ToObject(row.Properties),
		AddedByProfileID: row.AddedByProfileID,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:        vars.ToTimePtr(row.DeletedAt),
	}, nil
}

func (r *Repository) CreateProfileResource(
	ctx context.Context,
	id string,
	profileID string,
	kind string,
	isManaged bool,
	remoteID *string,
	publicID *string,
	url *string,
	title string,
	description *string,
	properties any,
	addedByProfileID string,
) (*profiles.ProfileResource, error) {
	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal properties: %w", err)
	}

	row, err := r.queries.CreateProfileResource(ctx, CreateProfileResourceParams{
		ID:               id,
		ProfileID:        profileID,
		Kind:             kind,
		IsManaged:        isManaged,
		RemoteID:         vars.ToSQLNullString(remoteID),
		PublicID:         vars.ToSQLNullString(publicID),
		URL:              vars.ToSQLNullString(url),
		Title:            title,
		Description:      vars.ToSQLNullString(description),
		Properties:       pqtype.NullRawMessage{RawMessage: propertiesJSON, Valid: true},
		AddedByProfileID: addedByProfileID,
	})
	if err != nil {
		return nil, err
	}

	return &profiles.ProfileResource{
		ID:               row.ID,
		ProfileID:        row.ProfileID,
		Kind:             row.Kind,
		IsManaged:        row.IsManaged,
		RemoteID:         vars.ToStringPtr(row.RemoteID),
		PublicID:         vars.ToStringPtr(row.PublicID),
		URL:              vars.ToStringPtr(row.URL),
		Title:            row.Title,
		Description:      vars.ToStringPtr(row.Description),
		Properties:       vars.ToObject(row.Properties),
		AddedByProfileID: row.AddedByProfileID,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:        vars.ToTimePtr(row.DeletedAt),
	}, nil
}

func (r *Repository) SoftDeleteProfileResource(
	ctx context.Context,
	id string,
) error {
	_, err := r.queries.SoftDeleteProfileResource(ctx, SoftDeleteProfileResourceParams{
		ID: id,
	})

	return err
}

func (r *Repository) UpdateProfileResourceProperties(
	ctx context.Context,
	id string,
	properties any,
) error {
	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	_, err = r.queries.UpdateProfileResourceProperties(ctx, UpdateProfileResourcePropertiesParams{
		ID:         id,
		Properties: pqtype.NullRawMessage{RawMessage: propertiesJSON, Valid: true},
	})

	return err
}

func (r *Repository) UpdateProfileMembershipProperties(
	ctx context.Context,
	id string,
	properties any,
) error {
	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	_, err = r.queries.UpdateProfileMembershipProperties(
		ctx,
		UpdateProfileMembershipPropertiesParams{
			ID:         id,
			Properties: pqtype.NullRawMessage{RawMessage: propertiesJSON, Valid: true},
		},
	)

	return err
}

func (r *Repository) GetManagedGitHubLinkByProfileID(
	ctx context.Context,
	profileID string,
) (*profiles.ManagedGitHubLink, error) {
	row, err := r.queries.GetManagedGitHubLinkByProfileID(
		ctx,
		GetManagedGitHubLinkByProfileIDParams{
			ProfileID: profileID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &profiles.ManagedGitHubLink{
		ID:              row.ID,
		ProfileID:       row.ProfileID,
		AuthAccessToken: row.AuthAccessToken.String,
	}, nil
}
