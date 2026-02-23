import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";
import type { ProfileResource } from "../types.ts";

export async function verifyTelegramRegisterCode(
  locale: string,
  slug: string,
  code: string,
): Promise<ProfileResource | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_resources/telegram/verify-code`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({ code }),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
