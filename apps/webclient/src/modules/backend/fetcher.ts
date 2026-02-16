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
 * Check if we are running in a browser environment.
 * Deno has globalThis.localStorage available on the server, so checking
 * localStorage alone is not enough. `document` is only present in browsers.
 */
function isBrowserEnvironment(): boolean {
  return typeof globalThis !== "undefined" &&
    typeof globalThis.document !== "undefined";
}

/**
 * Retrieve auth token from localStorage
 * Returns null when running server-side or if no token exists
 */
export function getAuthToken(): string | null {
  if (!isBrowserEnvironment()) {
    return null;
  }
  return globalThis.localStorage.getItem(AUTH_TOKEN_KEY);
}

/**
 * Store auth token and expiration time in localStorage
 * Silently skips when running server-side
 */
export function setAuthToken(token: string, expiresAt: number): void {
  if (!isBrowserEnvironment()) {
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
  if (!isBrowserEnvironment()) {
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
  if (!isBrowserEnvironment()) {
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
 * Build fetch options with credentials always included.
 * The session cookie (aya_session) is the primary auth mechanism,
 * so credentials must always be sent for cross-origin requests.
 */
function buildFetchOptions(
  requestInit: RequestInit,
  headers: HeadersInit,
): RequestInit {
  return {
    ...requestInit,
    credentials: "include",
    headers: {
      ...headers,
      ...(requestInit.headers ?? {}),
    },
  };
}

// Shared in-flight GET request cache for anonymous SSR and client-side requests.
// Safe to share because anonymous requests return identical public data.
const sharedInflightCache = new Map<string, Promise<unknown>>();

const SESSION_COOKIE_NAME = "aya_session";

/**
 * Get the appropriate dedup cache for the current context.
 *
 * - Client-side: shared module-level cache (single user, always safe)
 * - SSR without session cookie: shared module-level cache (anonymous, same data for all)
 * - SSR with session cookie: per-request scoped cache (response may be personalized)
 */
async function getInflightCache(): Promise<Map<string, Promise<unknown>>> {
  if (isBrowserEnvironment()) {
    return sharedInflightCache;
  }

  // SSR: check if the request has a session cookie
  try {
    const { requestContextBinder } = await import(
      "@/server/request-context-binder"
    );
    const ctx = requestContextBinder.getStore();

    if (
      ctx !== undefined &&
      ctx.cookieHeader !== undefined &&
      ctx.cookieHeader.includes(SESSION_COOKIE_NAME)
    ) {
      // Authenticated SSR — use per-request cache to prevent cross-user data leaks
      if (ctx.inflightGetRequests === undefined) {
        ctx.inflightGetRequests = new Map<string, Promise<unknown>>();
      }
      return ctx.inflightGetRequests;
    }
  } catch {
    // No request context available — fall through to shared cache
  }

  // Anonymous SSR — safe to share across requests
  return sharedInflightCache;
}

/**
 * Generic API fetcher with automatic token refresh and error handling
 */
export function fetcher<T>(
  locale: string,
  relativePath: string,
  requestInit: RequestInit = {},
): Promise<T | null> {
  const method = (requestInit.method ?? "GET").toUpperCase();

  if (method === "GET") {
    return fetcherWithDedup<T>(locale, relativePath, requestInit);
  }

  return fetcherInternal<T>(locale, relativePath, requestInit);
}

async function fetcherWithDedup<T>(
  locale: string,
  relativePath: string,
  requestInit: RequestInit,
): Promise<T | null> {
  const requestPath = `${locale}/${relativePath}`;
  const cache = await getInflightCache();

  const existing = cache.get(requestPath);
  if (existing !== undefined) {
    return existing as Promise<T | null>;
  }

  const promise = fetcherInternal<T>(locale, relativePath, requestInit).finally(() => {
    cache.delete(requestPath);
  });

  cache.set(requestPath, promise);
  return promise;
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

  // Forward the incoming request's cookies to backend API calls during SSR.
  // On the client, getStore() returns undefined and this is a no-op.
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

  const fetchOptions = buildFetchOptions(requestInit, headers);
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
        const retryOptions = buildFetchOptions(requestInit, headers);
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

  const uploadOptions: RequestInit = {
    method: "POST",
    credentials: "include",
    headers,
    body: formData,
  };

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
