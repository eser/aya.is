package discussions

import (
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
)

// Content length constraints.
const (
	MinContentLength = 3
	MaxContentLength = 10000
	MaxNestingDepth  = 6
	DefaultPageLimit = 25
	MaxPageLimit     = 100
)

// SortMode defines the comment sort order.
type SortMode string

const (
	SortHot    SortMode = "hot"
	SortNewest SortMode = "newest"
	SortOldest SortMode = "oldest"
)

// VoteDirection defines the vote direction.
type VoteDirection int

const (
	VoteUp   VoteDirection = 1
	VoteDown VoteDirection = -1
)

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string

// DefaultIDGenerator returns the default ULID-based ID generator.
func DefaultIDGenerator() string {
	return lib.IDsGenerateUnique()
}

// Thread represents a discussion thread anchored to a story or profile.
type Thread struct {
	CreatedAt    time.Time  `json:"created_at"`
	StoryID      *string    `json:"story_id"`
	ProfileID    *string    `json:"profile_id"`
	UpdatedAt    *time.Time `json:"updated_at"`
	ID           string     `json:"id"`
	CommentCount int        `json:"comment_count"`
	IsLocked     bool       `json:"is_locked"`
}

// Comment represents a discussion comment (supports nesting via ParentID).
type Comment struct {
	CreatedAt               time.Time  `json:"created_at"`
	ParentID                *string    `json:"parent_id"`
	AuthorProfileID         *string    `json:"author_profile_id"`
	AuthorProfileSlug       *string    `json:"author_profile_slug"`
	AuthorProfileTitle      *string    `json:"author_profile_title"`
	AuthorProfilePictureURI *string    `json:"author_profile_picture_uri"`
	UpdatedAt               *time.Time `json:"updated_at"`
	ID                      string     `json:"id"`
	ThreadID                string     `json:"thread_id"`
	AuthorUserID            string     `json:"-"`
	Content                 string     `json:"content"`
	Depth                   int        `json:"depth"`
	VoteScore               int        `json:"vote_score"`
	UpvoteCount             int        `json:"upvote_count"`
	DownvoteCount           int        `json:"downvote_count"`
	ReplyCount              int        `json:"reply_count"`
	ViewerVoteDirection     int        `json:"viewer_vote_direction"`
	IsPinned                bool       `json:"is_pinned"`
	IsHidden                bool       `json:"is_hidden"`
	IsEdited                bool       `json:"is_edited"`
}

// CommentVote represents a user's vote on a comment.
type CommentVote struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	CommentID string    `json:"comment_id"`
	UserID    string    `json:"user_id"`
	Direction int       `json:"direction"`
}

// CreateCommentParams holds parameters for creating a new comment.
type CreateCommentParams struct {
	StorySlug   *string
	ProfileSlug *string
	Locale      string
	UserID      string
	ParentID    *string
	Content     string
}

// EditCommentParams holds parameters for editing a comment.
type EditCommentParams struct {
	CommentID string
	UserID    string
	Content   string
}

// DeleteCommentParams holds parameters for deleting a comment.
type DeleteCommentParams struct {
	CommentID   string
	UserID      string
	ProfileSlug string
}

// VoteParams holds parameters for voting on a comment.
type VoteParams struct {
	CommentID string
	UserID    string
	Direction VoteDirection
}

// HideCommentParams holds parameters for hiding/showing a comment.
type HideCommentParams struct {
	CommentID   string
	UserID      string
	ProfileSlug string
	IsHidden    bool
}

// PinCommentParams holds parameters for pinning/unpinning a comment.
type PinCommentParams struct {
	CommentID   string
	UserID      string
	ProfileSlug string
	IsPinned    bool
}

// LockThreadParams holds parameters for locking/unlocking a thread.
type LockThreadParams struct {
	ThreadID    string
	UserID      string
	ProfileSlug string
	IsLocked    bool
}

// ListCommentsParams holds parameters for listing comments.
type ListCommentsParams struct {
	ParentID      *string
	ViewerUserID  *string
	ThreadID      string
	Locale        string
	Sort          SortMode
	Limit         int
	Offset        int
	IncludeHidden bool
}

// ListResponse wraps the thread and comments for API responses.
type ListResponse struct {
	Thread   *Thread    `json:"thread"`
	Comments []*Comment `json:"comments"`
}

// VoteResponse is returned after a vote operation.
type VoteResponse struct {
	VoteScore           int `json:"vote_score"`
	ViewerVoteDirection int `json:"viewer_vote_direction"`
}
