// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { StoryEx } from "../types";

export type GetStoriesData = StoryEx[];

export async function getStoriesByKinds(
  locale: string,
  kinds: string[],
): Promise<StoryEx[] | null> {
  const response = await fetcher<GetStoriesData>(
    locale,
    `/stories?filter_kind=${kinds.join(",")}`,
  );

  return response;
}
