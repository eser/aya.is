// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";

export type ConnectDevtoResponse = {
  status: string;
};

export async function connectDevto(
  locale: string,
  slug: string,
  url: string,
): Promise<ConnectDevtoResponse | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/connect/devto`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({ url }),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
