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
	LinkVisibilityContributors LinkVisibility = "contributors" // Contributors and above
	LinkVisibilityMaintainers  LinkVisibility = "maintainers"  // Maintainers and above
	LinkVisibilityLeads        LinkVisibility = "leads"        // Leads and above
	LinkVisibilityOwners       LinkVisibility = "owners"       // Owners only
)

// MembershipKind represents the type of membership a profile has with another.
type MembershipKind string

const (
	MembershipKindFollower    MembershipKind = "follower"
	MembershipKindSponsor     MembershipKind = "sponsor"
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
	MembershipKindContributor: 3,
	MembershipKindMaintainer:  4,
	MembershipKindLead:        5,
	MembershipKindOwner:       6,
}

// MinMembershipForVisibility maps visibility levels to minimum membership required.
var MinMembershipForVisibility = map[LinkVisibility]MembershipKind{
	LinkVisibilityPublic:       "",                        // no membership required
	LinkVisibilityFollowers:    MembershipKindFollower,    // followers+
	LinkVisibilitySponsors:     MembershipKindSponsor,     // sponsors+
	LinkVisibilityContributors: MembershipKindContributor, // contributors+
	LinkVisibilityMaintainers:  MembershipKindMaintainer,  // maintainers+
	LinkVisibilityLeads:        MembershipKindLead,        // leads+
	LinkVisibilityOwners:       MembershipKindOwner,       // owners only
}

type Profile struct {
	CreatedAt         time.Time  `json:"created_at"`
	Properties        any        `json:"properties"`
	CustomDomain      *string    `json:"custom_domain"`
	ProfilePictureURI *string    `json:"profile_picture_uri"`
	Pronouns          *string    `json:"pronouns"`
	UpdatedAt         *time.Time `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
	ID                string     `json:"id"`
	Slug              string     `json:"slug"`
	Kind              string     `json:"kind"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	Points            uint64     `json:"points"`
	HasTranslation    bool       `json:"has_translation"`
}

type ProfileWithChildren struct {
	*Profile

	Pages []*ProfilePageBrief `json:"pages"`
	Links []*ProfileLinkBrief `json:"links"`
}

type ProfilePage struct {
	CoverPictureURI *string    `json:"cover_picture_uri"`
	PublishedAt     *time.Time `json:"published_at"`
	ID              string     `json:"id"`
	Slug            string     `json:"slug"`
	Title           string     `json:"title"`
	Summary         string     `json:"summary"`
	Content         string     `json:"content"`
}

type ProfilePageBrief struct {
	ID              string  `json:"id"`
	Slug            string  `json:"slug"`
	CoverPictureURI *string `json:"cover_picture_uri"`
	Title           string  `json:"title"`
	Summary         string  `json:"summary"`
}

type ProfileLink struct {
	ID          string         `json:"id"`
	Kind        string         `json:"kind"`
	ProfileID   string         `json:"profile_id"`
	Order       int            `json:"order"`
	IsManaged   bool           `json:"is_managed"`
	IsVerified  bool           `json:"is_verified"`
	IsFeatured  bool           `json:"is_featured"`
	Visibility  LinkVisibility `json:"visibility"`
	RemoteID    *string        `json:"remote_id"`
	PublicID    *string        `json:"public_id"`
	URI         *string        `json:"uri"`
	Title       string         `json:"title"`       // From profile_link_tx
	Icon        *string        `json:"icon"`        // From profile_link_tx - custom emoticon or initials
	Group       *string        `json:"group"`       // From profile_link_tx
	Description *string        `json:"description"` // From profile_link_tx
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   *time.Time     `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at"`
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

// PendingGitHubConnection stores temporary OAuth data for account selection.
type PendingGitHubConnection struct {
	AccessToken string    `json:"access_token"`
	Scope       string    `json:"scope"`
	ProfileSlug string    `json:"profile_slug"`
	ProfileKind string    `json:"profile_kind"`
	Locale      string    `json:"locale"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type ProfileMembership struct {
	Properties    any        `json:"properties"`
	Profile       *Profile   `json:"profile"`
	MemberProfile *Profile   `json:"member_profile"`
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	ID            string     `json:"id"`
	Kind          string     `json:"kind"`
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
