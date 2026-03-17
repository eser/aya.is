// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";

export type CheckProfileSlugResponse = {
  available: boolean;
  message?: string;
  severity?: "error" | "warning" | "";
};

export async function checkProfileSlug(
  locale: string,
  slug: string,
): Promise<CheckProfileSlugResponse | null> {
  const result = await fetcher<CheckProfileSlugResponse>(
    locale,
    `/profiles/${slug}/_check`,
    {
      method: "GET",
    },
  );

  return result;
}
