package profiles

import (
	"context"
	"errors"
	"fmt"

	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
)

// Sentinel errors for AI content generation.
var (
	ErrFailedToGetProfileData  = errors.New("failed to get profile data")
	ErrFailedToGenerateContent = errors.New("failed to generate content")
	ErrFailedToCreatePage      = errors.New("failed to create generated page")
	ErrCVPageAlreadyExists     = errors.New("a page with slug 'cv' already exists")
	ErrNoLinkedInLinkFound     = errors.New("no LinkedIn link found on this profile")
)

// ContentGenerator defines the interface for AI-powered content generation.
// Implementations live in the adapter layer (e.g., HTTP adapter using aifx).
type ContentGenerator interface {
	GenerateCV(
		ctx context.Context,
		locale string,
		profileTitle string,
		profileDescription string,
		linkedInURL string,
		links []*ProfileLinkBrief,
		contributions []*ProfileMembership,
	) (title, summary, content string, err error)
}

// GenerateCVPageParams holds the parameters for AI-powered CV page generation.
type GenerateCVPageParams struct {
	UserID              string
	UserKind            string
	IndividualProfileID string
	ProfileSlug         string
	Locale              string
}

// GenerateCVPage orchestrates the full AI CV generation workflow:
// check permissions, deduct points, gather profile data, generate via AI, and create page.
func (s *Service) GenerateCVPage(
	ctx context.Context,
	params GenerateCVPageParams,
	generator ContentGenerator,
	pointsService *profile_points.Service,
) (*ProfilePage, error) {
	// Check authorization
	canEdit, permErr := s.HasUserAccessToProfile(
		ctx,
		params.UserID,
		params.ProfileSlug,
		MembershipKindMaintainer,
	)
	if permErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCheckPermissions, permErr)
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			params.UserID,
			params.ProfileSlug,
		)
	}

	// Check if "cv" slug already exists
	slugResult, slugErr := s.CheckPageSlugAvailability(
		ctx,
		params.Locale,
		params.ProfileSlug,
		"cv",
		nil,
		false,
	)
	if slugErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetProfileData, slugErr)
	}

	if !slugResult.Available && slugResult.Severity == SeverityError {
		return nil, ErrCVPageAlreadyExists
	}

	// Deduct points for content generation
	eventGenerateContent := profile_points.EventGenerateContent

	_, spendErr := pointsService.SpendPoints(ctx, profile_points.SpendParams{
		ActorID:         params.UserID,
		TargetProfileID: params.IndividualProfileID,
		Amount:          profile_points.CostGenerateContent,
		TriggeringEvent: &eventGenerateContent,
		Description:     "Generate CV page from profile data",
	})
	if spendErr != nil {
		return nil, spendErr //nolint:wrapcheck
	}

	// Fetch profile data (title, description, links)
	profileData, profileErr := s.GetBySlugEx(ctx, params.Locale, params.ProfileSlug)
	if profileErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetProfileData, profileErr)
	}

	// Fetch contributions (organizations the user is part of)
	contributions, contribErr := s.ListProfileContributionsBySlug(
		ctx,
		params.Locale,
		params.ProfileSlug,
		nil,
	)
	if contribErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetProfileData, contribErr)
	}

	// Extract LinkedIn URL from links
	linkedInURL := ""

	for _, link := range profileData.Links {
		if link.Kind == "linkedin" && link.URI != "" {
			linkedInURL = link.URI

			break
		}
	}

	if linkedInURL == "" {
		return nil, ErrNoLinkedInLinkFound
	}

	// Generate CV content via AI
	title, summary, content, genErr := generator.GenerateCV(
		ctx,
		params.Locale,
		profileData.Title,
		profileData.Description,
		linkedInURL,
		profileData.Links,
		contributions.Data,
	)
	if genErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGenerateContent, genErr)
	}

	// Create the page
	page, createErr := s.CreateProfilePage(
		ctx,
		params.UserID,
		params.UserKind,
		params.ProfileSlug,
		"cv",
		params.Locale,
		title,
		summary,
		content,
		nil,
		nil,
	)
	if createErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreatePage, createErr)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfilePageAIGenerated,
		EntityType: "profile_page",
		EntityID:   page.ID,
		ActorID:    &params.UserID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"locale":       params.Locale,
			"generator":    "cv_from_linkedin",
			"linkedin_url": linkedInURL,
		},
	})

	return page, nil
}
