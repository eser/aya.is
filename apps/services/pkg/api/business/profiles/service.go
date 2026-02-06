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
	ErrInsufficientAccess   = errors.New("insufficient access level")
	ErrNoMembershipFound    = errors.New("no membership found")
	ErrNoIndividualProfile  = errors.New("user has no individual profile")
	ErrProfileNotFound      = errors.New("profile not found")
	ErrInvalidURI           = errors.New("invalid URI")
	ErrInvalidURIPrefix     = errors.New("URI must start with allowed prefix")
	ErrSearchFailed         = errors.New("search failed")
)

// FallbackLocaleCode is used when the requested locale translation is not available.
const FallbackLocaleCode = "en"

// SupportedLocaleCodes contains all locales supported by the platform.
var SupportedLocaleCodes = map[string]bool{ //nolint:gochecknoglobals
	"ar": true, "de": true, "en": true, "es": true,
	"fr": true, "it": true, "ja": true, "ko": true,
	"nl": true, "pt-PT": true, "ru": true, "tr": true,
	"zh-CN": true,
}

// IsValidLocale checks whether the given locale code is supported.
func IsValidLocale(localeCode string) bool {
	return SupportedLocaleCodes[localeCode]
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
	GetCustomDomainByDomain(ctx context.Context, domain string) (*ProfileCustomDomain, error)
	ListCustomDomainsByProfileID(
		ctx context.Context,
		profileID string,
	) ([]*ProfileCustomDomain, error)
	CreateCustomDomain(
		ctx context.Context,
		id string,
		profileID string,
		domain string,
		defaultLocale *string,
	) error
	UpdateCustomDomain(ctx context.Context, id string, domain string, defaultLocale *string) error
	DeleteCustomDomain(ctx context.Context, id string) error
	CheckProfileSlugExists(ctx context.Context, slug string) (bool, error)
	CheckProfileSlugExistsIncludingDeleted(ctx context.Context, slug string) (bool, error)
	CheckPageSlugExistsIncludingDeleted(
		ctx context.Context,
		profileID string,
		pageSlug string,
	) (bool, error)
	GetProfileBasicByID(ctx context.Context, id string) (*ProfileBrief, error)
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
	ListProfileLinksByProfileIDForEditing(
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
		hideRelations *bool,
		hideLinks *bool,
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
	CreateProfileMembership(
		ctx context.Context,
		id string,
		profileID string,
		memberProfileID *string,
		kind string,
		properties map[string]any,
	) error
	GetProfileOwnershipForUser(
		ctx context.Context,
		userID string,
		profileSlug string,
	) (*ProfileOwnership, error)
	GetUserBasicInfo(
		ctx context.Context,
		userID string,
	) (*UserBasicInfo, error)
	GetUserProfilePermissions(
		ctx context.Context,
		userID string,
	) ([]*ProfilePermission, error)
	GetProfileTxByID(
		ctx context.Context,
		profileID string,
	) ([]*ProfileTx, error)
	// Profile Link methods
	GetProfileLink(ctx context.Context, localeCode string, id string) (*ProfileLink, error)
	CreateProfileLink(
		ctx context.Context,
		id string,
		kind string,
		profileID string,
		order int,
		uri *string,
		isFeatured bool,
		visibility LinkVisibility,
	) (*ProfileLink, error)
	UpdateProfileLink(
		ctx context.Context,
		id string,
		kind string,
		order int,
		uri *string,
		isFeatured bool,
		visibility LinkVisibility,
	) error
	GetMembershipBetweenProfiles(
		ctx context.Context,
		profileID string,
		memberProfileID string,
	) (MembershipKind, error)
	DeleteProfileLink(ctx context.Context, id string) error
	ListFeaturedProfileLinksByProfileID(
		ctx context.Context,
		localeCode string,
		profileID string,
	) ([]*ProfileLinkBrief, error)
	ListAllProfileLinksByProfileID(
		ctx context.Context,
		localeCode string,
		profileID string,
	) ([]*ProfileLinkBrief, error)
	UpsertProfileLinkTx(
		ctx context.Context,
		profileLinkID string,
		localeCode string,
		title string,
		icon *string,
		group *string,
		description *string,
	) error
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
		authScope string,
		accessToken string,
		accessTokenExpiresAt *sql.NullTime,
		refreshToken *string,
	) error
	GetMaxProfileLinkOrder(ctx context.Context, profileID string) (int, error)
	// Admin methods
	ListAllProfilesForAdmin(
		ctx context.Context,
		localeCode string,
		filterKind string,
		limit int,
		offset int,
	) ([]*Profile, error)
	CountAllProfilesForAdmin(
		ctx context.Context,
		filterKind string,
	) (int64, error)
	GetAdminProfileBySlug(
		ctx context.Context,
		localeCode string,
		slug string,
	) (*Profile, error)
	// Membership management methods
	ListProfileMembershipsForSettings(
		ctx context.Context,
		localeCode string,
		profileID string,
	) ([]*ProfileMembershipWithMember, error)
	GetProfileMembershipByID(
		ctx context.Context,
		id string,
	) (*ProfileMembership, error)
	GetProfileMembershipByProfileAndMember(
		ctx context.Context,
		profileID string,
		memberProfileID string,
	) (*ProfileMembership, error)
	UpdateProfileMembership(
		ctx context.Context,
		id string,
		kind string,
	) error
	DeleteProfileMembership(
		ctx context.Context,
		id string,
	) error
	CountProfileOwners(
		ctx context.Context,
		profileID string,
	) (int64, error)
	SearchUsersForMembership(
		ctx context.Context,
		localeCode string,
		profileID string,
		query string,
	) ([]*UserSearchResult, error)
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

// CanViewLink checks if a viewer has permission to see a link based on its visibility.
// If viewerProfileID is empty, only public links are visible.
//
// Visibility levels and required membership:
//   - "public": visible to everyone (no membership required)
//   - "followers": visible to followers and above
//   - "sponsors": visible to sponsors and above
//   - "contributors": visible to contributors and above
//   - "maintainers": visible to maintainers and above
//   - "leads": visible to leads and above
//   - "owners": visible to owners only
//
// NOTE: Currently, public API endpoints (GET /profiles/{slug}, GET /profiles/{slug}/links)
// don't pass a viewerProfileID, so only public links are returned. To enable visibility
// filtering for logged-in users, the HTTP routes would need to optionally detect the
// session and pass the viewer's profile ID.
func (s *Service) CanViewLink(
	ctx context.Context,
	link *ProfileLinkBrief,
	targetProfileID string,
	viewerProfileID string,
) bool {
	// Public links are always visible
	if link.Visibility == LinkVisibilityPublic || link.Visibility == "" {
		return true
	}

	// Anonymous viewers can only see public links
	if viewerProfileID == "" {
		return false
	}

	// Get viewer's membership with the target profile
	membershipKind, err := s.repo.GetMembershipBetweenProfiles(
		ctx,
		targetProfileID,
		viewerProfileID,
	)
	if err != nil {
		// No membership found
		return false
	}

	// Check if membership level is sufficient
	minRequired := MinMembershipForVisibility[link.Visibility]
	if minRequired == "" {
		// No minimum required (shouldn't happen for non-public)
		return true
	}

	viewerLevel := MembershipKindLevel[membershipKind]
	requiredLevel := MembershipKindLevel[minRequired]

	return viewerLevel >= requiredLevel
}

// FilterVisibleLinks filters a list of links to only include those visible to the viewer.
func (s *Service) FilterVisibleLinks(
	ctx context.Context,
	links []*ProfileLinkBrief,
	targetProfileID string,
	viewerProfileID string,
) []*ProfileLinkBrief {
	result := make([]*ProfileLinkBrief, 0, len(links))

	for _, link := range links {
		if s.CanViewLink(ctx, link, targetProfileID, viewerProfileID) {
			result = append(result, link)
		}
	}

	return result
}

func (s *Service) GetBasicByID(ctx context.Context, id string) (*ProfileBrief, error) {
	record, err := s.repo.GetProfileBasicByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return record, nil
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
	return s.GetBySlugExWithViewer(ctx, localeCode, slug, "")
}

// GetBySlugExWithViewer returns a profile with children, filtering links based on viewer's membership.
func (s *Service) GetBySlugExWithViewer(
	ctx context.Context,
	localeCode string,
	slug string,
	viewerProfileID string,
) (*ProfileWithChildren, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	record, err := s.repo.GetProfileByID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	// Try fallback locale if primary locale has no translation
	if record == nil && FallbackLocaleCode != localeCode {
		record, err = s.repo.GetProfileByID(ctx, FallbackLocaleCode, profileID)
		if err != nil {
			return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
		}
	}

	if record == nil {
		return nil, nil //nolint:nilnil
	}

	pages, err := s.repo.ListProfilePagesByProfileID(ctx, localeCode, record.ID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	// Try fallback locale for pages if none found
	if len(pages) == 0 && FallbackLocaleCode != localeCode {
		pages, err = s.repo.ListProfilePagesByProfileID(ctx, FallbackLocaleCode, record.ID)
		if err != nil {
			return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
		}
	}

	// Only include featured links for the profile sidebar
	links, err := s.repo.ListFeaturedProfileLinksByProfileID(ctx, localeCode, record.ID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	// Filter links based on viewer's membership
	filteredLinks := s.FilterVisibleLinks(ctx, links, profileID, viewerProfileID)

	result := &ProfileWithChildren{
		Profile: record,
		Pages:   pages,
		Links:   filteredLinks,
	}

	return result, nil
}

func (s *Service) GetByCustomDomain(
	ctx context.Context,
	localeCode string,
	domain string,
) (*Profile, *ProfileCustomDomain, error) {
	customDomain, err := s.repo.GetCustomDomainByDomain(ctx, domain)
	if err != nil {
		return nil, nil, fmt.Errorf("%w(custom_domain: %s): %w", ErrFailedToGetRecord, domain, err)
	}

	if customDomain == nil {
		return nil, nil, nil
	}

	record, err := s.repo.GetProfileByID(ctx, localeCode, customDomain.ProfileID)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"%w(profile_id: %s): %w",
			ErrFailedToGetRecord,
			customDomain.ProfileID,
			err,
		)
	}

	return record, customDomain, nil
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

// AdminProfileListResult holds the result of listing profiles for admin.
type AdminProfileListResult struct {
	Data   []*Profile `json:"data"`
	Total  int64      `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

// ListAllProfilesForAdmin lists all profiles for admin with pagination.
func (s *Service) ListAllProfilesForAdmin(
	ctx context.Context,
	localeCode string,
	filterKind string,
	limit int,
	offset int,
) (*AdminProfileListResult, error) {
	profiles, err := s.repo.ListAllProfilesForAdmin(ctx, localeCode, filterKind, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	total, err := s.repo.CountAllProfilesForAdmin(ctx, filterKind)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return &AdminProfileListResult{
		Data:   profiles,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// GetAdminProfileBySlug gets a single profile by slug for admin.
func (s *Service) GetAdminProfileBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
) (*Profile, error) {
	profile, err := s.repo.GetAdminProfileBySlug(ctx, localeCode, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	return profile, nil
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
	if len(pageSlug) < 2 {
		return &SlugAvailabilityResult{
			Available: false,
			Message:   "Slug must be at least 2 characters",
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
			existsDeleted, delErr := s.repo.CheckPageSlugExistsIncludingDeleted(
				ctx,
				profileID,
				pageSlug,
			)
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
	if len(slug) < 2 {
		return &SlugAvailabilityResult{
			Available: false,
			Message:   "Slug must be at least 2 characters",
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

// CreateProfileMembership creates a membership record linking a member profile to a profile.
// This establishes the relationship (e.g., owner, maintainer) between profiles.
func (s *Service) CreateProfileMembership(
	ctx context.Context,
	profileID string,
	memberProfileID *string,
	kind string,
) error {
	membershipID := s.idGenerator()

	err := s.repo.CreateProfileMembership(
		ctx,
		string(membershipID),
		profileID,
		memberProfileID,
		kind,
		nil, // No additional properties for now
	)
	if err != nil {
		return fmt.Errorf("%w: membership: %w", ErrFailedToCreateRecord, err)
	}

	return nil
}

// ensureProfileCanProfileAccess checks if an origin profile has the required
// membership level on a target profile. Returns nil if access is granted.
func (s *Service) ensureProfileCanProfileAccess(
	ctx context.Context,
	targetProfileID string,
	originProfileID string,
	requiredLevel MembershipKind,
) error {
	if targetProfileID == originProfileID {
		return nil
	}

	membershipKind, err := s.repo.GetMembershipBetweenProfiles(
		ctx,
		targetProfileID,
		originProfileID,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if membershipKind == "" {
		return fmt.Errorf("%w: %w", ErrInsufficientAccess, ErrNoMembershipFound)
	}

	if MembershipKindLevel[membershipKind] < MembershipKindLevel[requiredLevel] {
		return fmt.Errorf("%w", ErrInsufficientAccess)
	}

	return nil
}

// ensureUserCanProfileAccess checks if a user has the required access level on
// a target profile. Resolves user kind (admin bypass) and individual profile
// internally using cached lookups.
func (s *Service) ensureUserCanProfileAccess(
	ctx context.Context,
	targetProfileID string,
	originUserID string,
	requiredLevel MembershipKind,
) error {
	userInfo, err := s.repo.GetUserBasicInfo(ctx, originUserID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if userInfo.Kind == "admin" {
		return nil
	}

	if userInfo.IndividualProfileID == nil {
		return fmt.Errorf("%w: %w", ErrInsufficientAccess, ErrNoIndividualProfile)
	}

	return s.ensureProfileCanProfileAccess(
		ctx,
		targetProfileID,
		*userInfo.IndividualProfileID,
		requiredLevel,
	)
}

// HasUserAccessToProfile checks if a user has the required access level on a
// profile identified by slug. Returns (true, nil) if access is granted.
func (s *Service) HasUserAccessToProfile(
	ctx context.Context,
	userID string,
	profileSlug string,
	requiredLevel MembershipKind,
) (bool, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return false, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return false, nil
	}

	err = s.ensureUserCanProfileAccess(ctx, profileID, userID, requiredLevel)
	if err != nil {
		if errors.Is(err, ErrInsufficientAccess) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Update updates profile main fields (profile_picture_uri, pronouns, properties).
func (s *Service) Update(
	ctx context.Context,
	localeCode string,
	userID string,
	userKind string,
	profileSlug string,
	profilePictureURI *string,
	pronouns *string,
	properties map[string]any,
	hideRelations *bool,
	hideLinks *bool,
) (*Profile, error) {
	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return nil, err
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

	// Update the profile
	err = s.repo.UpdateProfile(
		ctx,
		profileID,
		profilePictureURI,
		pronouns,
		properties,
		hideRelations,
		hideLinks,
	)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToUpdateRecord, profileID, err)
	}

	// Return updated profile
	profile, err := s.repo.GetProfileByID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	return profile, nil
}

// UpdateTranslation updates profile translation fields (title, description).
func (s *Service) UpdateTranslation(
	ctx context.Context,
	userID string,
	userKind string,
	profileSlug string,
	localeCode string,
	title string,
	description string,
	properties map[string]any,
) error {
	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return err
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
	localeCode string,
	userID string,
	userKind string,
	profileSlug string,
	kind string,
	uri *string,
	title string,
	icon *string,
	group *string,
	description *string,
	isFeatured bool,
	visibility LinkVisibility,
) (*ProfileLink, error) {
	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return nil, err
	}

	// Get next order value
	existingLinks, err := s.repo.ListProfileLinksByProfileID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	order := len(existingLinks) + 1

	// Generate new link ID
	linkID := s.idGenerator()

	// Default visibility to public if not specified
	if visibility == "" {
		visibility = LinkVisibilityPublic
	}

	// Create the link
	link, err := s.repo.CreateProfileLink(
		ctx,
		string(linkID),
		kind,
		profileID,
		order,
		uri,
		isFeatured,
		visibility,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	// Create the translation for this link
	err = s.repo.UpsertProfileLinkTx(
		ctx,
		string(linkID),
		localeCode,
		title,
		icon,
		group,
		description,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	// Set the title from the translation
	link.Title = title
	link.Icon = icon
	link.Group = group
	link.Description = description

	return link, nil
}

// UpdateProfileLink updates an existing profile link with authorization check.
func (s *Service) UpdateProfileLink(
	ctx context.Context,
	localeCode string,
	userID string,
	userKind string,
	profileSlug string,
	linkID string,
	kind string,
	order int,
	uri *string,
	title string,
	icon *string,
	group *string,
	description *string,
	isFeatured bool,
	visibility LinkVisibility,
) (*ProfileLink, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return nil, err
	}

	// Verify the link exists and belongs to the profile
	existingLink, err := s.repo.GetProfileLink(ctx, localeCode, linkID)
	if err != nil {
		return nil, fmt.Errorf("%w(linkID: %s): %w", ErrFailedToGetRecord, linkID, err)
	}

	if existingLink == nil {
		return nil, fmt.Errorf("%w: link %s not found", ErrFailedToGetRecord, linkID)
	}

	if existingLink.ProfileID != profileID {
		return nil, fmt.Errorf(
			"%w: link %s does not belong to profile %s",
			ErrUnauthorized,
			linkID,
			profileSlug,
		)
	}

	// For managed links, preserve the original URI and kind
	// Users can only change order, visibility for managed links
	if existingLink.IsManaged {
		kind = existingLink.Kind
		uri = existingLink.URI
	}

	// Default visibility to public if not specified
	if visibility == "" {
		visibility = LinkVisibilityPublic
	}

	// Update the link
	err = s.repo.UpdateProfileLink(
		ctx,
		linkID,
		kind,
		order,
		uri,
		isFeatured,
		visibility,
	)
	if err != nil {
		return nil, fmt.Errorf("%w(linkID: %s): %w", ErrFailedToUpdateRecord, linkID, err)
	}

	// Update the translation (only for non-managed links, but icon/group/description can be updated for all)
	if !existingLink.IsManaged {
		err = s.repo.UpsertProfileLinkTx(ctx, linkID, localeCode, title, icon, group, description)
		if err != nil {
			return nil, fmt.Errorf("%w(linkID: %s): %w", ErrFailedToUpdateRecord, linkID, err)
		}
	} else {
		// For managed links, still allow updating icon, group and description
		err = s.repo.UpsertProfileLinkTx(ctx, linkID, localeCode, existingLink.Title, icon, group, description)
		if err != nil {
			return nil, fmt.Errorf("%w(linkID: %s): %w", ErrFailedToUpdateRecord, linkID, err)
		}
	}

	// Return updated link
	updatedLink, err := s.repo.GetProfileLink(ctx, localeCode, linkID)
	if err != nil {
		return nil, fmt.Errorf("%w(linkID: %s): %w", ErrFailedToGetRecord, linkID, err)
	}

	return updatedLink, nil
}

// DeleteProfileLink soft-deletes a profile link with authorization check.
func (s *Service) DeleteProfileLink(
	ctx context.Context,
	userID string,
	userKind string,
	profileSlug string,
	linkID string,
) error {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return err
	}

	// Verify the link exists and belongs to the profile
	existingLink, err := s.repo.GetProfileLink(ctx, "en", linkID)
	if err != nil {
		return fmt.Errorf("%w(linkID: %s): %w", ErrFailedToGetRecord, linkID, err)
	}

	if existingLink == nil {
		return fmt.Errorf("%w: link %s not found", ErrFailedToGetRecord, linkID)
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
	localeCode string,
	userID string,
	profileSlug string,
	linkID string,
) (*ProfileLink, error) {
	hasAccess, err := s.HasUserAccessToProfile(ctx, userID, profileSlug, MembershipKindMaintainer)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, fmt.Errorf("%w", ErrInsufficientAccess)
	}

	// Get the link
	link, err := s.repo.GetProfileLink(ctx, localeCode, linkID)
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
	localeCode string,
	userID string,
	userKind string,
	profileSlug string,
) ([]*ProfileLink, error) {
	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return nil, err
	}

	// Get all links for editing (this is for the settings page)
	briefLinks, err := s.repo.ListProfileLinksByProfileIDForEditing(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToListRecords, profileID, err)
	}

	// For now, we need to get each link individually to get the full data
	// This could be optimized with a dedicated repository method later
	links := make([]*ProfileLink, 0, len(briefLinks))

	for _, briefLink := range briefLinks {
		fullLink, err := s.repo.GetProfileLink(ctx, localeCode, briefLink.ID)
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
	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return nil, err
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
	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return nil, err
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
	userKind string,
	profileSlug string,
	pageID string,
	localeCode string,
	title string,
	summary string,
	content string,
) error {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return err
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
	userKind string,
	profileSlug string,
	pageID string,
) error {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return err
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
	// Get profile ID
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return nil, err
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
		return nil, fmt.Errorf("%w: %w", ErrSearchFailed, err)
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
	localeCode string,
	publicID string,
	uri string,
	title string,
	authScope string,
	accessToken string,
	accessTokenExpiresAt *sql.NullTime,
	refreshToken *string,
) error {
	// Update OAuth tokens
	err := s.repo.UpdateProfileLinkOAuthTokens(
		ctx, id, publicID, uri, authScope, accessToken, accessTokenExpiresAt, refreshToken,
	)
	if err != nil {
		return err
	}

	// Update/create translation for the title (icon, group, description are nil for OAuth links)
	return s.repo.UpsertProfileLinkTx(ctx, id, localeCode, title, nil, nil, nil)
}

// CreateOAuthProfileLink creates a new OAuth-connected profile link.
func (s *Service) CreateOAuthProfileLink(
	ctx context.Context,
	id string,
	kind string,
	profileID string,
	order int,
	localeCode string,
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
	// Create the OAuth profile link
	link, err := s.repo.CreateOAuthProfileLink(
		ctx, id, kind, profileID, order, remoteID, publicID, uri,
		authProvider, authScope, accessToken, accessTokenExpiresAt, refreshToken,
	)
	if err != nil {
		return nil, err
	}

	// Create translation for the title (icon, group, description are nil for OAuth links)
	err = s.repo.UpsertProfileLinkTx(ctx, id, localeCode, title, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return link, nil
}

// GetMaxProfileLinkOrder returns the maximum order value for profile links.
func (s *Service) GetMaxProfileLinkOrder(ctx context.Context, profileID string) (int, error) {
	return s.repo.GetMaxProfileLinkOrder(ctx, profileID)
}

// ListFeaturedLinksBySlug returns featured profile links visible to the viewer.
func (s *Service) ListFeaturedLinksBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
	viewerProfileID string,
) ([]*ProfileLinkBrief, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	links, err := s.repo.ListFeaturedProfileLinksByProfileID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	// Filter links based on viewer's membership
	return s.FilterVisibleLinks(ctx, links, profileID, viewerProfileID), nil
}

// ListAllLinksBySlug returns all profile links visible to the viewer.
func (s *Service) ListAllLinksBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
	viewerProfileID string,
) ([]*ProfileLinkBrief, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	links, err := s.repo.ListAllProfileLinksByProfileID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	// Filter links based on viewer's membership
	return s.FilterVisibleLinks(ctx, links, profileID, viewerProfileID), nil
}

// UpsertProfileLinkTranslation creates or updates a profile link translation.
func (s *Service) UpsertProfileLinkTranslation(
	ctx context.Context,
	userID string,
	userKind string,
	profileSlug string,
	linkID string,
	localeCode string,
	title string,
	icon *string,
	group *string,
	description *string,
) error {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return err
	}

	// Verify the link exists
	existingLink, err := s.repo.GetProfileLink(ctx, localeCode, linkID)
	if err != nil {
		return fmt.Errorf("%w(linkID: %s): %w", ErrFailedToGetRecord, linkID, err)
	}

	if existingLink == nil {
		return fmt.Errorf("%w: link %s not found", ErrFailedToGetRecord, linkID)
	}

	if existingLink.ProfileID != profileID {
		return fmt.Errorf(
			"%w: link %s does not belong to profile %s",
			ErrUnauthorized,
			linkID,
			profileSlug,
		)
	}

	// Upsert the translation
	err = s.repo.UpsertProfileLinkTx(ctx, linkID, localeCode, title, icon, group, description)
	if err != nil {
		return fmt.Errorf(
			"%w(linkID: %s, locale: %s): %w",
			ErrFailedToUpdateRecord,
			linkID,
			localeCode,
			err,
		)
	}

	return nil
}

// Profile Membership Management

var (
	ErrCannotRemoveLastOwner      = errors.New("cannot remove the last owner")
	ErrCannotRemoveIndividualSelf = errors.New(
		"cannot remove yourself from your individual profile",
	)
	ErrMembershipNotFound       = errors.New("membership not found")
	ErrInvalidMembershipKind    = errors.New("invalid membership kind")
	ErrCannotModifyOwnRole      = errors.New("cannot modify your own membership role")
	ErrCannotAssignHigherRole   = errors.New("cannot assign a role higher than your own")
	ErrCannotModifyHigherMember = errors.New("cannot modify a member with higher role than yours")
)

// membershipRoleLevel defines the hierarchy of membership roles.
// Higher number = higher privilege level.
var membershipRoleLevel = map[string]int{
	"follower":    1,
	"sponsor":     2,
	"contributor": 3,
	"maintainer":  4,
	"lead":        5,
	"owner":       6,
}

// ListMembershipsForSettings lists all memberships for a profile (for settings page).
func (s *Service) ListMembershipsForSettings(
	ctx context.Context,
	localeCode string,
	userID string,
	userKind string,
	profileSlug string,
) ([]*ProfileMembershipWithMember, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return nil, err
	}

	memberships, err := s.repo.ListProfileMembershipsForSettings(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToListRecords, profileID, err)
	}

	return memberships, nil
}

// UpdateMembership updates the kind of an existing membership.
func (s *Service) UpdateMembership( //nolint:cyclop,funlen
	ctx context.Context,
	userID string,
	userKind string,
	userIndividualProfileID *string,
	profileSlug string,
	membershipID string,
	newKind string,
) error {
	// Validate kind
	validKinds := map[string]bool{
		"owner": true, "lead": true, "maintainer": true,
		"contributor": true, "sponsor": true, "follower": true,
	}
	if !validKinds[newKind] {
		return ErrInvalidMembershipKind
	}

	// Get the membership first to check ownership
	membership, err := s.repo.GetProfileMembershipByID(ctx, membershipID)
	if err != nil {
		return fmt.Errorf("%w(membershipID: %s): %w", ErrFailedToGetRecord, membershipID, err)
	}

	if membership == nil {
		return ErrMembershipNotFound
	}

	// SECURITY: Prevent self-modification (users cannot change their own role)
	// Exception: Admins can modify anything
	if userKind != "admin" && userIndividualProfileID != nil {
		if membership.MemberProfileID != nil &&
			*membership.MemberProfileID == *userIndividualProfileID {
			return ErrCannotModifyOwnRole
		}
	}

	// Check authorization
	isAdmin := userKind == "admin"

	hasAccess, accessErr := s.HasUserAccessToProfile(
		ctx,
		userID,
		profileSlug,
		MembershipKindMaintainer,
	)
	if accessErr != nil {
		return accessErr
	}

	if !hasAccess {
		return fmt.Errorf("%w", ErrInsufficientAccess)
	}

	// SECURITY: Non-admin users can only assign roles at or below their own level
	// and can only modify members at or below their level
	if !isAdmin && userIndividualProfileID != nil {
		// Get the requesting user's membership level
		userMembership, membershipErr := s.repo.GetProfileMembershipByProfileAndMember(
			ctx,
			membership.ProfileID,
			*userIndividualProfileID,
		)
		if membershipErr != nil {
			return fmt.Errorf("%w: %w", ErrFailedToGetRecord, membershipErr)
		}

		userLevel := 0
		if userMembership != nil {
			userLevel = membershipRoleLevel[userMembership.Kind]
		}

		// Check: Cannot assign a role higher than your own
		newRoleLevel := membershipRoleLevel[newKind]
		if newRoleLevel > userLevel {
			return ErrCannotAssignHigherRole
		}

		// Check: Cannot modify a member who has a higher or equal role than you
		// (except for demoting yourself, which is already blocked above)
		targetCurrentLevel := membershipRoleLevel[membership.Kind]
		if targetCurrentLevel >= userLevel {
			return ErrCannotModifyHigherMember
		}
	}

	// Check if trying to change to 'owner' on individual profile - not allowed
	if newKind == "owner" {
		profile, profileErr := s.repo.GetProfileByID(ctx, "en", membership.ProfileID)
		if profileErr != nil {
			return fmt.Errorf(
				"%w(profileID: %s): %w",
				ErrFailedToGetRecord,
				membership.ProfileID,
				profileErr,
			)
		}

		if profile.Kind == "individual" {
			return fmt.Errorf(
				"%w: cannot set 'owner' on individual profiles",
				ErrInvalidMembershipKind,
			)
		}
	}

	// If changing from owner to something else, check we're not removing the last owner
	if membership.Kind == "owner" && newKind != "owner" {
		ownerCount, countErr := s.repo.CountProfileOwners(ctx, membership.ProfileID)
		if countErr != nil {
			return fmt.Errorf("%w: %w", ErrFailedToGetRecord, countErr)
		}

		if ownerCount <= 1 {
			return ErrCannotRemoveLastOwner
		}
	}

	// Update the membership
	err = s.repo.UpdateProfileMembership(ctx, membershipID, newKind)
	if err != nil {
		return fmt.Errorf("%w(membershipID: %s): %w", ErrFailedToUpdateRecord, membershipID, err)
	}

	return nil
}

// DeleteMembership deletes a membership with validation.
func (s *Service) DeleteMembership(
	ctx context.Context,
	userID string,
	userKind string,
	userIndividualProfileID *string,
	profileSlug string,
	membershipID string,
) error {
	// Check authorization
	hasAccess, accessErr := s.HasUserAccessToProfile(
		ctx,
		userID,
		profileSlug,
		MembershipKindMaintainer,
	)
	if accessErr != nil {
		return accessErr
	}

	if !hasAccess {
		return fmt.Errorf("%w", ErrInsufficientAccess)
	}

	// Get the membership
	membership, err := s.repo.GetProfileMembershipByID(ctx, membershipID)
	if err != nil {
		return fmt.Errorf("%w(membershipID: %s): %w", ErrFailedToGetRecord, membershipID, err)
	}

	if membership == nil {
		return ErrMembershipNotFound
	}

	// Get profile to check if it's an individual profile
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	profile, err := s.repo.GetProfileByID(ctx, "en", profileID)
	if err != nil {
		return fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	// Prevent removing self from individual profile if it matches user's individual_profile_id
	if profile != nil && profile.Kind == "individual" {
		if userIndividualProfileID != nil && *userIndividualProfileID == profileID {
			if membership.MemberProfileID != nil &&
				*membership.MemberProfileID == *userIndividualProfileID {
				return ErrCannotRemoveIndividualSelf
			}
		}
	}

	// If removing an owner, check we're not removing the last owner
	if membership.Kind == "owner" {
		ownerCount, err := s.repo.CountProfileOwners(ctx, membership.ProfileID)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
		}

		if ownerCount <= 1 {
			return ErrCannotRemoveLastOwner
		}
	}

	// Delete the membership
	err = s.repo.DeleteProfileMembership(ctx, membershipID)
	if err != nil {
		return fmt.Errorf("%w(membershipID: %s): %w", ErrFailedToDeleteRecord, membershipID, err)
	}

	return nil
}

// SearchUsersForMembership searches users for adding as members.
func (s *Service) SearchUsersForMembership(
	ctx context.Context,
	localeCode string,
	userID string,
	userKind string,
	profileSlug string,
	query string,
) ([]*UserSearchResult, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	if err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer); err != nil {
		return nil, err
	}

	results, err := s.repo.SearchUsersForMembership(ctx, localeCode, profileID, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return results, nil
}

// AddMembership adds a new membership to a profile.
func (s *Service) AddMembership( //nolint:cyclop,funlen
	ctx context.Context,
	userID string,
	userKind string,
	userIndividualProfileID *string,
	profileSlug string,
	memberProfileID string,
	kind string,
) error {
	// Validate kind
	validKinds := map[string]bool{
		"owner": true, "lead": true, "maintainer": true,
		"contributor": true, "sponsor": true, "follower": true,
	}
	if !validKinds[kind] {
		return ErrInvalidMembershipKind
	}

	// Check authorization
	isAdmin := userKind == "admin"

	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return ErrProfileNotFound
	}

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return accessErr
	}

	// SECURITY: Non-admin users can only assign roles at or below their own level
	if !isAdmin && userIndividualProfileID != nil {
		// Get the requesting user's membership level
		userMembership, membershipErr := s.repo.GetProfileMembershipByProfileAndMember(
			ctx,
			profileID,
			*userIndividualProfileID,
		)
		if membershipErr != nil {
			return fmt.Errorf("%w: %w", ErrFailedToGetRecord, membershipErr)
		}

		userLevel := 0
		if userMembership != nil {
			userLevel = membershipRoleLevel[userMembership.Kind]
		}

		// Check: Cannot assign a role higher than your own
		newRoleLevel := membershipRoleLevel[kind]
		if newRoleLevel > userLevel {
			return ErrCannotAssignHigherRole
		}
	}

	// Check if trying to add 'owner' to individual profile - not allowed
	// Individual profiles have implicit ownership through user.individual_profile_id
	if kind == "owner" {
		profile, profileErr := s.repo.GetProfileByID(ctx, "en", profileID)
		if profileErr != nil {
			return fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, profileErr)
		}

		if profile.Kind == "individual" {
			return fmt.Errorf(
				"%w: cannot add 'owner' to individual profiles",
				ErrInvalidMembershipKind,
			)
		}
	}

	// Create the membership
	membershipID := s.idGenerator()

	err = s.repo.CreateProfileMembership(
		ctx,
		string(membershipID),
		profileID,
		&memberProfileID,
		kind,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	return nil
}
