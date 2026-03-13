import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { CandidateStatus } from "../types";

export async function updateCandidateStatus(
  locale: string,
  slug: string,
  candidateId: string,
  status: CandidateStatus,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_candidates/${candidateId}/status`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({ status }),
    },
  );

  return response.ok;
}
