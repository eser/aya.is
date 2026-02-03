import { fetcher } from "../fetcher";
import type { SpotlightItem } from "../types";

// Re-export for backward compatibility
export type { SpotlightItem } from "../types";

export type GetSpotlightData = SpotlightItem[];

export async function getSpotlight(locale: string): Promise<GetSpotlightData | null> {
  return await fetcher<GetSpotlightData>(locale, `/site/spotlight`);
}
