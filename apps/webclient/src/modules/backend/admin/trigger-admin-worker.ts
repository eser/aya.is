import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export interface TriggerWorkerResult {
  name: string;
  triggered: boolean;
}

export async function triggerAdminWorker(
  name: string,
): Promise<TriggerWorkerResult | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/admin/workers/${encodeURIComponent(name)}/trigger`,
    {
      method: "POST",
      headers,
      credentials: "include",
    },
  );

  if (!response.ok) return null;
  return response.json();
}
