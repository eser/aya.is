import { fetcher } from "../fetcher";
import type { ProfileTranslation } from "../types";

export async function getProfileTranslations(
  locale: string,
  profileSlug: string
): Promise<ProfileTranslation[] | null> {
  return fetcher<ProfileTranslation[]>(`/${locale}/profiles/${profileSlug}/_tx`);
}
