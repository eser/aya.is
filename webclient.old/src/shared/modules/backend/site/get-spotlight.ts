import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { Profile } from "../profiles/types.ts";

export type GetSpotlightData = Profile[];

export async function getSpotlight() {
  const response = await fetcher<GetSpotlightData>(`/en/site/spotlight`);

  return response;
}
