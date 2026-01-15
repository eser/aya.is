import { betterAuth } from "better-auth";
import { ayaBackendPlugin } from "./aya-backend-plugin";

/**
 * Get the base URL for the auth server
 */
function getServerBaseURL(): string {
  const host = import.meta.env.VITE_HOST ?? "http://localhost:3000";
  return host;
}

/**
 * BetterAuth server configuration
 *
 * This configuration:
 * - Uses the custom aya-backend plugin to proxy auth to api.aya.is
 * - Does NOT use a local database
 * - Disables built-in providers (we use api.aya.is for GitHub OAuth)
 */
export const auth = betterAuth({
  // No database - api.aya.is is the source of truth
  // @ts-expect-error - database can be undefined for custom plugin-only auth
  database: undefined,

  // Base URL for the auth server
  baseURL: getServerBaseURL(),

  // Base path for auth endpoints
  basePath: "/api/auth",

  // Disable built-in email/password (using custom plugin)
  emailAndPassword: {
    enabled: false,
  },

  // Use custom plugin for api.aya.is proxy
  plugins: [ayaBackendPlugin()],

  // Session configuration (matches api.aya.is)
  session: {
    expiresIn: 60 * 60 * 24 * 7, // 7 days
    updateAge: 60 * 60 * 24, // 1 day
  },

  // Advanced options
  advanced: {
    // Don't generate secret since we're using external auth
    generateId: () => crypto.randomUUID(),
  },
});

export type Auth = typeof auth;
