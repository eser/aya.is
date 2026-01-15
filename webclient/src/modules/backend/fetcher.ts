import { getBackendUri } from "@/config";
import type { Result } from "./types";

// Token management utilities
export function getAuthToken(): string | null {
  if (typeof window !== "undefined" && window.localStorage) {
    return localStorage.getItem("auth_token");
  }
  return null;
}

export function setAuthToken(token: string, expiresAt: number): void {
  if (typeof window === "undefined") return;
  localStorage.setItem("auth_token", token);
  localStorage.setItem("auth_token_expires_at", expiresAt.toString());
}

export function clearAuthData(): void {
  if (typeof window === "undefined") return;
  localStorage.removeItem("auth_token");
  localStorage.removeItem("auth_token_expires_at");
  localStorage.removeItem("auth_session");
}

export function isTokenExpiringSoon(): boolean {
  if (typeof window === "undefined" || !window.localStorage) {
    return false;
  }

  const expiresAt = localStorage.getItem("auth_token_expires_at");
  if (expiresAt === null) {
    return false;
  }

  const expirationTime = parseInt(expiresAt, 10);
  const now = Date.now();
  const fiveMinutes = 5 * 60 * 1000;

  return expirationTime - now < fiveMinutes;
}

let refreshPromise: Promise<string | null> | null = null;

export async function refreshTokenRequest(locale: string = "en"): Promise<{ token: string; expiresAt: number } | null> {
  // Deduplicate concurrent refresh calls
  if (refreshPromise !== null) {
    const token = await refreshPromise;
    if (token !== null) {
      const expiresAt = localStorage.getItem("auth_token_expires_at");
      return { token, expiresAt: expiresAt !== null ? parseInt(expiresAt, 10) : Date.now() + 24 * 60 * 60 * 1000 };
    }
    return null;
  }

  const currentToken = getAuthToken();
  if (currentToken === null) return null;

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
        const newToken = data.data.token;
        const expiresAt = data.data.expires_at ?? Date.now() + 24 * 60 * 60 * 1000;
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
    const expiresAt = localStorage.getItem("auth_token_expires_at");
    return { token, expiresAt: expiresAt !== null ? parseInt(expiresAt, 10) : Date.now() + 24 * 60 * 60 * 1000 };
  }
  return null;
}

async function refreshToken(): Promise<string | null> {
  const currentToken = getAuthToken();
  if (currentToken === null) {
    return null;
  }

  try {
    const backendUri = getBackendUri();
    const response = await fetch(`${backendUri}/auth/refresh`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${currentToken}`,
      },
    });

    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    if (data.token) {
      if (typeof window !== "undefined" && window.localStorage) {
        localStorage.setItem("auth_token", data.token);
        if (data.expires_at) {
          localStorage.setItem("auth_token_expires_at", String(data.expires_at));
        }
      }
      return data.token;
    }

    return null;
  } catch {
    return null;
  }
}

export async function fetcher<T>(
  relativePath: string,
  requestInit: RequestInit = {}
): Promise<T | null> {
  const targetUrl = `${getBackendUri()}${relativePath}`;

  // Get auth token from localStorage (only available on client)
  let authToken = getAuthToken();

  // Check if token needs refresh
  if (authToken !== null && isTokenExpiringSoon()) {
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

  const response = await fetch(targetUrl, {
    ...requestInit,
    credentials: "include", // Send cookies for cross-domain SSO
    headers: {
      ...headers,
      ...(requestInit.headers ?? {}),
    },
  });

  // Handle authentication errors
  if (response.status === 401 || response.status === 405) {
    // If we have a token, try to refresh it and retry
    if (authToken !== null) {
      const newToken = await refreshToken();

      if (newToken !== null) {
        headers["Authorization"] = `Bearer ${newToken}`;

        const retryResponse = await fetch(targetUrl, {
          ...requestInit,
          credentials: "include",
          headers: {
            ...headers,
            ...(requestInit.headers ?? {}),
          },
        });

        if (retryResponse.status === 404) {
          return null;
        }

        if (retryResponse.status === 401 || retryResponse.status === 405) {
          // Still unauthorized after refresh, return null for unauthenticated state
          return null;
        }

        if (retryResponse.status >= 500) {
          throw new Error(`Internal server error: ${retryResponse.status}`);
        }

        const retryResult = (await retryResponse.json()) as Result<T>;

        if (retryResult.error !== undefined && retryResult.error !== null) {
          throw new Error(retryResult.error);
        }

        return retryResult.data;
      }
    }

    // No token or refresh failed - return null for unauthenticated state
    return null;
  }

  if (response.status === 404) {
    return null;
  }

  if (response.status >= 500) {
    throw new Error(`Internal server error: ${response.status}`);
  }

  const result = (await response.json()) as Result<T>;

  if (result.error !== undefined && result.error !== null) {
    throw new Error(result.error);
  }

  return result.data;
}

// Fetcher for file uploads (multipart/form-data)
export async function uploadFetcher<T>(
  relativePath: string,
  formData: FormData
): Promise<T | null> {
  const targetUrl = `${getBackendUri()}${relativePath}`;

  let authToken = getAuthToken();

  const headers: HeadersInit = {};

  if (authToken !== null) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  const response = await fetch(targetUrl, {
    method: "POST",
    credentials: "include",
    headers,
    body: formData,
  });

  if (response.status === 401 || response.status === 405) {
    return null;
  }

  if (response.status === 404) {
    return null;
  }

  if (response.status >= 500) {
    throw new Error(`Internal server error: ${response.status}`);
  }

  const result = (await response.json()) as Result<T>;

  if (result.error !== undefined && result.error !== null) {
    throw new Error(result.error);
  }

  return result.data;
}
