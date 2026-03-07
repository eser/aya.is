import { fetcher } from "../fetcher";
import type { DateProposal } from "../types";

export async function getDateProposals(
  locale: string,
  storySlug: string,
): Promise<DateProposal[] | null> {
  const response = await fetcher<DateProposal[]>(
    locale,
    `/stories/${storySlug}/date-proposals`,
  );

  return response;
}
