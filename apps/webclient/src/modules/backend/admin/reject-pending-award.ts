import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export interface RejectResult {
  status: string;
}

export async function rejectPendingAward(
  awardId: string,
  reason?: string,
): Promise<RejectResult | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/admin/points/pending/${awardId}/reject`,
    {
      method: "POST",
      headers,
      credentials: "include",
      body: JSON.stringify({ reason: reason ?? "" }),
    },
  );

  if (!response.ok) return null;
  return response.json();
}
