import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function rejectMailboxMessage(
  locale: string,
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
    `${getBackendUri()}/${locale}/mailbox/messages/${envelopeId}/reject`,
    {
      method: "POST",
      headers,
      credentials: "include",
    },
  );

  return response.ok;
}
