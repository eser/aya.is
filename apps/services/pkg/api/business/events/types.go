package events

// EventType identifies all events in the system.
// Used by both event_audit (recording what happened) and event_queue (async tasks).
type EventType string

// Story events.
const (
	StoryCreated            EventType = "story_created"
	StoryUpdated            EventType = "story_updated"
	StoryDeleted            EventType = "story_deleted"
	StoryPublished          EventType = "story_published"
	StoryUnpublished        EventType = "story_unpublished"
	StoryFeatured           EventType = "story_featured"
	StoryTranslationUpdated EventType = "story_translation_updated"
	StoryTranslationDeleted EventType = "story_translation_deleted"
	StoryAutoTranslated     EventType = "story_auto_translated"
)

// Profile events.
const (
	ProfileCreated            EventType = "profile_created"
	ProfileUpdated            EventType = "profile_updated"
	ProfileTranslationUpdated EventType = "profile_translation_updated"
)

// Profile page events.
const (
	ProfilePageCreated            EventType = "profile_page_created"
	ProfilePageUpdated            EventType = "profile_page_updated"
	ProfilePageDeleted            EventType = "profile_page_deleted"
	ProfilePageTranslationUpdated EventType = "profile_page_translation_updated"
	ProfilePageTranslationDeleted EventType = "profile_page_translation_deleted"
	ProfilePageAutoTranslated     EventType = "profile_page_auto_translated"
	ProfilePageAIGenerated        EventType = "profile_page_ai_generated"
)

// Profile link events.
const (
	ProfileLinkCreated EventType = "profile_link_created"
	ProfileLinkUpdated EventType = "profile_link_updated"
	ProfileLinkDeleted EventType = "profile_link_deleted"
)

// Profile membership events.
const (
	ProfileMembershipCreated EventType = "profile_membership_created"
	ProfileMembershipUpdated EventType = "profile_membership_updated"
	ProfileMembershipDeleted EventType = "profile_membership_deleted"
)

// Profile question events.
const (
	ProfileQuestionCreated      EventType = "profile_question_created"
	ProfileQuestionAnswered     EventType = "profile_question_answered"
	ProfileQuestionAnswerEdited EventType = "profile_question_answer_edited"
	ProfileQuestionVoted        EventType = "profile_question_voted"
	ProfileQuestionHidden       EventType = "profile_question_hidden"
)

// Discussion events.
const (
	DiscussionCommentCreated EventType = "discussion_comment_created"
	DiscussionCommentEdited  EventType = "discussion_comment_edited"
	DiscussionCommentDeleted EventType = "discussion_comment_deleted"
	DiscussionCommentVoted   EventType = "discussion_comment_voted"
	DiscussionCommentHidden  EventType = "discussion_comment_hidden"
	DiscussionCommentPinned  EventType = "discussion_comment_pinned"
	DiscussionThreadLocked   EventType = "discussion_thread_locked"
)

// Points events.
const (
	PointsGained        EventType = "points_gained"
	PointsSpent         EventType = "points_spent"
	PointsTransferred   EventType = "points_transferred"
	PointsAwardApproved EventType = "points_award_approved"
	PointsAwardRejected EventType = "points_award_rejected"
)

// Story interaction events.
const (
	StoryInteractionSet     EventType = "story_interaction_set"
	StoryInteractionRemoved EventType = "story_interaction_removed"
)

// Session events.
const (
	SessionCreated    EventType = "session_created"
	SessionTerminated EventType = "session_terminated"
)

// User events.
const (
	UserCreated EventType = "user_created"
	UserUpdated EventType = "user_updated"
)

// Telegram events.
const (
	TelegramInviteLinkGenerated EventType = "telegram_invite_link_generated"
)
