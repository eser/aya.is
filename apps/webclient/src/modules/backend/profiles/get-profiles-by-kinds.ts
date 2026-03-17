// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { Profile } from "../types";

export async function getProfilesByKinds(
  locale: string,
  kinds: string[],
  seed?: string,
  limit?: number,
  offset?: number,
  q?: string,
): Promise<Profile[] | null> {
  const kindsParam = kinds.join(",");
  const params = new URLSearchParams({ filter_kind: kindsParam });

  if (seed !== undefined) {
    params.set("seed", seed);
  }

  if (limit !== undefined) {
    params.set("limit", String(limit));
  }

  if (offset !== undefined) {
    params.set("offset", String(offset));
  }

  if (q !== undefined && q !== "") {
    params.set("filter_q", q);
  }

  return await fetcher<Profile[]>(locale, `/profiles?${params.toString()}`);
}
