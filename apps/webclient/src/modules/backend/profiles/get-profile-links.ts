import { fetcher } from "../fetcher";
import type { ProfileLink } from "../types";

export function getProfileLinks(
  locale: string,
  slug: string,
): Promise<ProfileLink[] | null> {
  return fetcher<ProfileLink[]>(locale, `/profiles/${slug}/links`);
}
