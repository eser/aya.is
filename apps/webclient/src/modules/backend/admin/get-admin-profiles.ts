import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { Profile } from "../types";

export interface AdminProfileListResult {
  data: Profile[];
  total: number;
  limit: number;
  offset: number;
}

export interface GetAdminProfilesParams {
  locale?: string;
  kind?: string;
  limit?: number;
  offset?: number;
}

export async function getAdminProfiles(
  params: GetAdminProfilesParams = {},
): Promise<AdminProfileListResult | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const searchParams = new URLSearchParams();

  if (params.locale !== undefined) {
    searchParams.set("locale", params.locale);
  }
  if (params.kind !== undefined && params.kind !== "") {
    searchParams.set("kind", params.kind);
  }
  if (params.limit !== undefined) {
    searchParams.set("limit", params.limit.toString());
  }
  if (params.offset !== undefined) {
    searchParams.set("offset", params.offset.toString());
  }

  const queryString = searchParams.toString();
  const url = `${getBackendUri()}/admin/profiles${queryString !== "" ? `?${queryString}` : ""}`;

  const response = await fetch(url, {
    method: "GET",
    headers,
    credentials: "include",
  });

  if (!response.ok) {
    return null;
  }

  return response.json();
}
