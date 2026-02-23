import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileTeam } from "../types";

export async function createProfileTeam(
  locale: string,
  slug: string,
  name: string,
  description?: string | null,
): Promise<ProfileTeam | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_teams`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({ name, description: description ?? null }),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data ?? null;
}
