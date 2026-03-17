// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { ProfilePage } from "../types";

export async function listProfilePages(
  locale: string,
  profileSlug: string,
): Promise<ProfilePage[] | null> {
  return await fetcher<ProfilePage[]>(locale, `/profiles/${profileSlug}/_pages`);
}
