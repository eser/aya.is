import { getBackendUri } from "@/config.ts";
import type { PreferencesResponse, SessionPreferences } from "./types.ts";

/**
 * Update session preferences on the server
 * Only updates the provided fields
 */
export async function updateSessionPreferences(
  locale: string,
  preferences: Partial<SessionPreferences>,
): Promise<SessionPreferences | null> {
  const backendUri = getBackendUri();
  const response = await fetch(
    `${backendUri}/${locale}/sessions/_current`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify(preferences),
    },
  );

  if (!response.ok) {
    console.error(
      "[Session] Failed to update preferences:",
      response.status,
    );
    return null;
  }

  const result = (await response.json()) as {
    data: PreferencesResponse | null;
    error: string | null;
  };

  if (result.error !== null || result.data === null) {
    console.error("[Session] Update preferences error:", result.error);
    return null;
  }

  return result.data.preferences;
}
