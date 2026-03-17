// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function unfollowProfile(
  locale: string,
  slug: string,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_follow`,
    {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
    },
  );

  return response.ok;
}
