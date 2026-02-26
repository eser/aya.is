package storage

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

func (r *Repository) CreateProfileMembershipReferral(
	ctx context.Context,
	id string,
	profileID string,
	referredProfileID string,
	referrerMembershipID string,
) (*profiles.ProfileMembershipReferral, error) {
	row, err := r.queries.CreateProfileMembershipReferral(
		ctx,
		CreateProfileMembershipReferralParams{
			ID:                   id,
			ProfileID:            profileID,
			ReferredProfileID:    referredProfileID,
			ReferrerMembershipID: referrerMembershipID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &profiles.ProfileMembershipReferral{
		ID:                   row.ID,
		ProfileID:            row.ProfileID,
		ReferredProfileID:    row.ReferredProfileID,
		ReferrerMembershipID: row.ReferrerMembershipID,
		Status:               profiles.ReferralStatus(row.Status),
		VoteCount:            int(row.VoteCount),
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) GetProfileMembershipReferralByID(
	ctx context.Context,
	id string,
) (*profiles.ProfileMembershipReferral, error) {
	row, err := r.queries.GetProfileMembershipReferralByID(
		ctx,
		GetProfileMembershipReferralByIDParams{
			ID: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return &profiles.ProfileMembershipReferral{
		ID:                   row.ID,
		ProfileID:            row.ProfileID,
		ReferredProfileID:    row.ReferredProfileID,
		ReferrerMembershipID: row.ReferrerMembershipID,
		Status:               profiles.ReferralStatus(row.Status),
		VoteCount:            int(row.VoteCount),
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) GetProfileMembershipReferralByProfileAndReferred(
	ctx context.Context,
	profileID string,
	referredProfileID string,
) (*profiles.ProfileMembershipReferral, error) {
	row, err := r.queries.GetProfileMembershipReferralByProfileAndReferred(
		ctx,
		GetProfileMembershipReferralByProfileAndReferredParams{
			ProfileID:         profileID,
			ReferredProfileID: referredProfileID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &profiles.ProfileMembershipReferral{
		ID:                   row.ID,
		ProfileID:            row.ProfileID,
		ReferredProfileID:    row.ReferredProfileID,
		ReferrerMembershipID: row.ReferrerMembershipID,
		Status:               profiles.ReferralStatus(row.Status),
		VoteCount:            int(row.VoteCount),
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) ListProfileMembershipReferralsByProfileID(
	ctx context.Context,
	localeCode string,
	profileID string,
	viewerMembershipID *string,
) ([]*profiles.ProfileMembershipReferral, error) {
	var viewerMembershipIDParam sql.NullString
	if viewerMembershipID != nil {
		viewerMembershipIDParam = sql.NullString{String: *viewerMembershipID, Valid: true}
	}

	rows, err := r.queries.ListProfileMembershipReferralsByProfileID(
		ctx,
		ListProfileMembershipReferralsByProfileIDParams{
			ViewerMembershipID: viewerMembershipIDParam,
			LocaleCode:         localeCode,
			ProfileID:          profileID,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.ProfileMembershipReferral, 0, len(rows))

	for _, row := range rows {
		avgScore, _ := strconv.ParseFloat(row.AverageScore, 64)

		var viewerVoteScore *int16

		if row.ViewerVoteScore != -1 {
			score := row.ViewerVoteScore
			viewerVoteScore = &score
		}

		referral := &profiles.ProfileMembershipReferral{
			ID:                   row.ID,
			ProfileID:            row.ProfileID,
			ReferredProfileID:    row.ReferredProfileID,
			ReferrerMembershipID: row.ReferrerMembershipID,
			Status:               profiles.ReferralStatus(row.Status),
			VoteCount:            int(row.VoteCount),
			CreatedAt:            row.CreatedAt,
			UpdatedAt:            vars.ToTimePtr(row.UpdatedAt),
			TotalVotes:           row.TotalVotes,
			AverageScore:         avgScore,
			ViewerVoteScore:      viewerVoteScore,
			ViewerVoteComment:    vars.ToStringPtr(row.ViewerVoteComment),
			ReferrerProfile: &profiles.ProfileBrief{
				Slug:              row.ReferrerProfileSlug,
				Kind:              row.ReferrerProfileKind,
				ProfilePictureURI: vars.ToStringPtr(row.ReferrerProfilePictureURI),
				Title:             strings.TrimRight(row.ReferrerProfileTitle, " "),
			},
			ReferredProfile: &profiles.ProfileBrief{
				Slug:              row.ReferredProfileSlug,
				Kind:              row.ReferredProfileKind,
				ProfilePictureURI: vars.ToStringPtr(row.ReferredProfilePictureURI),
				Title:             strings.TrimRight(row.ReferredProfileTitle, " "),
			},
		}

		result = append(result, referral)
	}

	return result, nil
}

func (r *Repository) UpsertReferralVote(
	ctx context.Context,
	id string,
	referralID string,
	voterMembershipID string,
	score int16,
	comment *string,
) (*profiles.ReferralVote, error) {
	var commentParam sql.NullString
	if comment != nil {
		commentParam = sql.NullString{String: *comment, Valid: true}
	}

	row, err := r.queries.UpsertReferralVote(ctx, UpsertReferralVoteParams{
		ID:                          id,
		ProfileMembershipReferralID: referralID,
		VoterMembershipID:           voterMembershipID,
		Score:                       score,
		Comment:                     commentParam,
	})
	if err != nil {
		return nil, err
	}

	return &profiles.ReferralVote{
		ID:                          row.ID,
		ProfileMembershipReferralID: row.ProfileMembershipReferralID,
		VoterMembershipID:           row.VoterMembershipID,
		Score:                       row.Score,
		Comment:                     vars.ToStringPtr(row.Comment),
		CreatedAt:                   row.CreatedAt,
		UpdatedAt:                   vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) ListReferralVotes(
	ctx context.Context,
	localeCode string,
	referralID string,
) ([]*profiles.ReferralVote, error) {
	rows, err := r.queries.ListReferralVotes(ctx, ListReferralVotesParams{
		LocaleCode: localeCode,
		ReferralID: referralID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.ReferralVote, 0, len(rows))

	for _, row := range rows {
		vote := &profiles.ReferralVote{
			ID:                          row.ID,
			ProfileMembershipReferralID: row.ProfileMembershipReferralID,
			VoterMembershipID:           row.VoterMembershipID,
			Score:                       row.Score,
			Comment:                     vars.ToStringPtr(row.Comment),
			CreatedAt:                   row.CreatedAt,
			UpdatedAt:                   vars.ToTimePtr(row.UpdatedAt),
			VoterProfile: &profiles.ProfileBrief{
				Slug:              row.VoterProfileSlug,
				Kind:              row.VoterProfileKind,
				ProfilePictureURI: vars.ToStringPtr(row.VoterProfilePictureURI),
				Title:             strings.TrimRight(row.VoterProfileTitle, " "),
			},
		}

		result = append(result, vote)
	}

	return result, nil
}

func (r *Repository) UpdateReferralVoteCount(
	ctx context.Context,
	referralID string,
) error {
	return r.queries.UpdateReferralVoteCount(ctx, UpdateReferralVoteCountParams{
		ID: referralID,
	})
}

func (r *Repository) InsertReferralTeam(
	ctx context.Context,
	id string,
	referralID string,
	teamID string,
) error {
	_, err := r.queries.InsertReferralTeam(ctx, InsertReferralTeamParams{
		ID:                          id,
		ProfileMembershipReferralID: referralID,
		ProfileTeamID:               teamID,
	})

	return err
}

func (r *Repository) ListReferralTeams(
	ctx context.Context,
	referralID string,
) ([]*profiles.ProfileTeam, error) {
	rows, err := r.queries.ListReferralTeams(ctx, ListReferralTeamsParams{
		ReferralID: referralID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.ProfileTeam, 0, len(rows))

	for _, row := range rows {
		result = append(result, &profiles.ProfileTeam{
			ID:        row.ID,
			ProfileID: row.ProfileID,
			Name:      row.Name,
			Description: func() *string {
				if row.Description.Valid {
					return &row.Description.String
				}

				return nil
			}(),
		})
	}

	return result, nil
}

func (r *Repository) GetReferralVoteBreakdown(
	ctx context.Context,
	referralID string,
) (map[int]int, error) {
	rows, err := r.queries.GetReferralVoteBreakdown(ctx, GetReferralVoteBreakdownParams{
		ReferralID: referralID,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[int]int, len(rows))

	for _, row := range rows {
		result[int(row.Score)] = int(row.Count)
	}

	return result, nil
}
