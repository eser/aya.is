import { fetcher } from "@/modules/backend/fetcher";

export interface AutoTranslateStoryResult {
  success: boolean;
  title: string;
  summary: string;
  content: string;
}

export async function autoTranslateStory(
  locale: string,
  profileSlug: string,
  storyId: string,
  targetLocale: string,
  sourceLocale: string,
): Promise<AutoTranslateStoryResult | null> {
  const response = await fetcher<AutoTranslateStoryResult>(
    locale,
    `/profiles/${profileSlug}/_stories/${storyId}/translations/${targetLocale}/auto-translate`,
    {
      method: "POST",
      body: JSON.stringify({ source_locale: sourceLocale }),
    },
  );
  return response;
}
