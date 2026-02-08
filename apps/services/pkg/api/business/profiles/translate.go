package profiles

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
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

// AutoTranslatePageParams holds the parameters for auto-translating a profile page.
type AutoTranslatePageParams struct {
	UserID              string
	UserKind            string
	IndividualProfileID string
	ProfileSlug         string
	PageID              string
	SourceLocale        string
	TargetLocale        string
}

// AutoTranslateProfilePage orchestrates the full auto-translate workflow for profile pages:
// check permissions, deduct points, get source content, translate via AI, and save.
func (s *Service) AutoTranslateProfilePage(
	ctx context.Context,
	params AutoTranslatePageParams,
	translator ContentTranslator,
	pointsService *profile_points.Service,
) error {
	// Check authorization
	canEdit, permErr := s.HasUserAccessToProfile(
		ctx,
		params.UserID,
		params.ProfileSlug,
		MembershipKindMaintainer,
	)
	if permErr != nil {
		return fmt.Errorf("failed to check permissions: %w", permErr)
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			params.UserID,
			params.ProfileSlug,
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
	title, summary, content, err := s.GetProfilePageTranslationContent(
		ctx,
		params.ProfileSlug,
		params.PageID,
		params.SourceLocale,
	)
	if err != nil {
		return fmt.Errorf("failed to get source content: %w", err)
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
	err = s.UpdateProfilePageTranslation(
		ctx,
		params.UserID,
		params.UserKind,
		params.ProfileSlug,
		params.PageID,
		params.TargetLocale,
		translatedTitle,
		translatedSummary,
		translatedContent,
	)
	if err != nil {
		return fmt.Errorf("failed to save translated content: %w", err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfilePageAutoTranslated,
		EntityType: "profile_page",
		EntityID:   params.PageID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"source_locale": params.SourceLocale,
			"target_locale": params.TargetLocale,
		},
	})

	return nil
}
