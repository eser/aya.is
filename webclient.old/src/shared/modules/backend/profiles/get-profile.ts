import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { Profile, ProfileLink, ProfilePage } from "./types.ts";

export type GetProfileData = Omit<Profile, "deleted_at"> & {
  links: Omit<
    ProfileLink,
    "profile_id" | "auth_token_scope" | "auth_token" | "auth_refresh_token" | "created_at" | "updated_at" | "deleted_at"
  >[];
  pages: Omit<
    ProfilePage,
    "profile_id" | "order" | "content" | "published_at" | "created_at" | "updated_at" | "deleted_at"
  >[];
};

export async function getProfile(locale: string, slug: string) {
  const response = await fetcher<GetProfileData>(`/${locale}/profiles/${slug}`);

  return response;
}
