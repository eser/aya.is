import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { StoryEx } from "../stories/types.ts";

export type GetProfileStoryData = StoryEx;

export async function getProfileStory(locale: string, slug: string, storySlug: string) {
  const response = await fetcher<GetProfileStoryData>(`/${locale}/profiles/${slug}/stories/${storySlug}`);

  return response;
}
