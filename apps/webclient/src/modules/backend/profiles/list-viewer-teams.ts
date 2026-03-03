import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileTeam } from "../types";

export async function listViewerTeams(
  locale: string,
  slug: string,
): Promise<ProfileTeam[] | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_my-teams`,
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
  const result = await response.json();
  return result.data ?? null;
}
