// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "@/modules/backend/fetcher";
import type { StoryEditData } from "@/modules/backend/types";

export async function getStoryForEdit(
  locale: string,
  profileSlug: string,
  storyIdOrSlug: string,
): Promise<StoryEditData | null> {
  const response = await fetcher<StoryEditData>(
    locale,
    `/profiles/${profileSlug}/_stories/${storyIdOrSlug}`,
  );
  return response;
}
