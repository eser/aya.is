package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	bulletinbiz "github.com/eser/aya.is/services/pkg/api/business/bulletin"
)

// nullStringFromPtr converts a *string to sql.NullString.
func nullStringFromPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{} //nolint:exhaustruct
	}

	return sql.NullString{String: *s, Valid: true}
}

// ptrFromNullString converts sql.NullString to *string.
func ptrFromNullString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}

	return &ns.String
}

type bulletinAdapter struct {
	repo *Repository
}

// NewBulletinRepository creates a new adapter that implements bulletin.Repository.
func NewBulletinRepository(repo *Repository) bulletinbiz.Repository {
	return &bulletinAdapter{repo: repo}
}

func (a *bulletinAdapter) GetActiveSubscriptionsForWindow(
	ctx context.Context,
	utcHour int,
) ([]*bulletinbiz.Subscription, error) {
	rows, err := a.repo.queries.GetActiveSubscriptionsForWindow(
		ctx,
		GetActiveSubscriptionsForWindowParams{
			UtcHour: int16(utcHour),
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*bulletinbiz.Subscription, 0, len(rows))

	for _, row := range rows {
		sub := &bulletinbiz.Subscription{
			ID:            row.ID,
			ProfileID:     row.ProfileID,
			ProfileSlug:   row.ProfileSlug,
			Channel:       bulletinbiz.ChannelKind(row.Channel),
			Frequency:     bulletinbiz.DigestFrequency(row.Frequency),
			PreferredTime: int(row.PreferredTime),
			DefaultLocale: strings.TrimRight(row.DefaultLocale, " "),
			CreatedAt:     row.CreatedAt,
		}

		if row.LastBulletinAt.Valid {
			t := row.LastBulletinAt.Time
			sub.LastBulletinAt = &t
		}

		if row.UpdatedAt.Valid {
			t := row.UpdatedAt.Time
			sub.UpdatedAt = &t
		}

		result = append(result, sub)
	}

	return result, nil
}

func (a *bulletinAdapter) GetFollowedProfileStoriesSince(
	ctx context.Context,
	subscriberProfileID string,
	localeCode string,
	since time.Time,
	maxStories int,
) ([]*bulletinbiz.DigestStory, error) {
	rows, err := a.repo.queries.GetFollowedProfileStoriesSince(
		ctx,
		GetFollowedProfileStoriesSinceParams{
			LocaleCode:          localeCode,
			SubscriberProfileID: sql.NullString{String: subscriberProfileID, Valid: true},
			Since:               sql.NullTime{Time: since, Valid: true},
			MaxStories:          int32(maxStories),
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*bulletinbiz.DigestStory, 0, len(rows))

	for _, row := range rows {
		story := &bulletinbiz.DigestStory{
			StoryID:                 row.StoryID,
			Slug:                    row.StorySlug,
			Kind:                    row.StoryKind,
			LocaleCode:              strings.TrimRight(row.StoryLocaleCode, " "),
			Title:                   row.StoryTitle,
			Summary:                 row.StorySummary,
			AuthorProfileID:         row.AuthorProfileID,
			AuthorSlug:              row.AuthorProfileSlug,
			AuthorTitle:             row.AuthorProfileTitle,
			AuthorProfilePictureURI: ptrFromNullString(row.AuthorProfilePictureURI),
		}

		story.StoryPictureURI = ptrFromNullString(row.StoryPictureURI)
		story.SummaryAI = ptrFromNullString(row.StorySummaryAi)

		if t, ok := row.PublishedAt.(time.Time); ok {
			story.PublishedAt = &t
		}

		result = append(result, story)
	}

	return result, nil
}

func (a *bulletinAdapter) UpdateLastBulletinAt(
	ctx context.Context,
	subscriptionID string,
) error {
	return a.repo.queries.UpdateBulletinSubscriptionLastSentAt(
		ctx,
		UpdateBulletinSubscriptionLastSentAtParams{
			ID: subscriptionID,
		},
	)
}

func (a *bulletinAdapter) CreateBulletinLog(
	ctx context.Context,
	log *bulletinbiz.BulletinLog,
) error {
	return a.repo.queries.CreateBulletinLog(ctx, CreateBulletinLogParams{
		ID:             log.ID,
		SubscriptionID: log.SubscriptionID,
		StoryCount:     int32(log.StoryCount),
		Status:         log.Status,
		ErrorMessage:   nullStringFromPtr(log.ErrorMessage),
	})
}

func (a *bulletinAdapter) UpsertStorySummaryAI(
	ctx context.Context,
	storyID string,
	localeCode string,
	summaryAI string,
) error {
	return a.repo.queries.UpsertStorySummaryAI(ctx, UpsertStorySummaryAIParams{
		StoryID:    storyID,
		LocaleCode: localeCode,
		SummaryAi:  sql.NullString{String: summaryAI, Valid: true},
	})
}

func (a *bulletinAdapter) GetSubscriptionsByProfileID(
	ctx context.Context,
	profileID string,
) ([]*bulletinbiz.Subscription, error) {
	rows, err := a.repo.queries.GetBulletinSubscriptionsByProfileID(
		ctx,
		GetBulletinSubscriptionsByProfileIDParams{
			ProfileID: profileID,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*bulletinbiz.Subscription, 0, len(rows))

	for _, row := range rows {
		sub := &bulletinbiz.Subscription{
			ID:            row.ID,
			ProfileID:     row.ProfileID,
			Channel:       bulletinbiz.ChannelKind(row.Channel),
			Frequency:     bulletinbiz.DigestFrequency(row.Frequency),
			PreferredTime: int(row.PreferredTime),
			CreatedAt:     row.CreatedAt,
		}

		if row.LastBulletinAt.Valid {
			t := row.LastBulletinAt.Time
			sub.LastBulletinAt = &t
		}

		if row.UpdatedAt.Valid {
			t := row.UpdatedAt.Time
			sub.UpdatedAt = &t
		}

		result = append(result, sub)
	}

	return result, nil
}

func (a *bulletinAdapter) GetSubscription(
	ctx context.Context,
	id string,
) (*bulletinbiz.Subscription, error) {
	row, err := a.repo.queries.GetBulletinSubscription(ctx, GetBulletinSubscriptionParams{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, bulletinbiz.ErrSubscriptionNotFound
		}

		return nil, err
	}

	sub := &bulletinbiz.Subscription{
		ID:            row.ID,
		ProfileID:     row.ProfileID,
		Channel:       bulletinbiz.ChannelKind(row.Channel),
		Frequency:     bulletinbiz.DigestFrequency(row.Frequency),
		PreferredTime: int(row.PreferredTime),
		CreatedAt:     row.CreatedAt,
	}

	if row.LastBulletinAt.Valid {
		t := row.LastBulletinAt.Time
		sub.LastBulletinAt = &t
	}

	if row.UpdatedAt.Valid {
		t := row.UpdatedAt.Time
		sub.UpdatedAt = &t
	}

	return sub, nil
}

func (a *bulletinAdapter) UpsertSubscription(
	ctx context.Context,
	sub *bulletinbiz.Subscription,
) error {
	_, err := a.repo.queries.UpsertBulletinSubscription(ctx, UpsertBulletinSubscriptionParams{
		ID:            sub.ID,
		ProfileID:     sub.ProfileID,
		Channel:       string(sub.Channel),
		Frequency:     string(sub.Frequency),
		PreferredTime: int16(sub.PreferredTime),
	})

	return err
}

func (a *bulletinAdapter) UpdateSubscriptionPreferences(
	ctx context.Context,
	id string,
	frequency bulletinbiz.DigestFrequency,
	preferredTime int,
) error {
	return a.repo.queries.UpdateBulletinSubscriptionPreferences(
		ctx,
		UpdateBulletinSubscriptionPreferencesParams{
			ID:            id,
			Frequency:     string(frequency),
			PreferredTime: int16(preferredTime),
		},
	)
}

func (a *bulletinAdapter) DeleteSubscription(
	ctx context.Context,
	id string,
) error {
	return a.repo.queries.DeleteBulletinSubscription(ctx, DeleteBulletinSubscriptionParams{
		ID: id,
	})
}

func (a *bulletinAdapter) DeleteSubscriptionsByProfileID(
	ctx context.Context,
	profileID string,
) error {
	return a.repo.queries.DeleteBulletinSubscriptionsByProfileID(
		ctx,
		DeleteBulletinSubscriptionsByProfileIDParams{
			ProfileID: profileID,
		},
	)
}

// bulletinEmailResolver implements bulletin.UserEmailResolver.
type bulletinEmailResolver struct {
	repo *Repository
}

// NewBulletinEmailResolver creates a new adapter that resolves profile IDs to emails.
func NewBulletinEmailResolver(repo *Repository) bulletinbiz.UserEmailResolver {
	return &bulletinEmailResolver{repo: repo}
}

func (r *bulletinEmailResolver) GetUserEmailByProfileID(
	ctx context.Context,
	profileID string,
) (string, error) {
	email, err := r.repo.queries.GetUserEmailByIndividualProfileID(
		ctx,
		GetUserEmailByIndividualProfileIDParams{
			ProfileID: sql.NullString{String: profileID, Valid: true},
		},
	)
	if err != nil {
		return "", err
	}

	if !email.Valid || email.String == "" {
		return "", bulletinbiz.ErrSubscriptionNotFound
	}

	return email.String, nil
}
