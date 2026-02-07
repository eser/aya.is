import { fetcher } from "@/modules/backend/fetcher";

export async function listProfilePageTranslationLocales(
  locale: string,
  profileSlug: string,
  pageId: string,
): Promise<string[] | null> {
  const response = await fetcher<string[]>(
    locale,
    `/profiles/${profileSlug}/_pages/${pageId}/_tx`,
    {
      method: "GET",
    },
  );

  return response;
}
