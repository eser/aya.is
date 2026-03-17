// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";

export async function deleteProfilePage(
  locale: string,
  profileSlug: string,
  pageId: string,
): Promise<{ success: boolean } | null> {
  return await fetcher<{ success: boolean }>(
    locale,
    `/profiles/${profileSlug}/_pages/${pageId}`,
    { method: "DELETE" },
  );
}
