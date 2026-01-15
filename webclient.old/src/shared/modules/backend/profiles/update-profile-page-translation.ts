import { fetcher } from "../fetcher.ts";

export type UpdateProfilePageTranslationParams = {
  title: string;
  summary: string;
  content: string;
};

export async function updateProfilePageTranslation(
  locale: string,
  profileSlug: string,
  pageId: string,
  translationLocale: string,
  params: UpdateProfilePageTranslationParams,
) {
  const url = `/${locale}/profiles/${profileSlug}/_pages/${pageId}/translations/${translationLocale}`;

  const response = await fetcher(url, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });

  return response;
}
