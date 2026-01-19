import { fetcher } from "@/modules/backend/fetcher";
import type { UpdateStoryInput, Story } from "@/modules/backend/types";

export type UpdateStoryData = Story;

export async function updateStory(
  locale: string,
  profileSlug: string,
  storyId: string,
  input: UpdateStoryInput,
): Promise<Story | null> {
  const response = await fetcher<UpdateStoryData>(
    `/${locale}/profiles/${profileSlug}/_stories/${storyId}`,
    {
      method: "PATCH",
      body: JSON.stringify(input),
    },
  );
  return response;
}
