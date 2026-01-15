import { fetcher } from "../fetcher";
import type { ProfilePage } from "../types";

export type CreateProfilePageRequest = {
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
  data: CreateProfilePageRequest
): Promise<ProfilePage | null> {
  return fetcher<ProfilePage>(`/${locale}/profiles/${profileSlug}/_pages`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}
