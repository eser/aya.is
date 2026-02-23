import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { GitHubRepo } from "../types";

export type ListGitHubReposResult = {
  repos: GitHubRepo[] | null;
  needsReauth: boolean;
};

export async function listGitHubRepos(
  locale: string,
  slug: string,
): Promise<ListGitHubReposResult> {
  const token = getAuthToken();
  if (token === null) return { repos: null, needsReauth: false };

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

  if (!response.ok) return { repos: null, needsReauth: false };
  const result = await response.json();

  if (result.needs_reauth === true) {
    return { repos: null, needsReauth: true };
  }

  return { repos: result.data, needsReauth: false };
}
