// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { ProfilePermissions } from "../types";

export async function getProfilePermissions(
  locale: string,
  slug: string,
): Promise<ProfilePermissions | null> {
  return await fetcher<ProfilePermissions>(locale, `/profiles/${slug}/_permissions`);
}
