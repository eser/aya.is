package stories

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

var (
	ErrFailedToGetRecord     = errors.New("failed to get record")
	ErrFailedToListRecords   = errors.New("failed to list records")
	ErrFailedToInsertRecord  = errors.New("failed to insert record")
	ErrFailedToUpdateRecord  = errors.New("failed to update record")
	ErrFailedToRemoveRecord  = errors.New("failed to remove record")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrStoryNotFound         = errors.New("story not found")
	ErrInvalidSlugPrefix     = errors.New("slug must start with YYYYMMDD of publish date")
	ErrInvalidURI            = errors.New("invalid URI")
	ErrInvalidURIPrefix      = errors.New("URI must start with allowed prefix")
	ErrHasActivePublications = errors.New(
		"story has active publications, unpublish from all profiles first",
	)
	ErrNoProfileAccess = errors.New(
		"user does not have membership access to the target profile",
	)
	ErrInsufficientProfileRole = errors.New(
		"user does not have sufficient role for this operation",
	)
	ErrManagedStory = errors.New(
		"this story is managed by an external sync and cannot be edited directly",
	)
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
// Returns severity based on publication state: error for published, warning for draft.
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
// Returns severity based on publication state: error for published, warning for draft.
func validateSlugDatePrefix(
	slug string,
	publishDate time.Time,
	isPublished bool,
) *slugValidationResult {
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
	GetProfileByID(
		ctx context.Context,
		localeCode string,
		id string,
	) (*profiles.Profile, error)
	GetStoryIDBySlug(ctx context.Context, slug string) (string, error)
	GetStoryIDBySlugForViewer(
		ctx context.Context,
		slug string,
		viewerUserID *string,
	) (string, error)
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
	ListStoriesOfPublicationForViewer(
		ctx context.Context,
		localeCode string,
		cursor *cursors.Cursor,
		viewerUserID *string,
	) (cursors.Cursored[[]*StoryWithChildren], error)
	ListStoriesByAuthorProfileID(
		ctx context.Context,
		localeCode string,
		authorProfileID string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*StoryWithChildren], error)
	ListStoriesByAuthorProfileIDForViewer(
		ctx context.Context,
		localeCode string,
		authorProfileID string,
		cursor *cursors.Cursor,
		viewerUserID *string,
	) (cursors.Cursored[[]*StoryWithChildren], error)
	// Story CRUD methods
	InsertStory(
		ctx context.Context,
		id string,
		authorProfileID string,
		slug string,
		kind string,
		storyPictureURI *string,
		properties map[string]any,
		isManaged bool,
		remoteID *string,
		visibility string,
		featDiscussions bool,
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
		isFeatured bool,
		publishedAt *time.Time,
		properties map[string]any,
	) error
	UpdateStory(
		ctx context.Context,
		id string,
		slug string,
		storyPictureURI *string,
		properties map[string]any,
		visibility string,
		featDiscussions *bool,
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
	DeleteStoryTx(ctx context.Context, storyID string, localeCode string) error
	ListStoryTxLocales(ctx context.Context, storyID string) ([]string, error)
	GetStoryForEdit(ctx context.Context, localeCode string, id string) (*StoryForEdit, error)
	GetStoryOwnershipForUser(
		ctx context.Context,
		userID string,
		storyID string,
	) (*StoryOwnership, error)
	// Publication management methods
	ListStoryPublications(
		ctx context.Context,
		localeCode string,
		storyID string,
	) ([]*StoryPublication, error)
	GetStoryPublicationProfileID(ctx context.Context, publicationID string) (string, error)
	UpdateStoryPublication(
		ctx context.Context,
		id string,
		isFeatured bool,
	) error
	RemoveStoryPublication(ctx context.Context, id string) error
	CountStoryPublications(ctx context.Context, storyID string) (int64, error)
	GetStoryFirstPublishedAt(ctx context.Context, storyID string) (*time.Time, error)
	GetUserMembershipForProfile(
		ctx context.Context,
		userID string,
		profileID string,
	) (string, error)
	InvalidateStorySlugCache(ctx context.Context, slug string) error
	ListActivityStories(
		ctx context.Context,
		localeCode string,
		filterAuthorProfileID *string,
	) ([]*StoryWithChildren, error)
}

type Service struct {
	logger       *logfx.Logger
	config       *Config
	repo         Repository
	auditService *events.AuditService
	idGenerator  RecordIDGenerator
}

func NewService(
	logger *logfx.Logger,
	config *Config,
	repo Repository,
	auditService *events.AuditService,
) *Service {
	return &Service{
		logger:       logger,
		config:       config,
		repo:         repo,
		auditService: auditService,
		idGenerator:  DefaultIDGenerator,
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

	if storyID == "" {
		return nil, nil //nolint:nilnil
	}

	record, err := s.repo.GetStoryByID(ctx, localeCode, storyID, nil)
	if err != nil {
		return nil, fmt.Errorf("%w(story_id: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	return record, nil
}

// GetBySlugForViewer returns a story by slug, respecting viewer permissions.
// - Stories with active publications are visible to everyone
// - Other stories are visible to admins, authors, and profile editors.
func (s *Service) GetBySlugForViewer(
	ctx context.Context,
	localeCode string,
	slug string,
	viewerUserID *string,
) (*StoryWithChildren, error) {
	storyID, err := s.repo.GetStoryIDBySlugForViewer(ctx, slug, viewerUserID)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	if storyID == "" {
		return nil, nil //nolint:nilnil
	}

	record, err := s.repo.GetStoryByID(ctx, localeCode, storyID, nil)
	if err != nil {
		return nil, fmt.Errorf("%w(story_id: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	return record, nil
}

// ResolveStorySlug resolves a story slug to its ID.
func (s *Service) ResolveStorySlug(ctx context.Context, slug string) (string, error) {
	storyID, err := s.repo.GetStoryIDBySlugIncludingDeleted(ctx, slug)
	if err != nil {
		return "", fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	return storyID, nil
}

// CheckSlugAvailability checks if a story slug is available and validates format requirements.
// It optionally excludes a specific story ID (for edit scenarios).
// When the story has publications, it validates the date prefix strictly.
// Parameters:
//   - slug: the slug to check
//   - excludeStoryID: story ID to exclude (for edit mode)
//   - storyID: optional story ID to check publication state
//   - publishedAt: publish date for date prefix validation
//   - includeDeleted: if true, also check against deleted stories
func (s *Service) CheckSlugAvailability(
	ctx context.Context,
	slug string,
	excludeStoryID *string,
	storyID *string,
	publishedAt *time.Time,
	includeDeleted bool,
) (*SlugAvailabilityResult, error) {
	// Determine if story is published by checking its publications
	isPublished := false

	if storyID != nil {
		count, err := s.repo.CountStoryPublications(ctx, *storyID)
		if err == nil && count > 0 {
			isPublished = true
		}
	}

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
	} else if storyID != nil {
		// Try to get the first published_at from publications
		firstPublishedAt, err := s.repo.GetStoryFirstPublishedAt(ctx, *storyID)
		if err == nil && firstPublishedAt != nil {
			datePrefix = *firstPublishedAt
		} else {
			datePrefix = time.Now()
		}
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
	existingStoryID, err := s.repo.GetStoryIDBySlug(ctx, slug)
	if err == nil && existingStoryID != "" {
		// If we're editing and the slug belongs to the same story, continue
		if excludeStoryID == nil || existingStoryID != *excludeStoryID {
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

func (s *Service) ListByAuthorProfileSlug(
	ctx context.Context,
	localeCode string,
	authorProfileSlug string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*StoryWithChildren], error) {
	authorProfileID, err := s.repo.GetProfileIDBySlug(ctx, authorProfileSlug)
	if err != nil {
		return cursors.Cursored[[]*StoryWithChildren]{}, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			authorProfileSlug,
			err,
		)
	}

	records, err := s.repo.ListStoriesByAuthorProfileID(
		ctx,
		localeCode,
		authorProfileID,
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

func (s *Service) ListByPublicationProfileSlugForViewer(
	ctx context.Context,
	localeCode string,
	publicationProfileSlug string,
	cursor *cursors.Cursor,
	viewerUserID *string,
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

	records, err := s.repo.ListStoriesOfPublicationForViewer(
		ctx,
		localeCode,
		cursor,
		viewerUserID,
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

func (s *Service) ListByAuthorProfileSlugForViewer(
	ctx context.Context,
	localeCode string,
	authorProfileSlug string,
	cursor *cursors.Cursor,
	viewerUserID *string,
) (cursors.Cursored[[]*StoryWithChildren], error) {
	authorProfileID, err := s.repo.GetProfileIDBySlug(ctx, authorProfileSlug)
	if err != nil {
		return cursors.Cursored[[]*StoryWithChildren]{}, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			authorProfileSlug,
			err,
		)
	}

	records, err := s.repo.ListStoriesByAuthorProfileIDForViewer(
		ctx,
		localeCode,
		authorProfileID,
		cursor,
		viewerUserID,
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

// GetForEdit retrieves a story for editing with its publications.
func (s *Service) GetForEdit(
	ctx context.Context,
	localeCode string,
	storyID string,
) (*StoryForEditWithPublications, error) {
	story, err := s.repo.GetStoryForEdit(ctx, localeCode, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	if story == nil {
		return nil, nil //nolint:nilnil
	}

	publications, err := s.repo.ListStoryPublications(ctx, localeCode, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToListRecords, storyID, err)
	}

	isFallback := strings.TrimRight(story.LocaleCode, " ") != localeCode

	return &StoryForEditWithPublications{
		StoryForEdit: story,
		IsFallback:   isFallback,
		Publications: publications,
	}, nil
}

// Create creates a new story with its translation. Optionally publishes to profiles.
func (s *Service) Create(
	ctx context.Context,
	userID string,
	userKind string,
	authorProfileSlug string,
	localeCode string,
	slug string,
	kind string,
	title string,
	summary string,
	content string,
	storyPictureURI *string,
	publishToProfileSlugs []string,
	properties map[string]any,
	visibility string,
) (*Story, error) {
	// Determine if the story will be published (for slug validation)
	isPublishing := len(publishToProfileSlugs) > 0

	status := "draft"

	if isPublishing {
		status = "published"
	}

	// Validate slug availability and date prefix (only block on errors, not warnings)
	var publishedAt *time.Time

	if isPublishing {
		now := time.Now()
		publishedAt = &now
	}

	slugResult, err := s.CheckSlugAvailability(ctx, slug, nil, nil, publishedAt, false)
	if err != nil {
		return nil, err
	}

	_ = status // used conceptually for slug validation

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

	// Fetch author profile to read discussion default
	authorProfile, err := s.repo.GetProfileByID(ctx, localeCode, authorProfileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, authorProfileID, err)
	}

	featDiscussions := false
	if authorProfile != nil {
		featDiscussions = authorProfile.OptionStoryDiscussionsByDefault
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
		storyPictureURI,
		properties,
		false,
		nil,
		visibility,
		featDiscussions,
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

	// Create publications for each target profile
	for _, profileSlug := range publishToProfileSlugs {
		profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
		if err != nil {
			return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
		}

		// Verify user has membership access to the target profile (contributor+)
		membershipKind, err := s.repo.GetUserMembershipForProfile(ctx, userID, profileID)
		if err != nil {
			return nil, fmt.Errorf(
				"%w(userID: %s, profileID: %s): %w",
				ErrFailedToGetRecord,
				userID,
				profileID,
				err,
			)
		}

		publishRoles := map[string]bool{
			"admin": true, "owner": true, "lead": true, "maintainer": true, "contributor": true,
		}
		if !publishRoles[membershipKind] {
			return nil, fmt.Errorf(
				"%w: user %s has no publish access to profile %s",
				ErrNoProfileAccess,
				userID,
				profileID,
			)
		}

		publicationID := s.idGenerator()
		now := time.Now()

		err = s.repo.InsertStoryPublication(
			ctx,
			string(publicationID),
			string(storyID),
			profileID,
			"original",
			false, // is_featured
			&now,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: publication: %w", ErrFailedToInsertRecord, err)
		}
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryCreated,
		EntityType: "story",
		EntityID:   string(storyID),
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"slug":                slug,
			"kind":                kind,
			"author_profile_slug": authorProfileSlug,
		},
	})

	return story, nil
}

// Update updates an existing story (slug, picture, and properties).
func (s *Service) Update(
	ctx context.Context,
	locale string,
	userID string,
	userKind string,
	storyID string,
	slug string,
	storyPictureURI *string,
	properties map[string]any,
	visibility string,
	featDiscussions *bool,
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

	// Check if story is managed (synced from external source)
	storyForEdit, err := s.repo.GetStoryForEdit(ctx, locale, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	if storyForEdit == nil {
		return nil, fmt.Errorf("%w: %s", ErrStoryNotFound, storyID)
	}

	if storyForEdit.IsManaged {
		return nil, ErrManagedStory
	}

	// Validate slug availability and date prefix (only block on errors, not warnings)
	slugResult, err := s.CheckSlugAvailability(ctx, slug, &storyID, &storyID, nil, false)
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
	err = s.repo.UpdateStory(
		ctx,
		storyID,
		slug,
		storyPictureURI,
		properties,
		visibility,
		featDiscussions,
	)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToUpdateRecord, storyID, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryUpdated,
		EntityType: "story",
		EntityID:   storyID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"slug": slug},
	})

	// Return updated story
	story, err := s.repo.GetStoryForEdit(ctx, locale, storyID)
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

	// Check if story is managed (synced from external source)
	storyForEdit, err := s.repo.GetStoryForEdit(ctx, localeCode, storyID)
	if err != nil {
		return fmt.Errorf("%w(storyID: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	if storyForEdit != nil && storyForEdit.IsManaged {
		return ErrManagedStory
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

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryTranslationUpdated,
		EntityType: "story",
		EntityID:   storyID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"locale_code": localeCode},
	})

	return nil
}

// DeleteTranslation deletes a specific translation for a story.
func (s *Service) DeleteTranslation(
	ctx context.Context,
	userID string,
	storyID string,
	localeCode string,
) error {
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

	err = s.repo.DeleteStoryTx(ctx, storyID, localeCode)
	if err != nil {
		return fmt.Errorf(
			"%w(storyID: %s, locale: %s): %w",
			ErrFailedToRemoveRecord,
			storyID,
			localeCode,
			err,
		)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryTranslationDeleted,
		EntityType: "story",
		EntityID:   storyID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"locale_code": localeCode},
	})

	return nil
}

// ListTranslationLocales returns locale codes that have translations for a story.
func (s *Service) ListTranslationLocales(
	ctx context.Context,
	storyID string,
) ([]string, error) {
	locales, err := s.repo.ListStoryTxLocales(ctx, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToListRecords, storyID, err)
	}

	// Trim whitespace from locale codes (char fields may be padded)
	result := make([]string, 0, len(locales))
	for _, l := range locales {
		result = append(result, strings.TrimSpace(l))
	}

	return result, nil
}

// GetTranslationContent returns the translation content for a specific locale.
func (s *Service) GetTranslationContent(
	ctx context.Context,
	storyID string,
	localeCode string,
) (string, string, string, error) {
	story, err := s.repo.GetStoryForEdit(ctx, localeCode, storyID)
	if err != nil {
		return "", "", "", fmt.Errorf("%w(storyID: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	if story == nil {
		return "", "", "", fmt.Errorf("%w: story %s not found", ErrStoryNotFound, storyID)
	}

	// Check if this is actual content for the requested locale (not fallback)
	actualLocale := strings.TrimSpace(story.LocaleCode)
	if actualLocale != localeCode {
		return "", "", "", fmt.Errorf(
			"%w: no translation for locale %s",
			ErrFailedToGetRecord,
			localeCode,
		)
	}

	return story.Title, story.Summary, story.Content, nil
}

// Delete soft-deletes a story. Fails if the story has active publications.
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

	// Check for active publications
	count, err := s.repo.CountStoryPublications(ctx, storyID)
	if err != nil {
		return fmt.Errorf("%w(storyID: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	if count > 0 {
		return ErrHasActivePublications
	}

	// Delete the story
	err = s.repo.RemoveStory(ctx, storyID)
	if err != nil {
		return fmt.Errorf("%w(storyID: %s): %w", ErrFailedToRemoveRecord, storyID, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryDeleted,
		EntityType: "story",
		EntityID:   storyID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
	})

	return nil
}

// AddPublication publishes a story to a profile.
func (s *Service) AddPublication(
	ctx context.Context,
	userID string,
	storyID string,
	profileID string,
	localeCode string,
	isFeatured bool,
) (*StoryPublication, error) {
	// Check story edit permission
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

	// Check user has membership access to the target profile (contributor+)
	membershipKind, err := s.repo.GetUserMembershipForProfile(ctx, userID, profileID)
	if err != nil {
		return nil, fmt.Errorf(
			"%w(userID: %s, profileID: %s): %w",
			ErrFailedToGetRecord,
			userID,
			profileID,
			err,
		)
	}

	publishRoles := map[string]bool{
		"admin": true, "owner": true, "lead": true, "maintainer": true, "contributor": true,
	}
	if !publishRoles[membershipKind] {
		return nil, fmt.Errorf(
			"%w: user %s has no publish access to profile %s",
			ErrNoProfileAccess,
			userID,
			profileID,
		)
	}

	// Check featured permission (maintainer+)
	if isFeatured {
		featureRoles := map[string]bool{
			"admin": true, "owner": true, "lead": true, "maintainer": true,
		}
		if !featureRoles[membershipKind] {
			return nil, fmt.Errorf(
				"%w: user %s cannot feature on profile %s (requires maintainer+)",
				ErrInsufficientProfileRole,
				userID,
				profileID,
			)
		}
	}

	// Create the publication
	publicationID := s.idGenerator()
	now := time.Now()

	err = s.repo.InsertStoryPublication(
		ctx,
		string(publicationID),
		storyID,
		profileID,
		"original",
		isFeatured,
		&now,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	// Invalidate the story slug cache so it becomes publicly findable
	// We need to get the story's slug first
	storyForEdit, err := s.repo.GetStoryForEdit(ctx, localeCode, storyID)
	if err == nil && storyForEdit != nil {
		_ = s.repo.InvalidateStorySlugCache(ctx, storyForEdit.Slug)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryPublished,
		EntityType: "story",
		EntityID:   storyID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"publication_id": string(publicationID),
			"profile_id":     profileID,
		},
	})

	// Return the created publication with profile info
	publications, err := s.repo.ListStoryPublications(ctx, localeCode, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	// Find the one we just created
	for _, pub := range publications {
		if pub.ID == string(publicationID) {
			return pub, nil
		}
	}

	// Fallback: return a basic publication object
	return &StoryPublication{
		ID:          string(publicationID),
		StoryID:     storyID,
		ProfileID:   profileID,
		Kind:        "original",
		IsFeatured:  isFeatured,
		PublishedAt: &now,
		CreatedAt:   now,
	}, nil
}

// RemovePublication unpublishes a story from a profile.
func (s *Service) RemovePublication(
	ctx context.Context,
	userID string,
	storyID string,
	publicationID string,
	localeCode string,
) error {
	// Check story edit permission
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

	// Check user has membership access to the publication's profile (contributor+)
	pubProfileID, err := s.repo.GetStoryPublicationProfileID(ctx, publicationID)
	if err != nil {
		return fmt.Errorf(
			"%w(publicationID: %s): %w",
			ErrFailedToGetRecord,
			publicationID,
			err,
		)
	}

	if pubProfileID != "" {
		membershipKind, err := s.repo.GetUserMembershipForProfile(ctx, userID, pubProfileID)
		if err != nil {
			return fmt.Errorf(
				"%w(userID: %s, profileID: %s): %w",
				ErrFailedToGetRecord,
				userID,
				pubProfileID,
				err,
			)
		}

		publishRoles := map[string]bool{
			"admin": true, "owner": true, "lead": true, "maintainer": true, "contributor": true,
		}
		if !publishRoles[membershipKind] {
			return fmt.Errorf(
				"%w: user %s has no access to profile %s",
				ErrNoProfileAccess,
				userID,
				pubProfileID,
			)
		}
	}

	// Remove the publication
	err = s.repo.RemoveStoryPublication(ctx, publicationID)
	if err != nil {
		return fmt.Errorf("%w(publicationID: %s): %w", ErrFailedToRemoveRecord, publicationID, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryUnpublished,
		EntityType: "story",
		EntityID:   storyID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"publication_id": publicationID},
	})

	// Invalidate slug cache
	storyForEdit, err := s.repo.GetStoryForEdit(ctx, localeCode, storyID)
	if err == nil && storyForEdit != nil {
		_ = s.repo.InvalidateStorySlugCache(ctx, storyForEdit.Slug)
	}

	return nil
}

// UpdatePublication updates a story publication's properties.
func (s *Service) UpdatePublication(
	ctx context.Context,
	userID string,
	storyID string,
	publicationID string,
	isFeatured bool,
) error {
	// Check story edit permission
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

	// Check user has membership access to the publication's profile (maintainer+ for featured)
	pubProfileID, err := s.repo.GetStoryPublicationProfileID(ctx, publicationID)
	if err != nil {
		return fmt.Errorf(
			"%w(publicationID: %s): %w",
			ErrFailedToGetRecord,
			publicationID,
			err,
		)
	}

	if pubProfileID != "" {
		membershipKind, err := s.repo.GetUserMembershipForProfile(ctx, userID, pubProfileID)
		if err != nil {
			return fmt.Errorf(
				"%w(userID: %s, profileID: %s): %w",
				ErrFailedToGetRecord,
				userID,
				pubProfileID,
				err,
			)
		}

		featureRoles := map[string]bool{
			"admin": true, "owner": true, "lead": true, "maintainer": true,
		}
		if !featureRoles[membershipKind] {
			return fmt.Errorf(
				"%w: user %s cannot toggle featured on profile %s (requires maintainer+)",
				ErrInsufficientProfileRole,
				userID,
				pubProfileID,
			)
		}
	}

	err = s.repo.UpdateStoryPublication(ctx, publicationID, isFeatured)
	if err != nil {
		return fmt.Errorf(
			"%w(publicationID: %s): %w",
			ErrFailedToUpdateRecord,
			publicationID,
			err,
		)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryFeatured,
		EntityType: "story",
		EntityID:   storyID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"publication_id": publicationID, "is_featured": isFeatured},
	})

	return nil
}

// ListActivities returns published activity stories sorted by activity_time_start.
func (s *Service) ListActivities(
	ctx context.Context,
	localeCode string,
	filterAuthorProfileID *string,
) ([]*StoryWithChildren, error) {
	records, err := s.repo.ListActivityStories(ctx, localeCode, filterAuthorProfileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return records, nil
}

// ListPublications returns all publications for a story.
func (s *Service) ListPublications(
	ctx context.Context,
	localeCode string,
	storyID string,
) ([]*StoryPublication, error) {
	publications, err := s.repo.ListStoryPublications(ctx, localeCode, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToListRecords, storyID, err)
	}

	return publications, nil
}
