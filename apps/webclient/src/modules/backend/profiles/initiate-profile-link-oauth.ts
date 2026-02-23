import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";

export type InitiateOAuthResponse = {
  auth_url: string;
};

export async function initiateProfileLinkOAuth(
  locale: string,
  slug: string,
  provider: string,
  scopeUpgrade?: string,
): Promise<InitiateOAuthResponse | null> {
  const token = getAuthToken();
  if (token === null) return null;

  let url =
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/connect/${provider}`;
  if (scopeUpgrade !== undefined) {
    url += `?scope_upgrade=${scopeUpgrade}`;
  }

  const response = await fetch(
    url,
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
