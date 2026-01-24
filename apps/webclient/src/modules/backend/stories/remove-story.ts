import { fetcher } from "@/modules/backend/fetcher";

export interface RemoveStoryResult {
  success: boolean;
  message: string;
}

export async function removeStory(
  locale: string,
  profileSlug: string,
  storyId: string,
): Promise<RemoveStoryResult | null> {
  const response = await fetcher<RemoveStoryResult>(
    locale,
    `/profiles/${profileSlug}/_stories/${storyId}`,
    {
      method: "DELETE",
    },
  );
  return response;
}
