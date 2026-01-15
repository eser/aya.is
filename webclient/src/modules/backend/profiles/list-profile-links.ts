import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileLink } from "../types";

export async function listProfileLinks(
  locale: string,
  slug: string
): Promise<ProfileLink[] | null> {
  const token = getAuthToken();
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (token !== null) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(`${getBackendUri()}/${locale}/profiles/${slug}/_links`, {
    method: "GET",
    headers,
    credentials: "include",
  });

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
