import { fetcher } from "@/modules/backend/fetcher";
import type { StoryEx } from "@/modules/backend/types";

export type GetProfileStoryData = StoryEx;

export async function getProfileStory(
  locale: string,
  slug: string,
  storySlug: string,
): Promise<StoryEx | null> {
  const response = await fetcher<GetProfileStoryData>(
    locale,
    `/profiles/${slug}/stories/${storySlug}`,
  );
  return response;
}
