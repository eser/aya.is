import { fetcher } from "@/shared/modules/backend/fetcher.ts";
import type { Profile } from "./types.ts";

export type UpdateProfileRequest = {
  profile_picture_uri?: string | null;
  pronouns?: string | null;
  // deno-lint-ignore no-explicit-any
  properties?: Record<string, any> | null;
};

export async function updateProfile(
  locale: string,
  slug: string,
  data: UpdateProfileRequest,
): Promise<Profile | null> {
  const response = await fetcher<Profile>(`/${locale}/profiles/${slug}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  return response;
}
