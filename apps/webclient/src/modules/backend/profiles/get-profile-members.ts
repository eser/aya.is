// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { ProfileMembership } from "../types";

export type GetProfileMembersData = ProfileMembership[];

export async function getProfileMembers(
  locale: string,
  slug: string,
): Promise<ProfileMembership[] | null> {
  const response = await fetcher<GetProfileMembersData>(
    locale,
    `/profiles/${slug}/members`,
  );

  return response;
}
