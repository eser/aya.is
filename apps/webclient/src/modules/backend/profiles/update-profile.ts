import { fetcher } from "../fetcher";
import type { Profile } from "../types";

export type UpdateProfileRequest = {
  profile_picture_uri?: string | null;
  pronouns?: string | null;
  properties?: Record<string, unknown> | null;
  feature_relations?: string;
  feature_links?: string;
  feature_qa?: string;
};

export async function updateProfile(
  locale: string,
  slug: string,
  data: UpdateProfileRequest,
): Promise<Profile | null> {
  return await fetcher<Profile>(locale, `/profiles/${slug}`, {
    method: "PATCH",
    body: JSON.stringify(data),
  });
}
