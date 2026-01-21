import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";

export type InitiateOAuthResponse = {
  auth_url: string;
};

export async function initiateProfileLinkOAuth(
  locale: string,
  slug: string,
  provider: string,
): Promise<InitiateOAuthResponse | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/connect/${provider}`,
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
