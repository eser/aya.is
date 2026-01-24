import { fetcher } from "../fetcher";

export interface SpotlightItem {
  icon: string;
  to: string;
  title: string;
}

export type GetSpotlightData = SpotlightItem[];

export async function getSpotlight(locale: string): Promise<GetSpotlightData | null> {
  return await fetcher<GetSpotlightData>(locale, `/site/spotlight`);
}
