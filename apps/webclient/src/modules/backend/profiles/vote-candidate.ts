import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { CandidateVote } from "../types";

export async function voteCandidate(
  locale: string,
  slug: string,
  candidateId: string,
  score: number,
  comment?: string | null,
): Promise<CandidateVote | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_candidates/${candidateId}/votes`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({
        score,
        comment: comment ?? null,
      }),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data ?? null;
}
