// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { Profile } from "../types";

export async function getProfile(
  locale: string,
  slug: string,
): Promise<Profile | null> {
  return await fetcher<Profile>(locale, `/profiles/${slug}`);
}
