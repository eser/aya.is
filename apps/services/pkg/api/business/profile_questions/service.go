package profile_questions

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func DefaultIDGenerator() string {
	return lib.IDsGenerateUnique()
}

// Service provides profile question operations.
type Service struct {
	logger         *logfx.Logger
	repo           Repository
	profileService *profiles.Service
	auditService   *events.AuditService
	idGenerator    IDGenerator
}

// NewService creates a new profile questions service.
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

// ListQuestions returns questions for a profile, stripping anonymous author info.
func (s *Service) ListQuestions(
	ctx context.Context,
	profileSlug string,
	viewerUserID *string,
	includeHidden bool,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*Question], error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return cursors.Cursored[[]*Question]{}, fmt.Errorf(
			"%w (slug: %s): %w",
			ErrFailedToGetRecord,
			profileSlug,
			err,
		)
	}

	hidden, err := s.repo.IsQAHidden(ctx, profileID)
	if err != nil {
		return cursors.Cursored[[]*Question]{}, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if hidden {
		return cursors.Cursored[[]*Question]{}, ErrQANotEnabled
	}

	result, err := s.repo.ListQuestions(
		ctx,
		profileID,
		viewerUserID,
		includeHidden,
		cursor,
	)
	if err != nil {
		return cursors.Cursored[[]*Question]{}, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	// Strip anonymous author information
	for _, q := range result.Data {
		if q.IsAnonymous {
			q.AuthorProfileID = nil
			q.AuthorProfileSlug = nil
		}
	}

	return result, nil
}

// CreateQuestion creates a new question on a profile.
func (s *Service) CreateQuestion(
	ctx context.Context,
	params CreateQuestionParams,
) (*Question, error) {
	contentLen := len([]rune(params.Content))
	if contentLen < MinContentLength {
		return nil, ErrContentTooShort
	}

	if contentLen > MaxContentLength {
		return nil, ErrContentTooLong
	}

	profileID, err := s.repo.GetProfileIDBySlug(ctx, params.ProfileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w (slug: %s): %w", ErrFailedToGetRecord, params.ProfileSlug, err)
	}

	hidden, err := s.repo.IsQAHidden(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if hidden {
		return nil, ErrQANotEnabled
	}

	id := s.idGenerator()

	question, err := s.repo.InsertQuestion(
		ctx,
		id,
		profileID,
		params.UserID,
		params.Content,
		params.IsAnonymous,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	// Strip anonymous author information from the response
	if question.IsAnonymous {
		question.AuthorProfileID = nil
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileQuestionCreated,
		EntityType: "profile_question",
		EntityID:   id,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_id":   profileID,
			"is_anonymous": params.IsAnonymous,
		},
	})

	return question, nil
}

// AnswerQuestion adds an answer to a question (requires contributor+ access).
func (s *Service) AnswerQuestion(ctx context.Context, params AnswerQuestionParams) error {
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

	question, err := s.repo.GetQuestion(ctx, params.QuestionID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrQuestionNotFound, err)
	}

	if question.AnswerContent != nil {
		return ErrQuestionAlreadyAnswered
	}

	err = s.repo.UpdateAnswer(
		ctx,
		params.QuestionID,
		params.AnswerContent,
		params.AnswerURI,
		params.AnswerKind,
		params.AnswererProfileID,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileQuestionAnswered,
		EntityType: "profile_question",
		EntityID:   params.QuestionID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_slug": params.ProfileSlug,
		},
	})

	return nil
}

// EditAnswer updates an existing answer on a question.
// Contributors can only edit answers they authored. Maintainers can edit all answers.
func (s *Service) EditAnswer(ctx context.Context, params AnswerQuestionParams) error {
	// Check minimum access: contributor+
	hasContributorAccess, err := s.profileService.HasUserAccessToProfile(
		ctx,
		params.UserID,
		params.ProfileSlug,
		profiles.MembershipKindContributor,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if !hasContributorAccess {
		return ErrInsufficientPermission
	}

	question, err := s.repo.GetQuestion(ctx, params.QuestionID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrQuestionNotFound, err)
	}

	if question.AnswerContent == nil {
		return ErrQuestionNotAnswered
	}

	// Contributors can only edit their own answers; maintainers can edit all
	if question.AnsweredByProfileID == nil ||
		*question.AnsweredByProfileID != params.AnswererProfileID {
		hasMaintainerAccess, maintainerErr := s.profileService.HasUserAccessToProfile(
			ctx,
			params.UserID,
			params.ProfileSlug,
			profiles.MembershipKindMaintainer,
		)
		if maintainerErr != nil {
			return fmt.Errorf("%w: %w", ErrFailedToGetRecord, maintainerErr)
		}

		if !hasMaintainerAccess {
			return ErrInsufficientPermission
		}
	}

	err = s.repo.UpdateAnswer(
		ctx,
		params.QuestionID,
		params.AnswerContent,
		params.AnswerURI,
		params.AnswerKind,
		params.AnswererProfileID,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileQuestionAnswerEdited,
		EntityType: "profile_question",
		EntityID:   params.QuestionID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_slug": params.ProfileSlug,
		},
	})

	return nil
}

// ToggleVote toggles a vote on a question. Returns true if a vote was added, false if removed.
func (s *Service) ToggleVote(ctx context.Context, params VoteParams) (bool, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, params.ProfileSlug)
	if err != nil {
		return false, fmt.Errorf("%w (slug: %s): %w", ErrFailedToGetRecord, params.ProfileSlug, err)
	}

	hidden, err := s.repo.IsQAHidden(ctx, profileID)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if hidden {
		return false, ErrQANotEnabled
	}

	_, err = s.repo.GetQuestion(ctx, params.QuestionID)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrQuestionNotFound, err)
	}

	existing, err := s.repo.GetVote(ctx, params.QuestionID, params.UserID)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if existing != nil {
		err = s.repo.DeleteVote(ctx, params.QuestionID, params.UserID)
		if err != nil {
			return false, fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
		}

		s.auditService.Record(ctx, events.AuditParams{
			EventType:  events.ProfileQuestionVoted,
			EntityType: "profile_question",
			EntityID:   params.QuestionID,
			ActorID:    &params.UserID,
			ActorKind:  events.ActorUser,
			Payload:    map[string]any{"action": "unvote"},
		})

		return false, nil
	}

	voteID := s.idGenerator()

	_, err = s.repo.InsertVote(ctx, voteID, params.QuestionID, params.UserID)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileQuestionVoted,
		EntityType: "profile_question",
		EntityID:   params.QuestionID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"action": "vote"},
	})

	return true, nil
}

// HideQuestion toggles the hidden state of a question (requires maintainer+ access).
func (s *Service) HideQuestion(ctx context.Context, params HideQuestionParams) error {
	hasAccess, err := s.profileService.HasUserAccessToProfile(
		ctx,
		params.UserID,
		params.ProfileSlug,
		profiles.MembershipKindMaintainer,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if !hasAccess {
		return ErrInsufficientPermission
	}

	err = s.repo.UpdateHidden(ctx, params.QuestionID, params.IsHidden)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileQuestionHidden,
		EntityType: "profile_question",
		EntityID:   params.QuestionID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_slug": params.ProfileSlug,
			"is_hidden":    params.IsHidden,
		},
	})

	return nil
}
