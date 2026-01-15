import { fetcher } from "../fetcher";
import type { CustomDomain } from "../types";

export async function getCustomDomain(
  host: string
): Promise<CustomDomain | null> {
  return fetcher<CustomDomain>(`/custom-domains/${host}`);
}
