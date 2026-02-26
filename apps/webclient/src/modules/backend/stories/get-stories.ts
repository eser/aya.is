import { fetcher } from "../fetcher";
import type { StoryEx } from "../types";

export async function getStories(
  locale: string,
): Promise<StoryEx[] | null> {
  const response = await fetcher<StoryEx[]>(
    locale,
    "/stories",
  );

  return response;
}
