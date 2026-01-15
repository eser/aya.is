import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { Story } from "../stories/types.ts";

export type getProfileStoriesData = Story[];

export async function getProfileStories(locale: string, slug: string) {
  const response = await fetcher<getProfileStoriesData>(`/${locale}/profiles/${slug}/stories`);

  return response;
}
