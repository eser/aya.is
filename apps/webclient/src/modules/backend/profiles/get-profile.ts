import { fetcher } from "../fetcher";
import type { Profile } from "../types";

export async function getProfile(
  locale: string,
  slug: string,
  fallbackLocale?: string,
): Promise<Profile | null> {
  const query = fallbackLocale !== undefined ? `?fallback_locale=${fallbackLocale}` : "";
  return await fetcher<Profile>(locale, `/profiles/${slug}${query}`);
}
