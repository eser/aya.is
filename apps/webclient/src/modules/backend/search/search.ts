// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { SearchResult } from "../types";

// Re-export for backward compatibility
export type { SearchResult } from "../types";

export async function search(
  locale: string,
  query: string,
  profileSlug?: string,
  limit: number = 20,
): Promise<SearchResult[] | null> {
  const params = new URLSearchParams({
    q: query,
    limit: String(limit),
  });

  if (profileSlug !== undefined && profileSlug !== "") {
    params.set("profile", profileSlug);
  }

  return await fetcher<SearchResult[]>(locale, `/search?${params}`);
}
