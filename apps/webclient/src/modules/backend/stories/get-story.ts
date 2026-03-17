// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "@/modules/backend/fetcher";
import type { StoryEx } from "@/modules/backend/types";

export type GetStoryData = StoryEx;

export async function getStory(
  locale: string,
  storyslug: string,
): Promise<StoryEx | null> {
  const response = await fetcher<GetStoryData>(
    locale,
    `/stories/${storyslug}`,
  );
  return response;
}
