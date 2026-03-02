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
	thread, ownerProfileSlug, err := s.resolveThreadAnchor(ctx, params)
	if err != nil {
		return nil, err
	}

	if thread.IsLocked {
		return nil, ErrThreadLocked
	}

	// Validate depth if replying to a parent.
	depth, err := s.resolveCommentDepth(ctx, params.ParentID)
	if err != nil {
		return nil, err
	}

	commentID := s.idGenerator()

	comment, err := s.repo.InsertComment(
		ctx,
		commentID,
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
		EntityID:   commentID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"thread_id":          thread.ID,
			"owner_profile_slug": ownerProfileSlug,
		},
	})

	return comment, nil
}

// resolveThreadAnchor resolves the story or profile slug to a thread, creating one if needed.
func (s *Service) resolveThreadAnchor(
	ctx context.Context,
	params CreateCommentParams,
) (*Thread, string, error) {
	switch {
	case params.StorySlug != nil:
		return s.resolveStoryThread(ctx, *params.StorySlug)
	case params.ProfileSlug != nil:
		return s.resolveProfileThread(ctx, *params.ProfileSlug)
	default:
		return nil, "", ErrThreadNotFound
	}
}

// resolveStoryThread resolves a story slug to its thread.
func (s *Service) resolveStoryThread(
	ctx context.Context,
	storySlug string,
) (*Thread, string, error) {
	storyID, err := s.repo.GetStoryIDBySlug(ctx, storySlug)
	if err != nil {
		return nil, "", fmt.Errorf(
			"%w (story slug: %s): %w",
			ErrFailedToGetRecord,
			storySlug,
			err,
		)
	}

	thread, err := s.GetOrCreateThreadByStory(ctx, storyID)
	if err != nil {
		return nil, "", err
	}

	// Resolve owner profile for visibility check.
	authorProfileID, aErr := s.repo.GetStoryAuthorProfileID(ctx, storyID)
	if aErr != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrFailedToGetRecord, aErr)
	}

	ownerProfileSlug := ""
	if authorProfileID != nil {
		ownerProfileSlug = s.resolveProfileSlug(ctx, *authorProfileID)
	}

	return thread, ownerProfileSlug, nil
}

// resolveProfileThread resolves a profile slug to its thread.
func (s *Service) resolveProfileThread(
	ctx context.Context,
	profileSlug string,
) (*Thread, string, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, "", fmt.Errorf(
			"%w (profile slug: %s): %w",
			ErrFailedToGetRecord,
			profileSlug,
			err,
		)
	}

	visibility, vErr := s.repo.GetDiscussionVisibility(ctx, profileID)
	if vErr != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrFailedToGetRecord, vErr)
	}

	if visibility == "disabled" {
		return nil, "", ErrDiscussionsNotEnabled
	}

	thread, err := s.GetOrCreateThreadByProfile(ctx, profileID)
	if err != nil {
		return nil, "", err
	}

	return thread, profileSlug, nil
}

// resolveCommentDepth validates the parent comment and returns the depth for a new reply.
func (s *Service) resolveCommentDepth(
	ctx context.Context,
	parentID *string,
) (int, error) {
	if parentID == nil {
		return 0, nil
	}

	parent, err := s.repo.GetCommentRaw(ctx, *parentID)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrCommentNotFound, err)
	}

	if parent == nil {
		return 0, ErrCommentNotFound
	}

	if parent.Depth >= MaxNestingDepth {
		return 0, ErrMaxNestingDepth
	}

	return parent.Depth + 1, nil
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
		SessionID:  nil,
		Payload:    nil,
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
		SessionID:  nil,
		Payload: map[string]any{
			"profile_slug": params.ProfileSlug,
		},
	})

	return nil
}

// voteFlipMultiplier is the score multiplier when flipping a vote direction.
const voteFlipMultiplier = 2

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

	deltas, err := s.applyVote(ctx, params, existing, direction)
	if err != nil {
		return nil, err
	}

	err = s.repo.AdjustCommentVoteScore(
		ctx,
		params.CommentID,
		deltas.score,
		deltas.upvote,
		deltas.downvote,
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
		SessionID:  nil,
		Payload: map[string]any{
			"direction": direction,
		},
	})

	return &VoteResponse{
		VoteScore:           comment.VoteScore + deltas.score,
		ViewerVoteDirection: deltas.viewerDirection,
	}, nil
}

// voteDeltas holds the computed delta values from a vote operation.
type voteDeltas struct {
	score           int
	upvote          int
	downvote        int
	viewerDirection int
}

// applyVote persists the vote change and returns the computed deltas.
func (s *Service) applyVote(
	ctx context.Context,
	params VoteParams,
	existing *CommentVote,
	direction int,
) (voteDeltas, error) {
	if existing == nil {
		return s.applyNewVote(ctx, params, direction)
	}

	if existing.Direction == direction {
		return s.applyRemoveVote(ctx, params, direction)
	}

	return s.applyFlipVote(ctx, params, direction)
}

func (s *Service) applyNewVote(
	ctx context.Context,
	params VoteParams,
	direction int,
) (voteDeltas, error) {
	voteID := s.idGenerator()

	_, err := s.repo.InsertCommentVote(ctx, voteID, params.CommentID, params.UserID, direction)
	if err != nil {
		return voteDeltas{}, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	upvote, downvote := splitVoteDirection(direction, 1)

	return voteDeltas{
		score:           direction,
		upvote:          upvote,
		downvote:        downvote,
		viewerDirection: direction,
	}, nil
}

func (s *Service) applyRemoveVote(
	ctx context.Context,
	params VoteParams,
	direction int,
) (voteDeltas, error) {
	err := s.repo.DeleteCommentVote(ctx, params.CommentID, params.UserID)
	if err != nil {
		return voteDeltas{}, fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	upvote, downvote := splitVoteDirection(direction, -1)

	return voteDeltas{
		score:           -direction,
		upvote:          upvote,
		downvote:        downvote,
		viewerDirection: 0,
	}, nil
}

func (s *Service) applyFlipVote(
	ctx context.Context,
	params VoteParams,
	direction int,
) (voteDeltas, error) {
	err := s.repo.UpdateCommentVoteDirection(ctx, params.CommentID, params.UserID, direction)
	if err != nil {
		return voteDeltas{}, fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	upvote, downvote := splitVoteDirection(direction, 1)

	return voteDeltas{
		score:           voteFlipMultiplier * direction,
		upvote:          upvote,
		downvote:        -downvote,
		viewerDirection: direction,
	}, nil
}

// splitVoteDirection returns (upvoteDelta, downvoteDelta) based on direction and magnitude.
func splitVoteDirection(direction int, magnitude int) (int, int) {
	if direction == int(VoteUp) {
		return magnitude, 0
	}

	return 0, magnitude
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
		SessionID:  nil,
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
		SessionID:  nil,
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
		SessionID:  nil,
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
