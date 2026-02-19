import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function acceptProfileEnvelope(
  locale: string,
  slug: string,
  envelopeId: string,
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
    `${getBackendUri()}/${locale}/profiles/${slug}/_envelopes/${envelopeId}/accept`,
    {
      method: "POST",
      headers,
      credentials: "include",
    },
  );

  return response.ok;
}
