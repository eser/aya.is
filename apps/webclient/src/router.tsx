import { createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import { DEFAULT_LOCALE, isValidLocale, predefinedSlugs, siteConfig, SUPPORTED_LOCALES } from "@/config";

import type { SupportedLocaleCode } from "@/config";
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

/**
 * Detect preferred locale from request context.
 * Same priority as getPreferredLocale() but works synchronously from context data.
 *
 * Priority: cookie → Accept-Language → domain default → DEFAULT_LOCALE
 */
function detectPreferredLocale(requestContext: RequestContext | undefined): SupportedLocaleCode {
  // 1. Check SITE_LOCALE cookie
  const cookieHeader = requestContext?.cookieHeader;
  if (cookieHeader !== undefined) {
    const match = cookieHeader.match(/(?:^|;\s*)SITE_LOCALE=([^;]+)/);
    if (match !== null && match[1] !== undefined && isValidLocale(match[1])) {
      return match[1];
    }
  }

  // 2. Check Accept-Language header (server-side only; on client, navigator.language)
  const acceptLanguage = requestContext?.acceptLanguageHeader;
  if (acceptLanguage !== undefined) {
    const matched = parseAcceptLanguageSync(acceptLanguage);
    if (matched !== null) {
      return matched;
    }
  } else if (typeof navigator !== "undefined" && navigator.language !== undefined) {
    // Client fallback: use navigator.language
    const lang = navigator.language.toLowerCase();
    const exactMatch = SUPPORTED_LOCALES.find((l) => l.toLowerCase() === lang);
    if (exactMatch !== undefined) {
      return exactMatch;
    }
    const prefix = lang.split("-")[0];
    if (prefix !== undefined) {
      const prefixMatch = SUPPORTED_LOCALES.find((l) => l.toLowerCase() === prefix);
      if (prefixMatch !== undefined) {
        return prefixMatch;
      }
    }
  }

  // 3. Domain's default locale
  const domainConfig = requestContext?.domainConfiguration;
  if (domainConfig?.type === "custom-domain" || domainConfig?.type === "main") {
    const domainLocale = domainConfig.defaultCulture;
    if (isValidLocale(domainLocale)) {
      return domainLocale;
    }
  }

  // 4. System default
  return DEFAULT_LOCALE;
}

/**
 * Synchronous Accept-Language parser (mirrors parseAcceptLanguage in get-locale.ts).
 */
function parseAcceptLanguageSync(header: string): SupportedLocaleCode | null {
  const entries = header
    .split(",")
    .map((part) => {
      const trimmed = part.trim();
      const [lang, ...params] = trimmed.split(";");
      const qParam = params.find((p) => p.trim().startsWith("q="));
      const quality = qParam !== undefined ? Number.parseFloat(qParam.trim().slice(2)) : 1.0;
      return { lang: lang?.trim() ?? "", quality };
    })
    .filter((entry) => entry.lang.length > 0)
    .sort((a, b) => b.quality - a.quality);

  for (const entry of entries) {
    const lang = entry.lang.toLowerCase();

    const exactMatch = SUPPORTED_LOCALES.find((locale) => locale.toLowerCase() === lang);
    if (exactMatch !== undefined) {
      return exactMatch;
    }

    const langPrefix = lang.split("-")[0];
    if (langPrefix !== undefined) {
      const prefixMatch = SUPPORTED_LOCALES.find((locale) => locale.toLowerCase() === langPrefix);
      if (prefixMatch !== undefined) {
        return prefixMatch;
      }
    }
  }

  return null;
}

// Create a new router instance
export async function getRouter() {
  const requestContext = await getRequestContext();

  const customDomainProfileSlug = (
    requestContext?.domainConfiguration.type === "custom-domain" &&
    requestContext.domainConfiguration.profileSlug
  ) || undefined;

  // Detect preferred locale for custom domain root path.
  // Uses: cookie → Accept-Language → domain default → DEFAULT_LOCALE
  const preferredLocale = detectPreferredLocale(requestContext);

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

          // Root path: use detected preferred locale
          // (cookie → Accept-Language → domain default → 'en')
          if (pathParts.length === 0) {
            pathParts.push(preferredLocale);
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
