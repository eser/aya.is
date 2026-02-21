import { fetcher } from "../fetcher";
import type { ContentVisibility, ProfilePage } from "../types";

export type UpdateProfilePageRequest = {
  slug: string;
  order: number;
  cover_picture_uri: string | null;
  published_at: string | null;
  visibility: ContentVisibility;
};

export async function updateProfilePage(
  locale: string,
  profileSlug: string,
  pageId: string,
  data: UpdateProfilePageRequest,
): Promise<ProfilePage | null> {
  return await fetcher<ProfilePage>(
    locale,
    `/profiles/${profileSlug}/_pages/${pageId}`,
    {
      method: "PATCH",
      body: JSON.stringify(data),
    },
  );
}
