import { auth } from "@/lib/auth/server";
import { createServerFileRoute } from "@tanstack/react-start/server";

/**
 * BetterAuth route handler for TanStack Start
 *
 * This catches all /api/auth/* requests and forwards them to BetterAuth
 */
export const ServerRoute = createServerFileRoute("/api/auth/$").methods({
  GET: ({ request }) => auth.handler(request),
  POST: ({ request }) => auth.handler(request),
});
