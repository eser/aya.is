import { fetcher } from "../fetcher";
import type { CandidateFormResponse } from "../types";

export async function getCandidateResponses(
  locale: string,
  slug: string,
  candidateId: string,
): Promise<CandidateFormResponse[] | null> {
  return await fetcher<CandidateFormResponse[]>(
    locale,
    `/profiles/${slug}/_candidates/${candidateId}/responses`,
  );
}
