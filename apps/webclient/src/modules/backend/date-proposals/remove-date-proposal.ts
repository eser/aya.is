// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
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
