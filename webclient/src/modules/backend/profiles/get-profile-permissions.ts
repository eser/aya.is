import { fetcher } from "../fetcher";
import type { ProfilePermissions } from "../types";

export async function getProfilePermissions(
  locale: string,
  slug: string
): Promise<ProfilePermissions | null> {
  return fetcher<ProfilePermissions>(`/${locale}/profiles/${slug}/permissions`);
}
