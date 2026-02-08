package events

// EventType identifies all events in the system.
// Used by both event_audit (recording what happened) and event_queue (async tasks).
type EventType string

// Story events.
const (
	StoryCreated            EventType = "story:created"
	StoryUpdated            EventType = "story:updated"
	StoryDeleted            EventType = "story:deleted"
	StoryPublished          EventType = "story:published"
	StoryUnpublished        EventType = "story:unpublished"
	StoryFeatured           EventType = "story:featured"
	StoryTranslationUpdated EventType = "story:translation_updated"
	StoryTranslationDeleted EventType = "story:translation_deleted"
	StoryAutoTranslated     EventType = "story:auto_translated"
)

// Profile events.
const (
	ProfileCreated            EventType = "profile:created"
	ProfileUpdated            EventType = "profile:updated"
	ProfileTranslationUpdated EventType = "profile:translation_updated"
)

// Profile page events.
const (
	PageCreated            EventType = "page:created"
	PageUpdated            EventType = "page:updated"
	PageDeleted            EventType = "page:deleted"
	PageTranslationUpdated EventType = "page:translation_updated"
	PageTranslationDeleted EventType = "page:translation_deleted"
	PageAutoTranslated     EventType = "page:auto_translated"
)

// Profile link events.
const (
	LinkCreated EventType = "link:created"
	LinkUpdated EventType = "link:updated"
	LinkDeleted EventType = "link:deleted"
)

// Membership events.
const (
	MembershipCreated EventType = "membership:created"
	MembershipUpdated EventType = "membership:updated"
	MembershipDeleted EventType = "membership:deleted"
)

// Points events.
const (
	PointsGained      EventType = "points:gained"
	PointsSpent       EventType = "points:spent"
	PointsTransferred EventType = "points:transferred"
	AwardApproved     EventType = "award:approved"
	AwardRejected     EventType = "award:rejected"
)

// Session events.
const (
	SessionCreated    EventType = "session:created"
	SessionTerminated EventType = "session:terminated"
)

// User events.
const (
	UserCreated EventType = "user:created"
	UserUpdated EventType = "user:updated"
)
