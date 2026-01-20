import { fetcher } from "../fetcher";
import type { Profile } from "../types";

export async function updateProfilePicture(
  locale: string,
  slug: string,
  profilePictureUri: string,
): Promise<Profile | null> {
  return await fetcher<Profile>(`/${locale}/profiles/${slug}`, {
    method: "PATCH",
    body: JSON.stringify({ profile_picture_uri: profilePictureUri }),
  });
}
