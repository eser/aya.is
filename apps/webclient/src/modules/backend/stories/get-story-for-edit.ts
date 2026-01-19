import { fetcher } from "@/modules/backend/fetcher";
import type { StoryEditData } from "@/modules/backend/types";

export async function getStoryForEdit(
  locale: string,
  profileSlug: string,
  storyId: string,
): Promise<StoryEditData | null> {
  const response = await fetcher<StoryEditData>(
    `/${locale}/profiles/${profileSlug}/_stories/${storyId}`,
  );
  return response;
}
