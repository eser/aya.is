import { fetcher } from "@/modules/backend/fetcher";
import type { StoryPermissions } from "@/modules/backend/types";

export async function getStoryPermissions(
  locale: string,
  profileSlug: string,
  storyId: string,
): Promise<StoryPermissions | null> {
  const response = await fetcher<StoryPermissions>(
    locale,
    `/profiles/${profileSlug}/_stories/${storyId}/_permissions`,
  );
  return response;
}
