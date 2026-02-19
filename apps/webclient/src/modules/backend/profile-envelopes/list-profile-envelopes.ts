import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileEnvelope } from "../types";

export async function listProfileEnvelopes(
  locale: string,
  slug: string,
  status?: string,
): Promise<ProfileEnvelope[] | null> {
  const token = getAuthToken();
  if (token === null) {
    return null;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  let url = `${getBackendUri()}/${locale}/profiles/${slug}/_envelopes`;
  if (status !== undefined && status !== "") {
    url += `?status=${encodeURIComponent(status)}`;
  }

  const response = await fetch(url, {
    method: "GET",
    headers,
    credentials: "include",
  });

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
