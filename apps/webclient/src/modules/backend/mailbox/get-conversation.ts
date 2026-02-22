import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ConversationDetail } from "../types";

export async function getConversation(
  locale: string,
  conversationId: string,
): Promise<ConversationDetail | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(
    `${getBackendUri()}/${locale}/mailbox/conversations/${conversationId}`,
    {
      method: "GET",
      headers,
      credentials: "include",
    },
  );

  if (!response.ok) {
    return null;
  }

  const result = await response.json();
  return result.data;
}
