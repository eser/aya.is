// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { DateProposalVoteResponse } from "../types";

export async function voteDateProposal(
  locale: string,
  storySlug: string,
  proposalId: string,
  direction: number,
): Promise<DateProposalVoteResponse | null> {
  const response = await fetcher<DateProposalVoteResponse>(
    locale,
    `/stories/${storySlug}/date-proposals/${proposalId}/vote`,
    {
      method: "POST",
      body: JSON.stringify({ direction }),
    },
  );

  return response;
}
