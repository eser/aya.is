// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { CandidateVote } from "../types";

export async function getCandidateVotes(
  locale: string,
  slug: string,
  candidateId: string,
): Promise<CandidateVote[] | null> {
  return await fetcher<CandidateVote[]>(
    locale,
    `/profiles/${slug}/_candidates/${candidateId}/votes`,
  );
}
