import { fetcher } from "@/modules/backend/fetcher";
import type { StoryPublication } from "@/modules/backend/types";

export interface AddStoryPublicationInput {
  profile_id: string;
  is_featured?: boolean;
}

export async function addStoryPublication(
  locale: string,
  profileSlug: string,
  storyId: string,
  input: AddStoryPublicationInput,
): Promise<StoryPublication | null> {
  const response = await fetcher<StoryPublication>(
    locale,
    `/profiles/${profileSlug}/_stories/${storyId}/publications`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );

  return response;
}
