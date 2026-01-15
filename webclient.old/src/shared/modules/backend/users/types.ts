export type UserKind = "regular" | "editor" | "admin";

export type User = {
  id: string;
  kind: UserKind;
  name: string;
  email: string;
  phone: string | null;
  github_handle: string | null;
  bsky_handle: string | null;
  x_handle: string | null;
  individual_profile_id: string | null;
  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
};
