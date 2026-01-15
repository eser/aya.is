import { fetcher } from "../fetcher.ts";
import type { ProfilePage } from "./types.ts";

export type ListProfilePagesData = Omit<ProfilePage, "order" | "created_at" | "updated_at" | "deleted_at">[];

export async function listProfilePages(locale: string, profileSlug: string) {
  const response = await fetcher<ListProfilePagesData>(`/${locale}/profiles/${profileSlug}/_pages`);

  return response;
}
