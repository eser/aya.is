import { fetcher } from "../fetcher.ts";

export type CreateProfilePageParams = {
  slug: string;
  title: string;
  summary: string;
  content: string;
  cover_picture_uri: string | null;
  published_at: string | null;
};

export async function createProfilePage(
  locale: string,
  profileSlug: string,
  params: CreateProfilePageParams,
) {
  const url = `/${locale}/profiles/${profileSlug}/_pages`;

  const response = await fetcher(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });

  return response;
}
