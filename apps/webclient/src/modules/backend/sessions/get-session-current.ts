// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { getBackendUri } from "@/config.ts";
import { getAuthToken } from "../fetcher.ts";
import type { SessionCurrentResponse } from "./types.ts";

/**
 * Get current session state via cookie-based check, with a Bearer fallback.
 * The session cookie is not sent when the frontend and the API are on different sites (e.g. localhost:3000 talking to
 * api.aya.is), so the stored token is sent as well; the backend accepts the cookie first and the Bearer token second.
 * Returns auth state, token, and preferences in a single response.
 * This is the primary session initialization call on app mount.
 */
export async function getSessionCurrent(locale: string): Promise<SessionCurrentResponse> {
  const backendUri = getBackendUri();
  const authToken = getAuthToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (authToken !== null) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  try {
    const response = await fetch(`${backendUri}/${locale}/sessions/_current`, {
      method: "GET",
      credentials: "include",
      headers,
    });

    if (!response.ok) {
      return { authenticated: false };
    }

    const result = await response.json();

    return result.data as SessionCurrentResponse;
  } catch {
    return { authenticated: false };
  }
}
