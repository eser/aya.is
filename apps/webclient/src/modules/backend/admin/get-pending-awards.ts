// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { CursoredResponse, PendingAward } from "../types";

export interface GetPendingAwardsParams {
  status?: "pending" | "approved" | "rejected";
  limit?: number;
  cursor?: string;
}

export async function getPendingAwards(
  params: GetPendingAwardsParams = {},
): Promise<CursoredResponse<PendingAward[]> | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const queryParams = new URLSearchParams();
  if (params.status !== undefined) {
    queryParams.set("status", params.status);
  }
  if (params.limit !== undefined) {
    queryParams.set("limit", params.limit.toString());
  }
  if (params.cursor !== undefined) {
    queryParams.set("cursor", params.cursor);
  }

  const queryString = queryParams.toString();
  const url = `${getBackendUri()}/admin/points/pending${queryString.length > 0 ? `?${queryString}` : ""}`;

  const response = await fetch(url, {
    method: "GET",
    headers,
    credentials: "include",
  });

  if (!response.ok) return null;
  return response.json();
}
