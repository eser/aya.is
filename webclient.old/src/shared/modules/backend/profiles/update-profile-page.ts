import { fetcher } from "../fetcher.ts";

export type UpdateProfilePageParams = {
  slug: string;
  order: number;
  cover_picture_uri: string | null;
  published_at: string | null;
};

export async function updateProfilePage(
  locale: string,
  profileSlug: string,
  pageId: string,
  params: UpdateProfilePageParams,
) {
  const response = await fetcher(`/${locale}/profiles/${profileSlug}/_pages/${pageId}`, {
    method: "PATCH",
    body: JSON.stringify(params),
  });

  return response;
}
