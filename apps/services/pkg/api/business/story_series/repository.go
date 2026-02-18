package story_series

import "context"

// Repository defines the storage operations for story series (port).
type Repository interface {
	// GetSeriesByID returns a series by its ID.
	GetSeriesByID(ctx context.Context, id string) (*StorySeries, error)

	// GetSeriesBySlug returns a series by its slug.
	GetSeriesBySlug(ctx context.Context, slug string) (*StorySeries, error)

	// ListSeries returns all series.
	ListSeries(ctx context.Context) ([]*StorySeries, error)

	// InsertSeries creates a new series.
	InsertSeries(
		ctx context.Context,
		id string,
		slug string,
		seriesPictureURI *string,
		title string,
		description string,
	) (*StorySeries, error)

	// UpdateSeries updates an existing series.
	UpdateSeries(
		ctx context.Context,
		id string,
		slug string,
		seriesPictureURI *string,
		title string,
		description string,
	) (int64, error)

	// RemoveSeries soft-deletes a series.
	RemoveSeries(ctx context.Context, id string) (int64, error)
}
