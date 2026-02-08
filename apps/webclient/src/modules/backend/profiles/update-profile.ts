import { fetcher } from "../fetcher";
import type { Profile } from "../types";

export type UpdateProfileRequest = {
  profile_picture_uri?: string | null;
  pronouns?: string | null;
  properties?: Record<string, unknown> | null;
  hide_relations?: boolean;
  hide_links?: boolean;
  hide_qa?: boolean;
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
