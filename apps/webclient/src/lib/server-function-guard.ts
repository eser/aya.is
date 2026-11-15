// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Pure request check for the TanStack Start server-function endpoint (kept free of env access so it is unit-testable).

/** The path prefix TanStack Start serves server functions under unless TSS_SERVER_FN_BASE overrides it. */
export const DEFAULT_SERVER_FN_BASE = "/_serverFn/";

/** The header TanStack Start's own client sends on every server-function call. */
const SERVER_FN_HEADER = "x-tsr-serverFn";

type HeaderReader = { get(name: string): string | null };

/**
 * Blocks requests to the server-function endpoint that did not come from the Start client.
 *
 * Reflected XSS in server-function responses (CVE-2026-102989, fixed in @tanstack/react-start 1.168.60 and
 * @tanstack/start-server-core 1.169.39) is triggered by a crafted link a victim opens in the browser. A link
 * navigation cannot set custom headers, so requiring the Start client's header, and refusing document
 * navigations outright, closes that path without changing the pinned TanStack set. Remove this guard once the
 * set is updated to the patched versions.
 */
export function isBlockedServerFunctionRequest(
  pathname: string,
  headers: HeaderReader,
  base: string = DEFAULT_SERVER_FN_BASE,
): boolean {
  if (!pathname.startsWith(base)) {
    return false;
  }

  if (headers.get(SERVER_FN_HEADER) !== "true") {
    return true;
  }

  return headers.get("sec-fetch-mode") === "navigate" || headers.get("sec-fetch-dest") === "document";
}
