import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function addReaction(
  locale: string,
  envelopeId: string,
  emoji: string,
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
    `${getBackendUri()}/${locale}/mailbox/messages/${envelopeId}/reactions`,
    {
      method: "POST",
      headers,
      credentials: "include",
      body: JSON.stringify({ emoji }),
    },
  );

  return response.ok;
}
