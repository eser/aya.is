import { fetcher } from "../fetcher.ts";

export type ProfileTranslation = {
  profile_id: string;
  locale_code: string;
  title: string;
  description: string;
  properties: Record<string, unknown> | null;
};

export async function getProfileTranslations(locale: string, profileSlug: string) {
  const url = `/${locale}/profiles/${profileSlug}/_tx`;

  const response = await fetcher<ProfileTranslation[]>(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return response;
}
