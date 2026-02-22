import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function unarchiveConversation(
  locale: string,
  conversationId: string,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) {
    return false;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/${locale}/mailbox/conversations/${conversationId}/unarchive`,
    {
      method: "POST",
      headers,
      credentials: "include",
    },
  );

  return response.ok;
}
