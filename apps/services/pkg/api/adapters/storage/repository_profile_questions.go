package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/profile_questions"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

// IsQAHidden checks whether Q&A is hidden for a profile.
func (r *Repository) IsQAHidden(ctx context.Context, profileID string) (bool, error) {
	hidden, err := r.queries.IsProfileQAHidden(ctx, IsProfileQAHiddenParams{
		ID: profileID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return hidden, nil
}

// GetQuestion returns a single question by ID.
func (r *Repository) GetQuestion(
	ctx context.Context,
	id string,
) (*profile_questions.Question, error) {
	row, err := r.queries.GetProfileQuestion(ctx, GetProfileQuestionParams{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, profile_questions.ErrQuestionNotFound
		}

		return nil, err
	}

	return r.rowToQuestion(row), nil
}

// ListQuestions returns questions for a profile with pagination.
func (r *Repository) ListQuestions(
	ctx context.Context,
	profileID string,
	localeCode string,
	viewerUserID *string,
	includeHidden bool,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profile_questions.Question], error) {
	rows, err := r.queries.ListProfileQuestionsByProfileID(
		ctx,
		ListProfileQuestionsByProfileIDParams{
			ViewerUserID:  vars.ToSQLNullString(viewerUserID),
			LocaleCode:    localeCode,
			ProfileID:     profileID,
			IncludeHidden: includeHidden,
		},
	)
	if err != nil {
		return cursors.Cursored[[]*profile_questions.Question]{}, err
	}

	result := make([]*profile_questions.Question, len(rows))
	for i, row := range rows {
		result[i] = r.listRowToQuestion(row)
	}

	return cursors.WrapResponseWithCursor[[]*profile_questions.Question](result, nil), nil
}

// InsertQuestion creates a new question record.
func (r *Repository) InsertQuestion(
	ctx context.Context,
	id string,
	profileID string,
	authorUserID string,
	content string,
	isAnonymous bool,
) (*profile_questions.Question, error) {
	row, err := r.queries.InsertProfileQuestion(ctx, InsertProfileQuestionParams{
		ID:           id,
		ProfileID:    profileID,
		AuthorUserID: authorUserID,
		Content:      content,
		IsAnonymous:  isAnonymous,
	})
	if err != nil {
		return nil, err
	}

	return r.rowToQuestion(row), nil
}

// UpdateAnswer sets the answer content on a question.
func (r *Repository) UpdateAnswer(
	ctx context.Context,
	questionID string,
	answerContent string,
	answerURI *string,
	answerKind *string,
	answeredByProfileID string,
) error {
	return r.queries.UpdateProfileQuestionAnswer(ctx, UpdateProfileQuestionAnswerParams{
		ID:            questionID,
		AnswerContent: sql.NullString{String: answerContent, Valid: true},
		AnswerURI:     vars.ToSQLNullString(answerURI),
		AnswerKind:    vars.ToSQLNullString(answerKind),
		AnsweredBy:    sql.NullString{String: answeredByProfileID, Valid: true},
	})
}

// UpdateHidden toggles the hidden state of a question.
func (r *Repository) UpdateHidden(ctx context.Context, questionID string, isHidden bool) error {
	return r.queries.UpdateProfileQuestionHidden(ctx, UpdateProfileQuestionHiddenParams{
		ID:       questionID,
		IsHidden: isHidden,
	})
}

// InsertVote creates a new vote record and increments the question's vote count.
func (r *Repository) InsertVote(
	ctx context.Context,
	voteID string,
	questionID string,
	userID string,
) (*profile_questions.Vote, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	queriesTx := r.queries.WithTx(tx)

	row, err := queriesTx.InsertProfileQuestionVote(ctx, InsertProfileQuestionVoteParams{
		ID:         voteID,
		QuestionID: questionID,
		UserID:     userID,
		Score:      1,
	})
	if err != nil {
		return nil, err
	}

	err = queriesTx.IncrementProfileQuestionVoteCount(ctx, IncrementProfileQuestionVoteCountParams{
		ID: questionID,
	})
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &profile_questions.Vote{
		ID:         row.ID,
		QuestionID: row.QuestionID,
		UserID:     row.UserID,
		Score:      int(row.Score),
		CreatedAt:  row.CreatedAt,
	}, nil
}

// DeleteVote removes a vote record and decrements the question's vote count.
func (r *Repository) DeleteVote(ctx context.Context, questionID string, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	queriesTx := r.queries.WithTx(tx)

	err = queriesTx.DeleteProfileQuestionVote(ctx, DeleteProfileQuestionVoteParams{
		QuestionID: questionID,
		UserID:     userID,
	})
	if err != nil {
		return err
	}

	err = queriesTx.DecrementProfileQuestionVoteCount(ctx, DecrementProfileQuestionVoteCountParams{
		ID: questionID,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetVote returns a vote by question and user.
func (r *Repository) GetVote(
	ctx context.Context,
	questionID string,
	userID string,
) (*profile_questions.Vote, error) {
	row, err := r.queries.GetProfileQuestionVote(ctx, GetProfileQuestionVoteParams{
		QuestionID: questionID,
		UserID:     userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &profile_questions.Vote{
		ID:         row.ID,
		QuestionID: row.QuestionID,
		UserID:     row.UserID,
		Score:      int(row.Score),
		CreatedAt:  row.CreatedAt,
	}, nil
}

// rowToQuestion converts a ProfileQuestion SQLC row to a domain Question.
func (r *Repository) rowToQuestion(row *ProfileQuestion) *profile_questions.Question {
	return &profile_questions.Question{
		ID:                  row.ID,
		ProfileID:           row.ProfileID,
		Content:             row.Content,
		AnswerContent:       vars.ToStringPtr(row.AnswerContent),
		AnswerURI:           vars.ToStringPtr(row.AnswerURI),
		AnswerKind:          vars.ToStringPtr(row.AnswerKind),
		AnsweredAt:          vars.ToTimePtr(row.AnsweredAt),
		AnsweredByProfileID: vars.ToStringPtr(row.AnsweredBy),
		VoteCount:           int(row.VoteCount),
		IsAnonymous:         row.IsAnonymous,
		IsHidden:            row.IsHidden,
		HasViewerVote:       false,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
	}
}

// listRowToQuestion converts a ListProfileQuestionsByProfileIDRow to a domain Question.
func (r *Repository) listRowToQuestion(
	row *ListProfileQuestionsByProfileIDRow,
) *profile_questions.Question {
	return &profile_questions.Question{
		ID:                     row.ID,
		ProfileID:              row.ProfileID,
		AuthorProfileID:        vars.ToStringPtr(row.AuthorProfileID),
		AuthorProfileSlug:      vars.ToStringPtr(row.AuthorProfileSlug),
		AuthorProfileTitle:     vars.ToStringPtr(row.AuthorProfileTitle),
		Content:                row.Content,
		AnswerContent:          vars.ToStringPtr(row.AnswerContent),
		AnswerURI:              vars.ToStringPtr(row.AnswerURI),
		AnswerKind:             vars.ToStringPtr(row.AnswerKind),
		AnsweredAt:             vars.ToTimePtr(row.AnsweredAt),
		AnsweredByProfileID:    vars.ToStringPtr(row.AnsweredBy),
		AnsweredByProfileSlug:  vars.ToStringPtr(row.AnsweredByProfileSlug),
		AnsweredByProfileTitle: vars.ToStringPtr(row.AnsweredByProfileTitle),
		VoteCount:              int(row.VoteCount),
		IsAnonymous:            row.IsAnonymous,
		IsHidden:               row.IsHidden,
		HasViewerVote:          row.HasViewerVote,
		CreatedAt:              row.CreatedAt,
		UpdatedAt:              vars.ToTimePtr(row.UpdatedAt),
	}
}
