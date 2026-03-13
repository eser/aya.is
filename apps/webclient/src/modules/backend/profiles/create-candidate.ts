import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileMembershipCandidate } from "../types";

export async function createCandidate(
  locale: string,
  slug: string,
  referredProfileSlug: string,
  teamIds: string[],
): Promise<ProfileMembershipCandidate | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_candidates`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({
        referred_profile_slug: referredProfileSlug,
        team_ids: teamIds,
      }),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data ?? null;
}
