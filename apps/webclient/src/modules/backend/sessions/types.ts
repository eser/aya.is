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

/**
 * Response from GET /sessions/_current consolidated endpoint
 */
export interface AccessibleProfile {
  id: string;
  slug: string;
  kind: string;
  title: string;
  profile_picture_uri?: string | null;
  membership_kind: string;
}

export interface SessionCurrentResponse {
  authenticated: boolean;
  token?: string;
  expires_at?: number;
  user?: {
    id: string;
    kind: string;
    name: string;
    email?: string;
    github_handle?: string;
    individual_profile_id?: string;
  };
  selected_profile?: {
    id: string;
    slug: string;
    kind: string;
    title?: string;
    description?: string;
    profile_picture_uri?: string;
    points?: number;
  };
  accessible_profiles?: AccessibleProfile[];
  total_pending_envelopes?: number;
  preferences?: SessionPreferences;
}
