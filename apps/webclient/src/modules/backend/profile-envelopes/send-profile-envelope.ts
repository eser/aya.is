// Send an envelope from a profile to a target profile

import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileEnvelope } from "../types";

export interface SendProfileEnvelopeParams {
  locale: string;
  senderSlug: string;
  targetProfileId: string;
  kind: string;
  title: string;
  description?: string;
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
    title: params.title,
  };

  if (params.description !== undefined && params.description !== "") {
    body.description = params.description;
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
    const text = await response.text();
    throw new Error(text || "Failed to send envelope");
  }

  const result = await response.json();
  return result.data;
}
