// Send an envelope from a profile to a target profile

import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileEnvelope } from "../types";

export interface SendProfileEnvelopeParams {
  locale: string;
  senderSlug: string;
  targetProfileId: string;
  kind: string;
  conversationTitle?: string;
  message: string;
  inviteCode?: string;
  properties?: Record<string, unknown>;
}

export async function sendProfileEnvelope(
  params: SendProfileEnvelopeParams,
): Promise<ProfileEnvelope | null> {
  const token = getAuthToken();

  if (token === null) {
    return null;
  }

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const url = `${getBackendUri()}/${params.locale}/profiles/${params.senderSlug}/_envelopes`;

  const body: Record<string, unknown> = {
    kind: params.kind,
    target_profile_id: params.targetProfileId,
    message: params.message,
  };

  if (params.conversationTitle !== undefined && params.conversationTitle !== "") {
    body.conversation_title = params.conversationTitle;
  }

  if (params.inviteCode !== undefined && params.inviteCode !== "") {
    body.invite_code = params.inviteCode;
  }

  if (params.properties !== undefined) {
    body.properties = params.properties;
  }

  const response = await fetch(url, {
    method: "POST",
    headers,
    credentials: "include",
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    const errorBody = await response.json().catch(() => null);
    const message = errorBody !== null && typeof errorBody === "object" && "error" in errorBody
      ? String(errorBody.error)
      : "Failed to send envelope";

    throw new Error(message);
  }

  const result = await response.json();
  return result.data;
}
