import { fetcher } from "../fetcher";
import type { ProfilePage } from "../types";

export async function listProfilePages(
  locale: string,
  profileSlug: string,
): Promise<ProfilePage[] | null> {
  return await fetcher<ProfilePage[]>(locale, `/profiles/${profileSlug}/_pages`);
}
