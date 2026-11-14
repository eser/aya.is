// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { StoryEx } from "../types";

export async function getProfileAuthoredStories(locale: string, slug: string): Promise<StoryEx[] | null> {
  return await fetcher<StoryEx[]>(locale, `/profiles/${slug}/stories-authored`);
}
