import { fetcher } from "../fetcher";
import type { ProfilePermissions } from "../types";

export async function getProfilePermissions(
  locale: string,
  slug: string,
): Promise<ProfilePermissions | null> {
  return await fetcher<ProfilePermissions>(`/${locale}/profiles/${slug}/_permissions`);
}
