import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { StoryEx } from "./types.ts";

export type GetStoryData = StoryEx;

export async function getStory(locale: string, storyslug: string) {
  const response = await fetcher<GetStoryData>(`/${locale}/stories/${storyslug}`);

  return response;
}
