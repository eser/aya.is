import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileEnvelope } from "../types";

export interface MailboxEnvelope extends ProfileEnvelope {
  owning_profile_slug: string;
  owning_profile_title: string;
  owning_profile_kind: string;
}

export async function listMailboxEnvelopes(
  locale: string,
  status?: string,
): Promise<MailboxEnvelope[] | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  let url = `${getBackendUri()}/${locale}/mailbox`;
  if (status !== undefined && status !== "") {
    url += `?status=${encodeURIComponent(status)}`;
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
  return result.data;
}
