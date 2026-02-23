package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/discussions"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

// GetDiscussionVisibility returns the discussions module visibility for a profile.
func (r *Repository) GetDiscussionVisibility(
	ctx context.Context,
	profileID string,
) (string, error) {
	visibility, err := r.queries.GetDiscussionVisibility(ctx, GetDiscussionVisibilityParams{
		ID: profileID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "public", nil
		}

		return "public", err
	}

	return visibility, nil
}

// GetThread returns a thread by ID.
func (r *Repository) GetThread(ctx context.Context, id string) (*discussions.Thread, error) {
	row, err := r.queries.GetDiscussionThread(ctx, GetDiscussionThreadParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return r.rowToThread(row), nil
}

// GetThreadByStoryID returns the thread for a story.
func (r *Repository) GetThreadByStoryID(
	ctx context.Context,
	storyID string,
) (*discussions.Thread, error) {
	row, err := r.queries.GetDiscussionThreadByStoryID(ctx, GetDiscussionThreadByStoryIDParams{
		StoryID: sql.NullString{String: storyID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return r.rowToThread(row), nil
}

// GetThreadByProfileID returns the thread for a profile.
func (r *Repository) GetThreadByProfileID(
	ctx context.Context,
	profileID string,
) (*discussions.Thread, error) {
	row, err := r.queries.GetDiscussionThreadByProfileID(ctx, GetDiscussionThreadByProfileIDParams{
		ProfileID: sql.NullString{String: profileID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return r.rowToThread(row), nil
}

// InsertThread creates a new thread.
func (r *Repository) InsertThread(
	ctx context.Context,
	id string,
	storyID *string,
	profileID *string,
) (*discussions.Thread, error) {
	row, err := r.queries.InsertDiscussionThread(ctx, InsertDiscussionThreadParams{
		ID:        id,
		StoryID:   vars.ToSQLNullString(storyID),
		ProfileID: vars.ToSQLNullString(profileID),
	})
	if err != nil {
		return nil, err
	}

	return r.rowToThread(row), nil
}

// UpdateThreadLocked toggles the locked state of a thread.
func (r *Repository) UpdateThreadLocked(ctx context.Context, threadID string, isLocked bool) error {
	return r.queries.UpdateDiscussionThreadLocked(ctx, UpdateDiscussionThreadLockedParams{
		ID:       threadID,
		IsLocked: isLocked,
	})
}

// IncrementThreadCommentCount increments the comment count of a thread.
func (r *Repository) IncrementThreadCommentCount(ctx context.Context, threadID string) error {
	return r.queries.IncrementDiscussionThreadCommentCount(
		ctx,
		IncrementDiscussionThreadCommentCountParams{
			ID: threadID,
		},
	)
}

// DecrementThreadCommentCount decrements the comment count of a thread.
func (r *Repository) DecrementThreadCommentCount(ctx context.Context, threadID string) error {
	return r.queries.DecrementDiscussionThreadCommentCount(
		ctx,
		DecrementDiscussionThreadCommentCountParams{
			ID: threadID,
		},
	)
}

// GetComment returns a comment by ID with author profile info.
func (r *Repository) GetComment(
	ctx context.Context,
	id string,
	locale string,
) (*discussions.Comment, error) {
	row, err := r.queries.GetDiscussionComment(ctx, GetDiscussionCommentParams{
		ID:         id,
		LocaleCode: locale,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return r.commentRowToComment(row), nil
}

// GetCommentRaw returns a comment by ID without JOINs.
func (r *Repository) GetCommentRaw(ctx context.Context, id string) (*discussions.Comment, error) {
	row, err := r.queries.GetDiscussionCommentRaw(ctx, GetDiscussionCommentRawParams{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return r.rawCommentRowToComment(row), nil
}

// ListTopLevelComments returns paginated top-level comments for a thread.
func (r *Repository) ListTopLevelComments(
	ctx context.Context,
	params discussions.ListCommentsParams,
) ([]*discussions.Comment, error) {
	rows, err := r.queries.ListTopLevelDiscussionComments(ctx, ListTopLevelDiscussionCommentsParams{
		ViewerUserID:  vars.ToSQLNullString(params.ViewerUserID),
		LocaleCode:    params.Locale,
		ThreadID:      params.ThreadID,
		IncludeHidden: params.IncludeHidden,
		SortMode:      string(params.Sort),
		PageLimit:     int32(params.Limit),
		PageOffset:    int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*discussions.Comment, len(rows))
	for i, row := range rows {
		result[i] = r.topLevelRowToComment(row)
	}

	return result, nil
}

// ListChildComments returns paginated child comments for a parent.
func (r *Repository) ListChildComments(
	ctx context.Context,
	params discussions.ListCommentsParams,
) ([]*discussions.Comment, error) {
	rows, err := r.queries.ListChildDiscussionComments(ctx, ListChildDiscussionCommentsParams{
		ViewerUserID:  vars.ToSQLNullString(params.ViewerUserID),
		LocaleCode:    params.Locale,
		ThreadID:      params.ThreadID,
		ParentID:      vars.ToSQLNullString(params.ParentID),
		IncludeHidden: params.IncludeHidden,
		PageLimit:     int32(params.Limit),
		PageOffset:    int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*discussions.Comment, len(rows))
	for i, row := range rows {
		result[i] = r.childRowToComment(row)
	}

	return result, nil
}

// InsertComment creates a new comment.
func (r *Repository) InsertComment(
	ctx context.Context,
	id, threadID string,
	parentID *string,
	authorUserID, content string,
	depth int,
) (*discussions.Comment, error) {
	row, err := r.queries.InsertDiscussionComment(ctx, InsertDiscussionCommentParams{
		ID:           id,
		ThreadID:     threadID,
		ParentID:     vars.ToSQLNullString(parentID),
		AuthorUserID: authorUserID,
		Content:      content,
		Depth:        int32(depth),
	})
	if err != nil {
		return nil, err
	}

	return r.rawCommentRowToComment(row), nil
}

// UpdateCommentContent updates the content of a comment.
func (r *Repository) UpdateCommentContent(ctx context.Context, commentID, content string) error {
	return r.queries.UpdateDiscussionCommentContent(ctx, UpdateDiscussionCommentContentParams{
		ID:      commentID,
		Content: content,
	})
}

// SoftDeleteComment soft-deletes a comment.
func (r *Repository) SoftDeleteComment(ctx context.Context, commentID string) error {
	return r.queries.SoftDeleteDiscussionComment(ctx, SoftDeleteDiscussionCommentParams{
		ID: commentID,
	})
}

// UpdateCommentHidden toggles the hidden state of a comment.
func (r *Repository) UpdateCommentHidden(
	ctx context.Context,
	commentID string,
	isHidden bool,
) error {
	return r.queries.UpdateDiscussionCommentHidden(ctx, UpdateDiscussionCommentHiddenParams{
		ID:       commentID,
		IsHidden: isHidden,
	})
}

// UpdateCommentPinned toggles the pinned state of a comment.
func (r *Repository) UpdateCommentPinned(
	ctx context.Context,
	commentID string,
	isPinned bool,
) error {
	return r.queries.UpdateDiscussionCommentPinned(ctx, UpdateDiscussionCommentPinnedParams{
		ID:       commentID,
		IsPinned: isPinned,
	})
}

// IncrementCommentReplyCount increments the reply count of a comment.
func (r *Repository) IncrementCommentReplyCount(ctx context.Context, commentID string) error {
	return r.queries.IncrementDiscussionCommentReplyCount(
		ctx,
		IncrementDiscussionCommentReplyCountParams{
			ID: commentID,
		},
	)
}

// DecrementCommentReplyCount decrements the reply count of a comment.
func (r *Repository) DecrementCommentReplyCount(ctx context.Context, commentID string) error {
	return r.queries.DecrementDiscussionCommentReplyCount(
		ctx,
		DecrementDiscussionCommentReplyCountParams{
			ID: commentID,
		},
	)
}

// GetCommentVote returns a vote by comment and user.
func (r *Repository) GetCommentVote(
	ctx context.Context,
	commentID, userID string,
) (*discussions.CommentVote, error) {
	row, err := r.queries.GetDiscussionCommentVote(ctx, GetDiscussionCommentVoteParams{
		CommentID: commentID,
		UserID:    userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &discussions.CommentVote{
		ID:        row.ID,
		CommentID: row.CommentID,
		UserID:    row.UserID,
		Direction: int(row.Direction),
		CreatedAt: row.CreatedAt,
	}, nil
}

// InsertCommentVote creates a new vote.
func (r *Repository) InsertCommentVote(
	ctx context.Context,
	voteID, commentID, userID string,
	direction int,
) (*discussions.CommentVote, error) {
	row, err := r.queries.InsertDiscussionCommentVote(ctx, InsertDiscussionCommentVoteParams{
		ID:        voteID,
		CommentID: commentID,
		UserID:    userID,
		Direction: int16(direction),
	})
	if err != nil {
		return nil, err
	}

	return &discussions.CommentVote{
		ID:        row.ID,
		CommentID: row.CommentID,
		UserID:    row.UserID,
		Direction: int(row.Direction),
		CreatedAt: row.CreatedAt,
	}, nil
}

// UpdateCommentVoteDirection updates the direction of a vote.
func (r *Repository) UpdateCommentVoteDirection(
	ctx context.Context,
	commentID, userID string,
	direction int,
) error {
	return r.queries.UpdateDiscussionCommentVoteDirection(
		ctx,
		UpdateDiscussionCommentVoteDirectionParams{
			CommentID: commentID,
			UserID:    userID,
			Direction: int16(direction),
		},
	)
}

// DeleteCommentVote removes a vote.
func (r *Repository) DeleteCommentVote(ctx context.Context, commentID, userID string) error {
	return r.queries.DeleteDiscussionCommentVote(ctx, DeleteDiscussionCommentVoteParams{
		CommentID: commentID,
		UserID:    userID,
	})
}

// AdjustCommentVoteScore atomically adjusts the vote score of a comment.
func (r *Repository) AdjustCommentVoteScore(
	ctx context.Context,
	commentID string,
	scoreDelta, upvoteDelta, downvoteDelta int,
) error {
	return r.queries.AdjustDiscussionCommentVoteScore(ctx, AdjustDiscussionCommentVoteScoreParams{
		ID:            commentID,
		ScoreDelta:    int32(scoreDelta),
		UpvoteDelta:   int32(upvoteDelta),
		DownvoteDelta: int32(downvoteDelta),
	})
}

// GetStoryAuthorProfileID returns the author profile ID of a story.
func (r *Repository) GetStoryAuthorProfileID(ctx context.Context, storyID string) (*string, error) {
	result, err := r.queries.GetStoryAuthorProfileID(ctx, GetStoryAuthorProfileIDParams{
		ID: storyID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return vars.ToStringPtr(result), nil
}

// rowToThread converts a DiscussionThread SQLC row to a domain Thread.
func (r *Repository) rowToThread(row *DiscussionThread) *discussions.Thread {
	return &discussions.Thread{
		ID:           row.ID,
		StoryID:      vars.ToStringPtr(row.StoryID),
		ProfileID:    vars.ToStringPtr(row.ProfileID),
		IsLocked:     row.IsLocked,
		CommentCount: int(row.CommentCount),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    vars.ToTimePtr(row.UpdatedAt),
	}
}

// rawCommentRowToComment converts a DiscussionComment (raw, no JOINs) to a domain Comment.
func (r *Repository) rawCommentRowToComment(row *DiscussionComment) *discussions.Comment {
	return &discussions.Comment{
		ID:            row.ID,
		ThreadID:      row.ThreadID,
		ParentID:      vars.ToStringPtr(row.ParentID),
		AuthorUserID:  row.AuthorUserID,
		Content:       row.Content,
		Depth:         int(row.Depth),
		VoteScore:     int(row.VoteScore),
		UpvoteCount:   int(row.UpvoteCount),
		DownvoteCount: int(row.DownvoteCount),
		ReplyCount:    int(row.ReplyCount),
		IsPinned:      row.IsPinned,
		IsHidden:      row.IsHidden,
		IsEdited:      row.IsEdited,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     vars.ToTimePtr(row.UpdatedAt),
	}
}

// commentRowToComment converts a GetDiscussionCommentRow (with JOINs) to a domain Comment.
func (r *Repository) commentRowToComment(row *GetDiscussionCommentRow) *discussions.Comment {
	return &discussions.Comment{
		ID:                      row.ID,
		ThreadID:                row.ThreadID,
		ParentID:                vars.ToStringPtr(row.ParentID),
		AuthorUserID:            row.AuthorUserID,
		AuthorProfileID:         vars.ToStringPtr(row.AuthorProfileID),
		AuthorProfileSlug:       vars.ToStringPtr(row.AuthorProfileSlug),
		AuthorProfileTitle:      vars.ToStringPtr(row.AuthorProfileTitle),
		AuthorProfilePictureURI: vars.ToStringPtr(row.AuthorProfilePictureURI),
		Content:                 row.Content,
		Depth:                   int(row.Depth),
		VoteScore:               int(row.VoteScore),
		UpvoteCount:             int(row.UpvoteCount),
		DownvoteCount:           int(row.DownvoteCount),
		ReplyCount:              int(row.ReplyCount),
		IsPinned:                row.IsPinned,
		IsHidden:                row.IsHidden,
		IsEdited:                row.IsEdited,
		CreatedAt:               row.CreatedAt,
		UpdatedAt:               vars.ToTimePtr(row.UpdatedAt),
	}
}

// topLevelRowToComment converts a ListTopLevelDiscussionCommentsRow to a domain Comment.
func (r *Repository) topLevelRowToComment(
	row *ListTopLevelDiscussionCommentsRow,
) *discussions.Comment {
	return &discussions.Comment{
		ID:                      row.ID,
		ThreadID:                row.ThreadID,
		ParentID:                vars.ToStringPtr(row.ParentID),
		AuthorUserID:            row.AuthorUserID,
		AuthorProfileID:         vars.ToStringPtr(row.AuthorProfileID),
		AuthorProfileSlug:       vars.ToStringPtr(row.AuthorProfileSlug),
		AuthorProfileTitle:      vars.ToStringPtr(row.AuthorProfileTitle),
		AuthorProfilePictureURI: vars.ToStringPtr(row.AuthorProfilePictureURI),
		Content:                 row.Content,
		Depth:                   int(row.Depth),
		VoteScore:               int(row.VoteScore),
		UpvoteCount:             int(row.UpvoteCount),
		DownvoteCount:           int(row.DownvoteCount),
		ReplyCount:              int(row.ReplyCount),
		IsPinned:                row.IsPinned,
		IsHidden:                row.IsHidden,
		IsEdited:                row.IsEdited,
		ViewerVoteDirection:     int(row.ViewerVoteDirection),
		CreatedAt:               row.CreatedAt,
		UpdatedAt:               vars.ToTimePtr(row.UpdatedAt),
	}
}

// childRowToComment converts a ListChildDiscussionCommentsRow to a domain Comment.
func (r *Repository) childRowToComment(row *ListChildDiscussionCommentsRow) *discussions.Comment {
	return &discussions.Comment{
		ID:                      row.ID,
		ThreadID:                row.ThreadID,
		ParentID:                vars.ToStringPtr(row.ParentID),
		AuthorUserID:            row.AuthorUserID,
		AuthorProfileID:         vars.ToStringPtr(row.AuthorProfileID),
		AuthorProfileSlug:       vars.ToStringPtr(row.AuthorProfileSlug),
		AuthorProfileTitle:      vars.ToStringPtr(row.AuthorProfileTitle),
		AuthorProfilePictureURI: vars.ToStringPtr(row.AuthorProfilePictureURI),
		Content:                 row.Content,
		Depth:                   int(row.Depth),
		VoteScore:               int(row.VoteScore),
		UpvoteCount:             int(row.UpvoteCount),
		DownvoteCount:           int(row.DownvoteCount),
		ReplyCount:              int(row.ReplyCount),
		IsPinned:                row.IsPinned,
		IsHidden:                row.IsHidden,
		IsEdited:                row.IsEdited,
		ViewerVoteDirection:     int(row.ViewerVoteDirection),
		CreatedAt:               row.CreatedAt,
		UpdatedAt:               vars.ToTimePtr(row.UpdatedAt),
	}
}
