import { fetcher } from "../fetcher";

export type UpdateProfilePageTranslationRequest = {
  title: string;
  summary: string;
  content: string;
};

export async function updateProfilePageTranslation(
  locale: string,
  profileSlug: string,
  pageId: string,
  translationLocale: string,
  data: UpdateProfilePageTranslationRequest
): Promise<{ success: boolean } | null> {
  return fetcher<{ success: boolean }>(
    `/${locale}/profiles/${profileSlug}/_pages/${pageId}/translations/${translationLocale}`,
    {
      method: "PATCH",
      body: JSON.stringify(data),
    }
  );
}
