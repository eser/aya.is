package story_interactions

import (
	"slices"
	"time"
)

// InteractionKind represents the type of interaction.
type InteractionKind string

// RSVP interaction kinds (mutually exclusive).
const (
	KindAttending    InteractionKind = "attending"
	KindInterested   InteractionKind = "interested"
	KindNotAttending InteractionKind = "not_attending"
)

// GetRSVPKinds returns all RSVP-related kinds for mutual exclusivity enforcement.
func GetRSVPKinds() []InteractionKind {
	return []InteractionKind{KindAttending, KindInterested, KindNotAttending}
}

// RSVPKindsCSV returns RSVP kinds as a comma-separated string for SQL queries.
func RSVPKindsCSV() string {
	return "attending,interested,not_attending"
}

// IsRSVPKind checks whether the given kind is an RSVP kind.
func IsRSVPKind(kind InteractionKind) bool {
	return slices.Contains(GetRSVPKinds(), kind)
}

// StoryInteraction represents a profile-to-story interaction.
type StoryInteraction struct {
	ID        string     `json:"id"`
	StoryID   string     `json:"story_id"`
	ProfileID string     `json:"profile_id"`
	Kind      string     `json:"kind"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// InteractionWithProfile extends StoryInteraction with profile display info.
type InteractionWithProfile struct {
	ID                string    `json:"id"`
	StoryID           string    `json:"story_id"`
	ProfileID         string    `json:"profile_id"`
	Kind              string    `json:"kind"`
	CreatedAt         time.Time `json:"created_at"`
	ProfileSlug       string    `json:"profile_slug"`
	ProfileTitle      string    `json:"profile_title"`
	ProfilePictureURI *string   `json:"profile_picture_uri"`
	ProfileKind       string    `json:"profile_kind"`
}

// InteractionCount represents the count of interactions for a specific kind.
type InteractionCount struct {
	Kind  string `json:"kind"`
	Count int64  `json:"count"`
}

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string
