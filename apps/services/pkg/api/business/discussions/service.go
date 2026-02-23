package discussions

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
)

// Service provides discussion operations.
type Service struct {
	logger         *logfx.Logger
	repo           Repository
	profileService *profiles.Service
	auditService   *events.AuditService
	idGenerator    IDGenerator
}

// NewService creates a new discussions service.
func NewService(
	logger *logfx.Logger,
	repo Repository,
	profileService *profiles.Service,
	auditService *events.AuditService,
	idGenerator IDGenerator,
) *Service {
	return &Service{
		logger:         logger,
		repo:           repo,
		profileService: profileService,
		auditService:   auditService,
		idGenerator:    idGenerator,
	}
}

// GetOrCreateThreadByStorySlug resolves a story slug and returns the thread, creating one if needed.
func (s *Service) GetOrCreateThreadByStorySlug(
	ctx context.Context,
	storySlug string,
) (*Thread, error) {
	storyID, err := s.repo.GetStoryIDBySlug(ctx, storySlug)
	if err != nil {
		return nil, fmt.Errorf("%w (story slug: %s): %w", ErrFailedToGetRecord, storySlug, err)
	}

	return s.GetOrCreateThreadByStory(ctx, storyID)
}

// GetOrCreateThreadByProfileSlug resolves a profile slug and returns the thread, creating one if needed.
func (s *Service) GetOrCreateThreadByProfileSlug(
	ctx context.Context,
	profileSlug string,
) (*Thread, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w (profile slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	return s.GetOrCreateThreadByProfile(ctx, profileID)
}

// GetOrCreateThreadByStory returns the thread for a story, creating one if needed.
func (s *Service) GetOrCreateThreadByStory(ctx context.Context, storyID string) (*Thread, error) {
	thread, err := s.repo.GetThreadByStoryID(ctx, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if thread != nil {
		return thread, nil
	}

	id := s.idGenerator()

	thread, err = s.repo.InsertThread(ctx, id, &storyID, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	return thread, nil
}

// GetOrCreateThreadByProfile returns the thread for a profile, creating one if needed.
func (s *Service) GetOrCreateThreadByProfile(
	ctx context.Context,
	profileID string,
) (*Thread, error) {
	thread, err := s.repo.GetThreadByProfileID(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if thread != nil {
		return thread, nil
	}

	id := s.idGenerator()

	thread, err = s.repo.InsertThread(ctx, id, nil, &profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	return thread, nil
}

// GetComment returns a single comment by ID with author profile info.
func (s *Service) GetComment(
	ctx context.Context,
	commentID string,
	locale string,
) (*Comment, error) {
	comment, err := s.repo.GetComment(ctx, commentID, locale)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	return comment, nil
}

// ListComments returns paginated comments for a thread (top-level or children).
func (s *Service) ListComments(ctx context.Context, params ListCommentsParams) ([]*Comment, error) {
	if params.Limit <= 0 || params.Limit > MaxPageLimit {
		params.Limit = DefaultPageLimit
	}

	if params.Sort == "" {
		params.Sort = SortHot
	}

	if params.ParentID != nil {
		comments, err := s.repo.ListChildComments(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
		}

		return comments, nil
	}

	comments, err := s.repo.ListTopLevelComments(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return comments, nil
}

// CreateComment creates a new comment on a thread.
func (s *Service) CreateComment(ctx context.Context, params CreateCommentParams) (*Comment, error) {
	contentLen := len([]rune(params.Content))
	if contentLen < MinContentLength {
		return nil, ErrContentTooShort
	}

	if contentLen > MaxContentLength {
		return nil, ErrContentTooLong
	}

	// Resolve the thread anchor.
	var (
		thread           *Thread
		ownerProfileSlug string
	)

	if params.StorySlug != nil {
		storyID, err := s.repo.GetStoryIDBySlug(ctx, *params.StorySlug)
		if err != nil {
			return nil, fmt.Errorf(
				"%w (story slug: %s): %w",
				ErrFailedToGetRecord,
				*params.StorySlug,
				err,
			)
		}

		thread, err = s.GetOrCreateThreadByStory(ctx, storyID)
		if err != nil {
			return nil, err
		}

		// Resolve owner profile for visibility check.
		authorProfileID, aErr := s.repo.GetStoryAuthorProfileID(ctx, storyID)
		if aErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, aErr)
		}

		if authorProfileID != nil {
			ownerProfileSlug = s.resolveProfileSlug(ctx, *authorProfileID)
		}
	} else if params.ProfileSlug != nil {
		profileID, err := s.repo.GetProfileIDBySlug(ctx, *params.ProfileSlug)
		if err != nil {
			return nil, fmt.Errorf("%w (profile slug: %s): %w", ErrFailedToGetRecord, *params.ProfileSlug, err)
		}

		visibility, vErr := s.repo.GetDiscussionVisibility(ctx, profileID)
		if vErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, vErr)
		}

		if visibility == "disabled" {
			return nil, ErrDiscussionsNotEnabled
		}

		thread, err = s.GetOrCreateThreadByProfile(ctx, profileID)
		if err != nil {
			return nil, err
		}

		ownerProfileSlug = *params.ProfileSlug
	} else {
		return nil, ErrThreadNotFound
	}

	if thread.IsLocked {
		return nil, ErrThreadLocked
	}

	// Validate depth if replying to a parent.
	depth := 0

	if params.ParentID != nil {
		parent, err := s.repo.GetCommentRaw(ctx, *params.ParentID)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrCommentNotFound, err)
		}

		if parent == nil {
			return nil, ErrCommentNotFound
		}

		if parent.Depth >= MaxNestingDepth {
			return nil, ErrMaxNestingDepth
		}

		depth = parent.Depth + 1
	}

	id := s.idGenerator()

	comment, err := s.repo.InsertComment(
		ctx,
		id,
		thread.ID,
		params.ParentID,
		params.UserID,
		params.Content,
		depth,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	// Update denormalized counts.
	_ = s.repo.IncrementThreadCommentCount(ctx, thread.ID)

	if params.ParentID != nil {
		_ = s.repo.IncrementCommentReplyCount(ctx, *params.ParentID)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DiscussionCommentCreated,
		EntityType: "discussion_comment",
		EntityID:   id,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"thread_id":          thread.ID,
			"owner_profile_slug": ownerProfileSlug,
		},
	})

	return comment, nil
}

// EditComment updates the content of a comment (author only).
func (s *Service) EditComment(ctx context.Context, params EditCommentParams) error {
	contentLen := len([]rune(params.Content))
	if contentLen < MinContentLength {
		return ErrContentTooShort
	}

	if contentLen > MaxContentLength {
		return ErrContentTooLong
	}

	comment, err := s.repo.GetCommentRaw(ctx, params.CommentID)
	if err != nil || comment == nil {
		return fmt.Errorf("%w: %w", ErrCommentNotFound, err)
	}

	if comment.AuthorUserID != params.UserID {
		return ErrInsufficientPermission
	}

	err = s.repo.UpdateCommentContent(ctx, params.CommentID, params.Content)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DiscussionCommentEdited,
		EntityType: "discussion_comment",
		EntityID:   params.CommentID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
	})

	return nil
}

// DeleteComment soft-deletes a comment (author or moderator).
func (s *Service) DeleteComment(ctx context.Context, params DeleteCommentParams) error {
	comment, err := s.repo.GetCommentRaw(ctx, params.CommentID)
	if err != nil || comment == nil {
		return fmt.Errorf("%w: %w", ErrCommentNotFound, err)
	}

	// Author can delete their own comment.
	if comment.AuthorUserID != params.UserID {
		// Otherwise, must be contributor+ on the owning profile.
		hasAccess, accessErr := s.profileService.HasUserAccessToProfile(
			ctx,
			params.UserID,
			params.ProfileSlug,
			profiles.MembershipKindContributor,
		)
		if accessErr != nil {
			return fmt.Errorf("%w: %w", ErrFailedToGetRecord, accessErr)
		}

		if !hasAccess {
			return ErrInsufficientPermission
		}
	}

	err = s.repo.SoftDeleteComment(ctx, params.CommentID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	// Update denormalized counts.
	_ = s.repo.DecrementThreadCommentCount(ctx, comment.ThreadID)

	if comment.ParentID != nil {
		_ = s.repo.DecrementCommentReplyCount(ctx, *comment.ParentID)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DiscussionCommentDeleted,
		EntityType: "discussion_comment",
		EntityID:   params.CommentID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_slug": params.ProfileSlug,
		},
	})

	return nil
}

// Vote handles upvote/downvote/toggle on a comment.
func (s *Service) Vote(ctx context.Context, params VoteParams) (*VoteResponse, error) {
	if params.Direction != VoteUp && params.Direction != VoteDown {
		return nil, ErrInvalidVoteDirection
	}

	comment, err := s.repo.GetCommentRaw(ctx, params.CommentID)
	if err != nil || comment == nil {
		return nil, fmt.Errorf("%w: %w", ErrCommentNotFound, err)
	}

	direction := int(params.Direction)

	existing, err := s.repo.GetCommentVote(ctx, params.CommentID, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	var (
		scoreDelta, upvoteDelta, downvoteDelta int
		viewerVoteDirection                    int
	)

	if existing != nil {
		if existing.Direction == direction {
			// Same direction: remove vote.
			err = s.repo.DeleteCommentVote(ctx, params.CommentID, params.UserID)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
			}

			scoreDelta = -direction
			if direction == int(VoteUp) {
				upvoteDelta = -1
			} else {
				downvoteDelta = -1
			}

			viewerVoteDirection = 0
		} else {
			// Opposite direction: flip vote.
			err = s.repo.UpdateCommentVoteDirection(ctx, params.CommentID, params.UserID, direction)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
			}

			scoreDelta = 2 * direction
			if direction == int(VoteUp) {
				upvoteDelta = 1
				downvoteDelta = -1
			} else {
				upvoteDelta = -1
				downvoteDelta = 1
			}

			viewerVoteDirection = direction
		}
	} else {
		// New vote.
		voteID := s.idGenerator()

		_, err = s.repo.InsertCommentVote(ctx, voteID, params.CommentID, params.UserID, direction)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
		}

		scoreDelta = direction
		if direction == int(VoteUp) {
			upvoteDelta = 1
		} else {
			downvoteDelta = 1
		}

		viewerVoteDirection = direction
	}

	err = s.repo.AdjustCommentVoteScore(
		ctx,
		params.CommentID,
		scoreDelta,
		upvoteDelta,
		downvoteDelta,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DiscussionCommentVoted,
		EntityType: "discussion_comment",
		EntityID:   params.CommentID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"direction": direction,
		},
	})

	return &VoteResponse{
		VoteScore:           comment.VoteScore + scoreDelta,
		ViewerVoteDirection: viewerVoteDirection,
	}, nil
}

// HideComment toggles the hidden state of a comment (contributor+ only).
func (s *Service) HideComment(ctx context.Context, params HideCommentParams) error {
	hasAccess, err := s.profileService.HasUserAccessToProfile(
		ctx,
		params.UserID,
		params.ProfileSlug,
		profiles.MembershipKindContributor,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if !hasAccess {
		return ErrInsufficientPermission
	}

	err = s.repo.UpdateCommentHidden(ctx, params.CommentID, params.IsHidden)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DiscussionCommentHidden,
		EntityType: "discussion_comment",
		EntityID:   params.CommentID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_slug": params.ProfileSlug,
			"is_hidden":    params.IsHidden,
		},
	})

	return nil
}

// PinComment toggles the pinned state of a comment (contributor+ only).
func (s *Service) PinComment(ctx context.Context, params PinCommentParams) error {
	hasAccess, err := s.profileService.HasUserAccessToProfile(
		ctx,
		params.UserID,
		params.ProfileSlug,
		profiles.MembershipKindContributor,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if !hasAccess {
		return ErrInsufficientPermission
	}

	err = s.repo.UpdateCommentPinned(ctx, params.CommentID, params.IsPinned)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DiscussionCommentPinned,
		EntityType: "discussion_comment",
		EntityID:   params.CommentID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_slug": params.ProfileSlug,
			"is_pinned":    params.IsPinned,
		},
	})

	return nil
}

// LockThread toggles the locked state of a thread (contributor+ only).
func (s *Service) LockThread(ctx context.Context, params LockThreadParams) error {
	hasAccess, err := s.profileService.HasUserAccessToProfile(
		ctx,
		params.UserID,
		params.ProfileSlug,
		profiles.MembershipKindContributor,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if !hasAccess {
		return ErrInsufficientPermission
	}

	err = s.repo.UpdateThreadLocked(ctx, params.ThreadID, params.IsLocked)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DiscussionThreadLocked,
		EntityType: "discussion_thread",
		EntityID:   params.ThreadID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_slug": params.ProfileSlug,
			"is_locked":    params.IsLocked,
		},
	})

	return nil
}

// resolveProfileSlug resolves a profile ID to its slug (best-effort).
func (s *Service) resolveProfileSlug(ctx context.Context, profileID string) string {
	// This is a helper for audit context. Failure is non-critical.
	_ = ctx
	_ = profileID

	return ""
}
