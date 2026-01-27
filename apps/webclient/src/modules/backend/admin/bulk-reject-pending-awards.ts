import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export interface BulkResultItem {
  id: string;
  status: string;
  error?: string;
}

export interface BulkRejectResult {
  results: BulkResultItem[];
}

export async function bulkRejectPendingAwards(
  ids: string[],
  reason?: string,
): Promise<BulkRejectResult | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/admin/points/pending/bulk-reject`,
    {
      method: "POST",
      headers,
      credentials: "include",
      body: JSON.stringify({ ids, reason: reason ?? "" }),
    },
  );

  if (!response.ok) return null;
  return response.json();
}
