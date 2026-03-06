package story_series

import "context"

// Repository defines the storage operations for story series (port).
type Repository interface {
	// GetSeriesByID returns a series by its ID with localized text.
	GetSeriesByID(ctx context.Context, localeCode string, id string) (*StorySeries, error)

	// GetSeriesBySlug returns a series by its slug with localized text.
	GetSeriesBySlug(ctx context.Context, localeCode string, slug string) (*StorySeries, error)

	// ListSeries returns all series with localized text.
	ListSeries(ctx context.Context, localeCode string) ([]*StorySeries, error)

	// InsertSeries creates a new series (base table only).
	InsertSeries(
		ctx context.Context,
		id string,
		slug string,
		seriesPictureURI *string,
	) error

	// UpsertSeriesTx creates or updates a series translation.
	UpsertSeriesTx(
		ctx context.Context,
		seriesID string,
		localeCode string,
		title string,
		description string,
	) error

	// UpdateSeries updates series base fields.
	UpdateSeries(
		ctx context.Context,
		id string,
		slug string,
		seriesPictureURI *string,
	) (int64, error)

	// RemoveSeries soft-deletes a series.
	RemoveSeries(ctx context.Context, id string) (int64, error)
}
