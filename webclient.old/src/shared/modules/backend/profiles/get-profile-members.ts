import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { ProfileMembership } from "./types.ts";

export type GetProfileMembersData = ProfileMembership[];

export async function getProfileMembers(locale: string, slug: string) {
  const response = await fetcher<GetProfileMembersData>(
    `/${locale}/profiles/${slug}/members`,
  );

  return response;
}
