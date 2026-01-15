// Custom Domain URL Rewriting Middleware
// This middleware detects custom domains and rewrites URLs to internal route patterns
// Example: eser.dev/tr/cv → /_cd/eser/tr/cv

import { defineEventHandler, getHeader, setHeader } from "h3";
import process from "node:process";

// Supported locales (must match config.ts)
const SUPPORTED_LOCALES = [
  "tr",
  "en",
  "fr",
  "de",
  "es",
  "pt-PT",
  "it",
  "nl",
  "ja",
  "ko",
  "ru",
  "zh-CN",
  "ar",
] as const;

// Default locale for custom domains when profile has no preference
const CUSTOM_DOMAIN_DEFAULT_LOCALE = "en";

// Main domains that should NOT be treated as custom domains
const MAIN_DOMAINS = [
  "localhost",
  "127.0.0.1",
  "aya.is",
  "www.aya.is",
];

// Cache for custom domain lookups (simple in-memory cache)
const domainCache = new Map<
  string,
  { slug: string; locale?: string | null } | null
>();
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes
const cacheTimestamps = new Map<string, number>();

function isValidLocale(locale: string): boolean {
  return SUPPORTED_LOCALES.includes(locale as typeof SUPPORTED_LOCALES[number]);
}

function isMainDomain(host: string): boolean {
  // Remove port for comparison
  const hostWithoutPort = host.split(":")[0];
  return MAIN_DOMAINS.some(
    (domain) => hostWithoutPort === domain || hostWithoutPort === `www.${domain}`,
  );
}

async function fetchCustomDomain(
  host: string,
): Promise<{ slug: string; locale?: string | null } | null> {
  // Check cache first
  const now = Date.now();
  const cached = domainCache.get(host);
  const cachedTime = cacheTimestamps.get(host);

  if (cached !== undefined && cachedTime && now - cachedTime < CACHE_TTL) {
    return cached;
  }

  try {
    // Get backend URI from environment
    const backendUri = process.env.BACKEND_URI ||
      process.env.PUBLIC_BACKEND_URI || "https://api.aya.is";
    const response = await fetch(
      `${backendUri}/custom-domains/${encodeURIComponent(host)}`,
    );

    if (!response.ok) {
      domainCache.set(host, null);
      cacheTimestamps.set(host, now);
      return null;
    }

    const data = await response.json();
    const result = { slug: data.slug, locale: data.locale };
    domainCache.set(host, result);
    cacheTimestamps.set(host, now);
    return result;
  } catch (error) {
    console.error(
      `[custom-domain] Failed to lookup custom domain: ${host}`,
      error,
    );
    domainCache.set(host, null);
    cacheTimestamps.set(host, now);
    return null;
  }
}

export default defineEventHandler(async (event) => {
  const host = getHeader(event, "host");

  if (!host) {
    return;
  }

  // Skip main domains
  if (isMainDomain(host)) {
    return;
  }

  // Skip static assets and API routes
  const path = event.path || "/";
  if (
    path.startsWith("/_") ||
    path.startsWith("/api/") ||
    path.startsWith("/assets/") ||
    path.includes(".")
  ) {
    return;
  }

  // Lookup custom domain from backend
  const customDomain = await fetchCustomDomain(host);

  if (!customDomain) {
    // Not a registered custom domain, let it pass through (will 404)
    return;
  }

  const { slug: profileSlug, locale: profileLocale } = customDomain;

  // Profile's preferred locale, fallback to default
  const defaultLocale = profileLocale ?? CUSTOM_DOMAIN_DEFAULT_LOCALE;

  // Parse incoming path
  const segments = path.split("/").filter(Boolean);

  let locale: string;
  let remainingPath: string;

  if (segments.length > 0 && isValidLocale(segments[0])) {
    // Explicit locale in URL: eser.dev/tr/cv
    locale = segments[0];
    remainingPath = segments.slice(1).join("/");
  } else {
    // No locale in URL: eser.dev/cv → use profile's default locale
    locale = defaultLocale;
    remainingPath = segments.join("/");
  }

  // Rewrite URL: /_cd/{profile}/{locale}/{path}
  const newPath = `/_cd/${profileSlug}/${locale}${remainingPath ? "/" + remainingPath : ""}`;

  // Update the request URL
  event.node.req.url = newPath;

  // Pass context via headers for components to read
  setHeader(event, "x-custom-domain", "true");
  setHeader(event, "x-custom-domain-host", host);
  setHeader(event, "x-custom-domain-profile", profileSlug);
  setHeader(event, "x-custom-domain-locale", locale);
  setHeader(event, "x-custom-domain-default-locale", defaultLocale);
});
