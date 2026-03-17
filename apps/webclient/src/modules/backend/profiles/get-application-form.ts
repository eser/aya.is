// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { ApplicationForm } from "../types";

export async function getApplicationForm(
  locale: string,
  slug: string,
): Promise<ApplicationForm | null> {
  return await fetcher<ApplicationForm>(
    locale,
    `/profiles/${slug}/_application-form`,
  );
}
