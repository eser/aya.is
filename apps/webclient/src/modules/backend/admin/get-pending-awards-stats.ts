// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { PendingAwardsStats } from "../types";

export async function getPendingAwardsStats(): Promise<PendingAwardsStats | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(`${getBackendUri()}/admin/points/stats`, {
    method: "GET",
    headers,
    credentials: "include",
  });

  if (!response.ok) return null;
  return response.json();
}
