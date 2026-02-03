import { fetcher } from "@/modules/backend/fetcher";
import type { StoryPublication } from "@/modules/backend/types";

export async function listStoryPublications(
  locale: string,
  profileSlug: string,
  storyId: string,
): Promise<StoryPublication[] | null> {
  const response = await fetcher<StoryPublication[]>(
    locale,
    `/profiles/${profileSlug}/_stories/${storyId}/publications`,
    {
      method: "GET",
    },
  );

  return response;
}
