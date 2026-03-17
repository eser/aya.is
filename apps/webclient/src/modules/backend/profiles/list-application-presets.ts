// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { ApplicationPreset } from "../types";

export async function listApplicationPresets(
  locale: string,
  slug: string,
): Promise<ApplicationPreset[] | null> {
  return await fetcher<ApplicationPreset[]>(
    locale,
    `/profiles/${slug}/_application-presets`,
  );
}
