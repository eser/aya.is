import { getBackendUri } from "@/config.ts";
import { getPOWChallenge, isPOWSolverSupported, solvePOW } from "../protection";
import type { CreateSessionResponse, SessionPreferences } from "./types.ts";

/**
 * Create a new session with optional preferences
 * Handles PoW challenge if required by the server
 */
export async function createSession(
  locale: string,
  preferences: SessionPreferences = {},
): Promise<CreateSessionResponse | null> {
  const backendUri = getBackendUri();

  // Get PoW challenge if required
  let powChallengeId: string | null = null;
  let nonce: string | null = null;

  if (isPOWSolverSupported()) {
    const challenge = await getPOWChallenge(locale);

    if (challenge !== null) {
      console.log(
        `[Session] Solving PoW challenge (difficulty: ${challenge.difficulty})...`,
      );
      try {
        nonce = await solvePOW(challenge.prefix, challenge.difficulty);
        powChallengeId = challenge.pow_challenge_id;
      } catch (error) {
        console.error("[Session] Failed to solve PoW challenge:", error);
        // Continue without PoW - server might accept if PoW is optional
      }
    }
  }

  // Build request body
  const body: Record<string, unknown> = {
    preferences,
  };

  if (powChallengeId !== null && nonce !== null) {
    body.pow_challenge_id = powChallengeId;
    body.nonce = nonce;
  }

  const response = await fetch(`${backendUri}/${locale}/sessions`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    console.error("[Session] Failed to create session:", response.status);
    return null;
  }

  const result = (await response.json()) as {
    data: CreateSessionResponse | null;
    error: string | null;
  };

  if (result.error !== null) {
    console.error("[Session] Create session error:", result.error);
    return null;
  }

  return result.data;
}
