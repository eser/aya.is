// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { ProfileTranslation } from "../types";

export async function getProfileTranslations(
  locale: string,
  profileSlug: string,
): Promise<ProfileTranslation[] | null> {
  return await fetcher<ProfileTranslation[]>(
    locale,
    `/profiles/${profileSlug}/_tx`,
  );
}
