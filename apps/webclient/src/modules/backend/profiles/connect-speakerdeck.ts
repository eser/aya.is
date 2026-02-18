import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";

export type ConnectSpeakerDeckResponse = {
  status: string;
};

export async function connectSpeakerDeck(
  locale: string,
  slug: string,
  url: string,
): Promise<ConnectSpeakerDeckResponse | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/connect/speakerdeck`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({ url }),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
