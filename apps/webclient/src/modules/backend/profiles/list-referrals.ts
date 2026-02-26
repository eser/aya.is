import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileMembershipReferral } from "../types";

export async function listReferrals(
  locale: string,
  slug: string,
): Promise<ProfileMembershipReferral[] | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_referrals`,
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
