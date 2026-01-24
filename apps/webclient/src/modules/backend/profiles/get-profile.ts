import { fetcher } from "../fetcher";
import type { Profile } from "../types";

export async function getProfile(
  locale: string,
  slug: string,
): Promise<Profile | null> {
  return await fetcher<Profile>(locale, `/profiles/${slug}`);
}
