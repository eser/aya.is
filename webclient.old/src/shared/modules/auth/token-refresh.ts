import { getBackendUri as _getBackendUri } from "@/shared/config.ts";
import { refreshTokenRequest } from "@/shared/modules/backend/auth/refresh-token.ts";

interface TokenRefreshResponse {
  token: string;
  expiresAt: number;
}

let refreshPromise: Promise<string | null> | null = null;

export async function refreshToken(): Promise<string | null> {
  // If a refresh is already in progress, return the same promise
  if (refreshPromise) {
    return refreshPromise;
  }

  refreshPromise = performTokenRefresh();

  try {
    const result = await refreshPromise;
    return result;
  } finally {
    refreshPromise = null;
  }
}

async function performTokenRefresh(): Promise<string | null> {
  try {
    const currentToken = globalThis.localStorage === undefined ? null : localStorage.getItem("auth_token");

    if (!currentToken) {
      return null;
    }

    // Get locale from URL or default to 'tr'
    const currentPath = globalThis.location?.pathname || "/tr";
    const locale = currentPath.split("/")[1] || "tr";

    const data = await refreshTokenRequest(locale, currentToken);

    if (data?.token) {
      // Store new token
      if (globalThis.localStorage !== undefined) {
        localStorage.setItem("auth_token", data.token);
        localStorage.setItem("auth_token_expires_at", data.expiresAt.toString());
      }
      return data.token;
    }

    return null;
  } catch (error) {
    console.error("Token refresh failed:", error);
    // Token is invalid or expired
    clearAuthData();
    return null;
  }
}

export function clearAuthData() {
  if (globalThis.localStorage !== undefined) {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("auth_token_expires_at");
    localStorage.removeItem("auth_session");
  }
}

export function isTokenExpiringSoon(): boolean {
  if (globalThis.localStorage === undefined) {
    return false;
  }

  const expiresAt = localStorage.getItem("auth_token_expires_at");
  if (!expiresAt) {
    return true; // If no expiration is stored, consider it expiring
  }

  const expirationTime = parseInt(expiresAt);
  const now = Date.now();
  const fiveMinutes = 5 * 60 * 1000;

  // Token expires in less than 5 minutes
  return (expirationTime - now) < fiveMinutes;
}

export function getTokenExpirationTime(): number | null {
  if (globalThis.localStorage === undefined) {
    return null;
  }

  const expiresAt = localStorage.getItem("auth_token_expires_at");
  if (!expiresAt) {
    return null;
  }

  return parseInt(expiresAt);
}

export function getToken(): string | null {
  if (globalThis.localStorage === undefined) {
    return null;
  }

  return localStorage.getItem("auth_token");
}
