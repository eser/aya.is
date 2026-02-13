import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export interface AdminWorkerStatus {
  name: string;
  is_running: boolean;
  is_enabled: boolean;
  last_run: string | null;
  next_run: string | null;
  last_error: string | null;
  success_count: number;
  skip_count: number;
  error_count: number;
  interval: string;
}

export async function getAdminWorkers(): Promise<AdminWorkerStatus[] | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(`${getBackendUri()}/admin/workers`, {
    method: "GET",
    headers,
    credentials: "include",
  });

  if (!response.ok) return null;
  return response.json();
}
