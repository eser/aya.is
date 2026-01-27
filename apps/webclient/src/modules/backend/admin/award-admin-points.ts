// Award points to a profile (admin only)

import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfilePointTransaction } from "../types";

export interface AwardAdminPointsParams {
  slug: string;
  amount: number;
  description: string;
}

export async function awardAdminPoints(
  params: AwardAdminPointsParams,
): Promise<ProfilePointTransaction | null> {
  const token = getAuthToken();

  if (token === null) {
    return null;
  }

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const url = `${getBackendUri()}/admin/profiles/${params.slug}/points`;

  const response = await fetch(url, {
    method: "POST",
    headers,
    credentials: "include",
    body: JSON.stringify({
      amount: params.amount,
      description: params.description,
    }),
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || "Failed to award points");
  }

  return response.json();
}
