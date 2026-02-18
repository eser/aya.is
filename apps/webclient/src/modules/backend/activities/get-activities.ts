import { fetcher } from "../fetcher";
import type { StoryEx } from "../types";

export type GetActivitiesData = StoryEx[];

export async function getActivities(
  locale: string,
): Promise<StoryEx[] | null> {
  const response = await fetcher<GetActivitiesData>(
    locale,
    `/activities`,
  );

  return response;
}
