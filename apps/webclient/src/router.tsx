import { createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";

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
            pathParts.splice(1, 0, customDomainProfileSlug);
            url.pathname = `/${pathParts.join("/")}`;

            return url;
          },
          output: ({ url }) => {
            // /en/{profileSlug}/about -> /en/about
            const pathParts = url.pathname.split("/").filter(Boolean);
            if (pathParts[1] === customDomainProfileSlug) {
              pathParts.splice(1, 1);
              url.pathname = `/${pathParts.join("/")}`;
            }

            return url;
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
