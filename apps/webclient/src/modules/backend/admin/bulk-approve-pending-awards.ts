import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export interface BulkResultItem {
  id: string;
  status: string;
  error?: string;
  transaction_id?: string;
}

export interface BulkApproveResult {
  results: BulkResultItem[];
}

export async function bulkApprovePendingAwards(
  ids: string[],
): Promise<BulkApproveResult | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/admin/points/pending/bulk-approve`,
    {
      method: "POST",
      headers,
      credentials: "include",
      body: JSON.stringify({ ids }),
    },
  );

  if (!response.ok) return null;
  return response.json();
}
