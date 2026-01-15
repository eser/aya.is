import { fetcher } from "../fetcher.ts";

export async function deleteProfilePage(
  locale: string,
  profileSlug: string,
  pageId: string,
) {
  const url = `/${locale}/profiles/${profileSlug}/_pages/${pageId}`;

  const response = await fetcher(url, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return response;
}
