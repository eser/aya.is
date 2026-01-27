import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { Profile } from "../types";

export async function getAdminProfile(
  locale: string,
  slug: string,
): Promise<Profile | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/admin/profiles/${slug}?locale=${locale}`,
    {
      method: "GET",
      headers,
      credentials: "include",
    },
  );

  if (!response.ok) {
    return null;
  }

  return response.json();
}
