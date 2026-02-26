package profiles

import (
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
)

type RecordID string

type RecordIDGenerator func() RecordID

func DefaultIDGenerator() RecordID {
	return RecordID(lib.IDsGenerateUnique())
}

// LinkVisibility defines who can see a profile link.
// Each level corresponds to a minimum membership level required to view.
type LinkVisibility string

const (
	LinkVisibilityPublic       LinkVisibility = "public"       // Everyone can see
	LinkVisibilityFollowers    LinkVisibility = "followers"    // Followers and above
	LinkVisibilitySponsors     LinkVisibility = "sponsors"     // Sponsors and above
	LinkVisibilityMembers      LinkVisibility = "members"      // Members and above
	LinkVisibilityContributors LinkVisibility = "contributors" // Contributors and above
	LinkVisibilityMaintainers  LinkVisibility = "maintainers"  // Maintainers and above
	LinkVisibilityLeads        LinkVisibility = "leads"        // Leads and above
	LinkVisibilityOwners       LinkVisibility = "owners"       // Owners only
)

// ModuleVisibility defines the visibility state of a profile module (e.g. Q&A, contributions).
type ModuleVisibility string

const (
	ModuleVisibilityPublic   ModuleVisibility = "public"   // Enabled and shown in navigation
	ModuleVisibilityHidden   ModuleVisibility = "hidden"   // Enabled but not shown in navigation
	ModuleVisibilityDisabled ModuleVisibility = "disabled" // Completely disabled, returns 404
)

// PageVisibility defines who can see a profile page.
type PageVisibility string

const (
	PageVisibilityPublic   PageVisibility = "public"   // Listed in sidebar, visible to all
	PageVisibilityUnlisted PageVisibility = "unlisted" // Accessible via direct link, not listed
	PageVisibilityPrivate  PageVisibility = "private"  // Only contributors+ and admins
)

// MembershipKind represents the type of membership a profile has with another.
type MembershipKind string

const (
	MembershipKindFollower    MembershipKind = "follower"
	MembershipKindSponsor     MembershipKind = "sponsor"
	MembershipKindMember      MembershipKind = "member"
	MembershipKindContributor MembershipKind = "contributor"
	MembershipKindMaintainer  MembershipKind = "maintainer"
	MembershipKindLead        MembershipKind = "lead"
	MembershipKindOwner       MembershipKind = "owner"
)

// MembershipKindLevel returns the privilege level of a membership kind.
// Higher values mean more privileges.
var MembershipKindLevel = map[MembershipKind]int{
	MembershipKindFollower:    1,
	MembershipKindSponsor:     2,
	MembershipKindMember:      3,
	MembershipKindContributor: 4,
	MembershipKindMaintainer:  5,
	MembershipKindLead:        6,
	MembershipKindOwner:       7,
}

// MinMembershipForVisibility maps visibility levels to minimum membership required.
var MinMembershipForVisibility = map[LinkVisibility]MembershipKind{
	LinkVisibilityPublic:       "",                        // no membership required
	LinkVisibilityFollowers:    MembershipKindFollower,    // followers+
	LinkVisibilitySponsors:     MembershipKindSponsor,     // sponsors+
	LinkVisibilityMembers:      MembershipKindMember,      // members+
	LinkVisibilityContributors: MembershipKindContributor, // contributors+
	LinkVisibilityMaintainers:  MembershipKindMaintainer,  // maintainers+
	LinkVisibilityLeads:        MembershipKindLead,        // leads+
	LinkVisibilityOwners:       MembershipKindOwner,       // owners only
}

type Profile struct {
	CreatedAt                       time.Time  `json:"created_at"`
	Properties                      any        `json:"properties"`
	ProfilePictureURI               *string    `json:"profile_picture_uri"`
	Pronouns                        *string    `json:"pronouns"`
	UpdatedAt                       *time.Time `json:"updated_at"`
	DeletedAt                       *time.Time `json:"deleted_at"`
	ID                              string     `json:"id"`
	Slug                            string     `json:"slug"`
	Kind                            string     `json:"kind"`
	LocaleCode                      string     `json:"locale_code"`
	Title                           string     `json:"title"`
	Description                     string     `json:"description"`
	DefaultLocale                   string     `json:"default_locale"`
	Points                          uint64     `json:"points"`
	HasTranslation                  bool       `json:"has_translation"`
	FeatureRelations                string     `json:"feature_relations"`   // Visibility of Members/Contributions module
	FeatureLinks                    string     `json:"feature_links"`       // Visibility of Links module
	FeatureQA                       string     `json:"feature_qa"`          // Visibility of Q&A module
	FeatureDiscussions              string     `json:"feature_discussions"` // Visibility of Discussions module
	OptionStoryDiscussionsByDefault bool       `json:"option_story_discussions_by_default"`
}

type ProfileCustomDomain struct {
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DefaultLocale *string    `json:"default_locale"`
	ID            string     `json:"id"`
	ProfileID     string     `json:"profile_id"`
	Domain        string     `json:"domain"`
}

type ProfileWithChildren struct {
	*Profile

	Pages []*ProfilePageBrief `json:"pages"`
	Links []*ProfileLinkBrief `json:"links"`
}

type ProfilePage struct {
	CoverPictureURI  *string        `json:"cover_picture_uri"`
	PublishedAt      *time.Time     `json:"published_at"`
	AddedByProfileID *string        `json:"added_by_profile_id"`
	AddedByProfile   *ProfileBrief  `json:"added_by_profile,omitempty"`
	ID               string         `json:"id"`
	Slug             string         `json:"slug"`
	LocaleCode       string         `json:"locale_code"`
	Title            string         `json:"title"`
	Summary          string         `json:"summary"`
	Content          string         `json:"content"`
	Visibility       PageVisibility `json:"visibility"`
	SortOrder        int32          `json:"sort_order"`
	CanRemove        bool           `json:"can_remove"`
}

type ProfilePageBrief struct {
	ID              string         `json:"id"`
	Slug            string         `json:"slug"`
	CoverPictureURI *string        `json:"cover_picture_uri"`
	Title           string         `json:"title"`
	Summary         string         `json:"summary"`
	Visibility      PageVisibility `json:"visibility"`
}

type ProfileLink struct {
	ID               string         `json:"id"`
	Kind             string         `json:"kind"`
	ProfileID        string         `json:"profile_id"`
	Order            int            `json:"order"`
	IsManaged        bool           `json:"is_managed"`
	IsVerified       bool           `json:"is_verified"`
	IsFeatured       bool           `json:"is_featured"`
	Visibility       LinkVisibility `json:"visibility"`
	RemoteID         *string        `json:"remote_id"`
	PublicID         *string        `json:"public_id"`
	URI              *string        `json:"uri"`
	Properties       map[string]any `json:"properties,omitempty"`
	Title            string         `json:"title"`       // From profile_link_tx
	Icon             *string        `json:"icon"`        // From profile_link_tx - custom emoticon or initials
	Group            *string        `json:"group"`       // From profile_link_tx
	Description      *string        `json:"description"` // From profile_link_tx
	AddedByProfileID *string        `json:"added_by_profile_id"`
	AddedByProfile   *ProfileBrief  `json:"added_by_profile,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        *time.Time     `json:"updated_at"`
	DeletedAt        *time.Time     `json:"deleted_at"`
	CanRemove        bool           `json:"can_remove"`
}

type ProfileLinkBrief struct {
	ID          string         `json:"id"`
	Kind        string         `json:"kind"`
	Order       int            `json:"order"`
	PublicID    string         `json:"public_id"`
	URI         string         `json:"uri"`
	Title       string         `json:"title"`       // From profile_link_tx
	Icon        string         `json:"icon"`        // From profile_link_tx - custom emoticon or initials
	Group       string         `json:"group"`       // From profile_link_tx
	Description string         `json:"description"` // From profile_link_tx
	IsManaged   bool           `json:"is_managed"`
	IsVerified  bool           `json:"is_verified"`
	IsFeatured  bool           `json:"is_featured"`
	Visibility  LinkVisibility `json:"visibility"`
}

// ProfileLinkState contains state for profile link OAuth flows.
// This extends the basic OAuth state with profile-specific data.
type ProfileLinkState struct {
	State          string    `json:"state"`
	ProfileSlug    string    `json:"profile_slug"`
	ProfileKind    string    `json:"profile_kind"`
	Locale         string    `json:"locale"`
	RedirectOrigin string    `json:"redirect_origin"`
	ExpiresAt      time.Time `json:"expires_at"`
}

// GitHubAccount represents a GitHub account (user or organization) for selection.
type GitHubAccount struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	HTMLURL     string `json:"html_url"`
	Type        string `json:"type"` // "User" or "Organization"
	Description string `json:"description,omitempty"`
}

// PendingOAuthConnection stores temporary OAuth data for account selection.
// Used by providers that support organization/page account selection (GitHub, LinkedIn).
type PendingOAuthConnection struct {
	Provider    string    `json:"provider"` // "github" or "linkedin"
	AccessToken string    `json:"access_token"`
	Scope       string    `json:"scope"`
	ProfileSlug string    `json:"profile_slug"`
	ProfileKind string    `json:"profile_kind"`
	Locale      string    `json:"locale"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// LinkedInAccount represents a LinkedIn account (personal or organization page) for selection.
type LinkedInAccount struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	VanityName  string `json:"vanity_name,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
	URI         string `json:"uri"`
	Type        string `json:"type"` // "Personal" or "Organization"
	Description string `json:"description,omitempty"`
}

type ProfileMembership struct {
	Properties      any        `json:"properties"`
	Profile         *Profile   `json:"profile"`
	MemberProfile   *Profile   `json:"member_profile"`
	StartedAt       *time.Time `json:"started_at"`
	FinishedAt      *time.Time `json:"finished_at"`
	ID              string     `json:"id"`
	ProfileID       string     `json:"profile_id"`
	MemberProfileID *string    `json:"member_profile_id"`
	Kind            string     `json:"kind"`
}

// ProfileBrief is a lightweight profile representation for lists and references.
type ProfileBrief struct {
	ID                string  `json:"id"`
	Slug              string  `json:"slug"`
	Kind              string  `json:"kind"`
	ProfilePictureURI *string `json:"profile_picture_uri"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
}

// ProfileTeam represents a team within a profile for organizing members.
type ProfileTeam struct {
	ID            string  `json:"id"`
	ProfileID     string  `json:"profile_id"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	MemberCount   int64   `json:"member_count"`
	ResourceCount int64   `json:"resource_count"`
}

// ProfileMembershipWithMember includes membership data with member profile details.
type ProfileMembershipWithMember struct {
	ID              string         `json:"id"`
	ProfileID       string         `json:"profile_id"`
	MemberProfileID *string        `json:"member_profile_id"`
	Kind            string         `json:"kind"`
	Properties      any            `json:"properties"`
	StartedAt       *time.Time     `json:"started_at"`
	FinishedAt      *time.Time     `json:"finished_at"`
	MemberProfile   *ProfileBrief  `json:"member_profile"`
	Teams           []*ProfileTeam `json:"teams"`
}

// UserSearchResult represents a user search result for adding memberships.
type UserSearchResult struct {
	UserID              string        `json:"user_id"`
	Email               string        `json:"email"`
	Name                *string       `json:"name"`
	IndividualProfileID *string       `json:"individual_profile_id"`
	Profile             *ProfileBrief `json:"profile"`
}

// ReferralStatus represents the current status of a membership referral.
type ReferralStatus string

const (
	ReferralStatusVoting                    ReferralStatus = "voting"
	ReferralStatusFrozen                    ReferralStatus = "frozen"
	ReferralStatusReferenceRejected         ReferralStatus = "reference_rejected"
	ReferralStatusInvitationPendingResponse ReferralStatus = "invitation_pending_response"
	ReferralStatusInvitationAccepted        ReferralStatus = "invitation_accepted"
	ReferralStatusInvitationRejected        ReferralStatus = "invitation_rejected"
)

// ProfileMembershipReferral represents a referral for a potential new member.
type ProfileMembershipReferral struct {
	ReferrerProfile      *ProfileBrief  `json:"referrer_profile"`
	ReferredProfile      *ProfileBrief  `json:"referred_profile"`
	Teams                []*ProfileTeam `json:"teams"`
	ViewerVoteComment    *string        `json:"viewer_vote_comment"`
	ViewerVoteScore      *int16         `json:"viewer_vote_score"`
	UpdatedAt            *time.Time     `json:"updated_at"`
	CreatedAt            time.Time      `json:"created_at"`
	ID                   string         `json:"id"`
	ProfileID            string         `json:"profile_id"`
	ReferredProfileID    string         `json:"referred_profile_id"`
	ReferrerMembershipID string         `json:"referrer_membership_id"`
	Status               ReferralStatus `json:"status"`
	AverageScore         float64        `json:"average_score"`
	TotalVotes           int64          `json:"total_votes"`
	VoteCount            int            `json:"vote_count"`
}

// ReferralVote represents a single vote on a referral.
type ReferralVote struct {
	VoterProfile                *ProfileBrief `json:"voter_profile"`
	Comment                     *string       `json:"comment"`
	UpdatedAt                   *time.Time    `json:"updated_at"`
	CreatedAt                   time.Time     `json:"created_at"`
	ID                          string        `json:"id"`
	ProfileMembershipReferralID string        `json:"profile_membership_referral_id"`
	VoterMembershipID           string        `json:"voter_membership_id"`
	Score                       int16         `json:"score"`
}

type ExternalPost struct {
	CreatedAt *time.Time `json:"created_at"` //nolint:tagliatelle
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Permalink string     `json:"permalink"`
}

type ProfileOwnership struct {
	ProfileID   string `json:"profile_id"`
	ProfileSlug string `json:"profile_slug"`
	ProfileKind string `json:"profile_kind"`
	UserKind    string `json:"user_kind"`
	CanEdit     bool   `json:"can_edit"`
}

// UserBriefInfo holds the minimum user information needed for access control checks.
type UserBriefInfo struct {
	Kind                string  `json:"kind"`
	IndividualProfileID *string `json:"individual_profile_id"`
}

type ProfilePermission struct {
	ProfileID      string `json:"profile_id"`
	ProfileSlug    string `json:"profile_slug"`
	ProfileKind    string `json:"profile_kind"`
	MembershipKind string `json:"membership_kind"`
	UserKind       string `json:"user_kind"`
}

type ProfileTx struct {
	ProfileID   string `json:"profile_id"`
	LocaleCode  string `json:"locale_code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Properties  any    `json:"properties"`
}

type ProfileLinkTx struct {
	ProfileLinkID string  `json:"profile_link_id"`
	LocaleCode    string  `json:"locale_code"`
	Title         string  `json:"title"`
	Icon          *string `json:"icon"`
	Group         *string `json:"group"`
	Description   *string `json:"description"`
}

// SpotlightItem represents an item in the spotlight section.
type SpotlightItem struct {
	Icon  string `json:"icon"`
	To    string `json:"to"`
	Title string `json:"title"`
}

// SearchResult represents a unified search result across profiles, stories, and pages.
type SearchResult struct {
	Type         string  `json:"type"`
	ID           string  `json:"id"`
	Slug         string  `json:"slug"`
	Title        string  `json:"title"`
	Summary      *string `json:"summary"`
	ImageURI     *string `json:"image_uri"`
	ProfileSlug  *string `json:"profile_slug"`
	ProfileTitle *string `json:"profile_title"`
	Kind         *string `json:"kind"`
	Rank         float32 `json:"rank"`
}

// ProfileResource represents an external resource linked to a profile (e.g. GitHub repo).
type ProfileResource struct {
	ID               string         `json:"id"`
	ProfileID        string         `json:"profile_id"`
	Kind             string         `json:"kind"`
	IsManaged        bool           `json:"is_managed"`
	RemoteID         *string        `json:"remote_id"`
	PublicID         *string        `json:"public_id"`
	URL              *string        `json:"url"`
	Title            string         `json:"title"`
	Description      *string        `json:"description"`
	Properties       any            `json:"properties"`
	AddedByProfileID string         `json:"added_by_profile_id"`
	AddedByProfile   *ProfileBrief  `json:"added_by_profile,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        *time.Time     `json:"updated_at"`
	DeletedAt        *time.Time     `json:"deleted_at"`
	CanRemove        bool           `json:"can_remove"`
	Teams            []*ProfileTeam `json:"teams"`
}

// ManagedGitHubLink holds the access token data for a managed GitHub profile link.
type ManagedGitHubLink struct {
	ID                   string  `json:"id"`
	ProfileID            string  `json:"profile_id"`
	AuthAccessToken      string  `json:"-"` // Never expose tokens via JSON
	AuthAccessTokenScope *string `json:"-"` // OAuth scope granted for this link
}
