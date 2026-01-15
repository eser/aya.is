import { fetcher } from "../fetcher";
import type { Story } from "../types";

export async function getProfileStories(
  locale: string,
  slug: string
): Promise<Story[] | null> {
  return fetcher<Story[]>(`/${locale}/profiles/${slug}/stories`);
}
