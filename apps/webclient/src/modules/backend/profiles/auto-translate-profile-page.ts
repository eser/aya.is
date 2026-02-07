import { fetcher } from "@/modules/backend/fetcher";

export interface AutoTranslateProfilePageResult {
  success: boolean;
  title: string;
  summary: string;
  content: string;
}

export async function autoTranslateProfilePage(
  locale: string,
  profileSlug: string,
  pageId: string,
  targetLocale: string,
  sourceLocale: string,
): Promise<AutoTranslateProfilePageResult | null> {
  const response = await fetcher<AutoTranslateProfilePageResult>(
    locale,
    `/profiles/${profileSlug}/_pages/${pageId}/translations/${targetLocale}/auto-translate`,
    {
      method: "POST",
      body: JSON.stringify({ source_locale: sourceLocale }),
    },
  );
  return response;
}
