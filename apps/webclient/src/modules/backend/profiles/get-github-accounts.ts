import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";

export type GitHubAccount = {
  id: string;
  login: string;
  name: string;
  avatar_url: string;
  html_url: string;
  type: "User" | "Organization";
  description?: string;
};

export type GitHubAccountsResponse = {
  accounts: GitHubAccount[];
  profile_kind: "individual" | "organization" | "product";
};

export async function getGitHubAccounts(
  locale: string,
  slug: string,
  pendingId: string,
): Promise<GitHubAccountsResponse | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/github/accounts?pending_id=${pendingId}`,
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
