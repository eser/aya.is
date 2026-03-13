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
