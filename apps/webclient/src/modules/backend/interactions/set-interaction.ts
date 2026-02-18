import { fetcher } from "../fetcher";
import type { StoryInteraction } from "../types";

export async function setInteraction(
  locale: string,
  storySlug: string,
  kind: string,
): Promise<StoryInteraction | null> {
  const response = await fetcher<StoryInteraction>(
    locale,
    `/stories/${storySlug}/interactions`,
    {
      method: "POST",
      body: JSON.stringify({ kind }),
    },
  );

  return response;
}
