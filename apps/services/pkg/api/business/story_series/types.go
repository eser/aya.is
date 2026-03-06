package story_series

import (
	"time"
)

// StorySeries represents a group of stories (article series, activity series, etc.).
type StorySeries struct {
	ID               string     `json:"id"`
	Slug             string     `json:"slug"`
	SeriesPictureURI *string    `json:"series_picture_uri"`
	LocaleCode       string     `json:"locale_code"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

// CreateParams holds parameters for creating a new series.
type CreateParams struct {
	Slug             string
	SeriesPictureURI *string
	LocaleCode       string
	Title            string
	Description      string
}

// UpdateParams holds parameters for updating series base fields.
type UpdateParams struct {
	Slug             string
	SeriesPictureURI *string
}

// TranslationParams holds parameters for upserting a series translation.
type TranslationParams struct {
	LocaleCode  string
	Title       string
	Description string
}

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string
