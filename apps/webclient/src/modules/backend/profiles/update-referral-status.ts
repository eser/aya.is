import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ReferralStatus } from "../types";

export async function updateReferralStatus(
  locale: string,
  slug: string,
  referralId: string,
  status: ReferralStatus,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_referrals/${referralId}/status`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({ status }),
    },
  );

  return response.ok;
}
