package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/story_interactions"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

func (r *Repository) UpsertInteraction(
	ctx context.Context,
	id string,
	storyID string,
	profileID string,
	kind string,
) (*story_interactions.StoryInteraction, error) {
	row, err := r.queries.UpsertStoryInteraction(ctx, UpsertStoryInteractionParams{
		ID:        id,
		StoryID:   storyID,
		ProfileID: profileID,
		Kind:      kind,
	})
	if err != nil {
		return nil, err
	}

	return &story_interactions.StoryInteraction{
		ID:        row.ID,
		StoryID:   row.StoryID,
		ProfileID: row.ProfileID,
		Kind:      row.Kind,
		CreatedAt: row.CreatedAt,
		UpdatedAt: vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) RemoveInteraction(
	ctx context.Context,
	storyID string,
	profileID string,
	kind string,
) (int64, error) {
	return r.queries.RemoveStoryInteraction(ctx, RemoveStoryInteractionParams{
		StoryID:   storyID,
		ProfileID: profileID,
		Kind:      kind,
	})
}

func (r *Repository) RemoveInteractionsByKinds(
	ctx context.Context,
	storyID string,
	profileID string,
	kindsCSV string,
) (int64, error) {
	return r.queries.RemoveStoryInteractionsByKinds(ctx, RemoveStoryInteractionsByKindsParams{
		StoryID:   storyID,
		ProfileID: profileID,
		Kinds:     kindsCSV,
	})
}

func (r *Repository) GetInteraction(
	ctx context.Context,
	storyID string,
	profileID string,
	kind string,
) (*story_interactions.StoryInteraction, error) {
	row, err := r.queries.GetStoryInteraction(ctx, GetStoryInteractionParams{
		StoryID:   storyID,
		ProfileID: profileID,
		Kind:      kind,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &story_interactions.StoryInteraction{
		ID:        row.ID,
		StoryID:   row.StoryID,
		ProfileID: row.ProfileID,
		Kind:      row.Kind,
		CreatedAt: row.CreatedAt,
		UpdatedAt: vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) ListInteractionsForProfile(
	ctx context.Context,
	storyID string,
	profileID string,
) ([]*story_interactions.StoryInteraction, error) {
	rows, err := r.queries.ListStoryInteractionsForProfile(
		ctx,
		ListStoryInteractionsForProfileParams{
			StoryID:   storyID,
			ProfileID: profileID,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*story_interactions.StoryInteraction, len(rows))
	for i, row := range rows {
		result[i] = &story_interactions.StoryInteraction{
			ID:        row.ID,
			StoryID:   row.StoryID,
			ProfileID: row.ProfileID,
			Kind:      row.Kind,
			CreatedAt: row.CreatedAt,
			UpdatedAt: vars.ToTimePtr(row.UpdatedAt),
		}
	}

	return result, nil
}

func (r *Repository) ListInteractions(
	ctx context.Context,
	localeCode string,
	storyID string,
	filterKind *string,
) ([]*story_interactions.InteractionWithProfile, error) {
	rows, err := r.queries.ListStoryInteractions(ctx, ListStoryInteractionsParams{
		LocaleCode: localeCode,
		StoryID:    storyID,
		FilterKind: vars.ToSQLNullString(filterKind),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*story_interactions.InteractionWithProfile, len(rows))
	for i, row := range rows {
		result[i] = &story_interactions.InteractionWithProfile{
			ID:                row.ID,
			StoryID:           row.StoryID,
			ProfileID:         row.ProfileID,
			Kind:              row.Kind,
			CreatedAt:         row.CreatedAt,
			ProfileSlug:       row.ProfileSlug,
			ProfileTitle:      row.ProfileTitle,
			ProfilePictureURI: vars.ToStringPtr(row.ProfilePictureURI),
			ProfileKind:       row.ProfileKind,
		}
	}

	return result, nil
}

func (r *Repository) CountInteractionsByKind(
	ctx context.Context,
	storyID string,
) ([]*story_interactions.InteractionCount, error) {
	rows, err := r.queries.CountStoryInteractionsByKind(
		ctx,
		CountStoryInteractionsByKindParams{StoryID: storyID},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*story_interactions.InteractionCount, len(rows))
	for i, row := range rows {
		result[i] = &story_interactions.InteractionCount{
			Kind:  row.Kind,
			Count: row.Count,
		}
	}

	return result, nil
}
