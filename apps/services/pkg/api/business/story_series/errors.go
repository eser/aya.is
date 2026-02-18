package story_series

import "errors"

var (
	ErrFailedToGetSeries    = errors.New("failed to get series")
	ErrFailedToListSeries   = errors.New("failed to list series")
	ErrFailedToCreateSeries = errors.New("failed to create series")
	ErrFailedToUpdateSeries = errors.New("failed to update series")
	ErrFailedToRemoveSeries = errors.New("failed to remove series")
	ErrSeriesNotFound       = errors.New("series not found")
)
