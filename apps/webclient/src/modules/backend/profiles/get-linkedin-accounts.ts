import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";

export type LinkedInAccount = {
  id: string;
  name: string;
  vanity_name?: string;
  logo_url?: string;
  uri: string;
  type: "Personal" | "Organization";
  description?: string;
};

export type LinkedInAccountsResponse = {
  accounts: LinkedInAccount[];
  profile_kind: "individual" | "organization" | "product";
};

export async function getLinkedInAccounts(
  locale: string,
  slug: string,
  pendingId: string,
): Promise<LinkedInAccountsResponse | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/linkedin/accounts?pending_id=${pendingId}`,
    {
      method: "GET",
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
