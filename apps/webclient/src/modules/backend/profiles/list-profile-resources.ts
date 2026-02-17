import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileResource } from "../types";

export async function listProfileResources(
  locale: string,
  slug: string,
): Promise<ProfileResource[] | null> {
  const token = getAuthToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (token !== null) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_resources`,
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
