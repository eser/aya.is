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
	// Check authorization
	canEdit, err := s.CanUserEditStory(ctx, params.UserID, params.StoryID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnauthorized, err)
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot edit story %s",
			ErrUnauthorized,
			params.UserID,
			params.StoryID,
		)
	}

	// Deduct points for auto-translation
	eventAutoTranslate := profile_points.EventAutoTranslate

	_, spendErr := pointsService.SpendPoints(ctx, profile_points.SpendParams{
		ActorID:         params.UserID,
		TargetProfileID: params.IndividualProfileID,
		Amount:          profile_points.CostAutoTranslate,
		TriggeringEvent: &eventAutoTranslate,
		Description:     "Auto-translate content",
	})
	if spendErr != nil {
		return spendErr //nolint:wrapcheck
	}

	// Get source content
	title, summary, content, err := s.GetTranslationContent(
		ctx,
		params.StoryID,
		params.SourceLocale,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetSourceContent, err)
	}

	// Translate via AI
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

	// Save translated content
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
		Payload: map[string]any{
			"source_locale": params.SourceLocale,
			"target_locale": params.TargetLocale,
		},
	})

	return nil
}
