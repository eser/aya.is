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

type Repository interface { //nolint:interfacebloat
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
		isManaged bool,
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
		seriesID *string,
		sortOrder *int32,
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
		isManaged bool,
	) error
	IsStoryTxManaged(ctx context.Context, storyID string, localeCode string) (bool, error)
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
	ListStoriesInSeries(
		ctx context.Context,
		localeCode string,
		seriesID string,
	) ([]*StoryWithChildren, error)
	// AI summarization methods
	GetUnsummarizedPublishedStories(
		ctx context.Context,
		maxItems int,
	) ([]*UnsummarizedStory, error)
	UpsertStorySummaryAI(
		ctx context.Context,
		storyID string,
		localeCode string,
		summaryAI string,
	) error
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
	isPublished := s.isStoryPublished(ctx, storyID)
	datePrefix := s.resolveSlugDatePrefix(ctx, publishedAt, storyID)

	// Validate format (errors block, warnings are collected)
	if result := checkSlugFormatError(slug, isPublished, datePrefix); result != nil {
		return result, nil
	}

	// Check uniqueness
	if result := s.checkSlugUniqueness(ctx, slug, excludeStoryID, includeDeleted); result != nil {
		return result, nil
	}

	// Return any warnings from format validation
	if result := checkSlugFormatWarning(slug, isPublished, datePrefix); result != nil {
		return result, nil
	}

	return &SlugAvailabilityResult{
		Available: true,
		Message:   "",
		Severity:  "",
	}, nil
}

// isStoryPublished checks whether a story has any active publications.
func (s *Service) isStoryPublished(ctx context.Context, storyID *string) bool {
	if storyID == nil {
		return false
	}

	count, err := s.repo.CountStoryPublications(ctx, *storyID)

	return err == nil && count > 0
}

// checkSlugFormatError returns an error-severity result if slug format is invalid, nil otherwise.
func checkSlugFormatError(
	slug string,
	isPublished bool,
	datePrefix time.Time,
) *SlugAvailabilityResult {
	if lengthResult := validateSlugLength(slug, isPublished); lengthResult != nil &&
		lengthResult.Severity == SeverityError {
		return &SlugAvailabilityResult{
			Available: false,
			Message:   lengthResult.Message,
			Severity:  lengthResult.Severity,
		}
	}

	if prefixResult := validateSlugDatePrefix(slug, datePrefix, isPublished); prefixResult != nil &&
		prefixResult.Severity == SeverityError {
		return &SlugAvailabilityResult{
			Available: false,
			Message:   prefixResult.Message,
			Severity:  prefixResult.Severity,
		}
	}

	return nil
}

// checkSlugFormatWarning returns a warning-severity result if slug format has warnings, nil otherwise.
func checkSlugFormatWarning(
	slug string,
	isPublished bool,
	datePrefix time.Time,
) *SlugAvailabilityResult {
	if lengthResult := validateSlugLength(slug, isPublished); lengthResult != nil &&
		lengthResult.Severity == SeverityWarning {
		return &SlugAvailabilityResult{
			Available: true,
			Message:   lengthResult.Message,
			Severity:  lengthResult.Severity,
		}
	}

	if prefixResult := validateSlugDatePrefix(slug, datePrefix, isPublished); prefixResult != nil &&
		prefixResult.Severity == SeverityWarning {
		return &SlugAvailabilityResult{
			Available: true,
			Message:   prefixResult.Message,
			Severity:  prefixResult.Severity,
		}
	}

	return nil
}

// checkSlugUniqueness checks if a slug is already used by active or deleted stories.
func (s *Service) checkSlugUniqueness(
	ctx context.Context,
	slug string,
	excludeStoryID *string,
	includeDeleted bool,
) *SlugAvailabilityResult {
	existingStoryID, err := s.repo.GetStoryIDBySlug(ctx, slug)
	if err == nil && existingStoryID != "" {
		if excludeStoryID == nil || existingStoryID != *excludeStoryID {
			return &SlugAvailabilityResult{
				Available: false,
				Message:   "This slug is already taken",
				Severity:  SeverityError,
			}
		}
	}

	if includeDeleted {
		deletedStoryID, delErr := s.repo.GetStoryIDBySlugIncludingDeleted(ctx, slug)
		if delErr == nil && deletedStoryID != "" {
			if excludeStoryID == nil || deletedStoryID != *excludeStoryID {
				return &SlugAvailabilityResult{
					Available: false,
					Message:   "This slug was previously used",
					Severity:  SeverityError,
				}
			}
		}
	}

	return nil
}

// resolveSlugDatePrefix determines the date prefix to use for slug validation.
func (s *Service) resolveSlugDatePrefix(
	ctx context.Context,
	publishedAt *time.Time,
	storyID *string,
) time.Time {
	switch {
	case publishedAt != nil:
		return *publishedAt
	case storyID != nil:
		firstPublishedAt, err := s.repo.GetStoryFirstPublishedAt(ctx, *storyID)
		if err == nil && firstPublishedAt != nil {
			return *firstPublishedAt
		}

		return time.Now()
	default:
		return time.Now()
	}
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
	featDiscussionsOverride *bool,
) (*Story, error) {
	validateErr := s.validateCreateInputs(
		ctx, slug, storyPictureURI, userKind, publishToProfileSlugs,
	)
	if validateErr != nil {
		return nil, validateErr
	}

	authorProfileID, featDiscussions, resolveErr := s.resolveAuthorAndFeatures(
		ctx, authorProfileSlug, localeCode, featDiscussionsOverride,
	)
	if resolveErr != nil {
		return nil, resolveErr
	}

	storyID := string(s.idGenerator())

	story, err := s.insertStoryWithTranslation(ctx, insertStoryParams{
		storyID:         storyID,
		authorProfileID: authorProfileID,
		slug:            slug,
		kind:            kind,
		localeCode:      localeCode,
		title:           title,
		summary:         summary,
		content:         content,
		storyPictureURI: storyPictureURI,
		properties:      properties,
		visibility:      visibility,
		featDiscussions: featDiscussions,
	})
	if err != nil {
		return nil, err
	}

	publishErr := s.publishAndAuditCreate(
		ctx, storyID, userID, slug, kind, authorProfileSlug, publishToProfileSlugs,
	)
	if publishErr != nil {
		return nil, publishErr
	}

	return story, nil
}

// resolveAuthorAndFeatures resolves the author profile ID and feature flags for story creation.
func (s *Service) resolveAuthorAndFeatures(
	ctx context.Context,
	authorProfileSlug, localeCode string,
	featDiscussionsOverride *bool,
) (string, bool, error) {
	authorProfileID, err := s.repo.GetProfileIDBySlug(ctx, authorProfileSlug)
	if err != nil {
		return "", false, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			authorProfileSlug,
			err,
		)
	}

	featDiscussions, featErr := s.resolveFeatDiscussions(
		ctx, localeCode, authorProfileID, featDiscussionsOverride,
	)
	if featErr != nil {
		return "", false, featErr
	}

	return authorProfileID, featDiscussions, nil
}

// publishAndAuditCreate publishes a newly created story to profiles and records an audit entry.
func (s *Service) publishAndAuditCreate(
	ctx context.Context,
	storyID, userID, slug, kind, authorProfileSlug string,
	publishToProfileSlugs []string,
) error {
	pubErr := s.createPublicationsForProfiles(ctx, storyID, userID, publishToProfileSlugs)
	if pubErr != nil {
		return pubErr
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryCreated,
		EntityType: "story",
		EntityID:   storyID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"slug":                slug,
			"kind":                kind,
			"author_profile_slug": authorProfileSlug,
		},
	})

	return nil
}

// insertStoryParams holds the parameters for inserting a story with its translation.
type insertStoryParams struct {
	storyID         string
	authorProfileID string
	slug            string
	kind            string
	localeCode      string
	title           string
	summary         string
	content         string
	storyPictureURI *string
	properties      map[string]any
	visibility      string
	featDiscussions bool
}

// insertStoryWithTranslation creates the story record and its localized translation.
func (s *Service) insertStoryWithTranslation(
	ctx context.Context,
	params insertStoryParams,
) (*Story, error) {
	story, err := s.repo.InsertStory(
		ctx,
		params.storyID,
		params.authorProfileID,
		params.slug,
		params.kind,
		params.storyPictureURI,
		params.properties,
		false,
		nil,
		params.visibility,
		params.featDiscussions,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	err = s.repo.InsertStoryTx(
		ctx,
		params.storyID,
		params.localeCode,
		params.title,
		params.summary,
		params.content,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: translation: %w", ErrFailedToInsertRecord, err)
	}

	return story, nil
}

// validateCreateInputs validates slug availability and picture URI for story creation.
func (s *Service) validateCreateInputs(
	ctx context.Context,
	slug string,
	storyPictureURI *string,
	userKind string,
	publishToProfileSlugs []string,
) error {
	isPublishing := len(publishToProfileSlugs) > 0

	var publishedAt *time.Time

	if isPublishing {
		now := time.Now()
		publishedAt = &now
	}

	slugResult, err := s.CheckSlugAvailability(ctx, slug, nil, nil, publishedAt, false)
	if err != nil {
		return err
	}

	if !slugResult.Available && slugResult.Severity == SeverityError {
		return fmt.Errorf("%w: %s", ErrInvalidSlugPrefix, slugResult.Message)
	}

	return s.validateStoryPictureURI(storyPictureURI, userKind)
}

// validateStoryPictureURI validates the story picture URI and checks prefix restrictions.
func (s *Service) validateStoryPictureURI(storyPictureURI *string, userKind string) error {
	pictureErr := validateOptionalURL(storyPictureURI)
	if pictureErr != nil {
		return pictureErr
	}

	// Non-admin users can only use URIs from our upload service
	if userKind != "admin" {
		return validateURIPrefixes(storyPictureURI, s.config.GetAllowedURIPrefixes())
	}

	return nil
}

// resolveFeatDiscussions determines the discussion feature flag for a new story.
func (s *Service) resolveFeatDiscussions(
	ctx context.Context,
	localeCode string,
	authorProfileID string,
	override *bool,
) (bool, error) {
	if override != nil {
		return *override, nil
	}

	authorProfile, err := s.repo.GetProfileByID(ctx, localeCode, authorProfileID)
	if err != nil {
		return false, fmt.Errorf(
			"%w(profileID: %s): %w",
			ErrFailedToGetRecord,
			authorProfileID,
			err,
		)
	}

	if authorProfile != nil {
		return authorProfile.OptionStoryDiscussionsByDefault, nil
	}

	return false, nil
}

// publishRoles defines the set of membership roles allowed to publish stories.
var publishRoles = map[string]bool{ //nolint:gochecknoglobals
	"admin": true, "owner": true, "lead": true, "maintainer": true, "contributor": true,
}

// createPublicationsForProfiles creates story publications for a list of profile slugs.
func (s *Service) createPublicationsForProfiles(
	ctx context.Context,
	storyID string,
	userID string,
	profileSlugs []string,
) error {
	for _, profileSlug := range profileSlugs {
		profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
		if err != nil {
			return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
		}

		membershipKind, err := s.repo.GetUserMembershipForProfile(ctx, userID, profileID)
		if err != nil {
			return fmt.Errorf(
				"%w(userID: %s, profileID: %s): %w",
				ErrFailedToGetRecord,
				userID,
				profileID,
				err,
			)
		}

		if !publishRoles[membershipKind] {
			return fmt.Errorf(
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
			storyID,
			profileID,
			"original",
			false, // is_featured
			&now,
			nil,
		)
		if err != nil {
			return fmt.Errorf("%w: publication: %w", ErrFailedToInsertRecord, err)
		}
	}

	return nil
}

// featureRoles defines the set of membership roles allowed to feature stories.
var featureRoles = map[string]bool{ //nolint:gochecknoglobals
	"admin": true, "owner": true, "lead": true, "maintainer": true,
}

// validatePublishAccess checks that a user has publish access and optionally featured permission.
func (s *Service) validatePublishAccess(
	ctx context.Context,
	userID string,
	profileID string,
	isFeatured bool,
) error {
	membershipKind, err := s.repo.GetUserMembershipForProfile(ctx, userID, profileID)
	if err != nil {
		return fmt.Errorf(
			"%w(userID: %s, profileID: %s): %w",
			ErrFailedToGetRecord,
			userID,
			profileID,
			err,
		)
	}

	if !publishRoles[membershipKind] {
		return fmt.Errorf(
			"%w: user %s has no publish access to profile %s",
			ErrNoProfileAccess,
			userID,
			profileID,
		)
	}

	if isFeatured && !featureRoles[membershipKind] {
		return fmt.Errorf(
			"%w: user %s cannot feature on profile %s (requires maintainer+)",
			ErrInsufficientProfileRole,
			userID,
			profileID,
		)
	}

	return nil
}

// validateUpdatePreconditions checks authorization, slug availability, and picture URI.
// Returns the existing story so callers can inspect IsManaged and preserve fields.
// For managed stories, slug and picture validation is skipped (those fields are locked).
func (s *Service) validateUpdatePreconditions(
	ctx context.Context,
	locale string,
	userID string,
	userKind string,
	storyID string,
	slug string,
	storyPictureURI *string,
) (*StoryForEdit, error) {
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

	storyForEdit, err := s.repo.GetStoryForEdit(ctx, locale, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	if storyForEdit == nil {
		return nil, fmt.Errorf("%w: %s", ErrStoryNotFound, storyID)
	}

	// Managed stories can only update visibility and discussions;
	// skip slug and picture URI validation since those fields are preserved.
	if storyForEdit.IsManaged {
		return storyForEdit, nil
	}

	slugResult, err := s.CheckSlugAvailability(ctx, slug, &storyID, &storyID, nil, false)
	if err != nil {
		return nil, err
	}

	if !slugResult.Available && slugResult.Severity == SeverityError {
		return nil, fmt.Errorf("%w: %s", ErrInvalidSlugPrefix, slugResult.Message)
	}

	return storyForEdit, s.validateStoryPictureURI(storyPictureURI, userKind)
}

// Update updates an existing story (slug, picture, and properties).
func (s *Service) Update( //nolint:funlen
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
	seriesID *string,
	sortOrder *int32,
) (*StoryForEdit, error) {
	// Validate authorization, slug, and picture
	existingStory, validateErr := s.validateUpdatePreconditions(
		ctx,
		locale,
		userID,
		userKind,
		storyID,
		slug,
		storyPictureURI,
	)
	if validateErr != nil {
		return nil, validateErr
	}

	// For managed stories, preserve managed fields — only visibility and discussions can change.
	if existingStory.IsManaged {
		slug = existingStory.Slug
		storyPictureURI = existingStory.StoryPictureURI

		if propsMap, ok := existingStory.Properties.(map[string]any); ok {
			properties = propsMap
		}
	}

	// Update the story
	err := s.repo.UpdateStory(
		ctx,
		storyID,
		slug,
		storyPictureURI,
		properties,
		visibility,
		featDiscussions,
		seriesID,
		sortOrder,
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
		SessionID:  nil,
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

	// If the specific translation is managed, skip the upsert silently
	txIsManaged, txErr := s.repo.IsStoryTxManaged(ctx, storyID, localeCode)
	if txErr == nil && txIsManaged {
		return nil
	}
	// If txErr (row doesn't exist) → this is a new translation → allow

	// Update the translation (use upsert to handle new locales)
	err = s.repo.UpsertStoryTx(ctx, storyID, localeCode, title, summary, content, false)
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
		SessionID:  nil,
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

	// Check if the translation is managed (cannot delete synced translations)
	txIsManaged, txErr := s.repo.IsStoryTxManaged(ctx, storyID, localeCode)
	if txErr == nil && txIsManaged {
		return ErrManagedStory
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
		SessionID:  nil,
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
		SessionID:  nil,
		Payload:    nil,
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
	// Check story edit and publish access
	authErr := s.authorizeStoryEdit(ctx, userID, storyID)
	if authErr != nil {
		return nil, authErr
	}

	accessErr := s.validatePublishAccess(ctx, userID, profileID, isFeatured)
	if accessErr != nil {
		return nil, accessErr
	}

	// Create the publication
	publicationID := s.idGenerator()
	now := time.Now()

	err := s.repo.InsertStoryPublication(
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

	s.invalidateSlugCache(ctx, localeCode, storyID)

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryPublished,
		EntityType: "story",
		EntityID:   storyID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"publication_id": string(publicationID),
			"profile_id":     profileID,
		},
	})

	return s.findOrBuildPublication(
		ctx, localeCode, storyID, string(publicationID), profileID, isFeatured, now,
	)
}

// invalidateSlugCache invalidates the story slug cache (best-effort).
func (s *Service) invalidateSlugCache(ctx context.Context, localeCode string, storyID string) {
	storyForEdit, err := s.repo.GetStoryForEdit(ctx, localeCode, storyID)
	if err == nil && storyForEdit != nil {
		_ = s.repo.InvalidateStorySlugCache(ctx, storyForEdit.Slug)
	}
}

// findOrBuildPublication returns the just-created publication with profile info, or a fallback.
func (s *Service) findOrBuildPublication(
	ctx context.Context,
	localeCode string,
	storyID string,
	publicationID string,
	profileID string,
	isFeatured bool,
	now time.Time,
) (*StoryPublication, error) {
	publications, err := s.repo.ListStoryPublications(ctx, localeCode, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	for _, pub := range publications {
		if pub.ID == publicationID {
			return pub, nil
		}
	}

	return &StoryPublication{
		ID:                publicationID,
		StoryID:           storyID,
		ProfileID:         profileID,
		ProfilePictureURI: nil,
		ProfileSlug:       "",
		ProfileTitle:      "",
		ProfileKind:       "",
		Kind:              "original",
		IsFeatured:        isFeatured,
		PublishedAt:       &now,
		CreatedAt:         now,
	}, nil
}

// validatePublicationEditAccess checks story edit permission and publication profile access.
// When requireMaintainer is true, it checks for maintainer+ role (for featured operations).
func (s *Service) validatePublicationEditAccess(
	ctx context.Context,
	userID string,
	storyID string,
	publicationID string,
	requireMaintainer bool,
) error {
	authErr := s.authorizeStoryEdit(ctx, userID, storyID)
	if authErr != nil {
		return authErr
	}

	pubProfileID, err := s.repo.GetStoryPublicationProfileID(ctx, publicationID)
	if err != nil {
		return fmt.Errorf(
			"%w(publicationID: %s): %w",
			ErrFailedToGetRecord,
			publicationID,
			err,
		)
	}

	if pubProfileID == "" {
		return nil
	}

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

	requiredRoles := publishRoles
	roleErr := ErrNoProfileAccess

	if requireMaintainer {
		requiredRoles = featureRoles
		roleErr = ErrInsufficientProfileRole
	}

	if !requiredRoles[membershipKind] {
		return fmt.Errorf(
			"%w: user %s insufficient access on profile %s",
			roleErr,
			userID,
			pubProfileID,
		)
	}

	return nil
}

// RemovePublication unpublishes a story from a profile.
func (s *Service) RemovePublication(
	ctx context.Context,
	userID string,
	storyID string,
	publicationID string,
	localeCode string,
) error {
	accessErr := s.validatePublicationEditAccess(ctx, userID, storyID, publicationID, false)
	if accessErr != nil {
		return accessErr
	}

	// Remove the publication
	err := s.repo.RemoveStoryPublication(ctx, publicationID)
	if err != nil {
		return fmt.Errorf("%w(publicationID: %s): %w", ErrFailedToRemoveRecord, publicationID, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.StoryUnpublished,
		EntityType: "story",
		EntityID:   storyID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload:    map[string]any{"publication_id": publicationID},
	})

	s.invalidateSlugCache(ctx, localeCode, storyID)

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
	accessErr := s.validatePublicationEditAccess(ctx, userID, storyID, publicationID, true)
	if accessErr != nil {
		return accessErr
	}

	err := s.repo.UpdateStoryPublication(ctx, publicationID, isFeatured)
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
		SessionID:  nil,
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

// ListBySeriesID returns all public stories in a series, ordered by sort_order then published_at.
func (s *Service) ListBySeriesID(
	ctx context.Context,
	localeCode string,
	seriesID string,
) ([]*StoryWithChildren, error) {
	records, err := s.repo.ListStoriesInSeries(ctx, localeCode, seriesID)
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

// GetUnsummarizedStories returns published story translations that have no AI summary.
func (s *Service) GetUnsummarizedStories(
	ctx context.Context,
	maxItems int,
) ([]*UnsummarizedStory, error) {
	stories, err := s.repo.GetUnsummarizedPublishedStories(ctx, maxItems)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return stories, nil
}

// PersistSummaryAI updates the AI-generated summary for a story translation.
func (s *Service) PersistSummaryAI(
	ctx context.Context,
	storyID string,
	localeCode string,
	summaryAI string,
) error {
	err := s.repo.UpsertStorySummaryAI(ctx, storyID, localeCode, summaryAI)
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
