// API Response Types

export type Result<T> = {
  data: T | null;
  error?: string | null;
};

// Profile Types
export interface Profile {
  id: string;
  slug: string;
  title: string;
  description?: string | null;
  pronouns?: string | null;
  kind: "individual" | "organization" | "product";
  profile_picture_uri?: string | null;
  links?: ProfileLink[];
  pages?: ProfilePage[];
  points: number;
  has_translation?: boolean;
  hide_relations?: boolean; // Hides Members/Contributions from sidebar
  hide_links?: boolean; // Hides Links from sidebar
  hide_qa?: boolean; // Hides Q&A section from sidebar
  created_at: string;
  updated_at: string;
}

export type ProfileLinkKind =
  | "website"
  | "github"
  | "x"
  | "linkedin"
  | "instagram"
  | "youtube"
  | "bsky"
  | "discord"
  | "telegram";

export type LinkVisibility =
  | "public"
  | "followers"
  | "sponsors"
  | "contributors"
  | "maintainers"
  | "leads"
  | "owners";

export interface ProfileLink {
  id: string;
  kind: ProfileLinkKind;
  profile_id: string;
  order: number;
  is_managed: boolean;
  is_verified: boolean;
  is_featured: boolean;
  visibility: LinkVisibility;
  remote_id?: string | null;
  public_id?: string | null;
  uri?: string | null;
  title: string;
  icon?: string | null; // Custom emoticon or initials - if not specified, kind-based icon is used
  group?: string | null;
  description?: string | null;
  created_at: string;
  updated_at?: string | null;
}

export type ProfilePage = {
  id: string;
  slug: string;
  locale_code?: string;
  title: string;
  summary: string;
  content: string;
  sort_order: number;
  cover_picture_uri?: string | null;
  published_at?: string | null;
};

export interface ProfileContribution {
  id: string;
  profile_id: string;
  target_profile_id: string;
  target_profile: Profile;
  role?: string | null;
  start_date?: string | null;
  end_date?: string | null;
}

export interface ProfileMember {
  id: string;
  profile_id: string;
  member_profile_id: string;
  member_profile: Profile;
  role?: string | null;
  start_date?: string | null;
  end_date?: string | null;
}

export type ProfileMembershipKind =
  | "follower"
  | "sponsor"
  | "contributor"
  | "maintainer"
  | "lead"
  | "owner";

export interface ProfileMembership {
  id: string;
  kind: ProfileMembershipKind;
  properties: {
    github?: {
      commits: number;
      prs: { resolved: number; total: number };
      issues: { resolved: number; total: number };
      stars: number;
    };
    videos?: number;
  };
  profile: Profile;
  member_profile: Profile;
}

// Type alias for MembershipKind (same as ProfileMembershipKind)
export type MembershipKind = ProfileMembershipKind;

// Brief profile for membership display
export interface ProfileBrief {
  id: string;
  slug: string;
  kind: string;
  profile_picture_uri?: string | null;
  title: string;
  description?: string | null;
}

// Membership with member profile details for settings
export interface ProfileMembershipWithMember {
  id: string;
  profile_id: string;
  member_profile_id?: string | null;
  kind: MembershipKind;
  properties?: Record<string, unknown> | null;
  started_at?: string | null;
  finished_at?: string | null;
  created_at: string;
  updated_at?: string | null;
  member_profile?: ProfileBrief | null;
}

// User search result for adding memberships
export interface UserSearchResult {
  user_id: string;
  email: string;
  name?: string | null;
  individual_profile_id?: string | null;
  profile?: ProfileBrief | null;
}

export interface ProfilePermissions {
  can_view: boolean;
  can_edit: boolean;
  can_delete: boolean;
  can_manage_members: boolean;
}

export interface ProfileTranslation {
  locale: string;
  title: string;
  description?: string | null;
}

// Story Types
export type StoryKind =
  | "status"
  | "announcement"
  | "article"
  | "news"
  | "content"
  | "presentation";

export interface Story {
  id: string;
  kind: StoryKind;
  slug: string | null;
  story_picture_uri: string | null;
  title: string | null;
  summary: string | null;
  content: string;
  author_profile_id: string | null;
  author_profile: Profile | null;
  published_at: string | null;
  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
}

export interface StoryEx extends Omit<Story, "author_profile_id" | "author_profile"> {
  author_profile: {
    id: string;
    slug: string;
    title: string;
    profile_picture_uri: string | null;
  } | null;
  publications: Profile[];
}

// User Types
export type UserKind = "admin" | "editor" | "regular" | "disabled";

export interface User {
  id: string;
  github_handle: string | null;
  email: string | null;
  name: string | null;
  profile_picture_uri: string | null;
  kind: UserKind;
  created_at: string;
}

// Session Types
export interface Session {
  id: string;
  user_id: string;
  profile?: Profile | null;
  email?: string | null;
  name?: string | null;
  avatar_url?: string | null;
}

export interface SessionPreferences {
  theme?: "light" | "dark" | "system";
  locale?: string;
  timezone?: string;
}

// Custom Domain Types
export interface CustomDomain {
  host: string;
  slug: string;
  profile?: Profile | null;
}

// Spotlight Types
export interface Spotlight {
  featured_profiles: Profile[];
  recent_stories: Story[];
}

// Create/Update Types
export interface CreateProfileInput {
  slug: string;
  title: string;
  description?: string;
  kind: "individual" | "organization" | "product";
}

export interface UpdateProfileInput {
  title?: string;
  description?: string;
  pronouns?: string;
}

export interface CreateProfileLinkInput {
  kind: string;
  uri: string;
  title?: string;
  visibility?: LinkVisibility;
}

export interface UpdateProfileLinkInput {
  kind?: string;
  uri?: string;
  title?: string;
  sort_order?: number;
  visibility?: LinkVisibility;
}

export interface CreateProfilePageInput {
  slug: string;
  title: string;
  content: string;
}

export interface UpdateProfilePageInput {
  title?: string;
  content?: string;
  sort_order?: number;
}

export interface UpdateProfileTranslationInput {
  locale: string;
  title: string;
  description?: string;
}

// Story CRUD Types
export interface InsertStoryInput {
  slug: string;
  kind: StoryKind;
  title: string;
  summary?: string;
  content: string;
  story_picture_uri?: string | null;
  publish_to_profiles?: string[];
}

export interface UpdateStoryInput {
  slug: string;
  story_picture_uri?: string | null;
}

export interface UpdateStoryTranslationInput {
  title: string;
  summary?: string;
  content: string;
}

export interface StoryPublication {
  id: string;
  story_id: string;
  profile_id: string;
  profile_slug: string;
  profile_title: string;
  profile_picture_uri: string | null;
  profile_kind: string;
  kind: string;
  is_featured: boolean;
  published_at: string | null;
  created_at: string;
}

export interface StoryEditData {
  id: string;
  kind: StoryKind;
  slug: string | null;
  locale_code: string;
  is_fallback: boolean;
  is_managed: boolean;
  story_picture_uri: string | null;
  title: string | null;
  summary: string | null;
  content: string;
  author_profile_id: string | null;
  author_profile_slug: string | null;
  created_at: string;
  updated_at: string | null;
  publications?: StoryPublication[];
}

export interface StoryPermissions {
  can_edit: boolean;
}

// Upload Types
export type UploadPurpose = "content-image" | "story-picture" | "profile-picture";

export interface GetPresignedURLRequest {
  filename: string;
  content_type: string;
  purpose: UploadPurpose;
}

export interface GetPresignedURLResponse {
  upload_url: string;
  key: string;
  public_url: string;
  expires_at: string;
}

// Profile Points Types
export type ProfilePointTransactionType = "GAIN" | "TRANSFER" | "SPEND";

export interface ProfilePointTransaction {
  id: string;
  target_profile_id: string;
  origin_profile_id: string | null;
  transaction_type: ProfilePointTransactionType;
  triggering_event: string | null;
  description: string;
  amount: number;
  balance_after: number;
  created_at: string;
}

// Pending Award Types
export type PendingAwardStatus = "pending" | "approved" | "rejected";

export interface PendingAwardProfileInfo {
  slug: string;
  title: string;
}

export interface PendingAward {
  id: string;
  target_profile_id: string;
  target_profile?: PendingAwardProfileInfo;
  triggering_event: string;
  description: string;
  amount: number;
  status: PendingAwardStatus;
  reviewed_by: string | null;
  reviewed_at: string | null;
  rejection_reason: string | null;
  metadata: Record<string, unknown> | null;
  created_at: string;
}

export interface PendingAwardsStats {
  total_pending: number;
  total_approved: number;
  total_rejected: number;
  points_awarded: number;
  by_event_type: Record<string, number>;
}

// Cursor Types
export interface CursoredResponse<T> {
  data: T;
  next_cursor: string | null;
}

// Spotlight Types (for command palette)
export interface SpotlightItem {
  icon: string;
  to: string;
  title: string;
}

// Search Result Types
export interface SearchResult {
  type: "profile" | "story" | "page";
  id: string;
  slug: string;
  title: string;
  summary: string | null;
  image_uri: string | null;
  profile_slug: string | null;
  profile_title: string | null;
  kind: string | null;
  rank: number;
}

// POW (Proof of Work) Types
export interface POWChallenge {
  pow_challenge_id: string;
  prefix: string;
  difficulty: number;
  expires_at: string;
}

export interface POWChallengeDisabled {
  enabled: false;
}

export interface POWSolution {
  nonce: string;
}

export interface POWSolverRequest {
  prefix: string;
  difficulty: number;
}

export interface POWSolverResponse {
  nonce: string;
  iterations: number;
  elapsed_ms: number;
}

// Profile Question Types
export interface ProfileQuestion {
  id: string;
  profile_id: string;
  author_user_id: string;
  content: string;
  author_name: string | null;
  author_slug: string | null;
  answer_content: string | null;
  answer_uri: string | null;
  answer_kind: string | null;
  answered_at: string | null;
  answered_by: string | null;
  vote_count: number;
  is_anonymous: boolean;
  is_hidden: boolean;
  has_viewer_vote: boolean;
  created_at: string;
  updated_at: string | null;
}
