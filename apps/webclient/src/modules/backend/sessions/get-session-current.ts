import { getBackendUri } from "@/config.ts";
import type { SessionCurrentResponse } from "./types.ts";

/**
 * Get current session state via cookie-based check.
 * Returns auth state, token, and preferences in a single response.
 * This is the primary session initialization call on app mount.
 */
export async function getSessionCurrent(
  locale: string,
): Promise<SessionCurrentResponse> {
  const backendUri = getBackendUri();

  try {
    const response = await fetch(`${backendUri}/${locale}/sessions/_current`, {
      method: "GET",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      return { authenticated: false };
    }

    const result = await response.json();

    return result.data as SessionCurrentResponse;
  } catch {
    return { authenticated: false };
  }
}
