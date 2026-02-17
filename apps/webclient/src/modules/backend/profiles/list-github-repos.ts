import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { GitHubRepo } from "../types";

export async function listGitHubRepos(
  locale: string,
  slug: string,
): Promise<GitHubRepo[] | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_resources/github/repos`,
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
