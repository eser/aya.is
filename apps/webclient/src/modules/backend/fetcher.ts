import { getBackendUri } from "@/config.ts";
import type { Result } from "./types.ts";

// Time constants for token management
const TOKEN_REFRESH_THRESHOLD_MS = 5 * 60 * 1000; // 5 minutes before expiry
const DEFAULT_TOKEN_EXPIRY_MS = 24 * 60 * 60 * 1000; // 24 hours fallback

// Storage keys
const AUTH_TOKEN_KEY = "auth_token";
const AUTH_TOKEN_EXPIRES_KEY = "auth_token_expires_at";
const AUTH_SESSION_KEY = "auth_session";

// HTTP status codes
const HTTP_STATUS_UNAUTHORIZED = 401;
const HTTP_STATUS_NOT_FOUND = 404;
const HTTP_STATUS_METHOD_NOT_ALLOWED = 405;
const HTTP_STATUS_SERVER_ERROR_THRESHOLD = 500;

/**
 * Check if localStorage is available (browser environment)
 * Uses globalThis for cross-runtime compatibility (Browser, Node, Deno)
 */
function isLocalStorageAvailable(): boolean {
  return typeof globalThis !== "undefined" &&
    typeof globalThis.localStorage !== "undefined";
}

/**
 * Retrieve auth token from localStorage
 * Returns null when running server-side or if no token exists
 */
export function getAuthToken(): string | null {
  if (!isLocalStorageAvailable()) {
    return null;
  }
  return globalThis.localStorage.getItem(AUTH_TOKEN_KEY);
}

/**
 * Store auth token and expiration time in localStorage
 * Silently skips when running server-side
 */
export function setAuthToken(token: string, expiresAt: number): void {
  if (!isLocalStorageAvailable()) {
    return;
  }
  globalThis.localStorage.setItem(AUTH_TOKEN_KEY, token);
  globalThis.localStorage.setItem(AUTH_TOKEN_EXPIRES_KEY, expiresAt.toString());
}

/**
 * Clear all authentication data from localStorage
 * Used on logout or when token refresh fails
 */
export function clearAuthData(): void {
  if (!isLocalStorageAvailable()) {
    return;
  }
  globalThis.localStorage.removeItem(AUTH_TOKEN_KEY);
  globalThis.localStorage.removeItem(AUTH_TOKEN_EXPIRES_KEY);
  globalThis.localStorage.removeItem(AUTH_SESSION_KEY);
}

/**
 * Check if the current token will expire within the refresh threshold
 * Triggers proactive refresh to avoid mid-request token expiration
 */
export function isTokenExpiringSoon(): boolean {
  if (!isLocalStorageAvailable()) {
    return false;
  }

  const expiresAt = globalThis.localStorage.getItem(AUTH_TOKEN_EXPIRES_KEY);
  if (expiresAt === null) {
    return false;
  }

  const expirationTime = Number(expiresAt);
  const now = Date.now();

  return expirationTime - now < TOKEN_REFRESH_THRESHOLD_MS;
}

// Singleton promise to deduplicate concurrent refresh requests
let refreshPromise: Promise<string | null> | null = null;

/**
 * Refresh the auth token via API call
 * Deduplicates concurrent refresh calls to prevent race conditions
 */
export async function refreshTokenRequest(
  locale: string,
): Promise<{ token: string; expiresAt: number } | null> {
  // Return existing refresh promise if one is in progress
  if (refreshPromise !== null) {
    const token = await refreshPromise;
    if (token !== null) {
      const expiresAt = globalThis.localStorage.getItem(AUTH_TOKEN_EXPIRES_KEY);
      const parsedExpiresAt = expiresAt !== null ? Number(expiresAt) : Date.now() + DEFAULT_TOKEN_EXPIRY_MS;
      return { token, expiresAt: parsedExpiresAt };
    }
    return null;
  }

  const currentToken = getAuthToken();
  if (currentToken === null) {
    return null;
  }

  refreshPromise = (async () => {
    try {
      const backendUri = getBackendUri();
      const response = await fetch(`${backendUri}/${locale}/auth/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${currentToken}`,
        },
        credentials: "include",
      });

      if (!response.ok) {
        clearAuthData();
        return null;
      }

      const data = await response.json();
      if (data.data?.token !== undefined) {
        const newToken = data.data.token as string;
        const expiresAt = (data.data.expires_at as number | undefined) ??
          Date.now() + DEFAULT_TOKEN_EXPIRY_MS;
        setAuthToken(newToken, expiresAt);
        return newToken;
      }
      return null;
    } catch {
      return null;
    } finally {
      refreshPromise = null;
    }
  })();

  const token = await refreshPromise;
  if (token !== null) {
    const expiresAt = globalThis.localStorage.getItem(AUTH_TOKEN_EXPIRES_KEY);
    const parsedExpiresAt = expiresAt !== null ? Number(expiresAt) : Date.now() + DEFAULT_TOKEN_EXPIRY_MS;
    return { token, expiresAt: parsedExpiresAt };
  }
  return null;
}

/**
 * Build fetch options with proper credentials handling
 * Credentials are only included when authenticated due to CORS wildcard limitation
 */
function buildFetchOptions(
  requestInit: RequestInit,
  headers: HeadersInit,
  includeCredentials: boolean,
): RequestInit {
  const options: RequestInit = {
    ...requestInit,
    headers: {
      ...headers,
      ...(requestInit.headers ?? {}),
    },
  };

  if (includeCredentials) {
    options.credentials = "include";
  }

  return options;
}

// In-flight GET request cache to deduplicate concurrent fetches to the same endpoint.
// Prevents multiple components or loaders from triggering duplicate network requests.
const inflightGetRequests = new Map<string, Promise<unknown>>();

/**
 * Generic API fetcher with automatic token refresh and error handling
 *
 * CORS Note: Credentials are only included when auth token is present
 * because the backend uses wildcard CORS origin which is incompatible
 * with credentials (browsers reject this combination per CORS spec)
 */
export function fetcher<T>(
  locale: string,
  relativePath: string,
  requestInit: RequestInit = {},
): Promise<T | null> {
  const method = (requestInit.method ?? "GET").toUpperCase();
  const requestPath = `${locale}/${relativePath}`;

  // Deduplicate concurrent GET requests to the same path
  if (method === "GET") {
    const existing = inflightGetRequests.get(requestPath);
    if (existing !== undefined) {
      return existing as Promise<T | null>;
    }

    const promise = fetcherInternal<T>(locale, relativePath, requestInit).finally(() => {
      inflightGetRequests.delete(requestPath);
    });

    inflightGetRequests.set(requestPath, promise);
    return promise;
  }

  return fetcherInternal<T>(locale, relativePath, requestInit);
}

async function fetcherInternal<T>(
  locale: string,
  relativePath: string,
  requestInit: RequestInit,
): Promise<T | null> {
  const backendUri = getBackendUri();
  const targetUrl = `${backendUri}/${locale}${relativePath}`;

  // Get auth token from localStorage (only available on client)
  const authToken = getAuthToken();

  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  if (authToken !== null) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  // During SSR, forward the incoming request's cookies to backend API calls.
  // This allows session-cookie-based auth to work during server-side rendering.
  if (!isLocalStorageAvailable()) {
    try {
      const { requestContextBinder } = await import(
        "@/server/request-context-binder"
      );
      const ctx = requestContextBinder.getStore();
      if (ctx !== undefined && ctx.cookieHeader !== undefined) {
        (headers as Record<string, string>)["Cookie"] = ctx.cookieHeader;
      }
    } catch {
      // Client bundle or no request context available — skip
    }
  }

  const isAuthenticated = authToken !== null;
  const fetchOptions = buildFetchOptions(requestInit, headers, isAuthenticated);
  const response = await fetch(targetUrl, fetchOptions);

  // Handle authentication errors with token refresh retry
  if (
    response.status === HTTP_STATUS_UNAUTHORIZED ||
    response.status === HTTP_STATUS_METHOD_NOT_ALLOWED
  ) {
    if (authToken !== null) {
      const tokenResponse = await refreshTokenRequest(locale);

      if (tokenResponse !== null) {
        headers["Authorization"] = `Bearer ${tokenResponse.token}`;
        const retryOptions = buildFetchOptions(requestInit, headers, true);
        const retryResponse = await fetch(targetUrl, retryOptions);

        if (retryResponse.status === HTTP_STATUS_NOT_FOUND) {
          return null;
        }

        if (
          retryResponse.status === HTTP_STATUS_UNAUTHORIZED ||
          retryResponse.status === HTTP_STATUS_METHOD_NOT_ALLOWED
        ) {
          // Still unauthorized after refresh - token is invalid
          return null;
        }

        if (retryResponse.status >= HTTP_STATUS_SERVER_ERROR_THRESHOLD) {
          throw new Error("Internal server error", {
            cause: { status: retryResponse.status, url: targetUrl },
          });
        }

        const retryResult = (await retryResponse.json()) as Result<T>;

        if (retryResult.error !== undefined && retryResult.error !== null) {
          throw new Error(retryResult.error, {
            cause: { url: targetUrl, response: retryResult },
          });
        }

        return retryResult.data;
      }
    }

    // No token or refresh failed - return null for unauthenticated state
    return null;
  }

  if (response.status === HTTP_STATUS_NOT_FOUND) {
    return null;
  }

  if (response.status >= HTTP_STATUS_SERVER_ERROR_THRESHOLD) {
    throw new Error("Internal server error", {
      cause: { status: response.status, url: targetUrl },
    });
  }

  const result = (await response.json()) as Result<T>;

  if (result.error !== undefined && result.error !== null) {
    throw new Error(result.error, {
      cause: { url: targetUrl, response: result },
    });
  }

  return result.data;
}

/**
 * Fetcher for file uploads using multipart/form-data
 * Does not set Content-Type header (browser sets it with boundary)
 */
export async function uploadFetcher<T>(
  locale: string,
  relativePath: string,
  formData: FormData,
): Promise<T | null> {
  const backendUri = getBackendUri();
  const targetUrl = `${backendUri}/${locale}${relativePath}`;

  const authToken = getAuthToken();

  const headers: HeadersInit = {};

  if (authToken !== null) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  const isAuthenticated = authToken !== null;
  const uploadOptions: RequestInit = {
    method: "POST",
    headers,
    body: formData,
  };

  if (isAuthenticated) {
    uploadOptions.credentials = "include";
  }

  const response = await fetch(targetUrl, uploadOptions);

  if (
    response.status === HTTP_STATUS_UNAUTHORIZED ||
    response.status === HTTP_STATUS_METHOD_NOT_ALLOWED
  ) {
    return null;
  }

  if (response.status === HTTP_STATUS_NOT_FOUND) {
    return null;
  }

  if (response.status >= HTTP_STATUS_SERVER_ERROR_THRESHOLD) {
    throw new Error("Internal server error", {
      cause: { status: response.status, url: targetUrl },
    });
  }

  const result = (await response.json()) as Result<T>;

  if (result.error !== undefined && result.error !== null) {
    throw new Error(result.error, {
      cause: { url: targetUrl, response: result },
    });
  }

  return result.data;
}
