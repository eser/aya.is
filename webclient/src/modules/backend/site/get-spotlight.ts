import { fetcher } from "../fetcher";
import type { Profile } from "../types";

export type GetSpotlightData = Profile[];

export async function getSpotlight(): Promise<GetSpotlightData | null> {
  return fetcher<GetSpotlightData>("/en/site/spotlight");
}
