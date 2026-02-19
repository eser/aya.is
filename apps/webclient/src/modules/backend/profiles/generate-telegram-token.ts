import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";

export type GenerateTelegramTokenResponse = {
  token: string;
  deep_link: string;
};

export async function generateTelegramToken(
  locale: string,
  slug: string,
): Promise<GenerateTelegramTokenResponse | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/telegram/generate-token`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
