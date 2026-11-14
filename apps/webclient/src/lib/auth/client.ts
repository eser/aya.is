// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { createAuthClient } from "better-auth/react";

/**
 * Get the base URL for auth client
 */
function getBaseURL(): string {
  if (typeof globalThis.location !== "undefined") {
    return `${globalThis.location.origin}/api/auth`;
  }
  // Server-side: use Vite env variable or default localhost
  const host = import.meta.env.VITE_HOST ?? "http://localhost:3000";
  return `${host}/api/auth`;
}

// Lazy-initialized auth client to avoid SSR issues
let _authClient: ReturnType<typeof createAuthClient> | null = null;

function getAuthClient() {
  if (_authClient === null) {
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

// better-auth's signIn is an object of methods (not callable), so each method is forwarded lazily
export const signIn = {
  email: (...args: Parameters<ReturnType<typeof createAuthClient>["signIn"]["email"]>) =>
    getAuthClient().signIn.email(...args),
  social: (...args: Parameters<ReturnType<typeof createAuthClient>["signIn"]["social"]>) =>
    getAuthClient().signIn.social(...args),
};

export function signOut(...args: Parameters<ReturnType<typeof createAuthClient>["signOut"]>) {
  return getAuthClient().signOut(...args);
}
