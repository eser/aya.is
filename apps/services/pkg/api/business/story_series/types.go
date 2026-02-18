package story_series

import (
	"time"
)

// StorySeries represents a group of stories (article series, activity series, etc.).
type StorySeries struct {
	ID               string     `json:"id"`
	Slug             string     `json:"slug"`
	SeriesPictureURI *string    `json:"series_picture_uri"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

// CreateParams holds parameters for creating a new series.
type CreateParams struct {
	Slug             string
	SeriesPictureURI *string
	Title            string
	Description      string
}

// UpdateParams holds parameters for updating a series.
type UpdateParams struct {
	Slug             string
	SeriesPictureURI *string
	Title            string
	Description      string
}

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string
