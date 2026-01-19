import { getBackendUri } from "@/config";

export interface SessionCheckResponse {
  authenticated: boolean;
  token?: string;
  expires_at?: number;
  user?: {
    id: string;
    name: string;
    email?: string;
    github_handle?: string;
    individual_profile_id?: string;
  };
  selected_profile?: {
    id: string;
    slug: string;
    kind: string;
    title: string;
    description?: string;
    profile_picture_uri?: string;
  };
}

/**
 * Check session via cookie for cross-domain SSO.
 * This endpoint reads from HttpOnly cookie instead of Authorization header.
 */
export async function checkSessionViaCookie(
  locale: string,
): Promise<SessionCheckResponse> {
  const backendUri = getBackendUri();

  try {
    const response = await fetch(`${backendUri}/${locale}/auth/session-check`, {
      method: "GET",
      credentials: "include", // CRITICAL: sends cookies cross-domain
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      return { authenticated: false };
    }

    const result = await response.json();

    return result.data as SessionCheckResponse;
  } catch {
    return { authenticated: false };
  }
}
