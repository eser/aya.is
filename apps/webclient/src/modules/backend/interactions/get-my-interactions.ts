// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { StoryInteraction } from "../types";

export type GetMyInteractionsData = StoryInteraction[];

export async function getMyInteractions(
  locale: string,
  storySlug: string,
): Promise<StoryInteraction[] | null> {
  const response = await fetcher<GetMyInteractionsData>(
    locale,
    `/stories/${storySlug}/interactions/me`,
  );

  return response;
}
