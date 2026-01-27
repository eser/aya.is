import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfilePointTransaction } from "../types";

export async function approvePendingAward(
  awardId: string,
): Promise<ProfilePointTransaction | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/admin/points/pending/${awardId}/approve`,
    {
      method: "POST",
      headers,
      credentials: "include",
    },
  );

  if (!response.ok) return null;
  return response.json();
}
