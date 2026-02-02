import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileMembershipWithMember } from "../types";

export async function listProfileMemberships(
  locale: string,
  slug: string,
): Promise<ProfileMembershipWithMember[] | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_memberships`,
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
