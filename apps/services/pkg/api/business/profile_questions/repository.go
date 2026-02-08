package profile_questions

import (
	"context"

	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

// Repository defines the storage operations for profile questions (port).
type Repository interface {
	// IsQAHidden checks whether Q&A is hidden for a profile.
	IsQAHidden(ctx context.Context, profileID string) (bool, error)

	// GetQuestion returns a single question by ID.
	GetQuestion(ctx context.Context, id string) (*Question, error)

	// ListQuestions returns questions for a profile with pagination.
	ListQuestions(
		ctx context.Context,
		profileID string,
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

	// UpdateAnswer sets the answer content on a question.
	UpdateAnswer(
		ctx context.Context,
		questionID string,
		answerContent string,
		answeredByUserID string,
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
