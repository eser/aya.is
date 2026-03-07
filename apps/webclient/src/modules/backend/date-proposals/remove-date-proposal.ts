import { fetcher } from "../fetcher";

export async function removeDateProposal(
  locale: string,
  storySlug: string,
  proposalId: string,
): Promise<boolean> {
  await fetcher(
    locale,
    `/stories/${storySlug}/date-proposals/${proposalId}`,
    {
      method: "DELETE",
    },
  );

  return true;
}
