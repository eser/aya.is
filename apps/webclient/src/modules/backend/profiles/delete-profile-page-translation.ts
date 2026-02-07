import { fetcher } from "@/modules/backend/fetcher";

export interface DeleteProfilePageTranslationResult {
  success: boolean;
}

export async function deleteProfilePageTranslation(
  locale: string,
  profileSlug: string,
  pageId: string,
  translationLocale: string,
): Promise<DeleteProfilePageTranslationResult | null> {
  const response = await fetcher<DeleteProfilePageTranslationResult>(
    locale,
    `/profiles/${profileSlug}/_pages/${pageId}/translations/${translationLocale}`,
    {
      method: "DELETE",
    },
  );
  return response;
}
