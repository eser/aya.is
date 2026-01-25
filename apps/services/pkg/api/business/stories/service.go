package stories

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

var (
	ErrFailedToGetRecord    = errors.New("failed to get record")
	ErrFailedToListRecords  = errors.New("failed to list records")
	ErrFailedToInsertRecord = errors.New("failed to insert record")
	ErrFailedToUpdateRecord = errors.New("failed to update record")
	ErrFailedToRemoveRecord = errors.New("failed to remove record")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrStoryNotFound        = errors.New("story not found")
	ErrInvalidSlugPrefix    = errors.New("slug must start with YYYYMMDD of publish date")
	ErrInvalidURI           = errors.New("invalid URI")
	ErrInvalidURIPrefix     = errors.New("URI must start with allowed prefix")
)

// Config holds the stories service configuration.
type Config struct {
	// AllowedURIPrefixes is a comma-separated list of allowed URI prefixes.
	AllowedURIPrefixes string `conf:"allowed_uri_prefixes" default:"https://objects.aya.is/"`
}

// GetAllowedURIPrefixes returns the allowed URI prefixes as a slice.
func (c *Config) GetAllowedURIPrefixes() []string {
	if c.AllowedURIPrefixes == "" {
		return nil
	}

	prefixes := strings.Split(c.AllowedURIPrefixes, ",")
	result := make([]string, 0, len(prefixes))

	for _, prefix := range prefixes {
		trimmed := strings.TrimSpace(prefix)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// slugValidationResult holds the result of slug format validation.
type slugValidationResult struct {
	Valid    bool
	Message  string
	Severity string
}

// validateSlugLength checks if the slug meets the minimum length requirement.
// Returns severity based on status: error for published, warning for draft.
func validateSlugLength(slug string, isPublished bool) *slugValidationResult {
	const minLength = 12 // 9 (YYYYMMDD-) + 3 (minimum content)

	if len(slug) >= minLength {
		return nil // Valid
	}

	severity := SeverityWarning
	if isPublished {
		severity = SeverityError
	}

	return &slugValidationResult{
		Valid:    false,
		Message:  fmt.Sprintf("Slug must be at least %d characters", minLength),
		Severity: severity,
	}
}

// validateSlugDatePrefix validates that the slug starts with YYYYMMDD- format
// matching the provided publish date.
// Returns severity based on status: error for published, warning for draft.
func validateSlugDatePrefix(slug string, publishDate time.Time, isPublished bool) *slugValidationResult {
	expectedPrefix := publishDate.Format("20060102") + "-"

	// Check if slug starts with the expected date prefix
	if strings.HasPrefix(slug, expectedPrefix) {
		return nil // Valid
	}

	severity := SeverityWarning
	if isPublished {
		severity = SeverityError
	}

	return &slugValidationResult{
		Valid:    false,
		Message:  "Slug must start with " + expectedPrefix,
		Severity: severity,
	}
}

// validateOptionalURL validates that a URL is either nil or a valid http/https URL.
func validateOptionalURL(uri *string) error {
	if uri == nil {
		return nil
	}

	parsedURL, err := url.ParseRequestURI(*uri)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidURI, *uri)
	}

	// Only accept http and https protocols
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("%w: URL must use http or https protocol: %s", ErrInvalidURI, *uri)
	}

	return nil
}

// validateURIPrefixes validates that a URI starts with one of the allowed prefixes.
// This is used to restrict non-admin users to only use URIs from our upload service.
func validateURIPrefixes(uri *string, allowedPrefixes []string) error {
	if uri == nil || *uri == "" {
		return nil
	}

	if len(allowedPrefixes) == 0 {
		return nil
	}

	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(*uri, prefix) {
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrInvalidURIPrefix, strings.Join(allowedPrefixes, ", "))
}

// Severity constants for slug availability results.
const (
	SeverityError   = "error"
	SeverityWarning = "warning"
)

// SlugAvailabilityResult holds the result of a slug availability check.
type SlugAvailabilityResult struct {
	Available bool   `json:"available"`
	Message   string `json:"message,omitempty"`
	Severity  string `json:"severity,omitempty"` // "error" | "warning" | ""
}

type Repository interface {
	GetProfileIDBySlug(ctx context.Context, slug string) (string, error)
	GetProfileByID(ctx context.Context, localeCode string, id string) (*profiles.Profile, error)
	GetStoryIDBySlug(ctx context.Context, slug string) (string, error)
	GetStoryIDBySlugIncludingDeleted(ctx context.Context, slug string) (string, error)
	GetStoryByID(
		ctx context.Context,
		localeCode string,
		id string,
		authorProfileID *string,
	) (*StoryWithChildren, error)
	ListStoriesOfPublication(
		ctx context.Context,
		localeCode string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*StoryWithChildren], error)
	// CRUD methods
	InsertStory(
		ctx context.Context,
		id string,
		authorProfileID string,
		slug string,
		kind string,
		status string,
		isFeatured bool,
		storyPictureURI *string,
		properties map[string]any,
		publishedAt *time.Time,
	) (*Story, error)
	InsertStoryTx(
		ctx context.Context,
		storyID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	InsertStoryPublication(
		ctx context.Context,
		id string,
		storyID string,
		profileID string,
		kind string,
		properties map[string]any,
	) error
	UpdateStory(
		ctx context.Context,
		id string,
		slug string,
		status string,
		isFeatured bool,
		storyPictureURI *string,
		publishedAt *time.Time,
	) error
	UpdateStoryTx(
		ctx context.Context,
		storyID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	UpsertStoryTx(
		ctx context.Context,
		storyID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	RemoveStory(ctx context.Context, id string) error
	GetStoryForEdit(ctx context.Context, localeCode string, id string) (*StoryForEdit, error)
	GetStoryOwnershipForUser(
		ctx context.Context,
		userID string,
		storyID string,
	) (*StoryOwnership, error)
}

type Service struct {
	logger      *logfx.Logger
	config      *Config
	repo        Repository
	idGenerator RecordIDGenerator
}

func NewService(logger *logfx.Logger, config *Config, repo Repository) *Service {
	return &Service{
		logger:      logger,
		config:      config,
		repo:        repo,
		idGenerator: DefaultIDGenerator,
	}
}

func (s *Service) GetByID(
	ctx context.Context,
	localeCode string,
	id string,
) (*StoryWithChildren, error) {
	record, err := s.repo.GetStoryByID(ctx, localeCode, id, nil)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return record, nil
}

func (s *Service) GetBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
) (*StoryWithChildren, error) {
	storyID, err := s.repo.GetStoryIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	record, err := s.repo.GetStoryByID(ctx, localeCode, storyID, nil)
	if err != nil {
		return nil, fmt.Errorf("%w(story_id: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	return record, nil
}

// CheckSlugAvailability checks if a story slug is available and validates format requirements.
// It optionally excludes a specific story ID (for edit scenarios).
// When status is "published" and publishedAt is provided, it also validates the date prefix.
// Parameters:
//   - slug: the slug to check
//   - excludeStoryID: story ID to exclude (for edit mode)
//   - status: "published" or "draft"
//   - publishedAt: publish date for date prefix validation
//   - includeDeleted: if true, also check against deleted stories
func (s *Service) CheckSlugAvailability(
	ctx context.Context,
	slug string,
	excludeStoryID *string,
	status string,
	publishedAt *time.Time,
	includeDeleted bool,
) (*SlugAvailabilityResult, error) {
	isPublished := status == "published"

	// Check minimum length (12 chars = 9 for date + 3 for content)
	if lengthResult := validateSlugLength(slug, isPublished); lengthResult != nil {
		// For errors, return immediately as unavailable
		if lengthResult.Severity == SeverityError {
			return &SlugAvailabilityResult{
				Available: false,
				Message:   lengthResult.Message,
				Severity:  lengthResult.Severity,
			}, nil
		}

		// For warnings, continue checking but remember to include the warning
		// We'll return this warning at the end if slug is otherwise available
	}

	// Validate slug date prefix
	var datePrefix time.Time
	if publishedAt != nil {
		datePrefix = *publishedAt
	} else {
		datePrefix = time.Now()
	}

	if prefixResult := validateSlugDatePrefix(slug, datePrefix, isPublished); prefixResult != nil {
		// For errors, return immediately as unavailable
		if prefixResult.Severity == SeverityError {
			return &SlugAvailabilityResult{
				Available: false,
				Message:   prefixResult.Message,
				Severity:  prefixResult.Severity,
			}, nil
		}

		// For warnings, continue checking but remember to include the warning
	}

	// Check if slug is already taken (active stories)
	storyID, err := s.repo.GetStoryIDBySlug(ctx, slug)
	if err == nil && storyID != "" {
		// If we're editing and the slug belongs to the same story, continue
		if excludeStoryID == nil || storyID != *excludeStoryID {
			return &SlugAvailabilityResult{
				Available: false,
				Message:   "This slug is already taken",
				Severity:  SeverityError,
			}, nil
		}
	}

	// Optionally check deleted stories
	if includeDeleted {
		deletedStoryID, err := s.repo.GetStoryIDBySlugIncludingDeleted(ctx, slug)
		if err == nil && deletedStoryID != "" {
			// If we're editing and it's the same story, that's fine
			if excludeStoryID == nil || deletedStoryID != *excludeStoryID {
				return &SlugAvailabilityResult{
					Available: false,
					Message:   "This slug was previously used",
					Severity:  SeverityError,
				}, nil
			}
		}
	}

	// Slug is available - but check if we have any warnings to return
	// Re-check for warnings
	if lengthResult := validateSlugLength(slug, isPublished); lengthResult != nil &&
		lengthResult.Severity == SeverityWarning {
		return &SlugAvailabilityResult{
			Available: true,
			Message:   lengthResult.Message,
			Severity:  lengthResult.Severity,
		}, nil
	}

	if prefixResult := validateSlugDatePrefix(slug, datePrefix, isPublished); prefixResult != nil &&
		prefixResult.Severity == SeverityWarning {
		return &SlugAvailabilityResult{
			Available: true,
			Message:   prefixResult.Message,
			Severity:  prefixResult.Severity,
		}, nil
	}

	return &SlugAvailabilityResult{
		Available: true,
	}, nil
}

func (s *Service) List(
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*StoryWithChildren], error) {
	records, err := s.repo.ListStoriesOfPublication(ctx, localeCode, cursor)
	if err != nil {
		return cursors.Cursored[[]*StoryWithChildren]{}, fmt.Errorf(
			"%w: %w",
			ErrFailedToListRecords,
			err,
		)
	}

	return records, nil
}

func (s *Service) ListByPublicationProfileSlug(
	ctx context.Context,
	localeCode string,
	publicationProfileSlug string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*StoryWithChildren], error) {
	publicationProfileID, err := s.repo.GetProfileIDBySlug(ctx, publicationProfileSlug)
	if err != nil {
		return cursors.Cursored[[]*StoryWithChildren]{}, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			publicationProfileSlug,
			err,
		)
	}

	cursor.Filters["publication_profile_id"] = publicationProfileID

	records, err := s.repo.ListStoriesOfPublication(
		ctx,
		localeCode,
		cursor,
	)
	if err != nil {
		return cursors.Cursored[[]*StoryWithChildren]{}, fmt.Errorf(
			"%w: %w",
			ErrFailedToListRecords,
			err,
		)
	}

	return records, nil
}

// CanUserEditStory checks if a user has permission to edit a story.
func (s *Service) CanUserEditStory(
	ctx context.Context,
	userID string,
	storyID string,
) (bool, error) {
	ownership, err := s.repo.GetStoryOwnershipForUser(ctx, userID, storyID)
	if err != nil {
		return false, fmt.Errorf(
			"%w(userID: %s, storyID: %s): %w",
			ErrFailedToGetRecord,
			userID,
			storyID,
			err,
		)
	}

	if ownership == nil {
		return false, nil
	}

	return ownership.CanEdit, nil
}

// GetForEdit retrieves a story for editing (raw content, no compilation).
func (s *Service) GetForEdit(
	ctx context.Context,
	localeCode string,
	storyID string,
) (*StoryForEdit, error) {
	story, err := s.repo.GetStoryForEdit(ctx, localeCode, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	return story, nil
}

// Create creates a new story with its translation and publication.
func (s *Service) Create(
	ctx context.Context,
	userID string,
	userKind string,
	authorProfileSlug string,
	publicationProfileSlug string,
	localeCode string,
	slug string,
	kind string,
	status string,
	title string,
	summary string,
	content string,
	storyPictureURI *string,
	isFeatured bool,
	publishedAt *time.Time,
) (*Story, error) {
	// Validate slug availability and date prefix (only block on errors, not warnings)
	slugResult, err := s.CheckSlugAvailability(ctx, slug, nil, status, publishedAt, false)
	if err != nil {
		return nil, err
	}

	if !slugResult.Available && slugResult.Severity == SeverityError {
		return nil, fmt.Errorf("%w: %s", ErrInvalidSlugPrefix, slugResult.Message)
	}

	// Validate story picture URI
	if err := validateOptionalURL(storyPictureURI); err != nil {
		return nil, err
	}

	// Non-admin users can only use URIs from our upload service
	if userKind != "admin" {
		err := validateURIPrefixes(storyPictureURI, s.config.GetAllowedURIPrefixes())
		if err != nil {
			return nil, err
		}
	}

	// Get author profile ID
	authorProfileID, err := s.repo.GetProfileIDBySlug(ctx, authorProfileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, authorProfileSlug, err)
	}

	// Get publication profile ID
	publicationProfileID, err := s.repo.GetProfileIDBySlug(ctx, publicationProfileSlug)
	if err != nil {
		return nil, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			publicationProfileSlug,
			err,
		)
	}

	// Generate new story ID
	storyID := s.idGenerator()

	// Create the main story record
	story, err := s.repo.InsertStory(
		ctx,
		string(storyID),
		authorProfileID,
		slug,
		kind,
		status,
		isFeatured,
		storyPictureURI,
		nil, // No additional properties for now
		publishedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	// Create the localized story data
	err = s.repo.InsertStoryTx(
		ctx,
		string(storyID),
		localeCode,
		title,
		summary,
		content,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: translation: %w", ErrFailedToInsertRecord, err)
	}

	// Create story publication (link story to profile)
	publicationID := s.idGenerator()

	err = s.repo.InsertStoryPublication(
		ctx,
		string(publicationID),
		string(storyID),
		publicationProfileID,
		"original", // publication kind
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: publication: %w", ErrFailedToInsertRecord, err)
	}

	return story, nil
}

// Update updates an existing story.
func (s *Service) Update(
	ctx context.Context,
	userID string,
	userKind string,
	storyID string,
	slug string,
	status string,
	isFeatured bool,
	storyPictureURI *string,
	publishedAt *time.Time,
) (*StoryForEdit, error) {
	// Check authorization
	canEdit, err := s.CanUserEditStory(ctx, userID, storyID)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit story %s",
			ErrUnauthorized,
			userID,
			storyID,
		)
	}

	// Validate slug availability and date prefix (only block on errors, not warnings)
	slugResult, err := s.CheckSlugAvailability(ctx, slug, &storyID, status, publishedAt, false)
	if err != nil {
		return nil, err
	}

	if !slugResult.Available && slugResult.Severity == SeverityError {
		return nil, fmt.Errorf("%w: %s", ErrInvalidSlugPrefix, slugResult.Message)
	}

	// Validate story picture URI
	if err := validateOptionalURL(storyPictureURI); err != nil {
		return nil, err
	}

	// Non-admin users can only use URIs from our upload service
	if userKind != "admin" {
		err := validateURIPrefixes(storyPictureURI, s.config.GetAllowedURIPrefixes())
		if err != nil {
			return nil, err
		}
	}

	// Update the story
	err = s.repo.UpdateStory(ctx, storyID, slug, status, isFeatured, storyPictureURI, publishedAt)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToUpdateRecord, storyID, err)
	}

	// Return updated story
	story, err := s.repo.GetStoryForEdit(ctx, "en", storyID)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	return story, nil
}

// UpdateTranslation updates story translation fields.
func (s *Service) UpdateTranslation(
	ctx context.Context,
	userID string,
	storyID string,
	localeCode string,
	title string,
	summary string,
	content string,
) error {
	// Check authorization
	canEdit, err := s.CanUserEditStory(ctx, userID, storyID)
	if err != nil {
		return err
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot edit story %s",
			ErrUnauthorized,
			userID,
			storyID,
		)
	}

	// Update the translation (use upsert to handle new locales)
	err = s.repo.UpsertStoryTx(ctx, storyID, localeCode, title, summary, content)
	if err != nil {
		return fmt.Errorf(
			"%w(storyID: %s, locale: %s): %w",
			ErrFailedToUpdateRecord,
			storyID,
			localeCode,
			err,
		)
	}

	return nil
}

// Delete soft-deletes a story.
func (s *Service) Delete(
	ctx context.Context,
	userID string,
	storyID string,
) error {
	// Check authorization
	canEdit, err := s.CanUserEditStory(ctx, userID, storyID)
	if err != nil {
		return err
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot delete story %s",
			ErrUnauthorized,
			userID,
			storyID,
		)
	}

	// Delete the story
	err = s.repo.RemoveStory(ctx, storyID)
	if err != nil {
		return fmt.Errorf("%w(storyID: %s): %w", ErrFailedToRemoveRecord, storyID, err)
	}

	return nil
}
