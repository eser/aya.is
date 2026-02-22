import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { Conversation } from "../types";

export interface ListConversationsResult {
  conversations: Conversation[];
  viewerHasTelegram: boolean;
}

export async function listConversations(
  locale: string,
  options?: { archived?: boolean },
): Promise<ListConversationsResult | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  let url = `${getBackendUri()}/${locale}/mailbox/conversations`;
  if (options?.archived === true) {
    url += "?archived=true";
  }

  const response = await fetch(url, {
    method: "GET",
    headers,
    credentials: "include",
  });

  if (!response.ok) {
    return null;
  }

  const result = await response.json();
  return {
    conversations: result.data ?? [],
    viewerHasTelegram: result.viewer_has_telegram === true,
  };
}
