package storage

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

// nullStringToString converts a sql.NullString to a plain string (empty if null).
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}

	return ""
}

func (r *Repository) CreateProfileMembershipCandidate(
	ctx context.Context,
	candidateID string,
	profileID string,
	referredProfileID string,
	referrerMembershipID *string,
	source string,
	applicantMessage *string,
) (*profiles.ProfileMembershipCandidate, error) {
	var referrerParam sql.NullString
	if referrerMembershipID != nil {
		referrerParam = sql.NullString{String: *referrerMembershipID, Valid: true}
	}

	var applicantMsgParam sql.NullString
	if applicantMessage != nil {
		applicantMsgParam = sql.NullString{String: *applicantMessage, Valid: true}
	}

	row, err := r.queries.CreateProfileMembershipCandidate(
		ctx,
		CreateProfileMembershipCandidateParams{
			ID:                   candidateID,
			ProfileID:            profileID,
			ReferredProfileID:    referredProfileID,
			ReferrerMembershipID: referrerParam,
			Source:               source,
			ApplicantMessage:     applicantMsgParam,
		},
	)
	if err != nil {
		return nil, err
	}

	return &profiles.ProfileMembershipCandidate{
		ID:                   row.ID,
		ProfileID:            row.ProfileID,
		ReferredProfileID:    row.ReferredProfileID,
		ReferrerMembershipID: nullStringToString(row.ReferrerMembershipID),
		Source:               row.Source,
		ApplicantMessage:     vars.ToStringPtr(row.ApplicantMessage),
		Status:               profiles.CandidateStatus(row.Status),
		VoteCount:            int(row.VoteCount),
		TotalVotes:           0,
		AverageScore:         0,
		ViewerVoteScore:      nil,
		ViewerVoteComment:    nil,
		ReferrerProfile:      nil,
		ReferredProfile:      nil,
		Teams:                nil,
		FormResponses:        nil,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) GetProfileMembershipCandidateByID(
	ctx context.Context,
	id string,
) (*profiles.ProfileMembershipCandidate, error) {
	row, err := r.queries.GetProfileMembershipCandidateByID(
		ctx,
		GetProfileMembershipCandidateByIDParams{
			ID: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return &profiles.ProfileMembershipCandidate{
		ID:                   row.ID,
		ProfileID:            row.ProfileID,
		ReferredProfileID:    row.ReferredProfileID,
		ReferrerMembershipID: nullStringToString(row.ReferrerMembershipID),
		Source:               row.Source,
		ApplicantMessage:     vars.ToStringPtr(row.ApplicantMessage),
		Status:               profiles.CandidateStatus(row.Status),
		VoteCount:            int(row.VoteCount),
		TotalVotes:           0,
		AverageScore:         0,
		ViewerVoteScore:      nil,
		ViewerVoteComment:    nil,
		ReferrerProfile:      nil,
		ReferredProfile:      nil,
		Teams:                nil,
		FormResponses:        nil,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) GetProfileMembershipCandidateByProfileAndReferred(
	ctx context.Context,
	profileID string,
	referredProfileID string,
) (*profiles.ProfileMembershipCandidate, error) {
	row, err := r.queries.GetProfileMembershipCandidateByProfileAndReferred(
		ctx,
		GetProfileMembershipCandidateByProfileAndReferredParams{
			ProfileID:         profileID,
			ReferredProfileID: referredProfileID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &profiles.ProfileMembershipCandidate{
		ID:                   row.ID,
		ProfileID:            row.ProfileID,
		ReferredProfileID:    row.ReferredProfileID,
		ReferrerMembershipID: nullStringToString(row.ReferrerMembershipID),
		Source:               row.Source,
		ApplicantMessage:     vars.ToStringPtr(row.ApplicantMessage),
		Status:               profiles.CandidateStatus(row.Status),
		VoteCount:            int(row.VoteCount),
		TotalVotes:           0,
		AverageScore:         0,
		ViewerVoteScore:      nil,
		ViewerVoteComment:    nil,
		ReferrerProfile:      nil,
		ReferredProfile:      nil,
		Teams:                nil,
		FormResponses:        nil,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

// candidateListRowToCandidate converts a ListProfileMembershipCandidatesByProfileIDRow
// to a business ProfileMembershipCandidate.
func candidateListRowToCandidate(
	row *ListProfileMembershipCandidatesByProfileIDRow,
) *profiles.ProfileMembershipCandidate {
	avgScore, _ := strconv.ParseFloat(row.AverageScore, 64)

	var viewerVoteScore *int16

	if row.ViewerVoteScore != -1 {
		score := row.ViewerVoteScore
		viewerVoteScore = &score
	}

	var referrerProfile *profiles.ProfileBrief
	if row.ReferrerProfileSlug.Valid {
		referrerProfile = &profiles.ProfileBrief{
			ID:                "",
			Slug:              row.ReferrerProfileSlug.String,
			Kind:              row.ReferrerProfileKind.String,
			ProfilePictureURI: vars.ToStringPtr(row.ReferrerProfilePictureURI),
			Title:             strings.TrimRight(row.ReferrerProfileTitle.String, " "),
			Description:       "",
		}
	}

	return &profiles.ProfileMembershipCandidate{
		ID:                   row.ID,
		ProfileID:            row.ProfileID,
		ReferredProfileID:    row.ReferredProfileID,
		ReferrerMembershipID: nullStringToString(row.ReferrerMembershipID),
		Source:               row.Source,
		ApplicantMessage:     vars.ToStringPtr(row.ApplicantMessage),
		Status:               profiles.CandidateStatus(row.Status),
		VoteCount:            int(row.VoteCount),
		TotalVotes:           row.TotalVotes,
		AverageScore:         avgScore,
		ViewerVoteScore:      viewerVoteScore,
		ViewerVoteComment:    vars.ToStringPtr(row.ViewerVoteComment),
		ReferrerProfile:      referrerProfile,
		ReferredProfile: &profiles.ProfileBrief{
			ID:                "",
			Slug:              row.ReferredProfileSlug,
			Kind:              row.ReferredProfileKind,
			ProfilePictureURI: vars.ToStringPtr(row.ReferredProfilePictureURI),
			Title:             strings.TrimRight(row.ReferredProfileTitle, " "),
			Description:       "",
		},
		Teams:         nil,
		FormResponses: nil,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     vars.ToTimePtr(row.UpdatedAt),
	}
}

func (r *Repository) ListProfileMembershipCandidatesByProfileID(
	ctx context.Context,
	localeCode string,
	profileID string,
	viewerMembershipID *string,
) ([]*profiles.ProfileMembershipCandidate, error) {
	var viewerMembershipIDParam sql.NullString
	if viewerMembershipID != nil {
		viewerMembershipIDParam = sql.NullString{String: *viewerMembershipID, Valid: true}
	}

	rows, err := r.queries.ListProfileMembershipCandidatesByProfileID(
		ctx,
		ListProfileMembershipCandidatesByProfileIDParams{
			ViewerMembershipID: viewerMembershipIDParam,
			LocaleCode:         localeCode,
			ProfileID:          profileID,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.ProfileMembershipCandidate, 0, len(rows))
	for _, row := range rows {
		result = append(result, candidateListRowToCandidate(row))
	}

	return result, nil
}

func (r *Repository) UpsertCandidateVote(
	ctx context.Context,
	voteID string,
	candidateID string,
	voterMembershipID string,
	score int16,
	comment *string,
) (*profiles.CandidateVote, error) {
	var commentParam sql.NullString
	if comment != nil {
		commentParam = sql.NullString{String: *comment, Valid: true}
	}

	row, err := r.queries.UpsertCandidateVote(ctx, UpsertCandidateVoteParams{
		ID:                voteID,
		CandidateID:       candidateID,
		VoterMembershipID: voterMembershipID,
		Score:             score,
		Comment:           commentParam,
	})
	if err != nil {
		return nil, err
	}

	return &profiles.CandidateVote{
		ID:                           row.ID,
		ProfileMembershipCandidateID: row.CandidateID,
		VoterMembershipID:            row.VoterMembershipID,
		Score:                        row.Score,
		Comment:                      vars.ToStringPtr(row.Comment),
		VoterProfile:                 nil,
		CreatedAt:                    row.CreatedAt,
		UpdatedAt:                    vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) ListCandidateVotes(
	ctx context.Context,
	localeCode string,
	candidateID string,
) ([]*profiles.CandidateVote, error) {
	rows, err := r.queries.ListCandidateVotes(ctx, ListCandidateVotesParams{
		LocaleCode:  localeCode,
		CandidateID: candidateID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.CandidateVote, 0, len(rows))

	for _, row := range rows {
		vote := &profiles.CandidateVote{
			ID:                           row.ID,
			ProfileMembershipCandidateID: row.CandidateID,
			VoterMembershipID:            row.VoterMembershipID,
			Score:                        row.Score,
			Comment:                      vars.ToStringPtr(row.Comment),
			CreatedAt:                    row.CreatedAt,
			UpdatedAt:                    vars.ToTimePtr(row.UpdatedAt),
			VoterProfile: &profiles.ProfileBrief{
				ID:                "",
				Slug:              row.VoterProfileSlug,
				Kind:              row.VoterProfileKind,
				ProfilePictureURI: vars.ToStringPtr(row.VoterProfilePictureURI),
				Title:             strings.TrimRight(row.VoterProfileTitle, " "),
				Description:       "",
			},
		}

		result = append(result, vote)
	}

	return result, nil
}

func (r *Repository) UpdateCandidateVoteCount(
	ctx context.Context,
	candidateID string,
) error {
	return r.queries.UpdateCandidateVoteCount(ctx, UpdateCandidateVoteCountParams{
		ID: candidateID,
	})
}

func (r *Repository) InsertCandidateTeam(
	ctx context.Context,
	id string,
	candidateID string,
	teamID string,
) error {
	_, err := r.queries.InsertCandidateTeam(ctx, InsertCandidateTeamParams{
		ID:            id,
		CandidateID:   candidateID,
		ProfileTeamID: teamID,
	})

	return err
}

func (r *Repository) ListCandidateTeams(
	ctx context.Context,
	candidateID string,
) ([]*profiles.ProfileTeam, error) {
	rows, err := r.queries.ListCandidateTeams(ctx, ListCandidateTeamsParams{
		CandidateID: candidateID,
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
			MemberCount:   0,
			ResourceCount: 0,
		})
	}

	return result, nil
}

func (r *Repository) UpdateCandidateStatus(
	ctx context.Context,
	candidateID string,
	profileID string,
	status profiles.CandidateStatus,
) error {
	return r.queries.UpdateCandidateStatus(ctx, UpdateCandidateStatusParams{
		Status:    string(status),
		ID:        candidateID,
		ProfileID: profileID,
	})
}

func (r *Repository) GetCandidateVoteBreakdown(
	ctx context.Context,
	candidateID string,
) (map[int]int, error) {
	rows, err := r.queries.GetCandidateVoteBreakdown(ctx, GetCandidateVoteBreakdownParams{
		CandidateID: candidateID,
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
