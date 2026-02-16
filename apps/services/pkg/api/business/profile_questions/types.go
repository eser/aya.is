package profile_questions

import (
	"time"
)

// Content length constraints.
const (
	MinContentLength = 10
	MaxContentLength = 2000
)

// Question represents a Q&A question on a profile.
type Question struct {
	ID                  string     `json:"id"`
	ProfileID           string     `json:"profile_id"`
	Content             string     `json:"content"`
	AuthorProfileID     *string    `json:"author_profile_id"`
	AuthorProfileSlug   *string    `json:"author_profile_slug"`
	AuthorProfileTitle  *string    `json:"author_profile_title"`
	AnswerContent       *string    `json:"answer_content"`
	AnswerURI           *string    `json:"answer_uri"`
	AnswerKind          *string    `json:"answer_kind"`
	AnsweredAt          *time.Time `json:"answered_at"`
	AnsweredByProfileID *string    `json:"answered_by_profile_id"`
	VoteCount           int        `json:"vote_count"`
	IsAnonymous         bool       `json:"is_anonymous"`
	IsHidden            bool       `json:"is_hidden"`
	HasViewerVote       bool       `json:"has_viewer_vote"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
}

// Vote represents a user's vote on a question.
type Vote struct {
	ID         string    `json:"id"`
	QuestionID string    `json:"question_id"`
	UserID     string    `json:"user_id"`
	Score      int       `json:"score"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateQuestionParams holds parameters for creating a new question.
type CreateQuestionParams struct {
	ProfileSlug string
	UserID      string
	Content     string
	IsAnonymous bool
}

// AnswerQuestionParams holds parameters for answering a question.
type AnswerQuestionParams struct {
	ProfileSlug       string
	QuestionID        string
	UserID            string
	AnswererProfileID string
	AnswerContent     string
	AnswerURI         *string
	AnswerKind        *string
}

// VoteParams holds parameters for toggling a vote on a question.
type VoteParams struct {
	ProfileSlug string
	QuestionID  string
	UserID      string
}

// HideQuestionParams holds parameters for hiding/showing a question.
type HideQuestionParams struct {
	ProfileSlug string
	QuestionID  string
	UserID      string
	IsHidden    bool
}

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string
