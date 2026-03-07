import { fetcher } from "../fetcher";

export async function finalizeDateProposal(
  locale: string,
  storySlug: string,
  proposalId: string,
): Promise<boolean> {
  await fetcher(
    locale,
    `/stories/${storySlug}/date-proposals/${proposalId}/finalize`,
    {
      method: "POST",
    },
  );

  return true;
}
