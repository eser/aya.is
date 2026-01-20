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
	PublishedAt     *time.Time `json:"published_at"`
	Properties      any        `json:"properties"`
	AuthorProfileID *string    `json:"author_profile_id"`
	StoryPictureURI *string    `json:"story_picture_uri"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	ID              string     `json:"id"`
	Slug            string     `json:"slug"`
	Kind            string     `json:"kind"`
	Status          string     `json:"status"`
	Title           string     `json:"title"`
	Summary         string     `json:"summary"`
	Content         string     `json:"content"`
	IsFeatured      bool       `json:"is_featured"`
}

type StoryWithChildren struct {
	*Story

	AuthorProfile *profiles.Profile   `json:"author_profile"`
	Publications  []*profiles.Profile `json:"publications"`
}

// StoryForEdit contains the raw story data for editing (without compiled content).
type StoryForEdit struct {
	CreatedAt       time.Time  `json:"created_at"`
	PublishedAt     *time.Time `json:"published_at"`
	AuthorProfileID *string    `json:"author_profile_id"`
	StoryPictureURI *string    `json:"story_picture_uri"`
	UpdatedAt       *time.Time `json:"updated_at"`
	ID              string     `json:"id"`
	Slug            string     `json:"slug"`
	Kind            string     `json:"kind"`
	Status          string     `json:"status"`
	Title           string     `json:"title"`
	Summary         string     `json:"summary"`
	Content         string     `json:"content"`
	IsFeatured      bool       `json:"is_featured"`
}

// StoryOwnership contains permission info for a story.
type StoryOwnership struct {
	ID              string  `json:"id"`
	Slug            string  `json:"slug"`
	AuthorProfileID *string `json:"author_profile_id"`
	UserKind        string  `json:"user_kind"`
	CanEdit         bool    `json:"can_edit"`
}
