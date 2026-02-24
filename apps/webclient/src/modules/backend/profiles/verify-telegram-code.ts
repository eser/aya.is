import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";

export type VerifyTelegramCodeResponse = {
  profile_id: string;
  profile_slug: string;
  telegram_user_id: number;
  telegram_username: string;
};

export async function verifyTelegramCode(
  locale: string,
  slug: string,
  code: string,
  visibility?: string,
): Promise<VerifyTelegramCodeResponse | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/telegram/verify-code`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({ code, visibility }),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
