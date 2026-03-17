// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { BulletinPreferences } from "./types";

export async function getBulletinPreferences(
  locale: string,
): Promise<BulletinPreferences | null> {
  return await fetcher<BulletinPreferences>(
    locale,
    "/bulletin/preferences",
  );
}
