import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { Profile } from "./types.ts";

export type GetProfilesData = Profile[];

export async function getProfilesByKinds(locale: string, kinds: string[]) {
  const response = await fetcher<GetProfilesData>(
    `/${locale}/profiles?filter_kind=${kinds.join(",")}`,
  );

  return response;
}
