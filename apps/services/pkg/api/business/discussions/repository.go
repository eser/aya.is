package discussions

import "context"

// Repository defines the storage operations for discussions (port).
type Repository interface {
	// GetDiscussionVisibility returns the discussions module visibility for a profile.
	GetDiscussionVisibility(ctx context.Context, profileID string) (string, error)

	// GetThread returns a thread by ID.
	GetThread(ctx context.Context, id string) (*Thread, error)

	// GetThreadByStoryID returns the thread for a story.
	GetThreadByStoryID(ctx context.Context, storyID string) (*Thread, error)

	// GetThreadByProfileID returns the thread for a profile.
	GetThreadByProfileID(ctx context.Context, profileID string) (*Thread, error)

	// InsertThread creates a new thread.
	InsertThread(
		ctx context.Context,
		id string,
		storyID *string,
		profileID *string,
	) (*Thread, error)

	// UpdateThreadLocked toggles the locked state of a thread.
	UpdateThreadLocked(ctx context.Context, threadID string, isLocked bool) error

	// IncrementThreadCommentCount increments the comment count of a thread.
	IncrementThreadCommentCount(ctx context.Context, threadID string) error

	// DecrementThreadCommentCount decrements the comment count of a thread.
	DecrementThreadCommentCount(ctx context.Context, threadID string) error

	// GetComment returns a comment by ID with author profile info.
	GetComment(ctx context.Context, id string, locale string) (*Comment, error)

	// GetCommentRaw returns a comment by ID without JOINs.
	GetCommentRaw(ctx context.Context, id string) (*Comment, error)

	// ListTopLevelComments returns paginated top-level comments for a thread.
	ListTopLevelComments(ctx context.Context, params ListCommentsParams) ([]*Comment, error)

	// ListChildComments returns paginated child comments for a parent.
	ListChildComments(ctx context.Context, params ListCommentsParams) ([]*Comment, error)

	// InsertComment creates a new comment.
	InsertComment(
		ctx context.Context,
		id, threadID string,
		parentID *string,
		authorUserID, content string,
		depth int,
	) (*Comment, error)

	// UpdateCommentContent updates the content of a comment.
	UpdateCommentContent(ctx context.Context, commentID, content string) error

	// SoftDeleteComment soft-deletes a comment.
	SoftDeleteComment(ctx context.Context, commentID string) error

	// UpdateCommentHidden toggles the hidden state of a comment.
	UpdateCommentHidden(ctx context.Context, commentID string, isHidden bool) error

	// UpdateCommentPinned toggles the pinned state of a comment.
	UpdateCommentPinned(ctx context.Context, commentID string, isPinned bool) error

	// IncrementCommentReplyCount increments the reply count of a comment.
	IncrementCommentReplyCount(ctx context.Context, commentID string) error

	// DecrementCommentReplyCount decrements the reply count of a comment.
	DecrementCommentReplyCount(ctx context.Context, commentID string) error

	// GetCommentVote returns a vote by comment and user.
	GetCommentVote(ctx context.Context, commentID, userID string) (*CommentVote, error)

	// InsertCommentVote creates a new vote.
	InsertCommentVote(
		ctx context.Context,
		voteID, commentID, userID string,
		direction int,
	) (*CommentVote, error)

	// UpdateCommentVoteDirection updates the direction of a vote.
	UpdateCommentVoteDirection(ctx context.Context, commentID, userID string, direction int) error

	// DeleteCommentVote removes a vote.
	DeleteCommentVote(ctx context.Context, commentID, userID string) error

	// AdjustCommentVoteScore atomically adjusts the vote score of a comment.
	AdjustCommentVoteScore(
		ctx context.Context,
		commentID string,
		scoreDelta, upvoteDelta, downvoteDelta int,
	) error

	// GetProfileIDBySlug resolves a profile slug to its ID.
	GetProfileIDBySlug(ctx context.Context, slug string) (string, error)

	// GetStoryIDBySlug resolves a story slug to its ID.
	GetStoryIDBySlug(ctx context.Context, slug string) (string, error)

	// GetStoryAuthorProfileID returns the author profile ID of a story.
	GetStoryAuthorProfileID(ctx context.Context, storyID string) (*string, error)
}
