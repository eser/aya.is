package profiles

import (
	"context"
	"database/sql" // TODO: replace sql.NullTime with *time.Time to remove database/sql dependency
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
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
	ErrDuplicateRecord      = errors.New("duplicate record")
	ErrInvalidInput         = errors.New("invalid input")
	ErrRelationsNotEnabled  = errors.New(
		"relations feature is not enabled for this profile",
	)
	ErrLinksNotEnabled               = errors.New("links feature is not enabled for this profile")
	ErrCannotDeleteTeamWithMembers   = errors.New("cannot delete team that has members")
	ErrCannotDeleteTeamWithResources = errors.New("cannot delete team that has resources")
	ErrReferralAlreadyExists         = errors.New("referral already exists for this profile")
	ErrCannotReferSelf               = errors.New("cannot refer yourself")
	ErrCannotReferExistingMember     = errors.New("cannot refer someone who is already a member")
	ErrReferralNotFound              = errors.New("referral not found")
	ErrInvalidVoteScore              = errors.New("vote score must be between 0 and 4")
	ErrReferralNotVoting             = errors.New("referral is not in voting status")
)

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
	ForbiddenSlugs string `conf:"forbidden_slugs" default:"about,admin,api,auth,communities,community,config,contact,contributions,dashboard,element,elements,events,faq,feed,guide,help,home,impressum,imprint,jobs,legal,login,logout,mailbox,new,news,null,organizations,orgs,people,policies,policy,privacy,product,products,profile,profiles,projects,register,root,search,services,settings,signin,signout,signup,site,stories,story,support,tag,tags,terms,tos,undefined,user,users,verify,wiki"`

	// DNSVerification holds the expected DNS targets for custom domain verification.
	DNSVerification DNSVerificationConfig `conf:"dns_verification"`
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
	GetFeatureRelationsVisibility(ctx context.Context, profileID string) (string, error)
	GetFeatureLinksVisibility(ctx context.Context, profileID string) (string, error)
	GetCustomDomainByDomain(ctx context.Context, domain string) (*ProfileCustomDomain, error)
	ListCustomDomainsByProfileID(
		ctx context.Context,
		profileID string,
	) ([]*ProfileCustomDomain, error)
	ListAllCustomDomains(ctx context.Context) ([]*ProfileCustomDomain, error)
	ListVerifiedCustomDomains(ctx context.Context) ([]*ProfileCustomDomain, error)
	CreateCustomDomain(
		ctx context.Context,
		id string,
		profileID string,
		domain string,
		defaultLocale *string,
	) error
	UpdateCustomDomain(ctx context.Context, id string, domain string, defaultLocale *string) error
	UpdateCustomDomainVerification(
		ctx context.Context,
		id string,
		status string,
		dnsVerifiedAt *time.Time,
		expiredAt *time.Time,
	) error
	UpdateCustomDomainWebserverSynced(ctx context.Context, id string, synced bool) error
	DeleteCustomDomain(ctx context.Context, id string) error
	CheckProfileSlugExists(ctx context.Context, slug string) (bool, error)
	CheckProfileSlugExistsIncludingDeleted(ctx context.Context, slug string) (bool, error)
	CheckPageSlugExistsIncludingDeleted(
		ctx context.Context,
		profileID string,
		pageSlug string,
	) (bool, error)
	GetProfileIdentifierByID(ctx context.Context, id string) (*ProfileBrief, error)
	GetProfileByID(
		ctx context.Context,
		localeCode string,
		id string,
	) (*Profile, error)
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
	ListProfilePagesByProfileIDForViewer(
		ctx context.Context,
		localeCode string,
		profileID string,
		viewerUserID *string,
	) ([]*ProfilePageBrief, error)
	GetProfilePageByProfileIDAndSlug(
		ctx context.Context,
		localeCode string,
		profileID string,
		pageSlug string,
	) (*ProfilePage, error)
	GetProfilePageByProfileIDAndSlugForViewer(
		ctx context.Context,
		localeCode string,
		profileID string,
		pageSlug string,
		viewerUserID *string,
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
		defaultLocale string,
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
		featureRelations *string,
		featureLinks *string,
		featureQA *string,
		featureDiscussions *string,
		optionStoryDiscussionsByDefault *bool,
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
	GetUserBriefInfo(
		ctx context.Context,
		userID string,
	) (*UserBriefInfo, error)
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
		addedByProfileID *string,
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
		addedByProfileID *string,
		visibility string,
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
		visibility string,
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
	DeleteProfilePageTx(ctx context.Context, profilePageID string, localeCode string) error
	ListProfilePageTxLocales(ctx context.Context, profilePageID string) ([]string, error)
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
	IsManagedProfileLinkRemoteIDInUse(
		ctx context.Context,
		kind string,
		remoteID string,
		excludeProfileID string,
	) (bool, error)
	ClearNonManagedProfileLinkRemoteID(
		ctx context.Context,
		profileID string,
		kind string,
		remoteID string,
	) error
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
		properties map[string]any,
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

	// Profile resources
	ListProfileResourcesByProfileID(
		ctx context.Context,
		profileID string,
	) ([]*ProfileResource, error)
	GetProfileResourceByID(
		ctx context.Context,
		id string,
	) (*ProfileResource, error)
	GetProfileResourceByRemoteID(
		ctx context.Context,
		profileID string,
		kind string,
		remoteID string,
	) (*ProfileResource, error)
	CreateProfileResource(
		ctx context.Context,
		id string,
		profileID string,
		kind string,
		isManaged bool,
		remoteID *string,
		publicID *string,
		url *string,
		title string,
		description *string,
		properties any,
		addedByProfileID string,
	) (*ProfileResource, error)
	SoftDeleteProfileResource(
		ctx context.Context,
		id string,
	) error
	UpdateProfileResourceProperties(
		ctx context.Context,
		id string,
		properties any,
	) error
	UpdateProfileMembershipProperties(
		ctx context.Context,
		id string,
		properties any,
	) error

	// Managed GitHub link
	GetManagedGitHubLinkByProfileID(
		ctx context.Context,
		profileID string,
	) (*ManagedGitHubLink, error)

	// Profile Team methods
	ListProfileTeamsWithMemberCount(
		ctx context.Context,
		profileID string,
	) ([]*ProfileTeam, error)
	CreateProfileTeam(
		ctx context.Context,
		id string,
		profileID string,
		name string,
		description *string,
	) (*ProfileTeam, error)
	UpdateProfileTeam(
		ctx context.Context,
		id string,
		name string,
		description *string,
	) error
	DeleteProfileTeam(
		ctx context.Context,
		id string,
	) error
	CountProfileTeamMembers(
		ctx context.Context,
		teamID string,
	) (int64, error)
	ListMembershipTeams(
		ctx context.Context,
		membershipID string,
	) ([]*ProfileTeam, error)
	SetMembershipTeams(
		ctx context.Context,
		membershipID string,
		teamIDs []string,
		idGenerator func() string,
	) error
	CountProfileTeamResources(
		ctx context.Context,
		teamID string,
	) (int64, error)
	ListResourceTeams(
		ctx context.Context,
		resourceID string,
	) ([]*ProfileTeam, error)
	SetResourceTeams(
		ctx context.Context,
		resourceID string,
		teamIDs []string,
		idGenerator func() string,
	) error

	// Referral methods
	CreateProfileMembershipReferral(
		ctx context.Context,
		id string,
		profileID string,
		referredProfileID string,
		referrerMembershipID string,
	) (*ProfileMembershipReferral, error)
	GetProfileMembershipReferralByID(
		ctx context.Context,
		id string,
	) (*ProfileMembershipReferral, error)
	GetProfileMembershipReferralByProfileAndReferred(
		ctx context.Context,
		profileID string,
		referredProfileID string,
	) (*ProfileMembershipReferral, error)
	ListProfileMembershipReferralsByProfileID(
		ctx context.Context,
		localeCode string,
		profileID string,
		viewerMembershipID *string,
	) ([]*ProfileMembershipReferral, error)
	UpsertReferralVote(
		ctx context.Context,
		id string,
		referralID string,
		voterMembershipID string,
		score int16,
		comment *string,
	) (*ReferralVote, error)
	ListReferralVotes(
		ctx context.Context,
		localeCode string,
		referralID string,
	) ([]*ReferralVote, error)
	UpdateReferralVoteCount(
		ctx context.Context,
		referralID string,
	) error
	InsertReferralTeam(
		ctx context.Context,
		id string,
		referralID string,
		teamID string,
	) error
	ListReferralTeams(
		ctx context.Context,
		referralID string,
	) ([]*ProfileTeam, error)
	GetReferralVoteBreakdown(
		ctx context.Context,
		referralID string,
	) (map[int]int, error)
}

// WebserverSyncer is the port for syncing custom domains to webserver infrastructure.
// Implementations: Coolify adapter (current), nginx, Vercel, etc.
type WebserverSyncer interface {
	GetCurrentDomains(ctx context.Context) ([]string, error)
	UpdateDomains(ctx context.Context, domains []string) error
	RestartApplication(ctx context.Context) error
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

func (s *Service) GetIdentifierByID(ctx context.Context, id string) (*ProfileBrief, error) {
	record, err := s.repo.GetProfileIdentifierByID(ctx, id)
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

	// Try fallback locale if primary locale has no translation
	if record == nil && localeCode != "en" {
		record, err = s.repo.GetProfileByID(ctx, "en", id)
		if err != nil {
			return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
		}
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
	if record == nil && localeCode != "en" {
		record, err = s.repo.GetProfileByID(ctx, "en", profileID)
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
	if len(pages) == 0 && localeCode != "en" {
		pages, err = s.repo.ListProfilePagesByProfileID(
			ctx,
			"en",
			record.ID,
		)
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

	record, err := s.repo.GetProfileByID(
		ctx,
		localeCode,
		customDomain.ProfileID,
	)
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

	page, err := s.repo.GetProfilePageByProfileIDAndSlug(
		ctx,
		localeCode,
		profileID,
		pageSlug,
	)
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

func (s *Service) ListPagesBySlugForViewer(
	ctx context.Context,
	localeCode string,
	slug string,
	viewerUserID *string,
) ([]*ProfilePageBrief, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	pages, err := s.repo.ListProfilePagesByProfileIDForViewer(
		ctx,
		localeCode,
		profileID,
		viewerUserID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	return pages, nil
}

func (s *Service) GetPageBySlugForViewer(
	ctx context.Context,
	localeCode string,
	slug string,
	pageSlug string,
	viewerUserID *string,
) (*ProfilePage, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	page, err := s.repo.GetProfilePageByProfileIDAndSlugForViewer(
		ctx,
		localeCode,
		profileID,
		pageSlug,
		viewerUserID,
	)
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

// GetBySlugExWithViewerUser returns a profile with children, filtering pages by viewer's access
// and links based on viewer's membership. Takes viewerUserID instead of viewerProfileID.
func (s *Service) GetBySlugExWithViewerUser( //nolint:cyclop
	ctx context.Context,
	localeCode string,
	slug string,
	viewerUserID *string,
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
	if record == nil && localeCode != "en" {
		record, err = s.repo.GetProfileByID(ctx, "en", profileID)
		if err != nil {
			return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
		}
	}

	if record == nil {
		return nil, nil //nolint:nilnil
	}

	pages, err := s.repo.ListProfilePagesByProfileIDForViewer(
		ctx,
		localeCode,
		record.ID,
		viewerUserID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	// Try fallback locale for pages if none found
	if len(pages) == 0 && localeCode != "en" {
		pages, err = s.repo.ListProfilePagesByProfileIDForViewer(
			ctx,
			"en",
			record.ID,
			viewerUserID,
		)
		if err != nil {
			return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
		}
	}

	// Only include featured links for the profile sidebar
	links, err := s.repo.ListFeaturedProfileLinksByProfileID(ctx, localeCode, record.ID)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	// Filter links based on viewer's membership (uses empty viewerProfileID for anonymous)
	filteredLinks := s.FilterVisibleLinks(ctx, links, profileID, "")

	result := &ProfileWithChildren{
		Profile: record,
		Pages:   pages,
		Links:   filteredLinks,
	}

	return result, nil
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

	page, err := s.repo.GetProfilePageByProfileIDAndSlug(
		ctx,
		localeCode,
		profileID,
		pageSlug,
	)
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

	visibility, err := s.repo.GetFeatureRelationsVisibility(ctx, profileID)
	if err != nil {
		return cursors.Cursored[[]*ProfileMembership]{}, fmt.Errorf(
			"%w: %w",
			ErrFailedToGetRecord,
			err,
		)
	}

	if visibility == "disabled" {
		return cursors.Cursored[[]*ProfileMembership]{}, ErrRelationsNotEnabled
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

	visibility, err := s.repo.GetFeatureRelationsVisibility(ctx, profileID)
	if err != nil {
		return cursors.Cursored[[]*ProfileMembership]{}, fmt.Errorf(
			"%w: %w",
			ErrFailedToGetRecord,
			err,
		)
	}

	if visibility == "disabled" {
		return cursors.Cursored[[]*ProfileMembership]{}, ErrRelationsNotEnabled
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
	userID string,
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

	// Create the main profile record with the request locale as default
	err := s.repo.CreateProfile(
		ctx,
		string(profileID),
		slug,
		kind,
		localeCode,
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

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileCreated,
		EntityType: "profile",
		EntityID:   string(profileID),
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
	})

	return profile, nil
}

// CreateProfileMembership creates a membership record linking a member profile to a profile.
// This establishes the relationship (e.g., owner, maintainer) between profiles.
func (s *Service) CreateProfileMembership(
	ctx context.Context,
	userID string,
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

	payload := map[string]any{"profile_id": profileID}
	if memberProfileID != nil {
		payload["member_profile_id"] = *memberProfileID
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileMembershipCreated,
		EntityType: "membership",
		EntityID:   string(membershipID),
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    payload,
	})

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
	userInfo, err := s.repo.GetUserBriefInfo(ctx, originUserID)
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

// GetProfilePermissions returns the viewer's edit permission and membership kind
// in a single query chain (slug→profileID, user→individualProfile, membership lookup).
func (s *Service) GetProfilePermissions(
	ctx context.Context,
	userID string,
	profileSlug string,
) (canEdit bool, viewerMembershipKind *string, err error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return false, nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return false, nil, nil
	}

	userInfo, err := s.repo.GetUserBriefInfo(ctx, userID)
	if err != nil {
		return false, nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	// Admins can always edit — still look up membership for badge display
	if userInfo.Kind == "admin" {
		if userInfo.IndividualProfileID != nil {
			kind, mkErr := s.repo.GetMembershipBetweenProfiles(
				ctx,
				profileID,
				*userInfo.IndividualProfileID,
			)
			if mkErr == nil && kind != "" {
				kindStr := string(kind)

				return true, &kindStr, nil
			}
		}

		return true, nil, nil
	}

	if userInfo.IndividualProfileID == nil {
		return false, nil, nil
	}

	// Own profile — can edit, no "membership" with self
	if profileID == *userInfo.IndividualProfileID {
		return true, nil, nil
	}

	// Single DB query: get membership kind between profiles
	kind, mkErr := s.repo.GetMembershipBetweenProfiles(
		ctx,
		profileID,
		*userInfo.IndividualProfileID,
	)
	if mkErr != nil || kind == "" {
		return false, nil, nil
	}

	kindStr := string(kind)
	canEdit = MembershipKindLevel[kind] >= MembershipKindLevel[MembershipKindMaintainer]

	return canEdit, &kindStr, nil
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
	featureRelations *string,
	featureLinks *string,
	featureQA *string,
	featureDiscussions *string,
	optionStoryDiscussionsByDefault *bool,
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

	// Validate profile picture URI (empty string means "remove picture")
	if profilePictureURI == nil || *profilePictureURI != "" {
		err := validateOptionalURL(profilePictureURI)
		if err != nil {
			return nil, err
		}

		// Non-admin users can only use URIs from allowed prefixes
		if userKind != "admin" {
			err := validateURIPrefixes(profilePictureURI, s.config.GetAllowedURIPrefixes())
			if err != nil {
				return nil, err
			}
		}
	}

	// Validate module visibility values
	for _, v := range []*string{featureRelations, featureLinks, featureQA, featureDiscussions} {
		if v != nil {
			switch ModuleVisibility(*v) {
			case ModuleVisibilityPublic, ModuleVisibilityHidden, ModuleVisibilityDisabled:
			default:
				return nil, ErrInvalidInput
			}
		}
	}

	// Update the profile
	err = s.repo.UpdateProfile(
		ctx,
		profileID,
		profilePictureURI,
		pronouns,
		properties,
		featureRelations,
		featureLinks,
		featureQA,
		featureDiscussions,
		optionStoryDiscussionsByDefault,
	)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToUpdateRecord, profileID, err)
	}

	// Return updated profile
	profile, err := s.repo.GetProfileByID(ctx, localeCode, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileUpdated,
		EntityType: "profile",
		EntityID:   profileID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
	})

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

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileTranslationUpdated,
		EntityType: "profile",
		EntityID:   profileID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"locale_code": localeCode},
	})

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

	// Get the user's individual profile for added_by tracking
	userInfo, userInfoErr := s.repo.GetUserBriefInfo(ctx, userID)

	var addedByProfileID *string
	if userInfoErr == nil && userInfo != nil && userInfo.IndividualProfileID != nil {
		addedByProfileID = userInfo.IndividualProfileID
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
		addedByProfileID,
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

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileLinkCreated,
		EntityType: "profile_link",
		EntityID:   string(linkID),
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"profile_id": profileID},
	})

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

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileLinkUpdated,
		EntityType: "profile_link",
		EntityID:   linkID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
	})

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

	// Three-tier authorization: admin > original adder > maintainer+
	canDelete := userKind == "admin"

	if !canDelete && existingLink.AddedByProfileID != nil && userID != "" {
		userInfo, upErr := s.repo.GetUserBriefInfo(ctx, userID)
		if upErr == nil && userInfo != nil && userInfo.IndividualProfileID != nil {
			if *userInfo.IndividualProfileID == *existingLink.AddedByProfileID {
				canDelete = true
			}
		}
	}

	if !canDelete {
		err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
		if err != nil {
			return err
		}

		canDelete = true
	}

	if !canDelete {
		return ErrUnauthorized
	}

	// Delete the link
	err = s.repo.DeleteProfileLink(ctx, linkID)
	if err != nil {
		return fmt.Errorf("%w(linkID: %s): %w", ErrFailedToDeleteRecord, linkID, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileLinkDeleted,
		EntityType: "profile_link",
		EntityID:   linkID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
	})

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

	// Determine if the current user can remove each link
	if userID != "" {
		userInfo, userErr := s.repo.GetUserBriefInfo(ctx, userID)

		for _, l := range links {
			canRemove := false

			// Site admins can always remove
			if userKind == "admin" {
				canRemove = true
			} else if userErr == nil && userInfo != nil && userInfo.IndividualProfileID != nil {
				// Check membership level
				membershipKind, mkErr := s.repo.GetMembershipBetweenProfiles(
					ctx, profileID, *userInfo.IndividualProfileID,
				)
				if mkErr == nil && membershipKind != "" {
					level, ok := MembershipKindLevel[membershipKind]
					if ok && level >= MembershipKindLevel[MembershipKindMaintainer] {
						canRemove = true
					}
				}
			}

			// The original adder can remove their own links
			if !canRemove && userErr == nil && userInfo != nil &&
				userInfo.IndividualProfileID != nil && l.AddedByProfileID != nil {
				if *userInfo.IndividualProfileID == *l.AddedByProfileID {
					canRemove = true
				}
			}

			l.CanRemove = canRemove
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
	visibility string,
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

	// Get the user's individual profile for added_by tracking
	userInfo, userInfoErr := s.repo.GetUserBriefInfo(ctx, userID)

	var addedByProfileID *string
	if userInfoErr == nil && userInfo != nil && userInfo.IndividualProfileID != nil {
		addedByProfileID = userInfo.IndividualProfileID
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
	existingPages, err := s.repo.ListProfilePagesByProfileID(
		ctx,
		localeCode,
		profileID,
	)
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
		addedByProfileID,
		visibility,
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
	fullPage, err := s.repo.GetProfilePageByProfileIDAndSlug(
		ctx,
		localeCode,
		profileID,
		slug,
	)
	if err != nil {
		return nil, fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, string(pageID), err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfilePageCreated,
		EntityType: "profile_page",
		EntityID:   string(pageID),
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"profile_id": profileID, "slug": slug},
	})

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
	visibility string,
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
	err = s.repo.UpdateProfilePage(
		ctx,
		pageID,
		slug,
		order,
		coverPictureURI,
		publishedAt,
		visibility,
	)
	if err != nil {
		return nil, fmt.Errorf("%w(pageID: %s): %w", ErrFailedToUpdateRecord, pageID, err)
	}

	// Return updated page
	updatedPage, err := s.repo.GetProfilePage(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, pageID, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfilePageUpdated,
		EntityType: "profile_page",
		EntityID:   pageID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
	})

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

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfilePageTranslationUpdated,
		EntityType: "profile_page",
		EntityID:   pageID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"locale_code": localeCode},
	})

	return nil
}

// DeleteProfilePageTranslation deletes a specific translation for a profile page.
func (s *Service) DeleteProfilePageTranslation(
	ctx context.Context,
	userID string,
	userKind string,
	profileSlug string,
	pageID string,
	localeCode string,
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

	err = s.repo.DeleteProfilePageTx(ctx, pageID, localeCode)
	if err != nil {
		return fmt.Errorf(
			"%w(pageID: %s, locale: %s): %w",
			ErrFailedToDeleteRecord,
			pageID,
			localeCode,
			err,
		)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfilePageTranslationDeleted,
		EntityType: "profile_page",
		EntityID:   pageID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload:    map[string]any{"locale_code": localeCode},
	})

	return nil
}

// ListProfilePageTranslationLocales returns locale codes that have translations for a page.
func (s *Service) ListProfilePageTranslationLocales(
	ctx context.Context,
	profileSlug string,
	pageID string,
) ([]string, error) {
	locales, err := s.repo.ListProfilePageTxLocales(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("%w(pageID: %s): %w", ErrFailedToListRecords, pageID, err)
	}

	result := make([]string, 0, len(locales))
	for _, l := range locales {
		result = append(result, strings.TrimSpace(l))
	}

	return result, nil
}

// GetProfilePageTranslationContent returns the translation content for a specific locale.
func (s *Service) GetProfilePageTranslationContent(
	ctx context.Context,
	profileSlug string,
	pageID string,
	localeCode string,
) (string, string, string, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return "", "", "", fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return "", "", "", ErrProfileNotFound
	}

	page, err := s.repo.GetProfilePageByProfileIDAndSlug(ctx, localeCode, profileID, "")
	if err != nil {
		// Try to get by ID directly
		pages, listErr := s.repo.ListProfilePagesByProfileID(ctx, localeCode, profileID)
		if listErr != nil {
			return "", "", "", fmt.Errorf(
				"%w(pageID: %s): %w",
				ErrFailedToGetRecord,
				pageID,
				listErr,
			)
		}

		for _, p := range pages {
			if p.ID == pageID {
				return p.Title, p.Summary, "", nil
			}
		}

		return "", "", "", fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, pageID, err)
	}

	if page == nil {
		return "", "", "", fmt.Errorf("%w: page not found", ErrFailedToGetRecord)
	}

	// Check if this is actual content for the requested locale
	actualLocale := strings.TrimSpace(page.LocaleCode)
	if actualLocale != localeCode {
		return "", "", "", fmt.Errorf(
			"%w: no translation for locale %s",
			ErrFailedToGetRecord,
			localeCode,
		)
	}

	return page.Title, page.Summary, page.Content, nil
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

	// Verify the page exists
	existingPage, err := s.repo.GetProfilePage(ctx, pageID)
	if err != nil {
		return fmt.Errorf("%w(pageID: %s): %w", ErrFailedToGetRecord, pageID, err)
	}

	if existingPage == nil {
		return fmt.Errorf("%w: page %s not found", ErrFailedToGetRecord, pageID)
	}

	// Three-tier authorization: admin > original adder > maintainer+
	canDelete := userKind == "admin"

	if !canDelete && existingPage.AddedByProfileID != nil && userID != "" {
		userInfo, upErr := s.repo.GetUserBriefInfo(ctx, userID)
		if upErr == nil && userInfo != nil && userInfo.IndividualProfileID != nil {
			if *userInfo.IndividualProfileID == *existingPage.AddedByProfileID {
				canDelete = true
			}
		}
	}

	if !canDelete {
		err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
		if err != nil {
			return err
		}

		canDelete = true
	}

	if !canDelete {
		return ErrUnauthorized
	}

	// Delete the page
	err = s.repo.DeleteProfilePage(ctx, pageID)
	if err != nil {
		return fmt.Errorf("%w(pageID: %s): %w", ErrFailedToDeleteRecord, pageID, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfilePageDeleted,
		EntityType: "profile_page",
		EntityID:   pageID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
	})

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
	fullPage, err := s.repo.GetProfilePageByProfileIDAndSlug(
		ctx,
		localeCode,
		profileID,
		page.Slug,
	)
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

// IsManagedProfileLinkRemoteIDInUse checks if a remote_id is already used by another
// profile's active link of the same kind (e.g., prevents connecting the same
// GitHub account to multiple profiles).
func (s *Service) IsManagedProfileLinkRemoteIDInUse(
	ctx context.Context,
	kind string,
	remoteID string,
	excludeProfileID string,
) (bool, error) {
	return s.repo.IsManagedProfileLinkRemoteIDInUse(ctx, kind, remoteID, excludeProfileID)
}

// ClearNonManagedProfileLinkRemoteID nulls out remote_id on non-managed links
// to avoid unique constraint violations when creating a new managed link.
func (s *Service) ClearNonManagedProfileLinkRemoteID(
	ctx context.Context,
	profileID string,
	kind string,
	remoteID string,
) error {
	return s.repo.ClearNonManagedProfileLinkRemoteID(ctx, profileID, kind, remoteID)
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
	properties map[string]any,
) (*ProfileLink, error) {
	// Create the OAuth profile link
	link, err := s.repo.CreateOAuthProfileLink(
		ctx, id, kind, profileID, order, remoteID, publicID, uri,
		authProvider, authScope, accessToken, accessTokenExpiresAt, refreshToken,
		properties,
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

	visibility, err := s.repo.GetFeatureLinksVisibility(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if visibility == "disabled" {
		return nil, ErrLinksNotEnabled
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
	"member":      3,
	"contributor": 4,
	"maintainer":  5,
	"lead":        6,
	"owner":       7,
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

	memberships, err := s.repo.ListProfileMembershipsForSettings(
		ctx,
		localeCode,
		profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToListRecords, profileID, err)
	}

	// Non-admin users should not see admin-only membership kinds (sponsor, follower)
	isAdmin := userKind == "admin"
	if !isAdmin {
		filtered := make([]*ProfileMembershipWithMember, 0, len(memberships))

		for _, m := range memberships {
			kind := MembershipKind(m.Kind)
			if kind == MembershipKindSponsor || kind == MembershipKindFollower {
				continue
			}

			filtered = append(filtered, m)
		}

		memberships = filtered
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
		"contributor": true, "member": true, "sponsor": true, "follower": true,
	}
	if !validKinds[newKind] {
		return ErrInvalidMembershipKind
	}

	// Check authorization
	isAdmin := userKind == "admin"

	// SECURITY: Only admins can assign sponsor or follower roles
	if !isAdmin && (newKind == "sponsor" || newKind == "follower") {
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
		} else if *userIndividualProfileID == membership.ProfileID {
			// Implicit owner of their own individual profile
			userLevel = membershipRoleLevel["owner"]
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

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileMembershipUpdated,
		EntityType: "membership",
		EntityID:   membershipID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_id":        membership.ProfileID,
			"member_profile_id": membership.MemberProfileID,
			"kind":              newKind,
			"last_properties": map[string]any{
				"kind": membership.Kind,
			},
		},
	})

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

	// If already a follower, perform a real delete.
	// Otherwise, demote to follower so they still follow the organization.
	if membership.Kind == string(MembershipKindFollower) {
		err = s.repo.DeleteProfileMembership(ctx, membershipID)
		if err != nil {
			return fmt.Errorf(
				"%w(membershipID: %s): %w",
				ErrFailedToDeleteRecord,
				membershipID,
				err,
			)
		}

		s.auditService.Record(ctx, events.AuditParams{
			EventType:  events.ProfileMembershipDeleted,
			EntityType: "membership",
			EntityID:   membershipID,
			ActorID:    &userID,
			ActorKind:  events.ActorUser,
			Payload: map[string]any{
				"profile_id":        membership.ProfileID,
				"member_profile_id": membership.MemberProfileID,
				"last_properties": map[string]any{
					"kind": membership.Kind,
				},
			},
		})
	} else {
		err = s.repo.UpdateProfileMembership(ctx, membershipID, string(MembershipKindFollower))
		if err != nil {
			return fmt.Errorf("%w(membershipID: %s): %w", ErrFailedToUpdateRecord, membershipID, err)
		}

		s.auditService.Record(ctx, events.AuditParams{
			EventType:  events.ProfileMembershipUpdated,
			EntityType: "membership",
			EntityID:   membershipID,
			ActorID:    &userID,
			ActorKind:  events.ActorUser,
			Payload: map[string]any{
				"profile_id":        membership.ProfileID,
				"member_profile_id": membership.MemberProfileID,
				"kind":              string(MembershipKindFollower),
				"last_properties": map[string]any{
					"kind": membership.Kind,
				},
			},
		})
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

	results, err := s.repo.SearchUsersForMembership(
		ctx,
		localeCode,
		profileID,
		query,
	)
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
		"contributor": true, "member": true, "sponsor": true, "follower": true,
	}
	if !validKinds[kind] {
		return ErrInvalidMembershipKind
	}

	// Check authorization
	isAdmin := userKind == "admin"

	// SECURITY: Only admins can assign sponsor or follower roles
	if !isAdmin && (kind == "sponsor" || kind == "follower") {
		return ErrInvalidMembershipKind
	}

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
		} else if *userIndividualProfileID == profileID {
			// Implicit owner of their own individual profile
			userLevel = membershipRoleLevel["owner"]
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

	// Check if the user already has a follower/sponsor membership that should be promoted
	existing, existingErr := s.repo.GetProfileMembershipByProfileAndMember(
		ctx,
		profileID,
		memberProfileID,
	)
	if existingErr != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, existingErr)
	}

	if existing != nil {
		// Existing membership found — promote it to the new kind
		err = s.repo.UpdateProfileMembership(ctx, existing.ID, kind)
		if err != nil {
			return fmt.Errorf("%w(membershipID: %s): %w", ErrFailedToUpdateRecord, existing.ID, err)
		}

		s.auditService.Record(ctx, events.AuditParams{
			EventType:  events.ProfileMembershipUpdated,
			EntityType: "membership",
			EntityID:   existing.ID,
			ActorID:    &userID,
			ActorKind:  events.ActorUser,
			Payload: map[string]any{
				"profile_id":        profileID,
				"member_profile_id": memberProfileID,
				"kind":              kind,
				"last_properties": map[string]any{
					"kind": existing.Kind,
				},
			},
		})

		return nil
	}

	// Create a new membership
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

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileMembershipCreated,
		EntityType: "membership",
		EntityID:   string(membershipID),
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_id":        profileID,
			"member_profile_id": memberProfileID,
			"kind":              kind,
		},
	})

	return nil
}

// FollowProfile creates a follower membership for the viewer on the given profile.
func (s *Service) FollowProfile(
	ctx context.Context,
	userID string,
	userIndividualProfileID string,
	profileSlug string,
) error {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return ErrProfileNotFound
	}

	// Cannot follow your own profile
	if userIndividualProfileID == profileID {
		return fmt.Errorf("%w: cannot follow your own profile", ErrInvalidMembershipKind)
	}

	// Check if already has a membership
	existing, err := s.repo.GetProfileMembershipByProfileAndMember(
		ctx,
		profileID,
		userIndividualProfileID,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if existing != nil {
		// Already has a relationship — no action needed
		return nil
	}

	membershipID := s.idGenerator()

	err = s.repo.CreateProfileMembership(
		ctx,
		string(membershipID),
		profileID,
		&userIndividualProfileID,
		string(MembershipKindFollower),
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileMembershipCreated,
		EntityType: "membership",
		EntityID:   string(membershipID),
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_id":        profileID,
			"member_profile_id": userIndividualProfileID,
			"kind":              string(MembershipKindFollower),
		},
	})

	return nil
}

// UnfollowProfile removes a follower membership. Only followers can unfollow;
// higher roles must be changed through the access settings.
func (s *Service) UnfollowProfile(
	ctx context.Context,
	userID string,
	userIndividualProfileID string,
	profileSlug string,
) error {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return ErrProfileNotFound
	}

	existing, err := s.repo.GetProfileMembershipByProfileAndMember(
		ctx,
		profileID,
		userIndividualProfileID,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if existing == nil {
		return ErrMembershipNotFound
	}

	// Only followers can self-unfollow
	if existing.Kind != string(MembershipKindFollower) {
		return fmt.Errorf("%w: only followers can unfollow", ErrInsufficientAccess)
	}

	err = s.repo.DeleteProfileMembership(ctx, existing.ID)
	if err != nil {
		return fmt.Errorf("%w(membershipID: %s): %w", ErrFailedToDeleteRecord, existing.ID, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileMembershipDeleted,
		EntityType: "membership",
		EntityID:   existing.ID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_id":        profileID,
			"member_profile_id": userIndividualProfileID,
			"last_properties": map[string]any{
				"kind": existing.Kind,
			},
		},
	})

	return nil
}

// ListProfileResources returns all resources associated with a profile,
// annotating each with whether the current user can remove it.
func (s *Service) ListProfileResources(
	ctx context.Context,
	locale string,
	userID string,
	userKind string,
	profileSlug string,
) ([]*ProfileResource, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	resources, err := s.repo.ListProfileResourcesByProfileID(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	// Populate teams for each resource
	for _, r := range resources {
		teams, teamsErr := s.repo.ListResourceTeams(ctx, r.ID)
		if teamsErr == nil && teams != nil {
			r.Teams = teams
		} else {
			r.Teams = []*ProfileTeam{}
		}
	}

	// Determine if the current user can remove each resource
	if userID != "" {
		userInfo, userErr := s.repo.GetUserBriefInfo(ctx, userID)

		for _, r := range resources {
			canRemove := false

			// Site admins can always remove
			if userKind == "admin" {
				canRemove = true
			} else if userErr == nil && userInfo != nil && userInfo.IndividualProfileID != nil {
				// Check membership level
				membershipKind, mkErr := s.repo.GetMembershipBetweenProfiles(
					ctx, profileID, *userInfo.IndividualProfileID,
				)
				if mkErr == nil && membershipKind != "" {
					level, ok := MembershipKindLevel[membershipKind]
					if ok && level >= MembershipKindLevel[MembershipKindMaintainer] {
						canRemove = true
					}
				}
			}

			// The original adder can remove their own resources
			if !canRemove && userErr == nil && userInfo != nil &&
				userInfo.IndividualProfileID != nil {
				if *userInfo.IndividualProfileID == r.AddedByProfileID {
					canRemove = true
				}
			}

			r.CanRemove = canRemove
		}
	}

	return resources, nil
}

// CreateProfileResource creates a new resource linked to a profile.
func (s *Service) CreateProfileResource(
	ctx context.Context,
	locale string,
	userID string,
	userKind string,
	profileSlug string,
	kind string,
	isManaged bool,
	remoteID *string,
	publicID *string,
	url *string,
	title string,
	description *string,
	properties any,
) (*ProfileResource, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return nil, ErrProfileNotFound
	}

	// Check permission - must be maintainer+ or admin
	if userKind != "admin" {
		err := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
		if err != nil {
			return nil, err
		}
	}

	// Check for duplicates
	if remoteID != nil {
		existing, _ := s.repo.GetProfileResourceByRemoteID(ctx, profileID, kind, *remoteID)
		if existing != nil {
			return nil, fmt.Errorf("%w: resource already exists", ErrDuplicateRecord)
		}
	}

	// Get the user's individual profile for added_by
	userInfo, err := s.repo.GetUserBriefInfo(ctx, userID)
	if err != nil || userInfo == nil || userInfo.IndividualProfileID == nil {
		return nil, fmt.Errorf("%w: could not find user profile", ErrFailedToGetRecord)
	}

	id := string(s.idGenerator())

	resource, err := s.repo.CreateProfileResource(
		ctx,
		id,
		profileID,
		kind,
		isManaged,
		remoteID,
		publicID,
		url,
		title,
		description,
		properties,
		*userInfo.IndividualProfileID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	return resource, nil
}

// DeleteProfileResource soft-deletes a profile resource with authorization check.
func (s *Service) DeleteProfileResource(
	ctx context.Context,
	locale string,
	userID string,
	userKind string,
	profileSlug string,
	resourceID string,
) error {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return ErrProfileNotFound
	}

	// Get the resource
	resource, err := s.repo.GetProfileResourceByID(ctx, resourceID)
	if err != nil {
		return fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, resourceID, err)
	}

	if resource == nil {
		return ErrProfileNotFound
	}

	// Check authorization
	canDelete := userKind == "admin"

	// Site admins can always delete

	// The original adder can delete
	if !canDelete && userID != "" {
		userInfo, upErr := s.repo.GetUserBriefInfo(ctx, userID)
		if upErr == nil && userInfo != nil && userInfo.IndividualProfileID != nil {
			if *userInfo.IndividualProfileID == resource.AddedByProfileID {
				canDelete = true
			}
		}
	}

	// Maintainers+ can delete
	if !canDelete {
		err = s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
		if err == nil {
			canDelete = true
		}
	}

	if !canDelete {
		return ErrUnauthorized
	}

	err = s.repo.SoftDeleteProfileResource(ctx, resourceID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToDeleteRecord, err)
	}

	return nil
}

// GetManagedGitHubLink returns the managed GitHub link for a profile.
func (s *Service) GetManagedGitHubLink(
	ctx context.Context,
	profileID string,
) (*ManagedGitHubLink, error) {
	link, err := s.repo.GetManagedGitHubLinkByProfileID(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToGetRecord, profileID, err)
	}

	return link, nil
}

// ListTeams lists all teams for a profile with member counts. Requires maintainer access.
func (s *Service) ListTeams(
	ctx context.Context,
	userID string,
	profileSlug string,
) ([]*ProfileTeam, error) {
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

	teams, err := s.repo.ListProfileTeamsWithMemberCount(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w(profileID: %s): %w", ErrFailedToListRecords, profileID, err)
	}

	return teams, nil
}

// CreateTeam creates a new team for a profile. Requires maintainer access.
func (s *Service) CreateTeam(
	ctx context.Context,
	userID string,
	profileSlug string,
	name string,
	description *string,
) (*ProfileTeam, error) {
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

	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("%w: team name is required", ErrInvalidInput)
	}

	id := string(s.idGenerator())

	team, err := s.repo.CreateProfileTeam(ctx, id, profileID, name, description)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	return team, nil
}

// UpdateTeam updates an existing team. Requires maintainer access.
func (s *Service) UpdateTeam(
	ctx context.Context,
	userID string,
	profileSlug string,
	teamID string,
	name string,
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

	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: team name is required", ErrInvalidInput)
	}

	err = s.repo.UpdateProfileTeam(ctx, teamID, name, description)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	return nil
}

// DeleteTeam deletes a team. Requires maintainer access. Fails if team has members.
func (s *Service) DeleteTeam(
	ctx context.Context,
	userID string,
	profileSlug string,
	teamID string,
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

	memberCount, err := s.repo.CountProfileTeamMembers(ctx, teamID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if memberCount > 0 {
		return ErrCannotDeleteTeamWithMembers
	}

	resourceCount, err := s.repo.CountProfileTeamResources(ctx, teamID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if resourceCount > 0 {
		return ErrCannotDeleteTeamWithResources
	}

	err = s.repo.DeleteProfileTeam(ctx, teamID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToDeleteRecord, err)
	}

	return nil
}

// SetMembershipTeams assigns teams to a membership. Requires maintainer access.
func (s *Service) SetMembershipTeams(
	ctx context.Context,
	userID string,
	profileSlug string,
	membershipID string,
	teamIDs []string,
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

	idGen := func() string { return string(s.idGenerator()) }

	err = s.repo.SetMembershipTeams(ctx, membershipID, teamIDs, idGen)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	return nil
}

// SetResourceTeams assigns teams to a resource. Requires maintainer access.
func (s *Service) SetResourceTeams(
	ctx context.Context,
	userID string,
	profileSlug string,
	resourceID string,
	teamIDs []string,
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

	idGen := func() string { return string(s.idGenerator()) }

	err = s.repo.SetResourceTeams(ctx, resourceID, teamIDs, idGen)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	return nil
}

// CreateReferral creates a new membership referral. The referrer must be member+ on the profile.
func (s *Service) CreateReferral( //nolint:cyclop
	ctx context.Context,
	userID string,
	profileSlug string,
	referredProfileSlug string,
	teamIDs []string,
) (*ProfileMembershipReferral, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	err = s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMember)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.repo.GetUserBriefInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if userInfo.IndividualProfileID == nil {
		return nil, fmt.Errorf("%w: %w", ErrInsufficientAccess, ErrNoIndividualProfile)
	}

	referrerMembership, err := s.repo.GetProfileMembershipByProfileAndMember(
		ctx, profileID, *userInfo.IndividualProfileID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNoMembershipFound, err)
	}

	referredProfileID, err := s.repo.GetProfileIDBySlug(ctx, referredProfileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: referred profile not found: %w", ErrProfileNotFound, err)
	}

	if *userInfo.IndividualProfileID == referredProfileID {
		return nil, ErrCannotReferSelf
	}

	existingMembership, _ := s.repo.GetProfileMembershipByProfileAndMember(
		ctx, profileID, referredProfileID,
	)
	if existingMembership != nil &&
		MembershipKindLevel[MembershipKind(existingMembership.Kind)] >= MembershipKindLevel[MembershipKindMember] {
		return nil, ErrCannotReferExistingMember
	}

	existingReferral, _ := s.repo.GetProfileMembershipReferralByProfileAndReferred(
		ctx, profileID, referredProfileID,
	)
	if existingReferral != nil {
		return nil, ErrReferralAlreadyExists
	}

	if len(teamIDs) > 0 {
		referrerTeams, teamsErr := s.repo.ListMembershipTeams(ctx, referrerMembership.ID)
		if teamsErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, teamsErr)
		}

		referrerTeamSet := make(map[string]bool, len(referrerTeams))
		for _, team := range referrerTeams {
			referrerTeamSet[team.ID] = true
		}

		for _, teamID := range teamIDs {
			if !referrerTeamSet[teamID] {
				return nil, fmt.Errorf(
					"%w: team %s does not belong to referrer",
					ErrInvalidInput,
					teamID,
				)
			}
		}
	}

	referralID := string(s.idGenerator())

	referral, err := s.repo.CreateProfileMembershipReferral(
		ctx, referralID, profileID, referredProfileID, referrerMembership.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	for _, teamID := range teamIDs {
		teamRecordID := string(s.idGenerator())

		insertErr := s.repo.InsertReferralTeam(ctx, teamRecordID, referralID, teamID)
		if insertErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, insertErr)
		}
	}

	teams, _ := s.repo.ListReferralTeams(ctx, referralID)
	referral.Teams = teams

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileReferralCreated,
		EntityType: "referral",
		EntityID:   referralID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"profile_id":             profileID,
			"referred_profile":       referredProfileSlug,
			"referrer_membership_id": referrerMembership.ID,
		},
	})

	return referral, nil
}

// ListReferrals lists all active referrals for a profile. Member+ only.
func (s *Service) ListReferrals(
	ctx context.Context,
	localeCode string,
	userID string,
	profileSlug string,
) ([]*ProfileMembershipReferral, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	err = s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMember)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.repo.GetUserBriefInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	var viewerMembershipID *string

	if userInfo.IndividualProfileID != nil {
		membership, membershipErr := s.repo.GetProfileMembershipByProfileAndMember(
			ctx, profileID, *userInfo.IndividualProfileID,
		)
		if membershipErr == nil && membership != nil {
			viewerMembershipID = &membership.ID
		}
	}

	referrals, err := s.repo.ListProfileMembershipReferralsByProfileID(
		ctx, localeCode, profileID, viewerMembershipID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	for _, referral := range referrals {
		teams, teamsErr := s.repo.ListReferralTeams(ctx, referral.ID)
		if teamsErr == nil {
			referral.Teams = teams
		}
	}

	return referrals, nil
}

// VoteOnReferral casts or updates a vote on a referral. Member+ only.
func (s *Service) VoteOnReferral(
	ctx context.Context,
	userID string,
	profileSlug string,
	referralID string,
	score int16,
	comment *string,
) (*ReferralVote, error) {
	const minScore, maxScore = 0, 4

	if score < minScore || score > maxScore {
		return nil, ErrInvalidVoteScore
	}

	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	err = s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMember)
	if err != nil {
		return nil, err
	}

	referral, err := s.repo.GetProfileMembershipReferralByID(ctx, referralID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReferralNotFound, err)
	}

	if referral.ProfileID != profileID {
		return nil, ErrReferralNotFound
	}

	if referral.Status != ReferralStatusVoting {
		return nil, ErrReferralNotVoting
	}

	userInfo, err := s.repo.GetUserBriefInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if userInfo.IndividualProfileID == nil {
		return nil, fmt.Errorf("%w: %w", ErrInsufficientAccess, ErrNoIndividualProfile)
	}

	voterMembership, err := s.repo.GetProfileMembershipByProfileAndMember(
		ctx, profileID, *userInfo.IndividualProfileID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNoMembershipFound, err)
	}

	voteID := string(s.idGenerator())

	vote, err := s.repo.UpsertReferralVote(
		ctx, voteID, referralID, voterMembership.ID, score, comment,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	_ = s.repo.UpdateReferralVoteCount(ctx, referralID)

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileReferralVoted,
		EntityType: "referral",
		EntityID:   referralID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		Payload: map[string]any{
			"score":   score,
			"vote_id": vote.ID,
		},
	})

	return vote, nil
}

// GetReferralVotes gets all votes for a referral. Member+ only.
func (s *Service) GetReferralVotes(
	ctx context.Context,
	localeCode string,
	userID string,
	profileSlug string,
	referralID string,
) ([]*ReferralVote, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	err = s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMember)
	if err != nil {
		return nil, err
	}

	referral, err := s.repo.GetProfileMembershipReferralByID(ctx, referralID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReferralNotFound, err)
	}

	if referral.ProfileID != profileID {
		return nil, ErrReferralNotFound
	}

	votes, err := s.repo.ListReferralVotes(ctx, localeCode, referralID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return votes, nil
}
