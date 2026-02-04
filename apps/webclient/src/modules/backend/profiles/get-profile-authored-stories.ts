import { fetcher } from "../fetcher";
import type { Story } from "../types";

export async function getProfileAuthoredStories(
  locale: string,
  slug: string,
): Promise<Story[] | null> {
  return await fetcher<Story[]>(locale, `/profiles/${slug}/stories-authored`);
}
