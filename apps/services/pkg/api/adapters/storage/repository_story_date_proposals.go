package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/story_date_proposals"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

func (r *Repository) InsertProposal(
	ctx context.Context,
	id string, //nolint:varnamelen
	storyID string,
	profileID string,
	datetimeStart time.Time,
	datetimeEnd *time.Time,
) (*story_date_proposals.DateProposal, error) {
	var endTime sql.NullTime
	if datetimeEnd != nil {
		endTime = sql.NullTime{Time: *datetimeEnd, Valid: true}
	}

	row, err := r.queries.InsertStoryDateProposal(ctx, InsertStoryDateProposalParams{
		ID:                id,
		StoryID:           storyID,
		ProposerProfileID: profileID,
		DatetimeStart:     datetimeStart,
		DatetimeEnd:       endTime,
	})
	if err != nil {
		return nil, err
	}

	return mapProposalRow(row), nil
}

func (r *Repository) GetProposal(
	ctx context.Context,
	id string,
) (*story_date_proposals.DateProposal, error) {
	row, err := r.queries.GetStoryDateProposal(ctx, GetStoryDateProposalParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return mapProposalRow(row), nil
}

func (r *Repository) ListProposals(
	ctx context.Context,
	localeCode string,
	storyID string,
	viewerProfileID *string,
) ([]*story_date_proposals.DateProposalWithProfile, error) {
	rows, err := r.queries.ListStoryDateProposals(ctx, ListStoryDateProposalsParams{
		LocaleCode:      localeCode,
		StoryID:         storyID,
		ViewerProfileID: vars.ToSQLNullString(viewerProfileID),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*story_date_proposals.DateProposalWithProfile, len(rows))
	for i, row := range rows {
		result[i] = &story_date_proposals.DateProposalWithProfile{
			ID:                        row.ID,
			StoryID:                   row.StoryID,
			ProposerProfileID:         row.ProposerProfileID,
			DatetimeStart:             row.DatetimeStart,
			DatetimeEnd:               vars.ToTimePtr(row.DatetimeEnd),
			IsFinalized:               row.IsFinalized,
			VoteScore:                 int(row.VoteScore),
			UpvoteCount:               int(row.UpvoteCount),
			DownvoteCount:             int(row.DownvoteCount),
			CreatedAt:                 row.CreatedAt,
			UpdatedAt:                 vars.ToTimePtr(row.UpdatedAt),
			ProposerProfileSlug:       row.ProposerProfileSlug,
			ProposerProfileTitle:      row.ProposerProfileTitle,
			ProposerProfilePictureURI: vars.ToStringPtr(row.ProposerProfilePictureURI),
			ProposerProfileKind:       row.ProposerProfileKind,
			ViewerVoteDirection:       int(row.ViewerVoteDirection),
		}
	}

	return result, nil
}

func (r *Repository) SoftDeleteProposal(ctx context.Context, id string) error {
	return r.queries.SoftDeleteStoryDateProposal(ctx, SoftDeleteStoryDateProposalParams{ID: id})
}

func (r *Repository) MarkProposalFinalized(ctx context.Context, id string) error {
	return r.queries.MarkStoryDateProposalFinalized(
		ctx,
		MarkStoryDateProposalFinalizedParams{ID: id},
	)
}

func (r *Repository) GetDateProposalVote(
	ctx context.Context,
	proposalID string,
	voterProfileID string,
) (*story_date_proposals.DateProposalVote, error) {
	row, err := r.queries.GetStoryDateProposalVote(ctx, GetStoryDateProposalVoteParams{
		ProposalID:     proposalID,
		VoterProfileID: voterProfileID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &story_date_proposals.DateProposalVote{
		ID:             row.ID,
		ProposalID:     row.ProposalID,
		VoterProfileID: row.VoterProfileID,
		Direction:      int(row.Direction),
		CreatedAt:      row.CreatedAt,
	}, nil
}

func (r *Repository) InsertDateProposalVote(
	ctx context.Context,
	id string, //nolint:varnamelen
	proposalID string,
	voterProfileID string,
	direction int,
) (*story_date_proposals.DateProposalVote, error) {
	row, err := r.queries.InsertStoryDateProposalVote(ctx, InsertStoryDateProposalVoteParams{
		ID:             id,
		ProposalID:     proposalID,
		VoterProfileID: voterProfileID,
		Direction:      int16(direction),
	})
	if err != nil {
		return nil, err
	}

	return &story_date_proposals.DateProposalVote{
		ID:             row.ID,
		ProposalID:     row.ProposalID,
		VoterProfileID: row.VoterProfileID,
		Direction:      int(row.Direction),
		CreatedAt:      row.CreatedAt,
	}, nil
}

func (r *Repository) UpdateDateProposalVoteDirection(
	ctx context.Context,
	proposalID string,
	voterProfileID string,
	direction int,
) error {
	return r.queries.UpdateStoryDateProposalVoteDirection(
		ctx,
		UpdateStoryDateProposalVoteDirectionParams{
			Direction:      int16(direction),
			ProposalID:     proposalID,
			VoterProfileID: voterProfileID,
		},
	)
}

func (r *Repository) DeleteDateProposalVote(
	ctx context.Context,
	proposalID string,
	voterProfileID string,
) error {
	return r.queries.DeleteStoryDateProposalVote(ctx, DeleteStoryDateProposalVoteParams{
		ProposalID:     proposalID,
		VoterProfileID: voterProfileID,
	})
}

func (r *Repository) DeleteAllVotesForProposal(ctx context.Context, proposalID string) error {
	return r.queries.DeleteAllVotesForProposal(ctx, DeleteAllVotesForProposalParams{
		ProposalID: proposalID,
	})
}

func (r *Repository) AdjustProposalVoteScore(
	ctx context.Context,
	proposalID string,
	scoreDelta int,
	upvoteDelta int,
	downvoteDelta int,
) error {
	return r.queries.AdjustStoryDateProposalVoteScore(ctx, AdjustStoryDateProposalVoteScoreParams{
		ScoreDelta:    int32(scoreDelta),
		UpvoteDelta:   int32(upvoteDelta),
		DownvoteDelta: int32(downvoteDelta),
		ID:            proposalID,
	})
}

// mapProposalRow converts a sqlc-generated StoryDateProposal to a business type.
func mapProposalRow(row *StoryDateProposal) *story_date_proposals.DateProposal {
	return &story_date_proposals.DateProposal{
		ID:                row.ID,
		StoryID:           row.StoryID,
		ProposerProfileID: row.ProposerProfileID,
		DatetimeStart:     row.DatetimeStart,
		DatetimeEnd:       vars.ToTimePtr(row.DatetimeEnd),
		IsFinalized:       row.IsFinalized,
		VoteScore:         int(row.VoteScore),
		UpvoteCount:       int(row.UpvoteCount),
		DownvoteCount:     int(row.DownvoteCount),
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         vars.ToTimePtr(row.UpdatedAt),
	}
}
