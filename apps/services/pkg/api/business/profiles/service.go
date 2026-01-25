package profiles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

var (
	ErrFailedToGetRecord    = errors.New("failed to get record")
	ErrFailedToListRecords  = errors.New("failed to list records")
	ErrFailedToCreateRecord = errors.New("failed to create record")
	ErrFailedToUpdateRecord = errors.New("failed to update record")
	ErrFailedToDeleteRecord = errors.New("failed to delete record")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrProfileNotFound      = errors.New("profile not found")
	ErrInvalidURI           = errors.New("invalid URI")
	ErrInvalidURIPrefix     = errors.New("URI must start with allowed prefix")
)

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

// Config holds the profiles service configuration.
type Config struct {
	// AllowedURIPrefixes is a comma-separated list of allowed URI prefixes.
	AllowedURIPrefixes string `conf:"allowed_uri_prefixes" default:"https://objects.aya.is/,https://avatars.githubusercontent.com/"`

	// ForbiddenSlugs is a comma-separated list of reserved slugs that cannot be used as profile slugs.
	ForbiddenSlugs string `conf:"forbidden_slugs" default:"about,admin,api,auth,communities,community,config,contact,contributions,dashboard,element,elements,events,faq,feed,guide,help,home,impressum,imprint,jobs,legal,login,logout,new,news,null,organizations,orgs,people,policies,policy,privacy,product,products,profile,profiles,projects,register,root,search,services,settings,signin,signout,signup,site,stories,story,support,tag,tags,terms,tos,undefined,user,users,verify,wiki"`
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

// GetForbiddenSlugs returns the forbidden slugs as a set.
func (c *Config) GetForbiddenSlugs() map[string]bool {
	if c.ForbiddenSlugs == "" {
		return nil
	}

	slugs := strings.Split(c.ForbiddenSlugs, ",")
	result := make(map[string]bool, len(slugs))

	for _, slug := range slugs {
		trimmed := strings.TrimSpace(slug)
		if trimmed != "" {
			result[trimmed] = true
		}
	}

	return result
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

type RecentPostsFetcher interface {
	GetRecentPostsByUsername(
		ctx context.Context,
		username string,
		userID string,
	) ([]*ExternalPost, error)
}

type Repository interface { //nolint:interfacebloat
	GetProfileIDBySlug(ctx context.Context, slug string) (string, error)
	GetProfileIDByCustomDomain(ctx context.Context, domain string) (*string, error)
	CheckProfileSlugExists(ctx context.Context, slug string) (bool, error)
	CheckProfileSlugExistsIncludingDeleted(ctx context.Context, slug string) (bool, error)
	CheckPageSlugExistsIncludingDeleted(
		ctx context.Context,
		profileID string,
		pageSlug string,
	) (bool, error)
	GetProfileByID(ctx context.Context, localeCode string, id string) (*Profile, error)
	ListProfiles(
		ctx context.Context,
		localeCode string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*Profile], error)
	// ListProfileLinksForKind(ctx context.Context, kind string) ([]*ProfileLink, error)
	ListProfilePagesByProfileID(
		ctx context.Context,
		localeCode string,
		profileID string,
	) ([]*ProfilePageBrief, error)
	GetProfilePageByProfileIDAndSlug(
		ctx context.Context,
		localeCode string,
		profileID string,
		pageSlug string,
	) (*ProfilePage, error)
	ListProfileLinksByProfileID(
		ctx context.Context,
		localeCode string,
		profileID string,
	) ([]*ProfileLinkBrief, error)
	ListProfileContributions(
		ctx context.Context,
		localeCode string,
		profileID string,
		kinds []string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*ProfileMembership], error)
	ListProfileMembers(
		ctx context.Context,
		localeCode string,
		profileID string,
		kinds []string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*ProfileMembership], error)
	GetProfileMembershipsByMemberProfileID(
		ctx context.Context,
		localeCode string,
		memberProfileID string,
	) ([]*ProfileMembership, error)
	CreateProfile(
		ctx context.Context,
		id string,
		slug string,
		kind string,
		customDomain *string,
		profilePictureURI *string,
		pronouns *string,
		properties map[string]any,
	) error
	CreateProfileTx(
		ctx context.Context,
		profileID string,
		localeCode string,
		title string,
		description string,
		properties map[string]any,
	) error
	UpdateProfile(
		ctx context.Context,
		id string,
		profilePictureURI *string,
		pronouns *string,
		properties map[string]any,
	) error
	UpdateProfileTx(
		ctx context.Context,
		profileID string,
		localeCode string,
		title string,
		description string,
		properties map[string]any,
	) error
	UpsertProfileTx(
		ctx context.Context,
		profileID string,
		localeCode string,
		title string,
		description string,
		properties map[string]any,
	) error
	GetProfileOwnershipForUser(
		ctx context.Context,
		userID string,
		profileSlug string,
	) (*ProfileOwnership, error)
	GetUserProfilePermissions(
		ctx context.Context,
		userID string,
	) ([]*ProfilePermission, error)
	GetProfileTxByID(
		ctx context.Context,
		profileID string,
	) ([]*ProfileTx, error)
	// Profile Link methods
	GetProfileLink(ctx context.Context, id string) (*ProfileLink, error)
	CreateProfileLink(
		ctx context.Context,
		id string,
		kind string,
		profileID string,
		order int,
		uri *string,
		title string,
		isHidden bool,
	) (*ProfileLink, error)
	UpdateProfileLink(
		ctx context.Context,
		id string,
		kind string,
		order int,
		uri *string,
		title string,
		isHidden bool,
	) error
	DeleteProfileLink(ctx context.Context, id string) error
	// Profile Page methods
	GetProfilePage(ctx context.Context, id string) (*ProfilePage, error)
	CreateProfilePage(
		ctx context.Context,
		id string,
		slug string,
		profileID string,
		order int,
		coverPictureURI *string,
		publishedAt *string,
	) (*ProfilePage, error)
	CreateProfilePageTx(
		ctx context.Context,
		profilePageID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	UpdateProfilePage(
		ctx context.Context,
		id string,
		slug string,
		order int,
		coverPictureURI *string,
		publishedAt *string,
	) error
	UpdateProfilePageTx(
		ctx context.Context,
		profilePageID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	UpsertProfilePageTx(
		ctx context.Context,
		profilePageID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	DeleteProfilePage(ctx context.Context, id string) error
	// Search methods
	Search(
		ctx context.Context,
		localeCode string,
		query string,
		profileSlug *string,
		limit int32,
	) ([]*SearchResult, error)
	// OAuth Profile Link methods
	GetProfileLinkByRemoteID(
		ctx context.Context,
		profileID string,
		kind string,
		remoteID string,
	) (*ProfileLink, error)
	CreateOAuthProfileLink(
		ctx context.Context,
		id string,
		kind string,
		profileID string,
		order int,
		remoteID string,
		publicID string,
		uri string,
		title string,
		authProvider string,
		authScope string,
		accessToken string,
		accessTokenExpiresAt *sql.NullTime,
		refreshToken *string,
	) (*ProfileLink, error)
	UpdateProfileLinkOAuthTokens(
		ctx context.Context,
		id string,
		publicID string,
		uri string,
		title string,
		authScope string,
		accessToken string,
		accessTokenExpiresAt *sql.NullTime,
		refreshToken *string,
	) error
	GetMaxProfileLinkOrder(ctx context.Context, profileID string) (int, error)
}

type Service struct {
	logger      *logfx.Logger
	config      *Config
	repo        Repository
	idGenerator RecordIDGenerator
}

func NewService(logger *logfx.Logger, config *Config, repo Repository) *Service {
	return &Service{logger: logger, config: config, repo: repo, idGenerator: DefaultIDGenerator}
}

func (s *Service) GetByID(ctx context.Context, localeCode string, id string) (*Profile, error) {
	record, err := s.repo.GetProfileByID(ctx, localeCode, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return record, nil
}

func (s *Service) GetBySlug(ctx context.Context, localeCode string, slug string) (*Profile, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	record, err := s.repo.GetProfileByID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	return record, nil
}

func (s *Service) GetBySlugEx(
	ctx context.Context,
	localeCode string,
	slug string,
) (*ProfileWithChildren, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	record, err := s.repo.GetProfileByID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	if record == nil {
		return nil, nil //nolint:nilnil
	}

	pages, err := s.repo.ListProfilePagesByProfileID(ctx, localeCode, record.ID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	links, err := s.repo.ListProfileLinksByProfileID(ctx, localeCode, record.ID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	result := &ProfileWithChildren{
		Profile: record,
		Pages:   pages,
		Links:   links,
	}

	return result, nil
}

func (s *Service) GetByCustomDomain(
	ctx context.Context,
	localeCode string,
	domain string,
) (*Profile, error) {
	profileID, err := s.repo.GetProfileIDByCustomDomain(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("%w(custom_domain: %s): %w", ErrFailedToGetRecord, domain, err)
	}

	if profileID == nil {
		return nil, nil //nolint:nilnil
	}

	record, err := s.repo.GetProfileByID(ctx, localeCode, *profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, *profileID, err)
	}

	return record, nil
}

func (s *Service) List(
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*Profile], error) {
	records, err := s.repo.ListProfiles(ctx, localeCode, cursor)
	if err != nil {
		return cursors.Cursored[[]*Profile]{}, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return records, nil
}

func (s *Service) ListPagesBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
) ([]*ProfilePageBrief, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	pages, err := s.repo.ListProfilePagesByProfileID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	return pages, nil
}

func (s *Service) GetPageBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
	pageSlug string,
) (*ProfilePage, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	page, err := s.repo.GetProfilePageByProfileIDAndSlug(ctx, localeCode, profileID, pageSlug)
	if err != nil {
		return nil, fmt.Errorf(
			"%w(profile_id: %s, page_slug: %s): %w",
			ErrFailedToGetRecord,
			profileID,
			pageSlug,
			err,
		)
	}

	return page, nil
}

// CheckPageSlugAvailability checks if a page slug is available within a profile.
// It optionally excludes a specific page ID (for edit scenarios).
func (s *Service) CheckPageSlugAvailability(
	ctx context.Context,
	localeCode string,
	profileSlug string,
	pageSlug string,
	excludePageID *string,
	includeDeleted bool,
) (*SlugAvailabilityResult, error) {
	// Check minimum length
	if len(pageSlug) < 3 {
		return &SlugAvailabilityResult{
			Available: false,
			Message:   "Slug must be at least 3 characters",
			Severity:  SeverityError,
		}, nil
	}

	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil || profileID == "" {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	page, err := s.repo.GetProfilePageByProfileIDAndSlug(ctx, localeCode, profileID, pageSlug)
	if err != nil || page == nil {
		// Page not found in active records, check deleted if requested
		if includeDeleted {
			existsDeleted, delErr := s.repo.CheckPageSlugExistsIncludingDeleted(ctx, profileID, pageSlug)
			if delErr != nil {
				return nil, delErr
			}

			if existsDeleted {
				return &SlugAvailabilityResult{
					Available: false,
					Message:   "This slug was previously used",
					Severity:  SeverityError,
				}, nil
			}
		}

		// Slug is available
		return &SlugAvailabilityResult{
			Available: true,
		}, nil
	}

	// If we're editing and the slug belongs to the same page, it's available
	if excludePageID != nil && page.ID == *excludePageID {
		return &SlugAvailabilityResult{
			Available: true,
		}, nil
	}

	return &SlugAvailabilityResult{
		Available: false,
		Message:   "This slug is already taken",
		Severity:  SeverityError,
	}, nil
}

func (s *Service) ListLinksBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
) ([]*ProfileLinkBrief, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	links, err := s.repo.ListProfileLinksByProfileID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	return links, nil
}

func (s *Service) ListProfileContributionsBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*ProfileMembership], error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return cursors.Cursored[[]*ProfileMembership]{}, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			slug,
			err,
		)
	}

	memberships, err := s.repo.ListProfileContributions(
		ctx,
		localeCode,
		profileID,
		[]string{"organization", "product"},
		cursor,
	)
	if err != nil {
		return cursors.Cursored[[]*ProfileMembership]{}, fmt.Errorf(
			"%w: %w",
			ErrFailedToListRecords,
			err,
		)
	}

	return memberships, nil
}

func (s *Service) ListProfileMembersBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*ProfileMembership], error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return cursors.Cursored[[]*ProfileMembership]{}, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			slug,
			err,
		)
	}

	memberships, err := s.repo.ListProfileMembers(
		ctx,
		localeCode,
		profileID,
		[]string{"organization", "individual"},
		cursor,
	)
	if err != nil {
		return cursors.Cursored[[]*ProfileMembership]{}, fmt.Errorf(
			"%w: %w",
			ErrFailedToListRecords,
			err,
		)
	}

	return memberships, nil
}

func (s *Service) Import(ctx context.Context, fetcher RecentPostsFetcher) error {
	// 	links, err := s.repo.ListProfileLinksForKind(ctx, "x")
	// 	if err != nil {
	// 		return fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	// 	}
	// 	for _, link := range links {
	// 		s.logger.InfoContext(ctx, "importing posts", "kind", link.Kind, "title", link.Title)
	// 		posts, err := fetcher.GetRecentPostsByUsername(ctx, link.RemoteID.String, link.AuthAccessToken)
	// 		if err != nil {
	// 			return fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	// 		}
	// 		s.logger.InfoContext(ctx, "posts imported", "kind", link.Kind, "title", link.Title, "posts", posts)
	// 	}
	return nil
}

func (s *Service) GetMembershipsByUserProfileID(
	ctx context.Context,
	localeCode string,
	userProfileID string,
) ([]*ProfileMembership, error) {
	memberships, err := s.repo.GetProfileMembershipsByMemberProfileID(
		ctx,
		localeCode,
		userProfileID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"%w(userProfileID: %s): %w",
			ErrFailedToGetRecord,
			userProfileID,
			err,
		)
	}

	return memberships, nil
}

func (s *Service) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	exists, err := s.repo.CheckProfileSlugExists(ctx, slug)
	if err != nil {
		return false, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	return exists, nil
}

func (s *Service) CheckSlugAvailability(
	ctx context.Context,
	slug string,
	includeDeleted bool,
) (*SlugAvailabilityResult, error) {
	// Check minimum length
	if len(slug) < 3 {
		return &SlugAvailabilityResult{
			Available: false,
			Message:   "Slug must be at least 3 characters",
			Severity:  SeverityError,
		}, nil
	}

	// Check forbidden slugs
	if s.config.GetForbiddenSlugs()[slug] {
		return &SlugAvailabilityResult{
			Available: false,
			Message:   "This slug is reserved",
			Severity:  SeverityError,
		}, nil
	}

	// Check if slug exists (active records)
	exists, err := s.CheckSlugExists(ctx, slug)
	if err != nil {
		return nil, err
	}

	if exists {
		return &SlugAvailabilityResult{
			Available: false,
			Message:   "This slug is already taken",
			Severity:  SeverityError,
		}, nil
	}

	// Check if slug was used by a deleted record (optional)
	if includeDeleted {
		existsDeleted, err := s.repo.CheckProfileSlugExistsIncludingDeleted(ctx, slug)
		if err != nil {
			return nil, err
		}

		if existsDeleted {
			return &SlugAvailabilityResult{
				Available: false,
				Message:   "This slug was previously used",
				Severity:  SeverityError,
			}, nil
		}
	}

	return &SlugAvailabilityResult{
		Available: true,
	}, nil
}

func (s *Service) Create(
	ctx context.Context,
	localeCode string,
	slug string,
	kind string,
	title string,
	description string,
	customDomain *string,
	profilePictureURI *string,
	pronouns *string,
	properties map[string]any,
) (*Profile, error) {
	// Generate new profile ID
	profileID := s.idGenerator()

	// Create the main profile record
	err := s.repo.CreateProfile(
		ctx,
		string(profileID),
		slug,
		kind,
		customDomain,
		profilePictureURI,
		pronouns,
		properties,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	// Create the localized profile data
	err = s.repo.CreateProfileTx(
		ctx,
		string(profileID),
		localeCode,
		title,
		description,
		nil, // No additional properties for profile_tx for now
	)
	if err != nil {
		return nil, fmt.Errorf("%w: translation: %w", ErrFailedToCreateRecord, err)
	}

	// Fetch and return the created profile
	profile, err := s.repo.GetProfileByID(ctx, localeCode, string(profileID))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	return profile, nil
}

// CanUserEditProfile checks if a user has permission to edit a profile.
func (s *Service) CanUserEditProfile(
	ctx context.Context,
	userID string,
	profileSlug string,
) (bool, error) {
	ownership, err := s.repo.GetProfileOwnershipForUser(ctx, userID, profileSlug)
	if err != nil {
		return false, fmt.Errorf(
			"%w(userID: %s, profileSlug: %s): %w",
			ErrFailedToGetRecord,
			userID,
			profileSlug,
			err,
		)
	}

	if ownership == nil {
		return false, nil
	}

	return ownership.CanEdit, nil
}

// Update updates profile main fields (profile_picture_uri, pronouns, properties).
func (s *Service) Update(
	ctx context.Context,
	userID string,
	userKind string,
	profileSlug string,
	profilePictureURI *string,
	pronouns *string,
	properties map[string]any,
) (*Profile, error) {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Validate profile picture URI
	if err := validateOptionalURL(profilePictureURI); err != nil {
		return nil, err
	}

	// Non-admin users can only use URIs from allowed prefixes
	if userKind != "admin" {
		err := validateURIPrefixes(profilePictureURI, s.config.GetAllowedURIPrefixes())
		if err != nil {
			return nil, err
		}
	}

	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	// Update the profile
	err = s.repo.UpdateProfile(ctx, profileID, profilePictureURI, pronouns, properties)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToUpdateRecord, profileID, err)
	}

	// Return updated profile (use default locale for now)
	profile, err := s.repo.GetProfileByID(ctx, "en", profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	return profile, nil
}

// UpdateTranslation updates profile translation fields (title, description).
func (s *Service) UpdateTranslation(
	ctx context.Context,
	userID string,
	profileSlug string,
	localeCode string,
	title string,
	description string,
	properties map[string]any,
) error {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return err
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	// Update the translation (use upsert to handle new locales)
	err = s.repo.UpsertProfileTx(ctx, profileID, localeCode, title, description, properties)
	if err != nil {
		return fmt.Errorf(
			"%w(profileID: %s, locale: %s): %w",
			ErrFailedToUpdateRecord,
			profileID,
			localeCode,
			err,
		)
	}

	return nil
}

// GetUserProfilePermissions returns all profiles the user has permissions for.
func (s *Service) GetUserProfilePermissions(
	ctx context.Context,
	userID string,
) ([]*ProfilePermission, error) {
	permissions, err := s.repo.GetUserProfilePermissions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w(userID: %s): %w", ErrFailedToGetRecord, userID, err)
	}

	return permissions, nil
}

// GetProfileTranslations returns all profile translations for a given profile with authorization check.
func (s *Service) GetProfileTranslations(
	ctx context.Context,
	profileSlug string,
) ([]*ProfileTx, error) {
	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	// Get all translations
	translations, err := s.repo.GetProfileTxByID(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	return translations, nil
}

// Profile Link Management

// CreateProfileLink creates a new profile link with authorization check.
func (s *Service) CreateProfileLink(
	ctx context.Context,
	userID string,
	profileSlug string,
	kind string,
	uri *string,
	title string,
	isHidden bool,
) (*ProfileLink, error) {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	// Get next order value
	existingLinks, err := s.repo.ListProfileLinksByProfileID(ctx, "en", profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	order := len(existingLinks) + 1

	// Generate new link ID
	linkID := s.idGenerator()

	// Create the link
	link, err := s.repo.CreateProfileLink(
		ctx,
		string(linkID),
		kind,
		profileID,
		order,
		uri,
		title,
		isHidden,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	return link, nil
}

// UpdateProfileLink updates an existing profile link with authorization check.
func (s *Service) UpdateProfileLink(
	ctx context.Context,
	userID string,
	profileSlug string,
	linkID string,
	kind string,
	order int,
	uri *string,
	title string,
	isHidden bool,
) (*ProfileLink, error) {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Verify the link exists and belongs to the profile
	existingLink, err := s.repo.GetProfileLink(ctx, linkID)
	if err != nil {
		return nil, fmt.Errorf("%w(linkID: %s): %w", ErrFailedToGetRecord, linkID, err)
	}

	if existingLink == nil {
		return nil, fmt.Errorf("%w: link %s not found", ErrFailedToGetRecord, linkID)
	}

	// Verify profile ownership through slug
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if existingLink.ProfileID != profileID {
		return nil, fmt.Errorf(
			"%w: link %s does not belong to profile %s",
			ErrUnauthorized,
			linkID,
			profileSlug,
		)
	}

	// For managed links, preserve the original title, URI, and kind
	// Users can only change order and visibility for managed links
	if existingLink.IsManaged {
		kind = existingLink.Kind
		title = existingLink.Title
		uri = existingLink.URI
	}

	// Update the link
	err = s.repo.UpdateProfileLink(ctx, linkID, kind, order, uri, title, isHidden)
	if err != nil {
		return nil, fmt.Errorf("%w(linkID: %s): %w", ErrFailedToUpdateRecord, linkID, err)
	}

	// Return updated link
	updatedLink, err := s.repo.GetProfileLink(ctx, linkID)
	if err != nil {
		return nil, fmt.Errorf("%w(linkID: %s): %w", ErrFailedToGetRecord, linkID, err)
	}

	return updatedLink, nil
}

// DeleteProfileLink soft-deletes a profile link with authorization check.
func (s *Service) DeleteProfileLink(
	ctx context.Context,
	userID string,
	profileSlug string,
	linkID string,
) error {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return err
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Verify the link exists and belongs to the profile
	existingLink, err := s.repo.GetProfileLink(ctx, linkID)
	if err != nil {
		return fmt.Errorf("%w(linkID: %s): %w", ErrFailedToGetRecord, linkID, err)
	}

	if existingLink == nil {
		return fmt.Errorf("%w: link %s not found", ErrFailedToGetRecord, linkID)
	}

	// Verify profile ownership through slug
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if existingLink.ProfileID != profileID {
		return fmt.Errorf(
			"%w: link %s does not belong to profile %s",
			ErrUnauthorized,
			linkID,
			profileSlug,
		)
	}

	// Delete the link
	err = s.repo.DeleteProfileLink(ctx, linkID)
	if err != nil {
		return fmt.Errorf("%w(linkID: %s): %w", ErrFailedToDeleteRecord, linkID, err)
	}

	return nil
}

// GetProfileLink retrieves a specific profile link with authorization check.
func (s *Service) GetProfileLink(
	ctx context.Context,
	userID string,
	profileSlug string,
	linkID string,
) (*ProfileLink, error) {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Get the link
	link, err := s.repo.GetProfileLink(ctx, linkID)
	if err != nil {
		return nil, fmt.Errorf("%w(linkID: %s): %w", ErrFailedToGetRecord, linkID, err)
	}

	if link == nil {
		return nil, fmt.Errorf("%w: link %s not found", ErrFailedToGetRecord, linkID)
	}

	// Verify profile ownership through slug
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if link.ProfileID != profileID {
		return nil, fmt.Errorf(
			"%w: link %s does not belong to profile %s",
			ErrUnauthorized,
			linkID,
			profileSlug,
		)
	}

	return link, nil
}

// ListProfileLinksBySlug retrieves all profile links for editing (includes hidden links).
func (s *Service) ListProfileLinksBySlug(
	ctx context.Context,
	userID string,
	profileSlug string,
) ([]*ProfileLink, error) {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	// Get brief links and convert them to full ProfileLink objects
	briefLinks, err := s.repo.ListProfileLinksByProfileID(ctx, "en", profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToListRecords, profileID, err)
	}

	// For now, we need to get each link individually to get the full data
	// This could be optimized with a dedicated repository method later
	links := make([]*ProfileLink, 0, len(briefLinks))

	for _, briefLink := range briefLinks {
		fullLink, err := s.repo.GetProfileLink(ctx, briefLink.ID)
		if err != nil {
			return nil, fmt.Errorf("%w(linkID: %s): %w", ErrFailedToGetRecord, briefLink.ID, err)
		}

		if fullLink != nil {
			links = append(links, fullLink)
		}
	}

	return links, nil
}

// Profile Page Management

// CreateProfilePage creates a new profile page with authorization check.
func (s *Service) CreateProfilePage(
	ctx context.Context,
	userID string,
	userKind string,
	profileSlug string,
	slug string,
	localeCode string,
	title string,
	summary string,
	content string,
	coverPictureURI *string,
	publishedAt *string,
) (*ProfilePage, error) {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Validate cover picture URI
	if err := validateOptionalURL(coverPictureURI); err != nil {
		return nil, err
	}

	// Non-admin users can only use URIs from allowed prefixes
	if userKind != "admin" {
		err := validateURIPrefixes(coverPictureURI, s.config.GetAllowedURIPrefixes())
		if err != nil {
			return nil, err
		}
	}

	// Validate slug availability
	slugResult, err := s.CheckPageSlugAvailability(ctx, localeCode, profileSlug, slug, nil, false)
	if err != nil {
		return nil, err
	}

	if !slugResult.Available && slugResult.Severity == SeverityError {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateRecord, slugResult.Message)
	}

	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	// Get next order value
	existingPages, err := s.repo.ListProfilePagesByProfileID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	order := len(existingPages) + 1

	// Generate new page ID
	pageID := s.idGenerator()

	// Create the page
	_, err = s.repo.CreateProfilePage(
		ctx,
		string(pageID),
		slug,
		profileID,
		order,
		coverPictureURI,
		publishedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	// Create the translation
	err = s.repo.CreateProfilePageTx(
		ctx,
		string(pageID),
		localeCode,
		title,
		summary,
		content,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	// Return the created page with translations
	fullPage, err := s.repo.GetProfilePageByProfileIDAndSlug(ctx, localeCode, profileID, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, string(pageID), err)
	}

	return fullPage, nil
}

// UpdateProfilePage updates an existing profile page with authorization check.
func (s *Service) UpdateProfilePage(
	ctx context.Context,
	userID string,
	userKind string,
	profileSlug string,
	pageID string,
	slug string,
	order int,
	coverPictureURI *string,
	publishedAt *string,
) (*ProfilePage, error) {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Validate cover picture URI
	if err := validateOptionalURL(coverPictureURI); err != nil {
		return nil, err
	}

	// Non-admin users can only use URIs from allowed prefixes
	if userKind != "admin" {
		err := validateURIPrefixes(coverPictureURI, s.config.GetAllowedURIPrefixes())
		if err != nil {
			return nil, err
		}
	}

	// Validate slug availability (exclude current page)
	slugResult, err := s.CheckPageSlugAvailability(ctx, "en", profileSlug, slug, &pageID, false)
	if err != nil {
		return nil, err
	}

	if !slugResult.Available && slugResult.Severity == SeverityError {
		return nil, fmt.Errorf("%w: %s", ErrFailedToUpdateRecord, slugResult.Message)
	}

	// Verify the page exists
	existingPage, err := s.repo.GetProfilePage(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, pageID, err)
	}

	if existingPage == nil {
		return nil, fmt.Errorf("%w: page %s not found", ErrFailedToGetRecord, pageID)
	}

	// Update the page
	err = s.repo.UpdateProfilePage(ctx, pageID, slug, order, coverPictureURI, publishedAt)
	if err != nil {
		return nil, fmt.Errorf("%w(pageID: %s): %w", ErrFailedToUpdateRecord, pageID, err)
	}

	// Return updated page
	updatedPage, err := s.repo.GetProfilePage(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, pageID, err)
	}

	return updatedPage, nil
}

// UpdateProfilePageTranslation updates profile page translation fields.
func (s *Service) UpdateProfilePageTranslation(
	ctx context.Context,
	userID string,
	profileSlug string,
	pageID string,
	localeCode string,
	title string,
	summary string,
	content string,
) error {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return err
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Verify the page exists
	existingPage, err := s.repo.GetProfilePage(ctx, pageID)
	if err != nil {
		return fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, pageID, err)
	}

	if existingPage == nil {
		return fmt.Errorf("%w: page %s not found", ErrFailedToGetRecord, pageID)
	}

	// Update the translation (use upsert to handle new locales)
	err = s.repo.UpsertProfilePageTx(ctx, pageID, localeCode, title, summary, content)
	if err != nil {
		return fmt.Errorf(
			"%w(pageID: %s, locale: %s): %w",
			ErrFailedToUpdateRecord,
			pageID,
			localeCode,
			err,
		)
	}

	return nil
}

// DeleteProfilePage soft-deletes a profile page with authorization check.
func (s *Service) DeleteProfilePage(
	ctx context.Context,
	userID string,
	profileSlug string,
	pageID string,
) error {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return err
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Verify the page exists
	existingPage, err := s.repo.GetProfilePage(ctx, pageID)
	if err != nil {
		return fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, pageID, err)
	}

	if existingPage == nil {
		return fmt.Errorf("%w: page %s not found", ErrFailedToGetRecord, pageID)
	}

	// Delete the page
	err = s.repo.DeleteProfilePage(ctx, pageID)
	if err != nil {
		return fmt.Errorf("%w(pageID: %s): %w", ErrFailedToDeleteRecord, pageID, err)
	}

	return nil
}

// GetProfilePage retrieves a specific profile page with authorization check.
func (s *Service) GetProfilePage(
	ctx context.Context,
	userID string,
	profileSlug string,
	pageID string,
	localeCode string,
) (*ProfilePage, error) {
	// Check authorization
	canEdit, err := s.CanUserEditProfile(ctx, userID, profileSlug)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit profile %s",
			ErrUnauthorized,
			userID,
			profileSlug,
		)
	}

	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	// Get the page
	page, err := s.repo.GetProfilePage(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, pageID, err)
	}

	if page == nil {
		return nil, fmt.Errorf("%w: page %s not found", ErrFailedToGetRecord, pageID)
	}

	// Get the full page with translations
	fullPage, err := s.repo.GetProfilePageByProfileIDAndSlug(ctx, localeCode, profileID, page.Slug)
	if err != nil {
		return nil, fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, pageID, err)
	}

	return fullPage, nil
}

// func (s *Service) Create(ctx context.Context, input *Profile) (*Profile, error) {
// 	record, err := s.repo.CreateProfile(ctx, input)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
// 	}

// 	return record, nil
// }

// Search performs a full-text search across profiles, stories, and profile pages.
// If profileSlug is provided, search is scoped to that profile only.
func (s *Service) Search(
	ctx context.Context,
	localeCode string,
	query string,
	profileSlug *string,
	limit int32,
) ([]*SearchResult, error) {
	if query == "" {
		return []*SearchResult{}, nil
	}

	results, err := s.repo.Search(ctx, localeCode, query, profileSlug, limit)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return results, nil
}

// GetProfileIDBySlug returns the profile ID for a given slug.
func (s *Service) GetProfileIDBySlug(ctx context.Context, slug string) (string, error) {
	return s.repo.GetProfileIDBySlug(ctx, slug)
}

// GetProfileLinkByRemoteID returns a profile link by its remote ID (e.g., YouTube channel ID).
func (s *Service) GetProfileLinkByRemoteID(
	ctx context.Context,
	profileID string,
	kind string,
	remoteID string,
) (*ProfileLink, error) {
	return s.repo.GetProfileLinkByRemoteID(ctx, profileID, kind, remoteID)
}

// UpdateProfileLinkOAuthTokens updates the OAuth tokens for an existing profile link.
func (s *Service) UpdateProfileLinkOAuthTokens(
	ctx context.Context,
	id string,
	publicID string,
	uri string,
	title string,
	authScope string,
	accessToken string,
	accessTokenExpiresAt *sql.NullTime,
	refreshToken *string,
) error {
	return s.repo.UpdateProfileLinkOAuthTokens(
		ctx, id, publicID, uri, title, authScope, accessToken, accessTokenExpiresAt, refreshToken,
	)
}

// CreateOAuthProfileLink creates a new OAuth-connected profile link.
func (s *Service) CreateOAuthProfileLink(
	ctx context.Context,
	id string,
	kind string,
	profileID string,
	order int,
	remoteID string,
	publicID string,
	uri string,
	title string,
	authProvider string,
	authScope string,
	accessToken string,
	accessTokenExpiresAt *sql.NullTime,
	refreshToken *string,
) (*ProfileLink, error) {
	return s.repo.CreateOAuthProfileLink(
		ctx, id, kind, profileID, order, remoteID, publicID, uri, title,
		authProvider, authScope, accessToken, accessTokenExpiresAt, refreshToken,
	)
}

// GetMaxProfileLinkOrder returns the maximum order value for profile links.
func (s *Service) GetMaxProfileLinkOrder(ctx context.Context, profileID string) (int, error) {
	return s.repo.GetMaxProfileLinkOrder(ctx, profileID)
}
