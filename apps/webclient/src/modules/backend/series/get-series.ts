// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { StorySeriesWithStories } from "../types";

export async function getSeries(
  locale: string,
  slug: string,
): Promise<StorySeriesWithStories | null> {
  const response = await fetcher<StorySeriesWithStories>(
    locale,
    `/series/${slug}`,
  );

  return response;
}
