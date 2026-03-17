package story_series

import (
	"time"
)

// StorySeries represents a group of stories (article series, activity series, etc.).
type StorySeries struct {
	CreatedAt        time.Time  `json:"created_at"`
	SeriesPictureURI *string    `json:"series_picture_uri"`
	UpdatedAt        *time.Time `json:"updated_at"`
	ID               string     `json:"id"`
	Slug             string     `json:"slug"`
	LocaleCode       string     `json:"locale_code"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
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
	SeriesPictureURI *string
	Slug             string
}

// TranslationParams holds parameters for upserting a series translation.
type TranslationParams struct {
	LocaleCode  string
	Title       string
	Description string
}

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string
