import { getBackendUri } from "@/config.ts";
import type { PreferencesResponse, SessionPreferences } from "./types.ts";

/**
 * Get session preferences from the server
 * Returns null if no session or error
 */
export async function getSessionPreferences(
  locale: string,
): Promise<SessionPreferences | null> {
  const backendUri = getBackendUri();
  const response = await fetch(
    `${backendUri}/${locale}/sessions/current/preferences`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    },
  );

  if (!response.ok) {
    // 401 is expected if no session exists
    if (response.status !== 401) {
      console.error(
        "[Session] Failed to get preferences:",
        response.status,
      );
    }
    return null;
  }

  const result = (await response.json()) as {
    data: PreferencesResponse | null;
    error: string | null;
  };

  if (result.error !== null || result.data === null) {
    return null;
  }

  return result.data.preferences;
}
