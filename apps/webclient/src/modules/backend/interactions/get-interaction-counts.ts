// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { InteractionCount } from "../types";

export type GetInteractionCountsData = InteractionCount[];

export async function getInteractionCounts(
  locale: string,
  storySlug: string,
): Promise<InteractionCount[] | null> {
  const response = await fetcher<GetInteractionCountsData>(
    locale,
    `/stories/${storySlug}/interactions/counts`,
  );

  return response;
}
