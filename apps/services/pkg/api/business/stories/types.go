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

// Valid story kinds.
const (
	KindStatus       = "status"
	KindAnnouncement = "announcement"
	KindArticle      = "article"
	KindNews         = "news"
	KindContent      = "content"
	KindPresentation = "presentation"
	KindActivity     = "activity"
)

// ActivityProperties holds activity-specific fields from the story properties JSONB.
type ActivityProperties struct {
	ActivityKind          string `json:"activity_kind"`
	ActivityTimeStart     string `json:"activity_time_start"`
	ActivityTimeEnd       string `json:"activity_time_end"`
	ExternalActivityURI   string `json:"external_activity_uri,omitempty"`
	ExternalAttendanceURI string `json:"external_attendance_uri,omitempty"`
	RSVPMode              string `json:"rsvp_mode"`
}

// Valid activity kinds.
const (
	ActivityKindMeetup     = "meetup"
	ActivityKindWorkshop   = "workshop"
	ActivityKindConference = "conference"
	ActivityKindBroadcast  = "broadcast"
	ActivityKindMeeting    = "meeting"
)

// Valid RSVP modes.
const (
	RSVPModeDisabled          = "disabled"
	RSVPModeManagedExternally = "managed_externally"
	RSVPModeEnabled           = "enabled"
)

type Story struct {
	CreatedAt       time.Time  `json:"created_at"`
	PublishedAt     *time.Time `json:"published_at"`
	Properties      any        `json:"properties"`
	AuthorProfileID *string    `json:"author_profile_id"`
	StoryPictureURI *string    `json:"story_picture_uri"`
	SeriesID        *string    `json:"series_id"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	ID              string     `json:"id"`
	Slug            string     `json:"slug"`
	Kind            string     `json:"kind"`
	Title           string     `json:"title"`
	Summary         string     `json:"summary"`
	Content         string     `json:"content"`
	IsManaged       bool       `json:"is_managed"`
}

type StoryWithChildren struct {
	*Story

	AuthorProfile *profiles.Profile   `json:"author_profile"`
	Publications  []*profiles.Profile `json:"publications"`
}

// StoryForEdit contains the raw story data for editing (without compiled content).
type StoryForEdit struct {
	CreatedAt         time.Time  `json:"created_at"`
	Properties        any        `json:"properties"`
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
	IsManaged         bool       `json:"is_managed"`
}

// StoryForEditWithPublications wraps StoryForEdit with its publications.
type StoryForEditWithPublications struct {
	*StoryForEdit

	IsFallback   bool                `json:"is_fallback"`
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
