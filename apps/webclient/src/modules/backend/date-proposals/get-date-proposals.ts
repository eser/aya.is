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
