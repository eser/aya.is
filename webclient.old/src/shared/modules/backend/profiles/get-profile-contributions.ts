import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { ProfileMembership } from "./types.ts";

export type GetProfileContributionsData = ProfileMembership[];

export async function getProfileContributions(locale: string, slug: string) {
  const response = await fetcher<GetProfileContributionsData>(
    `/${locale}/profiles/${slug}/contributions`,
  );

  return response;
}
