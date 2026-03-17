// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
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
