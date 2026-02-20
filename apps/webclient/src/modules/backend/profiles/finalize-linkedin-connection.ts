import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";
import type { LinkedInAccount } from "./get-linkedin-accounts.ts";

export async function finalizeLinkedInConnection(
  locale: string,
  slug: string,
  account: LinkedInAccount,
  pendingId: string,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/linkedin/finalize`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({
        pending_id: pendingId,
        account_id: account.id,
        name: account.name,
        vanity_name: account.vanity_name ?? "",
        uri: account.uri,
        type: account.type,
      }),
    },
  );

  return response.ok;
}
