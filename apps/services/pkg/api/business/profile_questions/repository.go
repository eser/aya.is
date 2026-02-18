package profile_questions

import (
	"context"

	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

// Repository defines the storage operations for profile questions (port).
type Repository interface {
	// GetQAVisibility returns the Q&A module visibility for a profile.
	GetQAVisibility(ctx context.Context, profileID string) (string, error)

	// GetQuestion returns a single question by ID.
	GetQuestion(ctx context.Context, id string) (*Question, error)

	// ListQuestions returns questions for a profile with pagination.
	ListQuestions(
		ctx context.Context,
		profileID string,
		localeCode string,
		viewerUserID *string,
		includeHidden bool,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*Question], error)

	// InsertQuestion creates a new question record.
	InsertQuestion(
		ctx context.Context,
		id string,
		profileID string,
		authorUserID string,
		content string,
		isAnonymous bool,
	) (*Question, error)

	// SetAnswer sets the initial answer on a question.
	SetAnswer(
		ctx context.Context,
		questionID string,
		answerContent string,
		answerURI *string,
		answerKind *string,
		answeredByProfileID string,
	) error

	// EditAnswer updates only the answer content, preserving the original answerer.
	EditAnswer(
		ctx context.Context,
		questionID string,
		answerContent string,
		answerURI *string,
		answerKind *string,
	) error

	// UpdateHidden toggles the hidden state of a question.
	UpdateHidden(ctx context.Context, questionID string, isHidden bool) error

	// InsertVote creates a new vote record.
	InsertVote(ctx context.Context, voteID string, questionID string, userID string) (*Vote, error)

	// DeleteVote removes a vote record.
	DeleteVote(ctx context.Context, questionID string, userID string) error

	// GetVote returns a vote by question and user.
	GetVote(ctx context.Context, questionID string, userID string) (*Vote, error)

	// GetProfileIDBySlug resolves a profile slug to its ID.
	GetProfileIDBySlug(ctx context.Context, slug string) (string, error)
}
