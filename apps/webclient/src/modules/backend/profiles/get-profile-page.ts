import { fetcher } from "../fetcher";
import type { ProfilePage } from "../types";

export async function getProfilePage(
  locale: string,
  slug: string,
  pageSlug: string,
): Promise<ProfilePage | null> {
  return await fetcher<ProfilePage>(`/${locale}/profiles/${slug}/pages/${pageSlug}`);
}
