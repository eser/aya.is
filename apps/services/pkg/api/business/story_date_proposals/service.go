package story_date_proposals

import (
	"context"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
)

// Service provides story date proposal operations.
type Service struct {
	logger        *logfx.Logger
	repo          Repository
	idGenerator   IDGenerator
	auditService  *events.AuditService
	storyProvider StoryProvider
	accessChecker AccessChecker
}

// NewService creates a new story date proposals service.
func NewService(
	logger *logfx.Logger,
	repo Repository,
	idGenerator IDGenerator,
	auditService *events.AuditService,
	storyProvider StoryProvider,
	accessChecker AccessChecker,
) *Service {
	return &Service{
		logger:        logger,
		repo:          repo,
		idGenerator:   idGenerator,
		auditService:  auditService,
		storyProvider: storyProvider,
		accessChecker: accessChecker,
	}
}

// validateStoryForProposals checks the story is an activity with undecided date.
func (s *Service) validateStoryForProposals(
	ctx context.Context,
	storyID string,
) (*stories.ActivityDateConfig, error) {
	config, err := s.storyProvider.GetActivityDateConfig(ctx, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrStoryNotFound, err)
	}

	if config.DateMode != stories.DateModeUndecided {
		return nil, ErrDateNotUndecided
	}

	return config, nil
}

// accessLevelToMembershipKind maps an access level string to a MembershipKind.
// Returns empty MembershipKind for "anyone" (no membership required).
func accessLevelToMembershipKind(access string) profiles.MembershipKind {
	switch access {
	case stories.DateAccessFollower:
		return profiles.MembershipKindFollower
	case stories.DateAccessMember:
		return profiles.MembershipKindMember
	case stories.DateAccessContributor:
		return profiles.MembershipKindContributor
	case stories.DateAccessMaintainer:
		return profiles.MembershipKindMaintainer
	default:
		return "" // "anyone" — no membership required
	}
}

// checkAccessLevel verifies that the user meets the minimum membership level
// against at least one of the access profile IDs (publication profiles, or
// author profile as fallback). The story author always has full access.
func (s *Service) checkAccessLevel(
	ctx context.Context,
	authorProfileID *string,
	accessProfileIDs []string,
	userProfileID string,
	requiredAccess string,
) error {
	requiredKind := accessLevelToMembershipKind(requiredAccess)
	if requiredKind == "" {
		// "anyone" — any authenticated user is allowed
		return nil
	}

	// Story author always has full access
	if authorProfileID != nil && *authorProfileID == userProfileID {
		return nil
	}

	if len(accessProfileIDs) == 0 {
		return ErrUnauthorized
	}

	levels := profiles.GetMembershipKindLevel()
	requiredLevel := levels[requiredKind]

	for _, profileID := range accessProfileIDs {
		membershipKind, err := s.accessChecker.GetMembershipKindBetween(
			ctx,
			profileID,
			userProfileID,
		)
		if err != nil || membershipKind == "" {
			continue
		}

		if levels[profiles.MembershipKind(membershipKind)] >= requiredLevel {
			return nil
		}
	}

	return ErrUnauthorized
}

// CreateProposal creates a new date proposal for an activity.
func (s *Service) CreateProposal(
	ctx context.Context,
	storyID string,
	profileID string,
	datetimeStart time.Time,
	datetimeEnd *time.Time,
) (*DateProposal, error) {
	config, err := s.validateStoryForProposals(ctx, storyID)
	if err != nil {
		return nil, err
	}

	accessErr := s.checkAccessLevel(
		ctx,
		config.AuthorProfileID,
		config.AccessProfileIDs,
		profileID,
		config.ProposalAccess,
	)
	if accessErr != nil {
		return nil, accessErr
	}

	id := s.idGenerator()

	proposal, err := s.repo.InsertProposal(ctx, id, storyID, profileID, datetimeStart, datetimeEnd)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DateProposalCreated,
		EntityType: "story_date_proposal",
		EntityID:   proposal.ID,
		ActorID:    &profileID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"story_id":       storyID,
			"datetime_start": datetimeStart,
		},
	})

	return proposal, nil
}

// ListProposals lists all proposals for a story with profile info and viewer permissions.
func (s *Service) ListProposals(
	ctx context.Context,
	localeCode string,
	storyID string,
	viewerProfileID *string,
) (*DateProposalListResponse, error) {
	proposals, err := s.repo.ListProposals(ctx, localeCode, storyID, viewerProfileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	canPropose, canVote := s.checkViewerAccess(ctx, storyID, viewerProfileID)

	return &DateProposalListResponse{
		Proposals:        proposals,
		ViewerCanPropose: canPropose,
		ViewerCanVote:    canVote,
	}, nil
}

// checkViewerAccess determines whether the viewer can propose dates and/or vote.
func (s *Service) checkViewerAccess( //nolint:nonamedreturns
	ctx context.Context,
	storyID string,
	viewerProfileID *string,
) (canPropose bool, canVote bool) {
	if viewerProfileID == nil {
		return false, false
	}

	config, err := s.validateStoryForProposals(ctx, storyID)
	if err != nil {
		return false, false
	}

	proposeErr := s.checkAccessLevel(
		ctx,
		config.AuthorProfileID,
		config.AccessProfileIDs,
		*viewerProfileID,
		config.ProposalAccess,
	)

	voteErr := s.checkAccessLevel(
		ctx,
		config.AuthorProfileID,
		config.AccessProfileIDs,
		*viewerProfileID,
		config.VoteAccess,
	)

	return proposeErr == nil, voteErr == nil
}

// RemoveProposal soft-deletes a proposal.
func (s *Service) RemoveProposal(
	ctx context.Context,
	proposalID string,
	profileID string,
) error {
	proposal, err := s.repo.GetProposal(ctx, proposalID)
	if err != nil || proposal == nil {
		return fmt.Errorf("%w: %w", ErrProposalNotFound, err)
	}

	if proposal.IsFinalized {
		return ErrCannotRemoveFinalized
	}

	// Only the proposer can remove their own proposal
	if proposal.ProposerProfileID != profileID {
		return ErrUnauthorized
	}

	_ = s.repo.DeleteAllVotesForProposal(ctx, proposalID)

	err = s.repo.SoftDeleteProposal(ctx, proposalID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRemoveRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DateProposalRemoved,
		EntityType: "story_date_proposal",
		EntityID:   proposalID,
		ActorID:    &profileID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"story_id": proposal.StoryID,
		},
	})

	return nil
}

// RemoveProposalAsAdmin soft-deletes a proposal (maintainer+ / admin).
func (s *Service) RemoveProposalAsAdmin(
	ctx context.Context,
	proposalID string,
	actorProfileID string,
) error {
	proposal, err := s.repo.GetProposal(ctx, proposalID)
	if err != nil || proposal == nil {
		return fmt.Errorf("%w: %w", ErrProposalNotFound, err)
	}

	if proposal.IsFinalized {
		return ErrCannotRemoveFinalized
	}

	_ = s.repo.DeleteAllVotesForProposal(ctx, proposalID)

	err = s.repo.SoftDeleteProposal(ctx, proposalID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRemoveRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DateProposalRemoved,
		EntityType: "story_date_proposal",
		EntityID:   proposalID,
		ActorID:    &actorProfileID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"story_id": proposal.StoryID,
			"admin":    true,
		},
	})

	return nil
}

// Vote handles agree/disagree/toggle on a date proposal.
func (s *Service) Vote( //nolint:funlen
	ctx context.Context,
	proposalID string,
	profileID string,
	direction VoteDirection,
) (*VoteResponse, error) {
	if direction != VoteAgree && direction != VoteDisagree {
		return nil, ErrInvalidVoteDirection
	}

	proposal, err := s.repo.GetProposal(ctx, proposalID)
	if err != nil || proposal == nil {
		return nil, fmt.Errorf("%w: %w", ErrProposalNotFound, err)
	}

	// Check vote access level
	config, validateErr := s.validateStoryForProposals(ctx, proposal.StoryID)
	if validateErr != nil {
		return nil, validateErr
	}

	accessErr := s.checkAccessLevel(
		ctx,
		config.AuthorProfileID,
		config.AccessProfileIDs,
		profileID,
		config.VoteAccess,
	)
	if accessErr != nil {
		return nil, accessErr
	}

	dir := int(direction)

	existing, err := s.repo.GetDateProposalVote(ctx, proposalID, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	deltas, err := s.applyVote(ctx, proposalID, profileID, existing, dir)
	if err != nil {
		return nil, err
	}

	err = s.repo.AdjustProposalVoteScore(
		ctx,
		proposalID,
		deltas.score,
		deltas.upvote,
		deltas.downvote,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DateProposalVoted,
		EntityType: "story_date_proposal",
		EntityID:   proposalID,
		ActorID:    &profileID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"direction": dir,
		},
	})

	return &VoteResponse{
		VoteScore:           proposal.VoteScore + deltas.score,
		ViewerVoteDirection: deltas.viewerDirection,
	}, nil
}

// applyVote persists the vote change and returns the computed deltas.
func (s *Service) applyVote(
	ctx context.Context,
	proposalID string,
	profileID string,
	existing *DateProposalVote,
	direction int,
) (voteDeltas, error) {
	if existing == nil {
		return s.applyNewVote(ctx, proposalID, profileID, direction)
	}

	if existing.Direction == direction {
		return s.applyRemoveVote(ctx, proposalID, profileID, direction)
	}

	return s.applyFlipVote(ctx, proposalID, profileID, direction)
}

func (s *Service) applyNewVote(
	ctx context.Context,
	proposalID string,
	profileID string,
	direction int,
) (voteDeltas, error) {
	voteID := s.idGenerator()

	_, err := s.repo.InsertDateProposalVote(ctx, voteID, proposalID, profileID, direction)
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
	proposalID string,
	profileID string,
	direction int,
) (voteDeltas, error) {
	err := s.repo.DeleteDateProposalVote(ctx, proposalID, profileID)
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
	proposalID string,
	profileID string,
	direction int,
) (voteDeltas, error) {
	err := s.repo.UpdateDateProposalVoteDirection(ctx, proposalID, profileID, direction)
	if err != nil {
		return voteDeltas{}, fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	// Flipping means: +1 to the new direction's count, -1 from the old direction's count.
	// direction=+1 (agree): upvote +1, downvote -1
	// direction=-1 (disagree): upvote -1, downvote +1
	return voteDeltas{
		score:           voteFlipMultiplier * direction,
		upvote:          direction,
		downvote:        -direction,
		viewerDirection: direction,
	}, nil
}

// splitVoteDirection returns (upvoteDelta, downvoteDelta) based on direction and magnitude.
func splitVoteDirection(direction int, magnitude int) (int, int) {
	if direction == int(VoteAgree) {
		return magnitude, 0
	}

	return 0, magnitude
}

// FinalizeProposal finalizes a proposal and sets the activity date.
func (s *Service) FinalizeProposal(
	ctx context.Context,
	proposalID string,
	storyID string,
	actorProfileID string,
) error {
	// Validate story is still undecided
	_, err := s.validateStoryForProposals(ctx, storyID)
	if err != nil {
		return err
	}

	// Get proposal
	proposal, err := s.repo.GetProposal(ctx, proposalID)
	if err != nil || proposal == nil {
		return fmt.Errorf("%w: %w", ErrProposalNotFound, err)
	}

	if proposal.StoryID != storyID {
		return ErrProposalNotFound
	}

	// Mark proposal as finalized
	err = s.repo.MarkProposalFinalized(ctx, proposalID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	// Delegate date finalization to the stories domain
	err = s.storyProvider.FinalizeActivityDate(
		ctx,
		storyID,
		proposal.DatetimeStart,
		proposal.DatetimeEnd,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.DateProposalFinalized,
		EntityType: "story_date_proposal",
		EntityID:   proposalID,
		ActorID:    &actorProfileID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"story_id":       storyID,
			"datetime_start": proposal.DatetimeStart,
		},
	})

	return nil
}
