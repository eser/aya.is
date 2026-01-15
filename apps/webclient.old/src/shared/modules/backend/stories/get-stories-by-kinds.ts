import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { StoryEx } from "./types.ts";

export type GetStoriesData = StoryEx[];

export async function getStoriesByKinds(locale: string, kinds: string[]) {
  const response = await fetcher<GetStoriesData>(
    `/${locale}/stories?filter_kind=${kinds.join(",")}`,
  );

  return response;
}
