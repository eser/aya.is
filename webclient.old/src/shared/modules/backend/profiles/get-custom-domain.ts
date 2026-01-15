import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { Profile } from "./types.ts";

export type GetCustomDomainData = Profile;

export async function getCustomDomain(domain: string) {
  const response = await fetcher<GetCustomDomainData>(`/en/site/custom-domains/${domain}`);

  return response;
}
