package stories

import (
	"context"
	"errors"
	"fmt"

	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
)

// Sentinel errors for auto-translation.
var (
	ErrTranslationNotAvailable       = errors.New("AI translation not available")
	ErrNoIndividualProfile           = errors.New("user has no individual profile")
	ErrFailedToGetSourceContent      = errors.New("failed to get source content")
	ErrFailedToSaveTranslatedContent = errors.New("failed to save translated content")
)

// ContentTranslator defines the interface for AI-powered content translation.
// Implementations live in the adapter layer (e.g., HTTP adapter using aifx).
type ContentTranslator interface {
	Translate(
		ctx context.Context,
		sourceLocale, targetLocale string,
		title, summary, content string,
	) (translatedTitle, translatedSummary, translatedContent string, err error)
}

// AutoTranslateStoryParams holds the parameters for auto-translating a story.
type AutoTranslateStoryParams struct {
	UserID              string
	IndividualProfileID string
	StoryID             string
	SourceLocale        string
	TargetLocale        string
}

// AutoTranslateStory orchestrates the full auto-translate workflow:
// check permissions, deduct points, get source content, translate via AI, and save.
func (s *Service) AutoTranslateStory(
	ctx context.Context,
	params AutoTranslateStoryParams,
	translator ContentTranslator,
	pointsService *profile_points.Service,
) error {
	authErr := s.authorizeStoryEdit(ctx, params.UserID, params.StoryID)
	if authErr != nil {
		return authErr
	}

	// Deduct points for auto-translation
	spendErr := s.deductTranslationPoints(ctx, pointsService, params)
	if spendErr != nil {
		return spendErr
	}

	// Translate and save
	return s.translateAndSave(ctx, translator, params)
}

// authorizeStoryEdit checks that a user can edit the given story.
func (s *Service) authorizeStoryEdit(ctx context.Context, userID string, storyID string) error {
	canEdit, err := s.CanUserEditStory(ctx, userID, storyID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnauthorized, err)
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot edit story %s",
			ErrUnauthorized,
			userID,
			storyID,
		)
	}

	return nil
}

// deductTranslationPoints deducts points for auto-translation.
func (s *Service) deductTranslationPoints(
	ctx context.Context,
	pointsService *profile_points.Service,
	params AutoTranslateStoryParams,
) error {
	eventAutoTranslate := profile_points.EventAutoTranslate

	_, err := pointsService.SpendPoints(ctx, profile_points.SpendParams{
		ActorID:         params.UserID,
		TargetProfileID: params.IndividualProfileID,
		Amount:          profile_points.CostAutoTranslate,
		TriggeringEvent: &eventAutoTranslate,
		Description:     "Auto-translate content",
	})

	return err //nolint:wrapcheck
}

// translateAndSave gets source content, translates via AI, and saves the result.
func (s *Service) translateAndSave(
	ctx context.Context,
	translator ContentTranslator,
	params AutoTranslateStoryParams,
) error {
	title, summary, content, err := s.GetTranslationContent(
		ctx,
		params.StoryID,
		params.SourceLocale,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetSourceContent, err)
	}

	translatedTitle, translatedSummary, translatedContent, err := translator.Translate(
		ctx,
		params.SourceLocale,
		params.TargetLocale,
		title,
		summary,
		content,
	)
	if err != nil {
		return err //nolint:wrapcheck
	}

	err = s.UpdateTranslation(
		ctx,
		params.UserID,
		params.StoryID,
		params.TargetLocale,
		translatedTitle,
		translatedSummary,
		translatedContent,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSaveTranslatedContent, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryAutoTranslated,
		EntityType: "story",
		EntityID:   params.StoryID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"source_locale": params.SourceLocale,
			"target_locale": params.TargetLocale,
		},
	})

	return nil
}
