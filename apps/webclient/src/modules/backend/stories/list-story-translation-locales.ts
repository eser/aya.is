import { fetcher } from "@/modules/backend/fetcher";

export async function listStoryTranslationLocales(
  locale: string,
  profileSlug: string,
  storyId: string,
): Promise<string[] | null> {
  const response = await fetcher<string[]>(
    locale,
    `/profiles/${profileSlug}/_stories/${storyId}/_tx`,
    {
      method: "GET",
    },
  );

  return response;
}
