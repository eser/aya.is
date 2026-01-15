import { fetcher } from "../fetcher";

export async function deleteProfilePage(
  locale: string,
  profileSlug: string,
  pageId: string
): Promise<{ success: boolean } | null> {
  return fetcher<{ success: boolean }>(
    `/${locale}/profiles/${profileSlug}/_pages/${pageId}`,
    { method: "DELETE" }
  );
}
