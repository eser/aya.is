import { createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import { DEFAULT_LOCALE, predefinedSlugs, siteConfig } from "@/config";

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

  // read from AsyncLocalStorage (server-side)
  const { requestContextBinder } = await import("./server/request-context-binder");
  const requestContext = requestContextBinder.getStore();

  return requestContext;
}

// Create a new router instance
export async function getRouter() {
  const requestContext = await getRequestContext();

  const customDomainProfileSlug = (
    requestContext?.domainConfiguration.type === "custom-domain" &&
    requestContext.domainConfiguration.profileSlug
  ) || undefined;

  const defaultLocale = (
    (
      requestContext?.domainConfiguration.type === "main" ||
      requestContext?.domainConfiguration.type === "custom-domain"
    ) && requestContext?.domainConfiguration?.defaultCulture
  ) || DEFAULT_LOCALE;

  const router = createRouter({
    routeTree,
    context: { requestContext },
    defaultPreload: "intent",
    scrollRestoration: true,

    rewrite: customDomainProfileSlug !== undefined
      ? {
        input: ({ url }) => {
          // /en/about -> /en/{profileSlug}/about
          // / -> /tr/{profileSlug} (adds default locale for empty paths)
          const pathParts = url.pathname.split("/").filter(Boolean);

          // If path is empty, prepend the default locale
          if (pathParts.length === 0) {
            pathParts.push(defaultLocale);
          }

          // Skip profile slug injection for system routes (auth, api, etc.)
          // These routes exist at /$locale/auth/* and shouldn't be nested under profile
          if (predefinedSlugs.includes(pathParts[1])) {
            return url;
          }

          pathParts.splice(1, 0, customDomainProfileSlug);
          url.pathname = `/${pathParts.join("/")}`;

          return url;
        },
        output: ({ url }) => {
          // /tr/{profileSlug}/about -> /tr/about (always keeps locale)
          const pathParts = url.pathname.split("/").filter(Boolean);

          // Skip profile slug injection for system routes (auth, api, etc.)
          // These routes exist at /$locale/auth/* and shouldn't be nested under profile
          if (predefinedSlugs.includes(pathParts[1])) {
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
