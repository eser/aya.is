// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { LiveStreamInfo } from "../types";

export async function getLiveNow(locale: string): Promise<LiveStreamInfo[] | null> {
  return await fetcher<LiveStreamInfo[]>(locale, `/site/live-now`);
}
