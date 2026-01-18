import { getBackendUri } from "@/config.ts";
import type { POWChallenge, POWChallengeDisabled } from "./types.ts";

/**
 * Get a PoW challenge from the server
 * Returns null if PoW is disabled or an error occurred
 */
export async function getPOWChallenge(
  locale: string,
): Promise<POWChallenge | null> {
  const backendUri = getBackendUri();
  const response = await fetch(
    `${backendUri}/${locale}/protection/pow-challenges`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    },
  );

  if (!response.ok) {
    console.error("[Protection] Failed to get PoW challenge:", response.status);
    return null;
  }

  const result = (await response.json()) as {
    data: POWChallenge | POWChallengeDisabled | null;
    error: string | null;
  };

  if (result.error !== null) {
    console.error("[Protection] PoW challenge error:", result.error);
    return null;
  }

  if (result.data === null) {
    return null;
  }

  // Check if PoW is disabled
  if ("enabled" in result.data && result.data.enabled === false) {
    return null;
  }

  return result.data as POWChallenge;
}
