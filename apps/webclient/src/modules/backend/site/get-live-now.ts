import { fetcher } from "../fetcher";
import type { LiveStreamInfo } from "../types";

export async function getLiveNow(locale: string): Promise<LiveStreamInfo[] | null> {
  return await fetcher<LiveStreamInfo[]>(locale, `/site/live-now`);
}
