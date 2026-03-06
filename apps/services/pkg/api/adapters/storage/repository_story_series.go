package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/eser/aya.is/services/pkg/api/business/story_series"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

func (r *Repository) GetSeriesByID(
	ctx context.Context,
	localeCode string,
	id string,
) (*story_series.StorySeries, error) {
	row, err := r.queries.GetStorySeriesByID(ctx, GetStorySeriesByIDParams{
		LocaleCode: localeCode,
		ID:         id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return mapStorySeriesFromIDRow(row), nil
}

func (r *Repository) GetSeriesBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
) (*story_series.StorySeries, error) {
	row, err := r.queries.GetStorySeriesBySlug(ctx, GetStorySeriesBySlugParams{
		LocaleCode: localeCode,
		Slug:       slug,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return mapStorySeriesFromSlugRow(row), nil
}

func (r *Repository) ListSeries(
	ctx context.Context,
	localeCode string,
) ([]*story_series.StorySeries, error) {
	rows, err := r.queries.ListStorySeries(ctx, ListStorySeriesParams{LocaleCode: localeCode})
	if err != nil {
		return nil, err
	}

	result := make([]*story_series.StorySeries, len(rows))
	for i, row := range rows {
		result[i] = mapStorySeriesFromListRow(row)
	}

	return result, nil
}

func (r *Repository) InsertSeries(
	ctx context.Context,
	seriesID string,
	slug string,
	seriesPictureURI *string,
) error {
	_, err := r.queries.InsertStorySeries(ctx, InsertStorySeriesParams{
		ID:               seriesID,
		Slug:             slug,
		SeriesPictureURI: vars.ToSQLNullString(seriesPictureURI),
	})

	return err
}

func (r *Repository) UpsertSeriesTx(
	ctx context.Context,
	seriesID string,
	localeCode string,
	title string,
	description string,
) error {
	return r.queries.UpsertStorySeriesTx(ctx, UpsertStorySeriesTxParams{
		StorySeriesID: seriesID,
		LocaleCode:    localeCode,
		Title:         title,
		Description:   description,
	})
}

func (r *Repository) UpdateSeries(
	ctx context.Context,
	seriesID string,
	slug string,
	seriesPictureURI *string,
) (int64, error) {
	return r.queries.UpdateStorySeries(ctx, UpdateStorySeriesParams{
		ID:               seriesID,
		Slug:             slug,
		SeriesPictureURI: vars.ToSQLNullString(seriesPictureURI),
	})
}

func (r *Repository) RemoveSeries(ctx context.Context, id string) (int64, error) {
	return r.queries.RemoveStorySeries(ctx, RemoveStorySeriesParams{ID: id})
}

func mapStorySeriesFromIDRow(row *GetStorySeriesByIDRow) *story_series.StorySeries {
	return &story_series.StorySeries{
		ID:               row.ID,
		Slug:             row.Slug,
		SeriesPictureURI: vars.ToStringPtr(row.SeriesPictureURI),
		LocaleCode:       strings.TrimRight(row.LocaleCode, " "),
		Title:            row.Title,
		Description:      row.Description,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
	}
}

func mapStorySeriesFromSlugRow(row *GetStorySeriesBySlugRow) *story_series.StorySeries {
	return &story_series.StorySeries{
		ID:               row.ID,
		Slug:             row.Slug,
		SeriesPictureURI: vars.ToStringPtr(row.SeriesPictureURI),
		LocaleCode:       strings.TrimRight(row.LocaleCode, " "),
		Title:            row.Title,
		Description:      row.Description,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
	}
}

func mapStorySeriesFromListRow(row *ListStorySeriesRow) *story_series.StorySeries {
	return &story_series.StorySeries{
		ID:               row.ID,
		Slug:             row.Slug,
		SeriesPictureURI: vars.ToStringPtr(row.SeriesPictureURI),
		LocaleCode:       strings.TrimRight(row.LocaleCode, " "),
		Title:            row.Title,
		Description:      row.Description,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
	}
}
