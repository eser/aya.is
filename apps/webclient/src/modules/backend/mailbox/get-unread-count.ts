import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function getUnreadCount(
  locale: string,
): Promise<number> {
  const token = getAuthToken();
  if (token === null) {
    return 0;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/${locale}/mailbox/unread-count`,
    {
      method: "GET",
      headers,
      credentials: "include",
    },
  );

  if (!response.ok) {
    return 0;
  }

  const result = await response.json();
  return result.data ?? 0;
}
