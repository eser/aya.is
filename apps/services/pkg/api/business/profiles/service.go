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
	ErrCandidateAlreadyExists        = errors.New("candidate already exists for this profile")
	ErrCannotReferSelf               = errors.New("cannot refer yourself")
	ErrCannotReferExistingMember     = errors.New("cannot refer someone who is already a member")
	ErrCannotReferNonIndividual      = errors.New("only individual profiles can be referred")
	ErrCandidateNotFound             = errors.New("candidate not found")
	ErrInvalidVoteScore              = errors.New("vote score must be between 0 and 4")
	ErrCandidateNotVoting            = errors.New("candidate is not in voting status")
	ErrInvalidStatusTransition       = errors.New("invalid candidate status transition")
	ErrApplicationsNotEnabled        = errors.New(
		"applications feature is not enabled for this profile",
	)
	ErrNoApplicationForm = errors.New(
		"no active application form configured for this profile",
	)
	ErrAlreadyApplied             = errors.New("you have already applied to this organization")
	ErrMissingRequiredField       = errors.New("a required form field is missing")
	ErrInvalidFieldType           = errors.New("invalid field type")
	ErrInvalidResponsesVisibility = errors.New(
		"responses visibility must be 'members' or 'leads'",
	)
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

// String constants used across the service.
const (
	UserKindAdmin         = "admin"
	ProfileKindIndividual = "individual"
)

// minSlugLength is the minimum allowed length for slugs.
const minSlugLength = 2

// SlugAvailabilityResult holds the result of a slug availability check.
type SlugAvailabilityResult struct {
	Available bool   `json:"available"`
	Message   string `json:"message,omitempty"`
	Severity  string `json:"severity,omitempty"` // "error" | "warning" | ""
}

// Config holds the profiles service configuration.
type Config struct {
	// AllowedURIPrefixes is a comma-separated list of allowed URI prefixes.
	AllowedURIPrefixes string `conf:"allowed_uri_prefixes" default:"https://objects.aya.is/,https://avatars.githubusercontent.com/"` //nolint:lll

	// ForbiddenSlugs is a comma-separated list of reserved slugs
	// that cannot be used as profile slugs.
	ForbiddenSlugs string `conf:"forbidden_slugs" default:"about,admin,api,auth,communities,community,config,contact,contributions,dashboard,element,elements,events,faq,feed,guide,help,home,impressum,imprint,jobs,legal,login,logout,mailbox,new,news,null,organizations,orgs,people,policies,policy,privacy,product,products,profile,profiles,projects,register,root,search,services,settings,signin,signout,signup,site,stories,story,support,tag,tags,terms,tos,undefined,user,users,verify,wiki"` //nolint:lll

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
		featureReferrals *string,
		featureApplications *string,
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
	ListOnlineProfileLinks(
		ctx context.Context,
		localeCode string,
	) ([]*LiveStreamInfo, error)
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
	InvalidateMembershipKindCache(
		ctx context.Context,
		profileID string,
		memberProfileID string,
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

	// Candidate methods
	CreateProfileMembershipCandidate(
		ctx context.Context,
		id string,
		profileID string,
		referredProfileID string,
		referrerMembershipID *string,
		source string,
		applicantMessage *string,
	) (*ProfileMembershipCandidate, error)
	GetProfileMembershipCandidateByID(
		ctx context.Context,
		id string,
	) (*ProfileMembershipCandidate, error)
	GetProfileMembershipCandidateByProfileAndReferred(
		ctx context.Context,
		profileID string,
		referredProfileID string,
	) (*ProfileMembershipCandidate, error)
	ListProfileMembershipCandidatesByProfileID(
		ctx context.Context,
		localeCode string,
		profileID string,
		viewerMembershipID *string,
	) ([]*ProfileMembershipCandidate, error)
	UpsertCandidateVote(
		ctx context.Context,
		id string,
		candidateID string,
		voterMembershipID string,
		score int16,
		comment *string,
	) (*CandidateVote, error)
	ListCandidateVotes(
		ctx context.Context,
		localeCode string,
		candidateID string,
	) ([]*CandidateVote, error)
	UpdateCandidateVoteCount(
		ctx context.Context,
		candidateID string,
	) error
	InsertCandidateTeam(
		ctx context.Context,
		id string,
		candidateID string,
		teamID string,
	) error
	ListCandidateTeams(
		ctx context.Context,
		candidateID string,
	) ([]*ProfileTeam, error)
	GetCandidateVoteBreakdown(
		ctx context.Context,
		candidateID string,
	) (map[int]int, error)
	UpdateCandidateStatus(
		ctx context.Context,
		candidateID string,
		profileID string,
		status CandidateStatus,
	) error
	UpdateCandidateApplicantMessage(
		ctx context.Context,
		candidateID string,
		applicantMessage *string,
	) error
	SoftDeleteCandidate(ctx context.Context, candidateID string) error

	// Application form methods
	GetActiveApplicationForm(
		ctx context.Context,
		profileID string,
	) (*ApplicationForm, error)
	GetApplicationFormByProfileID(
		ctx context.Context,
		profileID string,
	) (*ApplicationForm, string, error) // returns form, featureApplications, error
	CreateApplicationForm(
		ctx context.Context,
		formID string,
		profileID string,
		presetKey *string,
		responsesVisibility string,
	) (*ApplicationForm, error)
	UpdateApplicationForm(
		ctx context.Context,
		formID string,
		presetKey *string,
		responsesVisibility string,
	) error
	DeactivateApplicationForms(ctx context.Context, profileID string) error
	ListApplicationFormFields(
		ctx context.Context,
		formID string,
	) ([]*ApplicationFormField, error)
	CreateApplicationFormField(
		ctx context.Context,
		fieldID string,
		formID string,
		label string,
		fieldType string,
		isRequired bool,
		sortOrder int,
		placeholder *string,
	) (*ApplicationFormField, error)
	DeleteApplicationFormFields(ctx context.Context, formID string) error
	CreateCandidateResponse(
		ctx context.Context,
		responseID string,
		candidateID string,
		formFieldID string,
		value string,
	) error
	ListCandidateResponses(
		ctx context.Context,
		candidateID string,
	) ([]*CandidateFormResponse, error)
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
	minRequired := GetMinMembershipForVisibility()[link.Visibility]
	if minRequired == "" {
		// No minimum required (shouldn't happen for non-public)
		return true
	}

	viewerLevel := GetMembershipKindLevel()[membershipKind]
	requiredLevel := GetMembershipKindLevel()[minRequired]

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

func (s *Service) GetByID(
	ctx context.Context,
	localeCode string,
	id string, //nolint:varnamelen
) (*Profile, error) {
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
func (s *Service) GetBySlugExWithViewer( //nolint:cyclop
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
func (s *Service) GetBySlugExWithViewerUser( //nolint:cyclop,funlen
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

// checkDeletedPageSlug checks whether a page slug was previously used by a deleted page.
func (s *Service) checkDeletedPageSlug(
	ctx context.Context,
	profileID string,
	pageSlug string,
) (*SlugAvailabilityResult, error) {
	existsDeleted, delErr := s.repo.CheckPageSlugExistsIncludingDeleted(
		ctx,
		profileID,
		pageSlug,
	)
	if delErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, delErr)
	}

	if existsDeleted {
		return &SlugAvailabilityResult{
			Available: false,
			Message:   "This slug was previously used",
			Severity:  SeverityError,
		}, nil
	}

	return nil, nil //nolint:nilnil
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
	if len(pageSlug) < minSlugLength {
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
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if page == nil {
		return s.resolvePageSlugForMissingPage(ctx, profileID, pageSlug, includeDeleted)
	}

	// If we're editing and the slug belongs to the same page, it's available
	if excludePageID != nil && page.ID == *excludePageID {
		return &SlugAvailabilityResult{
			Available: true,
			Message:   "",
			Severity:  "",
		}, nil
	}

	return &SlugAvailabilityResult{
		Available: false,
		Message:   "This slug is already taken",
		Severity:  SeverityError,
	}, nil
}

// resolvePageSlugForMissingPage handles slug availability when no active page exists.
func (s *Service) resolvePageSlugForMissingPage(
	ctx context.Context,
	profileID string,
	pageSlug string,
	includeDeleted bool,
) (*SlugAvailabilityResult, error) {
	if includeDeleted {
		result, err := s.checkDeletedPageSlug(ctx, profileID, pageSlug)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	// Slug is available
	return &SlugAvailabilityResult{
		Available: true,
		Message:   "",
		Severity:  "",
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

// resolveProfileIDWithRelationsCheck resolves a profile slug to an ID and
// verifies that the relations feature is enabled.
func (s *Service) resolveProfileIDWithRelationsCheck(
	ctx context.Context,
	slug string,
) (string, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return "", fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	visibility, err := s.repo.GetFeatureRelationsVisibility(ctx, profileID)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if visibility == string(ModuleVisibilityDisabled) {
		return "", ErrRelationsNotEnabled
	}

	return profileID, nil
}

func (s *Service) ListProfileContributionsBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*ProfileMembership], error) {
	profileID, err := s.resolveProfileIDWithRelationsCheck(ctx, slug)
	if err != nil {
		return cursors.Cursored[[]*ProfileMembership]{}, err
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
	profileID, err := s.resolveProfileIDWithRelationsCheck(ctx, slug)
	if err != nil {
		return cursors.Cursored[[]*ProfileMembership]{}, err
	}

	memberships, err := s.repo.ListProfileMembers(
		ctx,
		localeCode,
		profileID,
		[]string{"organization", ProfileKindIndividual},
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
	if len(slug) < minSlugLength {
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
		existsDeleted, delErr := s.repo.CheckProfileSlugExistsIncludingDeleted(ctx, slug)
		if delErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, delErr)
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
		Message:   "",
		Severity:  "",
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
		SessionID:  nil,
		Payload:    nil,
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
		SessionID:  nil,
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

	if GetMembershipKindLevel()[membershipKind] < GetMembershipKindLevel()[requiredLevel] {
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

	if userInfo.Kind == UserKindAdmin {
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
func (s *Service) GetProfilePermissions( //nolint:cyclop,nonamedreturns
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
	if userInfo.Kind == UserKindAdmin {
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
	if mkErr != nil {
		return false, nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, mkErr)
	}

	if kind == "" {
		return false, nil, nil
	}

	kindStr := string(kind)
	canEdit = GetMembershipKindLevel()[kind] >= GetMembershipKindLevel()[MembershipKindMaintainer]

	return canEdit, &kindStr, nil
}

// Update updates profile main fields (profile_picture_uri, pronouns, properties).
func (s *Service) Update( //nolint:cyclop,funlen
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
	featureReferrals *string,
	featureApplications *string,
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
	}

	// Validate profile picture URI (empty string means "remove picture")
	if profilePictureURI == nil || *profilePictureURI != "" {
		urlErr := validateOptionalURL(profilePictureURI)
		if urlErr != nil {
			return nil, urlErr
		}

		// Non-admin users can only use URIs from allowed prefixes
		if userKind != UserKindAdmin {
			prefixErr := validateURIPrefixes(profilePictureURI, s.config.GetAllowedURIPrefixes())
			if prefixErr != nil {
				return nil, prefixErr
			}
		}
	}

	// Validate module visibility values
	for _, v := range []*string{featureRelations, featureLinks, featureQA, featureDiscussions, featureReferrals, featureApplications} {
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
		featureReferrals,
		featureApplications,
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
		SessionID:  nil,
		Payload:    nil,
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return accessErr
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
		SessionID:  nil,
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
func (s *Service) CreateProfileLink( //nolint:cyclop,funlen
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
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
		SessionID:  nil,
		Payload:    map[string]any{"profile_id": profileID},
	})

	return link, nil
}

// UpdateProfileLink updates an existing profile link with authorization check.
func (s *Service) UpdateProfileLink( //nolint:cyclop,funlen
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
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
		SessionID:  nil,
		Payload:    nil,
	})

	return updatedLink, nil
}

// DeleteProfileLink soft-deletes a profile link with authorization check.
func (s *Service) DeleteProfileLink( //nolint:cyclop,funlen
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
	canDelete := userKind == UserKindAdmin

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
		SessionID:  nil,
		Payload:    nil,
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

// annotateLinksCanRemove sets the CanRemove flag on each link based on user permissions.
func (s *Service) annotateLinksCanRemove(
	ctx context.Context,
	links []*ProfileLink,
	profileID string,
	userID string,
	userKind string,
) {
	if userID == "" {
		return
	}

	userInfo, userErr := s.repo.GetUserBriefInfo(ctx, userID)

	for _, link := range links {
		link.CanRemove = s.canUserRemoveLink(
			ctx, link, profileID, userKind, userInfo, userErr,
		)
	}
}

// canUserRemoveLink determines if a user can remove a specific profile link.
func (s *Service) canUserRemoveLink( //nolint:cyclop
	ctx context.Context,
	link *ProfileLink,
	profileID string,
	userKind string,
	userInfo *UserBriefInfo,
	userErr error,
) bool {
	// Site admins can always remove
	if userKind == UserKindAdmin {
		return true
	}

	if userErr != nil || userInfo == nil || userInfo.IndividualProfileID == nil {
		return false
	}

	// Check membership level
	membershipKind, mkErr := s.repo.GetMembershipBetweenProfiles(
		ctx, profileID, *userInfo.IndividualProfileID,
	)
	if mkErr == nil && membershipKind != "" {
		level, ok := GetMembershipKindLevel()[membershipKind]
		if ok && level >= GetMembershipKindLevel()[MembershipKindMaintainer] {
			return true
		}
	}

	// The original adder can remove their own links
	if link.AddedByProfileID != nil &&
		*userInfo.IndividualProfileID == *link.AddedByProfileID {
		return true
	}

	return false
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
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
		fullLink, linkErr := s.repo.GetProfileLink(ctx, localeCode, briefLink.ID)
		if linkErr != nil {
			return nil, fmt.Errorf(
				"%w(linkID: %s): %w",
				ErrFailedToGetRecord,
				briefLink.ID,
				linkErr,
			)
		}

		if fullLink != nil {
			links = append(links, fullLink)
		}
	}

	s.annotateLinksCanRemove(ctx, links, profileID, userID, userKind)

	return links, nil
}

// Profile Page Management

// CreateProfilePage creates a new profile page with authorization check.
func (s *Service) CreateProfilePage( //nolint:cyclop,funlen
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
	}

	// Get the user's individual profile for added_by tracking
	userInfo, userInfoErr := s.repo.GetUserBriefInfo(ctx, userID)

	var addedByProfileID *string
	if userInfoErr == nil && userInfo != nil && userInfo.IndividualProfileID != nil {
		addedByProfileID = userInfo.IndividualProfileID
	}

	// Validate cover picture URI
	coverErr := validateOptionalURL(coverPictureURI)
	if coverErr != nil {
		return nil, coverErr
	}

	// Non-admin users can only use URIs from allowed prefixes
	if userKind != UserKindAdmin {
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
		SessionID:  nil,
		Payload:    map[string]any{"profile_id": profileID, "slug": slug},
	})

	return fullPage, nil
}

// UpdateProfilePage updates an existing profile page with authorization check.
func (s *Service) UpdateProfilePage( //nolint:cyclop,funlen
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
	}

	// Validate cover picture URI
	coverErr := validateOptionalURL(coverPictureURI)
	if coverErr != nil {
		return nil, coverErr
	}

	// Non-admin users can only use URIs from allowed prefixes
	if userKind != UserKindAdmin {
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
		SessionID:  nil,
		Payload:    nil,
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return accessErr
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
		SessionID:  nil,
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return accessErr
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
		SessionID:  nil,
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
func (s *Service) DeleteProfilePage( //nolint:cyclop,funlen
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
	canDelete := userKind == UserKindAdmin

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
		SessionID:  nil,
		Payload:    nil,
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
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
	profileID, err := s.repo.GetProfileIDBySlug(ctx, slug)
	if err != nil {
		return "", fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	return profileID, nil
}

// GetProfileLinkByRemoteID returns a profile link by its remote ID (e.g., YouTube channel ID).
func (s *Service) GetProfileLinkByRemoteID(
	ctx context.Context,
	profileID string,
	kind string,
	remoteID string,
) (*ProfileLink, error) {
	link, err := s.repo.GetProfileLinkByRemoteID(ctx, profileID, kind, remoteID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	return link, nil
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
	inUse, err := s.repo.IsManagedProfileLinkRemoteIDInUse(ctx, kind, remoteID, excludeProfileID)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	return inUse, nil
}

// ClearNonManagedProfileLinkRemoteID nulls out remote_id on non-managed links
// to avoid unique constraint violations when creating a new managed link.
func (s *Service) ClearNonManagedProfileLinkRemoteID(
	ctx context.Context,
	profileID string,
	kind string,
	remoteID string,
) error {
	err := s.repo.ClearNonManagedProfileLinkRemoteID(ctx, profileID, kind, remoteID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	return nil
}

// UpdateProfileLinkOAuthTokens updates the OAuth tokens for an existing profile link.
func (s *Service) UpdateProfileLinkOAuthTokens(
	ctx context.Context,
	linkID string,
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
		ctx, linkID, publicID, uri, authScope, accessToken, accessTokenExpiresAt, refreshToken,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	// Update/create translation for the title (icon, group, description are nil for OAuth links)
	txErr := s.repo.UpsertProfileLinkTx(ctx, linkID, localeCode, title, nil, nil, nil)
	if txErr != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, txErr)
	}

	return nil
}

// CreateOAuthProfileLink creates a new OAuth-connected profile link.
func (s *Service) CreateOAuthProfileLink(
	ctx context.Context,
	linkID string,
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
		ctx, linkID, kind, profileID, order, remoteID, publicID, uri,
		authProvider, authScope, accessToken, accessTokenExpiresAt, refreshToken,
		properties,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	// Create translation for the title (icon, group, description are nil for OAuth links)
	txErr := s.repo.UpsertProfileLinkTx(ctx, linkID, localeCode, title, nil, nil, nil)
	if txErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, txErr)
	}

	return link, nil
}

// GetMaxProfileLinkOrder returns the maximum order value for profile links.
func (s *Service) GetMaxProfileLinkOrder(ctx context.Context, profileID string) (int, error) {
	maxOrder, err := s.repo.GetMaxProfileLinkOrder(ctx, profileID)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	return maxOrder, nil
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

// ListOnlineProfileLinks returns all currently live profile links across all profiles.
func (s *Service) ListOnlineProfileLinks(
	ctx context.Context,
	localeCode string,
) ([]*LiveStreamInfo, error) {
	streams, err := s.repo.ListOnlineProfileLinks(ctx, localeCode)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return streams, nil
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return accessErr
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

// getMembershipRoleLevel returns the hierarchy of membership roles.
// Higher number = higher privilege level.
func getMembershipRoleLevel() map[string]int {
	return map[string]int{
		string(MembershipKindFollower):    GetMembershipKindLevel()[MembershipKindFollower],
		string(MembershipKindSponsor):     GetMembershipKindLevel()[MembershipKindSponsor],
		string(MembershipKindMember):      GetMembershipKindLevel()[MembershipKindMember],
		string(MembershipKindContributor): GetMembershipKindLevel()[MembershipKindContributor],
		string(MembershipKindMaintainer):  GetMembershipKindLevel()[MembershipKindMaintainer],
		string(MembershipKindLead):        GetMembershipKindLevel()[MembershipKindLead],
		string(MembershipKindOwner):       GetMembershipKindLevel()[MembershipKindOwner],
	}
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
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
	isAdmin := userKind == UserKindAdmin
	if !isAdmin {
		filtered := make([]*ProfileMembershipWithMember, 0, len(memberships))

		for _, membership := range memberships {
			kind := MembershipKind(membership.Kind)
			if kind == MembershipKindSponsor || kind == MembershipKindFollower {
				continue
			}

			filtered = append(filtered, membership)
		}

		memberships = filtered
	}

	return memberships, nil
}

// UpdateMembership updates the kind of an existing membership.
func (s *Service) UpdateMembership( //nolint:cyclop,funlen,gocognit
	ctx context.Context,
	userID string,
	userKind string,
	userIndividualProfileID *string,
	profileSlug string,
	membershipID string,
	newKind string,
) error {
	roleLevel := getMembershipRoleLevel()

	// Validate kind
	validKinds := map[string]bool{
		string(MembershipKindOwner): true, string(MembershipKindLead): true, string(MembershipKindMaintainer): true,
		string(MembershipKindContributor): true, string(MembershipKindMember): true,
		string(MembershipKindSponsor): true, string(MembershipKindFollower): true,
	}
	if !validKinds[newKind] {
		return ErrInvalidMembershipKind
	}

	// Check authorization
	isAdmin := userKind == UserKindAdmin

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
	if userKind != UserKindAdmin && userIndividualProfileID != nil {
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
	if !isAdmin && userIndividualProfileID != nil { //nolint:nestif
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
			userLevel = roleLevel[userMembership.Kind]
		} else if *userIndividualProfileID == membership.ProfileID {
			// Implicit owner of their own individual profile
			userLevel = roleLevel[string(MembershipKindOwner)]
		}

		// Check: Cannot assign a role higher than your own
		newRoleLevel := roleLevel[newKind]
		if newRoleLevel > userLevel {
			return ErrCannotAssignHigherRole
		}

		// Check: Cannot modify a member who has a higher or equal role than you
		// (except for demoting yourself, which is already blocked above)
		targetCurrentLevel := roleLevel[membership.Kind]
		if targetCurrentLevel >= userLevel {
			return ErrCannotModifyHigherMember
		}
	}

	// Check if trying to change to 'owner' on individual profile - not allowed
	if newKind == string(MembershipKindOwner) {
		profile, profileErr := s.repo.GetProfileByID(ctx, "en", membership.ProfileID)
		if profileErr != nil {
			return fmt.Errorf(
				"%w(profileID: %s): %w",
				ErrFailedToGetRecord,
				membership.ProfileID,
				profileErr,
			)
		}

		if profile.Kind == ProfileKindIndividual {
			return fmt.Errorf(
				"%w: cannot set 'owner' on individual profiles",
				ErrInvalidMembershipKind,
			)
		}
	}

	// If changing from owner to something else, check we're not removing the last owner
	if membership.Kind == string(MembershipKindOwner) && newKind != string(MembershipKindOwner) {
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

	if membership.MemberProfileID != nil {
		_ = s.repo.InvalidateMembershipKindCache(
			ctx,
			membership.ProfileID,
			*membership.MemberProfileID,
		)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileMembershipUpdated,
		EntityType: "membership",
		EntityID:   membershipID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
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
func (s *Service) DeleteMembership( //nolint:cyclop,funlen
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
	if profile != nil && profile.Kind == ProfileKindIndividual {
		if userIndividualProfileID != nil && *userIndividualProfileID == profileID {
			if membership.MemberProfileID != nil &&
				*membership.MemberProfileID == *userIndividualProfileID {
				return ErrCannotRemoveIndividualSelf
			}
		}
	}

	// If removing an owner, check we're not removing the last owner
	if membership.Kind == string(MembershipKindOwner) {
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

		if membership.MemberProfileID != nil {
			_ = s.repo.InvalidateMembershipKindCache(
				ctx,
				membership.ProfileID,
				*membership.MemberProfileID,
			)
		}

		s.auditService.Record(ctx, events.AuditParams{
			EventType:  events.ProfileMembershipDeleted,
			EntityType: "membership",
			EntityID:   membershipID,
			ActorID:    &userID,
			ActorKind:  events.ActorUser,
			SessionID:  nil,
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

		if membership.MemberProfileID != nil {
			_ = s.repo.InvalidateMembershipKindCache(ctx, membership.ProfileID, *membership.MemberProfileID)
		}

		s.auditService.Record(ctx, events.AuditParams{
			EventType:  events.ProfileMembershipUpdated,
			EntityType: "membership",
			EntityID:   membershipID,
			ActorID:    &userID,
			ActorKind:  events.ActorUser,
			SessionID:  nil,
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
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
) (string, error) {
	roleLevel := getMembershipRoleLevel()

	// Validate kind
	validKinds := map[string]bool{
		string(MembershipKindOwner): true, string(MembershipKindLead): true, string(MembershipKindMaintainer): true,
		string(MembershipKindContributor): true, string(MembershipKindMember): true,
		string(MembershipKindSponsor): true, string(MembershipKindFollower): true,
	}
	if !validKinds[kind] {
		return "", ErrInvalidMembershipKind
	}

	// Check authorization
	isAdmin := userKind == UserKindAdmin

	// SECURITY: Only admins can assign sponsor or follower roles
	if !isAdmin && (kind == "sponsor" || kind == "follower") {
		return "", ErrInvalidMembershipKind
	}

	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return "", fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, profileSlug, err)
	}

	if profileID == "" {
		return "", ErrProfileNotFound
	}

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return "", accessErr
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
			return "", fmt.Errorf("%w: %w", ErrFailedToGetRecord, membershipErr)
		}

		userLevel := 0
		if userMembership != nil {
			userLevel = roleLevel[userMembership.Kind]
		} else if *userIndividualProfileID == profileID {
			// Implicit owner of their own individual profile
			userLevel = roleLevel[string(MembershipKindOwner)]
		}

		// Check: Cannot assign a role higher than your own
		newRoleLevel := roleLevel[kind]
		if newRoleLevel > userLevel {
			return "", ErrCannotAssignHigherRole
		}
	}

	// Check if trying to add 'owner' to individual profile - not allowed
	// Individual profiles have implicit ownership through user.individual_profile_id
	if kind == string(MembershipKindOwner) {
		profile, profileErr := s.repo.GetProfileByID(ctx, "en", profileID)
		if profileErr != nil {
			return "", fmt.Errorf(
				"%w(profileID: %s): %w",
				ErrFailedToGetRecord,
				profileID,
				profileErr,
			)
		}

		if profile.Kind == ProfileKindIndividual {
			return "", fmt.Errorf(
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
		return "", fmt.Errorf("%w: %w", ErrFailedToGetRecord, existingErr)
	}

	if existing != nil {
		// Existing membership found — promote it to the new kind
		err = s.repo.UpdateProfileMembership(ctx, existing.ID, kind)
		if err != nil {
			return "", fmt.Errorf(
				"%w(membershipID: %s): %w",
				ErrFailedToUpdateRecord,
				existing.ID,
				err,
			)
		}

		_ = s.repo.InvalidateMembershipKindCache(ctx, profileID, memberProfileID)

		s.auditService.Record(ctx, events.AuditParams{
			EventType:  events.ProfileMembershipUpdated,
			EntityType: "membership",
			EntityID:   existing.ID,
			ActorID:    &userID,
			ActorKind:  events.ActorUser,
			SessionID:  nil,
			Payload: map[string]any{
				"profile_id":        profileID,
				"member_profile_id": memberProfileID,
				"kind":              kind,
				"last_properties": map[string]any{
					"kind": existing.Kind,
				},
			},
		})

		return existing.ID, nil
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
		return "", fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileMembershipCreated,
		EntityType: "membership",
		EntityID:   string(membershipID),
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"profile_id":        profileID,
			"member_profile_id": memberProfileID,
			"kind":              kind,
		},
	})

	return string(membershipID), nil
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
		SessionID:  nil,
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

	_ = s.repo.InvalidateMembershipKindCache(ctx, profileID, userIndividualProfileID)

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileMembershipDeleted,
		EntityType: "membership",
		EntityID:   existing.ID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
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

// annotateResourcesCanRemove sets the CanRemove flag on each resource based on user permissions.
func (s *Service) annotateResourcesCanRemove(
	ctx context.Context,
	resources []*ProfileResource,
	profileID string,
	userID string,
	userKind string,
) {
	if userID == "" {
		return
	}

	userInfo, userErr := s.repo.GetUserBriefInfo(ctx, userID)

	for _, resource := range resources {
		resource.CanRemove = s.canUserRemoveResource(
			ctx, resource, profileID, userKind, userInfo, userErr,
		)
	}
}

// canUserRemoveResource determines if a user can remove a specific profile resource.
func (s *Service) canUserRemoveResource(
	ctx context.Context,
	resource *ProfileResource,
	profileID string,
	userKind string,
	userInfo *UserBriefInfo,
	userErr error,
) bool {
	// Site admins can always remove
	if userKind == UserKindAdmin {
		return true
	}

	if userErr != nil || userInfo == nil || userInfo.IndividualProfileID == nil {
		return false
	}

	// Check membership level
	membershipKind, mkErr := s.repo.GetMembershipBetweenProfiles(
		ctx, profileID, *userInfo.IndividualProfileID,
	)
	if mkErr == nil && membershipKind != "" {
		level, ok := GetMembershipKindLevel()[membershipKind]
		if ok && level >= GetMembershipKindLevel()[MembershipKindMaintainer] {
			return true
		}
	}

	// The original adder can remove their own resources
	if *userInfo.IndividualProfileID == resource.AddedByProfileID {
		return true
	}

	return false
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
	for _, resource := range resources {
		teams, teamsErr := s.repo.ListResourceTeams(ctx, resource.ID)
		if teamsErr == nil && teams != nil {
			resource.Teams = teams
		} else {
			resource.Teams = []*ProfileTeam{}
		}
	}

	s.annotateResourcesCanRemove(ctx, resources, profileID, userID, userKind)

	return resources, nil
}

// CreateProfileResource creates a new resource linked to a profile.
func (s *Service) CreateProfileResource( //nolint:cyclop,funlen
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
	if userKind != UserKindAdmin {
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
func (s *Service) DeleteProfileResource( //nolint:cyclop
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
	canDelete := userKind == UserKindAdmin

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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return nil, accessErr
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return accessErr
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return accessErr
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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return accessErr
	}

	idGen := func() string { return string(s.idGenerator()) }

	err = s.repo.SetMembershipTeams(ctx, membershipID, teamIDs, idGen)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileMembershipTeamsUpdated,
		EntityType: "membership",
		EntityID:   membershipID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"profile_id": profileID,
			"team_ids":   teamIDs,
		},
	})

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

	accessErr := s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if accessErr != nil {
		return accessErr
	}

	idGen := func() string { return string(s.idGenerator()) }

	err = s.repo.SetResourceTeams(ctx, resourceID, teamIDs, idGen)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	return nil
}

// CreateCandidate creates a new membership candidate. The referrer must be member+ on the profile.
func (s *Service) CreateCandidate( //nolint:cyclop,funlen
	ctx context.Context,
	userID string,
	profileSlug string,
	referredProfileSlug string,
	teamIDs []string,
) (*ProfileMembershipCandidate, error) {
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

	// Only individual profiles can be referred
	referredBrief, err := s.repo.GetProfileIdentifierByID(ctx, referredProfileID)
	if err != nil || referredBrief == nil {
		return nil, fmt.Errorf("%w: referred profile not found", ErrProfileNotFound)
	}

	if referredBrief.Kind != "individual" {
		return nil, ErrCannotReferNonIndividual
	}

	if *userInfo.IndividualProfileID == referredProfileID {
		return nil, ErrCannotReferSelf
	}

	existingMembership, _ := s.repo.GetProfileMembershipByProfileAndMember(
		ctx, profileID, referredProfileID,
	)
	if existingMembership != nil &&
		GetMembershipKindLevel()[MembershipKind(existingMembership.Kind)] >= GetMembershipKindLevel()[MembershipKindMember] {
		return nil, ErrCannotReferExistingMember
	}

	existingCandidate, _ := s.repo.GetProfileMembershipCandidateByProfileAndReferred(
		ctx, profileID, referredProfileID,
	)
	if existingCandidate != nil {
		return nil, ErrCandidateAlreadyExists
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

	candidateID := string(s.idGenerator())

	candidate, err := s.repo.CreateProfileMembershipCandidate(
		ctx, candidateID, profileID, referredProfileID, &referrerMembership.ID,
		string(CandidateSourceReferral), nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	for _, teamID := range teamIDs {
		teamRecordID := string(s.idGenerator())

		insertErr := s.repo.InsertCandidateTeam(ctx, teamRecordID, candidateID, teamID)
		if insertErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, insertErr)
		}
	}

	teams, _ := s.repo.ListCandidateTeams(ctx, candidateID)
	candidate.Teams = teams

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileCandidateCreated,
		EntityType: "candidate",
		EntityID:   candidateID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"profile_id":             profileID,
			"referred_profile":       referredProfileSlug,
			"referrer_membership_id": referrerMembership.ID,
			"source":                 string(CandidateSourceReferral),
		},
	})

	return candidate, nil
}

// ListCandidates lists all active candidates for a profile. Member+ only.
func (s *Service) ListCandidates(
	ctx context.Context,
	localeCode string,
	userID string,
	profileSlug string,
) ([]*ProfileMembershipCandidate, error) {
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

	candidates, err := s.repo.ListProfileMembershipCandidatesByProfileID(
		ctx, localeCode, profileID, viewerMembershipID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	for _, candidate := range candidates {
		teams, teamsErr := s.repo.ListCandidateTeams(ctx, candidate.ID)
		if teamsErr == nil {
			candidate.Teams = teams
		}
	}

	return candidates, nil
}

// VoteCandidate casts or updates a vote on a candidate. Member+ only.
func (s *Service) VoteCandidate( //nolint:cyclop,funlen
	ctx context.Context,
	userID string,
	profileSlug string,
	candidateID string,
	score int16,
	comment *string,
) (*CandidateVote, error) {
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

	candidate, err := s.repo.GetProfileMembershipCandidateByID(ctx, candidateID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCandidateNotFound, err)
	}

	if candidate.ProfileID != profileID {
		return nil, ErrCandidateNotFound
	}

	if candidate.Status != CandidateStatusVoting {
		return nil, ErrCandidateNotVoting
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

	vote, err := s.repo.UpsertCandidateVote(
		ctx, voteID, candidateID, voterMembership.ID, score, comment,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	_ = s.repo.UpdateCandidateVoteCount(ctx, candidateID)

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileCandidateVoted,
		EntityType: "candidate",
		EntityID:   candidateID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"score":   score,
			"vote_id": vote.ID,
		},
	})

	return vote, nil
}

// GetCandidateVotes gets all votes for a candidate. Member+ only.
func (s *Service) GetCandidateVotes(
	ctx context.Context,
	localeCode string,
	userID string,
	profileSlug string,
	candidateID string,
) ([]*CandidateVote, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	err = s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMember)
	if err != nil {
		return nil, err
	}

	candidate, err := s.repo.GetProfileMembershipCandidateByID(ctx, candidateID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCandidateNotFound, err)
	}

	if candidate.ProfileID != profileID {
		return nil, ErrCandidateNotFound
	}

	votes, err := s.repo.ListCandidateVotes(ctx, localeCode, candidateID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return votes, nil
}

// validCandidateTransitions defines allowed status transitions.
var validCandidateTransitions = map[CandidateStatus][]CandidateStatus{ //nolint:gochecknoglobals
	CandidateStatusVoting: {
		CandidateStatusFrozen,
		CandidateStatusReferenceRejected,
		CandidateStatusInvitationPendingResponse,
		CandidateStatusApplicationAccepted,
	},
	CandidateStatusFrozen: {
		CandidateStatusVoting,
		CandidateStatusReferenceRejected,
	},
	// invitation_pending_response, reference_rejected, invitation_accepted,
	// invitation_rejected, application_accepted are terminal or managed by mailbox.
}

func isValidCandidateTransition(from, target CandidateStatus) bool {
	allowed, ok := validCandidateTransitions[from]
	if !ok {
		return false
	}

	for _, status := range allowed {
		if status == target {
			return true
		}
	}

	return false
}

// UpdateCandidateStatus changes the status of a candidate. Requires lead+ access.
func (s *Service) UpdateCandidateStatus(
	ctx context.Context,
	userID string,
	profileSlug string,
	candidateID string,
	newStatus CandidateStatus,
) error {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	err = s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindMaintainer)
	if err != nil {
		return err
	}

	candidate, err := s.repo.GetProfileMembershipCandidateByID(ctx, candidateID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCandidateNotFound, err)
	}

	if candidate.ProfileID != profileID {
		return ErrCandidateNotFound
	}

	if !isValidCandidateTransition(candidate.Status, newStatus) {
		return fmt.Errorf(
			"%w: cannot transition from %s to %s",
			ErrInvalidStatusTransition,
			candidate.Status,
			newStatus,
		)
	}

	err = s.repo.UpdateCandidateStatus(ctx, candidateID, profileID, newStatus)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileCandidateStatusChanged,
		EntityType: "candidate",
		EntityID:   candidateID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"profile_id": profileID,
			"old_status": string(candidate.Status),
			"new_status": string(newStatus),
		},
	})

	return nil
}

// GetCandidateByID returns a candidate by its ID. No access check.
func (s *Service) GetCandidateByID(
	ctx context.Context,
	candidateID string,
) (*ProfileMembershipCandidate, error) {
	candidate, err := s.repo.GetProfileMembershipCandidateByID(ctx, candidateID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCandidateNotFound, err)
	}

	return candidate, nil
}

// GetBySlugInternal returns a profile by slug using the default locale. No access check.
func (s *Service) GetBySlugInternal(ctx context.Context, slug string) (*Profile, error) {
	return s.GetBySlug(ctx, "en", slug)
}

// UpdateCandidateStatusInternal changes the status of a candidate without access checks.
// Used by system callbacks (e.g. mailbox acceptance/rejection).
func (s *Service) UpdateCandidateStatusInternal(
	ctx context.Context,
	candidateID string,
	profileID string,
	newStatus CandidateStatus,
) error {
	err := s.repo.UpdateCandidateStatus(ctx, candidateID, profileID, newStatus)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileCandidateStatusChanged,
		EntityType: "candidate",
		EntityID:   candidateID,
		ActorID:    nil,
		ActorKind:  events.ActorSystem,
		SessionID:  nil,
		Payload: map[string]any{
			"profile_id": profileID,
			"new_status": string(newStatus),
			"internal":   true,
		},
	})

	return nil
}

// EnsureMembershipFromCandidateInternal handles the membership+teams logic when
// a candidate invitation is accepted. No access checks (system-initiated).
//
// Membership: keeps existing member+ memberships, upgrades follower/sponsor to member, or creates new.
// Teams: merges existing membership teams with the candidate's suggested teams (union, no duplicates).
func (s *Service) EnsureMembershipFromCandidateInternal(
	ctx context.Context,
	profileID string,
	memberProfileID string,
	candidateID string,
) (string, error) {
	membershipID, err := s.ensureMinMemberMembership(ctx, profileID, memberProfileID)
	if err != nil {
		return "", err
	}

	teamsErr := s.mergeCandidateTeams(ctx, membershipID, candidateID)
	if teamsErr != nil {
		return membershipID, teamsErr
	}

	return membershipID, nil
}

// ensureMinMemberMembership ensures a membership at member+ level exists between
// profileID and memberProfileID. Returns the membership ID.
func (s *Service) ensureMinMemberMembership(
	ctx context.Context,
	profileID string,
	memberProfileID string,
) (string, error) {
	existing, _ := s.repo.GetProfileMembershipByProfileAndMember(ctx, profileID, memberProfileID)

	if existing != nil {
		levels := GetMembershipKindLevel()
		if levels[MembershipKind(existing.Kind)] < levels[MembershipKindMember] {
			err := s.repo.UpdateProfileMembership(ctx, existing.ID, string(MembershipKindMember))
			if err != nil {
				return "", fmt.Errorf("%w: upgrade membership: %w", ErrFailedToUpdateRecord, err)
			}

			_ = s.repo.InvalidateMembershipKindCache(ctx, profileID, memberProfileID)

			s.recordSystemMembershipAudit(
				ctx,
				events.ProfileMembershipUpdated,
				existing.ID,
				map[string]any{
					"profile_id": profileID, "member_profile_id": memberProfileID,
					"old_kind": existing.Kind, "new_kind": string(MembershipKindMember), "source": "candidate",
				},
			)
		}

		return existing.ID, nil
	}

	newID := string(s.idGenerator())

	err := s.repo.CreateProfileMembership(
		ctx,
		newID,
		profileID,
		&memberProfileID,
		string(MembershipKindMember),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("%w: create membership: %w", ErrFailedToCreateRecord, err)
	}

	s.recordSystemMembershipAudit(ctx, events.ProfileMembershipCreated, newID, map[string]any{
		"profile_id": profileID, "member_profile_id": memberProfileID,
		"kind": string(MembershipKindMember), "source": "candidate",
	})

	return newID, nil
}

// mergeCandidateTeams merges a candidate's suggested teams into the membership's existing teams.
func (s *Service) mergeCandidateTeams(
	ctx context.Context,
	membershipID string,
	candidateID string,
) error {
	candidateTeams, _ := s.repo.ListCandidateTeams(ctx, candidateID)
	if len(candidateTeams) == 0 {
		return nil
	}

	existingTeams, _ := s.repo.ListMembershipTeams(ctx, membershipID)

	existingSet := make(map[string]struct{}, len(existingTeams))
	for _, t := range existingTeams {
		existingSet[t.ID] = struct{}{}
	}

	var teamsAdded []string

	for _, t := range candidateTeams {
		if _, exists := existingSet[t.ID]; !exists {
			teamsAdded = append(teamsAdded, t.ID)
		}
	}

	if len(teamsAdded) == 0 {
		return nil
	}

	mergedIDs := make([]string, 0, len(existingTeams)+len(teamsAdded))
	for _, t := range existingTeams {
		mergedIDs = append(mergedIDs, t.ID)
	}

	mergedIDs = append(mergedIDs, teamsAdded...)

	idGen := func() string { return string(s.idGenerator()) }

	err := s.repo.SetMembershipTeams(ctx, membershipID, mergedIDs, idGen)
	if err != nil {
		return fmt.Errorf("%w: set teams: %w", ErrFailedToUpdateRecord, err)
	}

	s.recordSystemMembershipAudit(
		ctx,
		events.ProfileMembershipUpdated,
		membershipID,
		map[string]any{
			"teams_added": teamsAdded, "source": "candidate", "candidate_id": candidateID,
		},
	)

	return nil
}

// RecordCandidateInvitationSent records an audit event when a candidate invitation is sent.
func (s *Service) RecordCandidateInvitationSent(
	ctx context.Context,
	userID string,
	candidateID string,
	profileID string,
	referredProfileID string,
) {
	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileCandidateInvitationSent,
		EntityType: "candidate",
		EntityID:   candidateID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"profile_id":          profileID,
			"referred_profile_id": referredProfileID,
		},
	})
}

// recordSystemMembershipAudit records a system-initiated membership audit event.
func (s *Service) recordSystemMembershipAudit(
	ctx context.Context,
	eventType events.EventType,
	entityID string,
	payload map[string]any,
) {
	s.auditService.Record(ctx, events.AuditParams{
		EventType:  eventType,
		EntityType: "membership",
		EntityID:   entityID,
		ActorID:    nil,
		ActorKind:  events.ActorSystem,
		SessionID:  nil,
		Payload:    payload,
	})
}

// GetMembershipKindBetween returns the membership kind of memberProfileID on profileID.
// Returns empty string and no error when no membership exists.
func (s *Service) GetMembershipKindBetween(
	ctx context.Context,
	profileID string,
	memberProfileID string,
) (string, error) {
	kind, err := s.repo.GetMembershipBetweenProfiles(ctx, profileID, memberProfileID)
	if err != nil {
		return "", nil //nolint:nilerr // no membership is not an error
	}

	return string(kind), nil
}

// --- Application Form Methods ---

// GetApplicationForm returns the active application form for a profile.
// Public access — anyone can see the form to know what questions will be asked.
func (s *Service) GetApplicationForm(
	ctx context.Context,
	profileSlug string,
) (*ApplicationForm, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	form, featureApplications, err := s.repo.GetApplicationFormByProfileID(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNoApplicationForm, err)
	}

	if featureApplications == string(ModuleVisibilityDisabled) {
		return nil, ErrApplicationsNotEnabled
	}

	return form, nil
}

// UpsertApplicationForm creates or updates the application form for a profile.
// Requires lead+ access. If a preset key is provided, preset fields are used as initial fields
// but the org can customize them afterward.
func (s *Service) UpsertApplicationForm( //nolint:cyclop,gocognit,funlen
	ctx context.Context,
	userID string,
	profileSlug string,
	presetKey *string,
	fields []ApplicationFormFieldInput,
	responsesVisibility string,
) (*ApplicationForm, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	err = s.ensureUserCanProfileAccess(ctx, profileID, userID, MembershipKindLead)
	if err != nil {
		return nil, err
	}

	if responsesVisibility != "members" && responsesVisibility != "leads" {
		return nil, ErrInvalidResponsesVisibility
	}

	// Validate field types
	validFieldTypes := map[string]bool{"short_text": true, "long_text": true, "url": true}
	for _, field := range fields {
		if !validFieldTypes[field.FieldType] {
			return nil, fmt.Errorf("%w: %s", ErrInvalidFieldType, field.FieldType)
		}
	}

	// If preset is given and no explicit fields, use preset fields (always English for DB storage)
	if presetKey != nil && len(fields) == 0 {
		preset, ok := ApplicationPresets[*presetKey]
		if ok {
			fields = make([]ApplicationFormFieldInput, 0, len(preset.Fields))
			for i, presetField := range preset.Fields {
				var placeholder *string
				if presetField.Placeholder != "" {
					placeholder = &presetField.Placeholder
				}

				fields = append(fields, ApplicationFormFieldInput{
					Label:       presetField.Label,
					FieldType:   presetField.FieldType,
					IsRequired:  presetField.IsRequired,
					SortOrder:   i,
					Placeholder: placeholder,
				})
			}
		}
	}

	// Check if an active form already exists
	existingForm, _ := s.repo.GetActiveApplicationForm(ctx, profileID)

	if existingForm != nil { //nolint:nestif
		// Update existing form
		err = s.repo.UpdateApplicationForm(ctx, existingForm.ID, presetKey, responsesVisibility)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, err)
		}

		// Replace fields: delete old, create new
		err = s.repo.DeleteApplicationFormFields(ctx, existingForm.ID)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToDeleteRecord, err)
		}

		for _, field := range fields {
			fieldID := string(s.idGenerator())

			_, fieldErr := s.repo.CreateApplicationFormField(
				ctx, fieldID, existingForm.ID,
				field.Label, field.FieldType, field.IsRequired,
				field.SortOrder, field.Placeholder,
			)
			if fieldErr != nil {
				return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, fieldErr)
			}
		}
	} else {
		// Create new form
		formID := string(s.idGenerator())

		_, err = s.repo.CreateApplicationForm(ctx, formID, profileID, presetKey, responsesVisibility)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
		}

		for _, field := range fields {
			fieldID := string(s.idGenerator())

			_, fieldErr := s.repo.CreateApplicationFormField(
				ctx, fieldID, formID,
				field.Label, field.FieldType, field.IsRequired,
				field.SortOrder, field.Placeholder,
			)
			if fieldErr != nil {
				return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, fieldErr)
			}
		}
	}

	// Re-fetch final form with all fields populated
	form, err := s.repo.GetActiveApplicationForm(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	return form, nil
}

// CreateApplication creates a self-nominated candidate (application).
// Requires: authenticated user, feature_applications enabled.
// If no application form is configured, accepts the application with just the message.
func (s *Service) CreateApplication( //nolint:cyclop,funlen
	ctx context.Context,
	userID string,
	profileSlug string,
	applicantMessage *string,
	formResponses map[string]string,
) (*ProfileMembershipCandidate, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	// Try to get application form (also returns feature flag)
	form, featureApplications, formErr := s.repo.GetApplicationFormByProfileID(ctx, profileID)
	if formErr != nil {
		// No form configured — check feature flag from profile directly
		profile, profileErr := s.repo.GetProfileByID(ctx, "en", profileID)
		if profileErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, profileErr)
		}

		featureApplications = profile.FeatureApplications
		form = nil
	}

	if featureApplications == string(ModuleVisibilityDisabled) {
		return nil, ErrApplicationsNotEnabled
	}

	// Get the applicant's individual profile
	userInfo, err := s.repo.GetUserBriefInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if userInfo.IndividualProfileID == nil {
		return nil, fmt.Errorf("%w: %w", ErrInsufficientAccess, ErrNoIndividualProfile)
	}

	applicantProfileID := *userInfo.IndividualProfileID

	// Cannot apply if already a member (member+)
	existingMembership, _ := s.repo.GetProfileMembershipByProfileAndMember(
		ctx, profileID, applicantProfileID,
	)
	if existingMembership != nil &&
		GetMembershipKindLevel()[MembershipKind(existingMembership.Kind)] >= GetMembershipKindLevel()[MembershipKindMember] {
		return nil, ErrCannotReferExistingMember
	}

	// Check if user already has a pending candidate record
	existingCandidate, _ := s.repo.GetProfileMembershipCandidateByProfileAndReferred(
		ctx, profileID, applicantProfileID,
	)

	if existingCandidate != nil && existingCandidate.Source == string(CandidateSourceApplication) {
		// Already applied via application — reject
		return nil, ErrAlreadyApplied
	}

	// Validate required fields (only when a form is configured)
	if form != nil {
		for _, field := range form.Fields {
			if field.IsRequired {
				value, exists := formResponses[field.ID]
				if !exists || value == "" {
					return nil, fmt.Errorf("%w: %s", ErrMissingRequiredField, field.Label)
				}
			}
		}
	}

	var candidate *ProfileMembershipCandidate

	var candidateID string

	if existingCandidate != nil && existingCandidate.Source == string(CandidateSourceReferral) {
		// Referred candidate is now also applying — update with their message
		candidateID = existingCandidate.ID

		msgErr := s.repo.UpdateCandidateApplicantMessage(ctx, candidateID, applicantMessage)
		if msgErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToUpdateRecord, msgErr)
		}

		candidate = existingCandidate
		candidate.ApplicantMessage = applicantMessage
	} else {
		// New application — create candidate record
		candidateID = string(s.idGenerator())

		var createErr error

		candidate, createErr = s.repo.CreateProfileMembershipCandidate(
			ctx, candidateID, profileID, applicantProfileID, nil,
			string(CandidateSourceApplication), applicantMessage,
		)
		if createErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, createErr)
		}
	}

	// Save form responses (only when a form is configured)
	if form != nil {
		for fieldID, value := range formResponses {
			responseID := string(s.idGenerator())

			responseErr := s.repo.CreateCandidateResponse(
				ctx,
				responseID,
				candidateID,
				fieldID,
				value,
			)
			if responseErr != nil {
				return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, responseErr)
			}
		}
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.ProfileCandidateCreated,
		EntityType: "candidate",
		EntityID:   candidateID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
		SessionID:  nil,
		Payload: map[string]any{
			"profile_id": profileID,
			"source":     string(CandidateSourceApplication),
		},
	})

	return candidate, nil
}

// GetMyApplication returns the current user's candidate record for a profile (if any).
// Used for applicant status tracking on the frontend.
func (s *Service) GetMyApplication(
	ctx context.Context,
	userID string,
	profileSlug string,
) (*ProfileMembershipCandidate, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	userInfo, err := s.repo.GetUserBriefInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	if userInfo.IndividualProfileID == nil {
		return nil, fmt.Errorf("%w: %w", ErrInsufficientAccess, ErrNoIndividualProfile)
	}

	candidate, err := s.repo.GetProfileMembershipCandidateByProfileAndReferred(
		ctx, profileID, *userInfo.IndividualProfileID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCandidateNotFound, err)
	}

	return candidate, nil
}

// GetCandidateFormResponses returns form responses for a candidate.
// Access is gated by the form's responses_visibility setting:
// "members" → member+ can see, "leads" → lead+ only.
func (s *Service) GetCandidateFormResponses(
	ctx context.Context,
	userID string,
	profileSlug string,
	candidateID string,
) ([]*CandidateFormResponse, error) {
	profileID, err := s.repo.GetProfileIDBySlug(ctx, profileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrProfileNotFound, err)
	}

	// Determine minimum access level from form's visibility setting
	requiredLevel := MembershipKindMember

	form, _ := s.repo.GetActiveApplicationForm(ctx, profileID)
	if form != nil && form.ResponsesVisibility == "leads" {
		requiredLevel = MembershipKindLead
	}

	err = s.ensureUserCanProfileAccess(ctx, profileID, userID, requiredLevel)
	if err != nil {
		return nil, err
	}

	responses, err := s.repo.ListCandidateResponses(ctx, candidateID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return responses, nil
}
