// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { DateProposalListResponse } from "../types";

export async function getDateProposals(
  locale: string,
  storySlug: string,
): Promise<DateProposalListResponse | null> {
  const response = await fetcher<DateProposalListResponse>(
    locale,
    `/stories/${storySlug}/date-proposals`,
  );

  return response;
}
