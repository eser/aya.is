import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfilePointTransaction } from "../types";

export async function listProfilePointTransactions(
  locale: string,
  slug: string,
): Promise<ProfilePointTransaction[] | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_points/transactions`,
    {
      method: "GET",
      headers,
      credentials: "include",
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
