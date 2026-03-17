// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "@/modules/backend/fetcher";
import type { StoryEx } from "@/modules/backend/types";

export type GetProfileStoryData = StoryEx;

export async function getProfileStory(
  locale: string,
  slug: string,
  storySlug: string,
): Promise<StoryEx | null> {
  const response = await fetcher<GetProfileStoryData>(
    locale,
    `/profiles/${slug}/stories/${storySlug}`,
  );
  return response;
}
