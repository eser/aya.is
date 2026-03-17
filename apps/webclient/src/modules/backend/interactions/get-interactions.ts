// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { InteractionWithProfile } from "../types";

export type GetInteractionsData = InteractionWithProfile[];

export async function getInteractions(
  locale: string,
  storySlug: string,
  kind?: string,
): Promise<InteractionWithProfile[] | null> {
  const query = kind !== undefined ? `?kind=${kind}` : "";

  const response = await fetcher<GetInteractionsData>(
    locale,
    `/stories/${storySlug}/interactions${query}`,
  );

  return response;
}
