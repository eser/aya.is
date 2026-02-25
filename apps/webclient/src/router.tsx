import { createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import { predefinedSlugs, siteConfig } from "@/config";

import type { RequestContext } from "@/request-context";

// globalThis injection for client-side hydration
declare global {
  var __REQUEST_CONTEXT__: RequestContext | undefined;
}

async function getRequestContext(): Promise<RequestContext | undefined> {
  // read injected (client-side)
  if (globalThis.__REQUEST_CONTEXT__ !== undefined) {
    return globalThis.__REQUEST_CONTEXT__;
  }

  // Server-side only: import.meta.env.SSR is replaced at build time
  // This allows Vite to dead-code eliminate this branch from client bundle
  if (import.meta.env.SSR) {
    const { requestContextBinder } = await import("./server/request-context-binder");
    return requestContextBinder.getStore();
  }

  return undefined;
}

// Create a new router instance
export async function getRouter() {
  const requestContext = await getRequestContext();

  const customDomainProfileSlug = (
    requestContext?.domainConfiguration.type === "custom-domain" &&
    requestContext.domainConfiguration.profileSlug
  ) || undefined;

  const router = createRouter({
    routeTree,
    context: { requestContext },
    defaultPreload: "intent",
    scrollRestoration: true,

    rewrite: customDomainProfileSlug !== undefined
      ? {
        input: ({ url }) => {
          // /en/about -> /en/{profileSlug}/about
          const pathParts = url.pathname.split("/").filter(Boolean);

          // Root path: let routes/index.tsx detect preferred locale via
          // cookie → Accept-Language → domain default → fallback chain
          if (pathParts.length === 0) {
            return url;
          }

          // System routes (e.g., /auth/callback) are not profile-scoped
          if (predefinedSlugs.includes(pathParts[0])) {
            return url;
          }

          pathParts.splice(1, 0, customDomainProfileSlug);
          url.pathname = `/${pathParts.join("/")}`;

          return url;
        },
        output: ({ url }) => {
          // /tr/{profileSlug}/about -> /tr/about (always keeps locale)
          const pathParts = url.pathname.split("/").filter(Boolean);

          // System routes are not profile-scoped
          if (predefinedSlugs.includes(pathParts[0])) {
            return url;
          }

          // Remove profile slug from position 1
          if (pathParts[1] === customDomainProfileSlug) {
            pathParts.splice(1, 1);
            url.pathname = `/${pathParts.join("/")}`;

            return url;
          }

          const newUrl = new URL(siteConfig.host);
          newUrl.pathname = `/${pathParts.join("/")}`;
          return newUrl;
        },
      }
      : undefined,
  });

  return router;
}

declare module "@tanstack/react-router" {
  interface Register {
    router: Awaited<ReturnType<typeof getRouter>>;
  }
}
