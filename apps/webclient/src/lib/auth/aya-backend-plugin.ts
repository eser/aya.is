import type { BetterAuthPlugin } from "better-auth";
import { createAuthEndpoint } from "better-auth/plugins";
import { DEFAULT_LOCALE, getBackendUri } from "@/config";

const BEARER_PREFIX = "Bearer ";

/**
 * Extracts the token from a Bearer authorization header.
 * Returns null if the header is missing or doesn't start with "Bearer ".
 */
function extractBearerToken(authHeader: string | null | undefined): string | null {
  if (authHeader === null || authHeader === undefined) {
    return null;
  }

  if (!authHeader.startsWith(BEARER_PREFIX)) {
    return null;
  }

  const token = authHeader.slice(BEARER_PREFIX.length);
  return token.length > 0 ? token : null;
}

/**
 * Custom BetterAuth plugin that proxies auth operations to api.aya.is
 *
 * This plugin:
 * - Proxies GitHub OAuth to the api.aya.is backend
 * - Fetches sessions from api.aya.is
 * - Handles token refresh via api.aya.is
 * - Does NOT use a local database (api.aya.is is the source of truth)
 */
export const ayaBackendPlugin = () => {
  return {
    id: "aya-backend",

    endpoints: {
      // Proxy GitHub login to api.aya.is
      githubLogin: createAuthEndpoint(
        "/github/login",
        { method: "GET" },
        (ctx) => {
          const backendUri = getBackendUri();
          const url = new URL(ctx.request?.url ?? "");
          const locale = url.searchParams.get("locale") ?? DEFAULT_LOCALE;
          const redirectUri = url.searchParams.get("redirect_uri") ?? "/";

          // Redirect to api.aya.is OAuth endpoint
          const loginUrlObj = new URL(`${backendUri}/${locale}/auth/github/login`);
          loginUrlObj.searchParams.set("redirect_uri", redirectUri);
          const loginUrl = loginUrlObj.toString();

          return new Response(null, {
            status: 302,
            headers: {
              Location: loginUrl,
            },
          });
        },
      ),

      // Handle callback from api.aya.is (receives auth_token in query)
      githubCallback: createAuthEndpoint(
        "/github/callback",
        { method: "GET" },
        async (ctx) => {
          const url = new URL(ctx.request?.url ?? "");
          const authToken = url.searchParams.get("auth_token");
          const locale = url.searchParams.get("locale") ?? DEFAULT_LOCALE;

          if (authToken === null) {
            return ctx.json({ error: "No auth token provided" }, {
              status: 400,
            });
          }

          // Fetch session from api.aya.is
          const backendUri = getBackendUri();
          try {
            const sessionResponse = await fetch(
              `${backendUri}/${locale}/auth/session`,
              {
                headers: { Authorization: `Bearer ${authToken}` },
              },
            );

            if (!sessionResponse.ok) {
              return ctx.json(
                { error: "Failed to fetch session" },
                { status: 401 },
              );
            }

            const session = await sessionResponse.json();

            return ctx.json({
              token: authToken,
              session: session.data || session,
            });
          } catch {
            return ctx.json(
              { error: "Failed to authenticate" },
              { status: 500 },
            );
          }
        },
      ),

      // Get current session from api.aya.is
      getSession: createAuthEndpoint(
        "/session",
        { method: "GET" },
        async (ctx) => {
          const token = extractBearerToken(ctx.request?.headers.get("authorization"));

          if (token === null) {
            return ctx.json({ session: null });
          }

          const url = new URL(ctx.request?.url ?? "");
          const locale = url.searchParams.get("locale") ?? DEFAULT_LOCALE;
          const backendUri = getBackendUri();

          try {
            const response = await fetch(
              `${backendUri}/${locale}/auth/session`,
              {
                headers: { Authorization: `Bearer ${token}` },
              },
            );

            if (!response.ok) {
              return ctx.json({ session: null });
            }

            const data = await response.json();
            return ctx.json({ session: data.data || data });
          } catch {
            return ctx.json({ session: null });
          }
        },
      ),

      // Refresh token via api.aya.is
      refreshToken: createAuthEndpoint(
        "/refresh",
        { method: "POST" },
        async (ctx) => {
          const token = extractBearerToken(ctx.request?.headers.get("authorization"));

          if (token === null) {
            return ctx.json({ error: "No token provided" }, { status: 401 });
          }

          const backendUri = getBackendUri();

          try {
            const body = await ctx.request?.json().catch(() => ({}));
            const locale = body?.locale ?? DEFAULT_LOCALE;

            const response = await fetch(
              `${backendUri}/${locale}/auth/refresh`,
              {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                  Authorization: `Bearer ${token}`,
                },
              },
            );

            if (!response.ok) {
              return ctx.json({ error: "Refresh failed" }, { status: 401 });
            }

            const data = await response.json();
            return ctx.json(data.data || data);
          } catch {
            return ctx.json({ error: "Refresh failed" }, { status: 401 });
          }
        },
      ),

      // Logout via api.aya.is
      logout: createAuthEndpoint(
        "/logout",
        { method: "POST" },
        async (ctx) => {
          const token = extractBearerToken(ctx.request?.headers.get("authorization"));
          const backendUri = getBackendUri();

          if (token !== null) {
            try {
              const body = await ctx.request?.json().catch(() => ({}));
              const locale = body?.locale ?? DEFAULT_LOCALE;

              await fetch(`${backendUri}/${locale}/auth/logout`, {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                  Authorization: `Bearer ${token}`,
                },
              });
            } catch {
              // Ignore logout errors
            }
          }

          return ctx.json({ success: true });
        },
      ),
    },
  } satisfies BetterAuthPlugin;
};
