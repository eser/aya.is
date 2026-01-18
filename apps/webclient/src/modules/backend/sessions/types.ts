// Session Preference Types

/**
 * Session preferences that can be stored server-side
 */
export interface SessionPreferences {
  theme?: "light" | "dark" | "system";
  locale?: string;
  timezone?: string;
}

/**
 * Response from session creation
 */
export interface CreateSessionResponse {
  session_id: string;
  preferences: SessionPreferences;
}

/**
 * Response from preferences endpoints
 */
export interface PreferencesResponse {
  preferences: SessionPreferences;
}
