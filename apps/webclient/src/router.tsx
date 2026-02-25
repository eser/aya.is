import { createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import { DEFAULT_LOCALE, predefinedSlugs, siteConfig } from "@/config";
import i18n from "@/modules/i18n/i18n";

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

  const defaultLocale = (
    (
      requestContext?.domainConfiguration.type === "main" ||
      requestContext?.domainConfiguration.type === "custom-domain"
    ) && requestContext?.domainConfiguration?.defaultCulture
  ) || DEFAULT_LOCALE;

  const router = createRouter({
    routeTree,
    context: { requestContext, i18nInstance: i18n },
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
