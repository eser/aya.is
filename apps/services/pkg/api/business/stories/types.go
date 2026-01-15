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
