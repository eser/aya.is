// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export interface ToggleWorkerResult {
  name: string;
  is_enabled: boolean;
}

export async function toggleAdminWorker(
  name: string,
): Promise<ToggleWorkerResult | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/admin/workers/${encodeURIComponent(name)}/toggle`,
    {
      method: "POST",
      headers,
      credentials: "include",
    },
  );

  if (!response.ok) return null;
  return response.json();
}
