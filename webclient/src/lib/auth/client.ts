import { createAuthClient } from "better-auth/react";

/**
 * Get the base URL for auth client
 */
function getBaseURL(): string {
  if (typeof window !== "undefined") {
    return `${window.location.origin}/api/auth`;
  }
  // Server-side: use environment variable or default localhost
  const host = process.env.PUBLIC_HOST || "http://localhost:3000";
  return `${host}/api/auth`;
}

// Lazy-initialized auth client to avoid SSR issues
let _authClient: ReturnType<typeof createAuthClient> | null = null;

function getAuthClient() {
  if (!_authClient) {
    _authClient = createAuthClient({
      baseURL: getBaseURL(),
    });
  }
  return _authClient;
}

/**
 * BetterAuth client - lazily initialized
 */
export const authClient = {
  get instance() {
    return getAuthClient();
  },
};

// Export auth methods for use in components
// These are re-exported from the lazy client
export function useSession() {
  return getAuthClient().useSession();
}

export function signIn(...args: Parameters<ReturnType<typeof createAuthClient>["signIn"]>) {
  return getAuthClient().signIn(...args);
}

export function signOut(...args: Parameters<ReturnType<typeof createAuthClient>["signOut"]>) {
  return getAuthClient().signOut(...args);
}
