package story_date_proposals

import (
	"context"
	"time"
)

// Repository defines the storage operations for story date proposals (port).
type Repository interface { //nolint:interfacebloat
	// InsertProposal creates a new date proposal.
	InsertProposal(
		ctx context.Context,
		id string,
		storyID string,
		profileID string,
		datetimeStart time.Time,
		datetimeEnd *time.Time,
	) (*DateProposal, error)

	// GetProposal returns a proposal by ID.
	GetProposal(ctx context.Context, id string) (*DateProposal, error)

	// ListProposals returns all proposals for a story with profile info and viewer's vote.
	ListProposals(
		ctx context.Context,
		localeCode string,
		storyID string,
		viewerProfileID *string,
	) ([]*DateProposalWithProfile, error)

	// SoftDeleteProposal soft-deletes a proposal.
	SoftDeleteProposal(ctx context.Context, id string) error

	// MarkProposalFinalized sets is_finalized=true on a proposal.
	MarkProposalFinalized(ctx context.Context, id string) error

	// GetDateProposalVote returns a vote by proposal and voter.
	GetDateProposalVote(
		ctx context.Context,
		proposalID string,
		voterProfileID string,
	) (*DateProposalVote, error)

	// InsertDateProposalVote creates a new vote.
	InsertDateProposalVote(
		ctx context.Context,
		id string,
		proposalID string,
		voterProfileID string,
		direction int,
	) (*DateProposalVote, error)

	// UpdateDateProposalVoteDirection updates the direction of an existing vote.
	UpdateDateProposalVoteDirection(
		ctx context.Context,
		proposalID string,
		voterProfileID string,
		direction int,
	) error

	// DeleteDateProposalVote removes a vote.
	DeleteDateProposalVote(ctx context.Context, proposalID string, voterProfileID string) error

	// DeleteAllVotesForProposal removes all votes for a proposal.
	DeleteAllVotesForProposal(ctx context.Context, proposalID string) error

	// AdjustProposalVoteScore atomically adjusts vote score.
	AdjustProposalVoteScore(
		ctx context.Context,
		proposalID string,
		scoreDelta int,
		upvoteDelta int,
		downvoteDelta int,
	) error
}
