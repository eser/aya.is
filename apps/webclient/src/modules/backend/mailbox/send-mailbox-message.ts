import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { MailboxEnvelope } from "../types";

export interface SendMailboxMessageParams {
  locale: string;
  senderProfileSlug: string;
  targetProfileSlug: string;
  kind?: string;
  conversationTitle?: string;
  message: string;
  replyToId?: string;
}

export async function sendMailboxMessage(
  params: SendMailboxMessageParams,
): Promise<MailboxEnvelope | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const body: Record<string, unknown> = {
    sender_profile_slug: params.senderProfileSlug,
    target_profile_slug: params.targetProfileSlug,
    kind: params.kind ?? "message",
    message: params.message,
  };

  if (params.conversationTitle !== undefined && params.conversationTitle !== "") {
    body.conversation_title = params.conversationTitle;
  }

  if (params.replyToId !== undefined) {
    body.reply_to_id = params.replyToId;
  }

  const response = await fetch(
    `${getBackendUri()}/${params.locale}/mailbox/messages`,
    {
      method: "POST",
      headers,
      credentials: "include",
      body: JSON.stringify(body),
    },
  );

  if (!response.ok) {
    const errorBody = await response.json().catch(() => null);
    const message = errorBody !== null && typeof errorBody === "object" && "error" in errorBody
      ? String(errorBody.error)
      : "Failed to send message";

    throw new Error(message);
  }

  const result = await response.json();
  return result.data;
}
