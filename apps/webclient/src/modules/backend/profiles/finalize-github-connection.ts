import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";
import type { GitHubAccount } from "./get-github-accounts.ts";

export type FinalizeGitHubConnectionInput = {
  pending_id: string;
  account_id: string;
  login: string;
  name: string;
  html_url: string;
  type: string;
};

export async function finalizeGitHubConnection(
  locale: string,
  slug: string,
  account: GitHubAccount,
  pendingId: string,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/github/finalize`,
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
        login: account.login,
        name: account.name,
        html_url: account.html_url,
        type: account.type,
      }),
    },
  );

  return response.ok;
}
