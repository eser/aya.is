import { getBackendUri } from "@/shared/config.ts";
import { isTokenExpiringSoon, refreshToken } from "@/shared/modules/auth/token-refresh.ts";
import type { Result } from "./types.ts";

export async function fetcher<T>(relativePath: string, requestInit: RequestInit = {}): Promise<T | null> {
  const targetUrl = `${getBackendUri()}${relativePath}`;
  console.log("targetUrl", targetUrl);

  // Get auth token from localStorage (only available on client)
  let authToken: string | null = null;
  if (globalThis.localStorage !== undefined) {
    authToken = localStorage.getItem("auth_token");
  }

  // Check if token needs refresh
  if (authToken && isTokenExpiringSoon()) {
    const newToken = await refreshToken();

    if (newToken !== null) {
      authToken = newToken;
    }
  }

  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  if (authToken !== null) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  const request = await fetch(targetUrl, {
    ...requestInit,
    headers: {
      ...headers,
      ...(requestInit.headers ?? {}),
    },
  });

  // Handle authentication errors
  if (request.status === 401 || request.status === 405) {
    // If we have a token, try to refresh it and retry
    if (authToken) {
      const newToken = await refreshToken();

      if (newToken !== null) {
        headers["Authorization"] = `Bearer ${newToken}`;

        const retryRequest = await fetch(targetUrl, { headers });

        if (retryRequest.status === 404) {
          return null;
        }

        if (retryRequest.status === 401 || retryRequest.status === 405) {
          // Still unauthorized after refresh, return null for unauthenticated state
          return null;
        }

        if (retryRequest.status >= 500) {
          throw new Error(`Internal server error: ${retryRequest.status}`);
        }

        const retryResult = (await retryRequest.json()) as Result<T>;

        if (retryResult.error !== undefined && retryResult.error !== null) {
          throw new Error(retryResult.error);
        }

        return retryResult.data;
      }
    }

    // No token or refresh failed - return null for unauthenticated state
    return null;
  }

  if (request.status === 404) {
    return null;
  }

  if (request.status >= 500) {
    throw new Error(`Internal server error: ${request.status}`);
  }

  const result = (await request.json()) as Result<T>;

  if (result.error !== undefined && result.error !== null) {
    throw new Error(result.error);
  }

  return result.data;
}
