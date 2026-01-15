export type ProfileKind = "individual" | "product" | "organization";

export type Profile = {
  id: string;
  kind: ProfileKind;
  slug: string;
  custom_domain: string | null;
  profile_picture_uri: string | null;
  pronouns: string | null;
  title: string;
  description: string;

  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;

  properties: Record<string, unknown> | null;
};

export type ProfileLinkKind = "website" | "github" | "bsky" | "x";

export type ProfileLink = {
  id: string;
  kind: ProfileLinkKind;

  profile_id: string;
  order: number;

  is_managed: boolean;
  is_verified: boolean;

  remote_id: string | null;
  public_id: string | null;
  uri: string | null;
  title: string;

  auth_provider: string | null;
  auth_access_token_scope: string | null;
  auth_access_token: string | null;
  auth_access_token_expires_at: string | null;
  auth_refresh_token: string | null;
  auth_refresh_token_expires_at: string | null;

  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
};

export type ProfilePage = {
  id: string;
  slug: string;

  profile_id: string;
  order: number;

  cover_picture_uri: string | null;
  title: string;
  summary: string;
  content: string;

  published_at: string | null;

  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
};

export type ProfileMembershipKind = "follower" | "sponsor" | "contributor" | "maintainer" | "lead" | "owner";

export type ProfileMembership = {
  id: string;
  kind: ProfileMembershipKind;
  properties: {
    github?: {
      commits: number;
      prs: {
        resolved: number;
        total: number;
      };
      issues: {
        resolved: number;
        total: number;
      };
      stars: number;
    };
    videos?: number;
  };

  profile: Profile;
  member_profile: Profile;
};
