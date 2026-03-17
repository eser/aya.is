// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { ProfileMembershipCandidate } from "../types";

export async function listCandidates(
  locale: string,
  slug: string,
): Promise<ProfileMembershipCandidate[] | null> {
  return await fetcher<ProfileMembershipCandidate[]>(
    locale,
    `/profiles/${slug}/_candidates`,
  );
}
