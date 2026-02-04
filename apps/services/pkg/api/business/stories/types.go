package stories

import (
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
)

type RecordID string

type RecordIDGenerator func() RecordID

func DefaultIDGenerator() RecordID {
	return RecordID(lib.IDsGenerateUnique())
}

type Story struct {
	CreatedAt       time.Time  `json:"created_at"`
	Properties      any        `json:"properties"`
	AuthorProfileID *string    `json:"author_profile_id"`
	StoryPictureURI *string    `json:"story_picture_uri"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	ID              string     `json:"id"`
	Slug            string     `json:"slug"`
	Kind            string     `json:"kind"`
	Title           string     `json:"title"`
	Summary         string     `json:"summary"`
	Content         string     `json:"content"`
}

type StoryWithChildren struct {
	*Story

	AuthorProfile *profiles.Profile   `json:"author_profile"`
	Publications  []*profiles.Profile `json:"publications"`
}

// StoryForEdit contains the raw story data for editing (without compiled content).
type StoryForEdit struct {
	CreatedAt         time.Time  `json:"created_at"`
	AuthorProfileID   *string    `json:"author_profile_id"`
	AuthorProfileSlug *string    `json:"author_profile_slug"`
	StoryPictureURI   *string    `json:"story_picture_uri"`
	UpdatedAt         *time.Time `json:"updated_at"`
	ID                string     `json:"id"`
	Slug              string     `json:"slug"`
	Kind              string     `json:"kind"`
	LocaleCode        string     `json:"locale_code"`
	Title             string     `json:"title"`
	Summary           string     `json:"summary"`
	Content           string     `json:"content"`
}

// StoryForEditWithPublications wraps StoryForEdit with its publications.
type StoryForEditWithPublications struct {
	*StoryForEdit

	Publications []*StoryPublication `json:"publications"`
}

// StoryPublication represents a publication of a story to a profile.
type StoryPublication struct {
	CreatedAt         time.Time  `json:"created_at"`
	PublishedAt       *time.Time `json:"published_at"`
	ProfilePictureURI *string    `json:"profile_picture_uri"`
	ID                string     `json:"id"`
	StoryID           string     `json:"story_id"`
	ProfileID         string     `json:"profile_id"`
	ProfileSlug       string     `json:"profile_slug"`
	ProfileTitle      string     `json:"profile_title"`
	ProfileKind       string     `json:"profile_kind"`
	Kind              string     `json:"kind"`
	IsFeatured        bool       `json:"is_featured"`
}

// StoryOwnership contains permission info for a story.
type StoryOwnership struct {
	ID              string  `json:"id"`
	Slug            string  `json:"slug"`
	AuthorProfileID *string `json:"author_profile_id"`
	UserKind        string  `json:"user_kind"`
	CanEdit         bool    `json:"can_edit"`
}
