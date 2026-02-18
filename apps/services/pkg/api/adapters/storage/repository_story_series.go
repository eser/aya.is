package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/story_series"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

func (r *Repository) GetSeriesByID(
	ctx context.Context,
	id string,
) (*story_series.StorySeries, error) {
	row, err := r.queries.GetStorySeriesByID(ctx, GetStorySeriesByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return mapStorySeries(row), nil
}

func (r *Repository) GetSeriesBySlug(
	ctx context.Context,
	slug string,
) (*story_series.StorySeries, error) {
	row, err := r.queries.GetStorySeriesBySlug(ctx, GetStorySeriesBySlugParams{Slug: slug})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return mapStorySeries(row), nil
}

func (r *Repository) ListSeries(ctx context.Context) ([]*story_series.StorySeries, error) {
	rows, err := r.queries.ListStorySeries(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*story_series.StorySeries, len(rows))
	for i, row := range rows {
		result[i] = mapStorySeries(row)
	}

	return result, nil
}

func (r *Repository) InsertSeries(
	ctx context.Context,
	id string,
	slug string,
	seriesPictureURI *string,
	title string,
	description string,
) (*story_series.StorySeries, error) {
	row, err := r.queries.InsertStorySeries(ctx, InsertStorySeriesParams{
		ID:               id,
		Slug:             slug,
		SeriesPictureURI: vars.ToSQLNullString(seriesPictureURI),
		Title:            title,
		Description:      description,
	})
	if err != nil {
		return nil, err
	}

	return mapStorySeries(row), nil
}

func (r *Repository) UpdateSeries(
	ctx context.Context,
	id string,
	slug string,
	seriesPictureURI *string,
	title string,
	description string,
) (int64, error) {
	return r.queries.UpdateStorySeries(ctx, UpdateStorySeriesParams{
		ID:               id,
		Slug:             slug,
		SeriesPictureURI: vars.ToSQLNullString(seriesPictureURI),
		Title:            title,
		Description:      description,
	})
}

func (r *Repository) RemoveSeries(ctx context.Context, id string) (int64, error) {
	return r.queries.RemoveStorySeries(ctx, RemoveStorySeriesParams{ID: id})
}

func mapStorySeries(row *StorySeries) *story_series.StorySeries {
	return &story_series.StorySeries{
		ID:               row.ID,
		Slug:             row.Slug,
		SeriesPictureURI: vars.ToStringPtr(row.SeriesPictureURI),
		Title:            row.Title,
		Description:      row.Description,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        vars.ToTimePtr(row.UpdatedAt),
	}
}
