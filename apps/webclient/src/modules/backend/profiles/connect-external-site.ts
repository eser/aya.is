import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";

export type ConnectExternalSiteResponse = {
  status: string;
};

export async function connectExternalSite(
  locale: string,
  slug: string,
  system: string,
  url: string,
  siteUrl: string,
): Promise<ConnectExternalSiteResponse | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/connect/external-site`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({ system, url, site_url: siteUrl }),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
