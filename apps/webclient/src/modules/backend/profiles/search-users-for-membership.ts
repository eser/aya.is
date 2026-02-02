import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { UserSearchResult } from "../types";

export async function searchUsersForMembership(
  locale: string,
  slug: string,
  query: string,
): Promise<UserSearchResult[] | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const encodedQuery = encodeURIComponent(query);
  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_memberships/search?q=${encodedQuery}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
    },
  );

  if (!response.ok) return null;
  return await response.json();
}
