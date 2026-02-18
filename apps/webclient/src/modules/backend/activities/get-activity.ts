import { fetcher } from "../fetcher";
import type { StoryEx } from "../types";

export type GetActivityData = StoryEx;

export async function getActivity(
  locale: string,
  slug: string,
): Promise<StoryEx | null> {
  const response = await fetcher<GetActivityData>(
    locale,
    `/activities/${slug}`,
  );

  return response;
}
