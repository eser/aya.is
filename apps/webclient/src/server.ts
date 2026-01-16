// Server Entry Point - handles custom domain URL rewriting
import process from "node:process";
import handler, { createServerEntry } from "@tanstack/react-start/server-entry";
import { DEFAULT_LOCALE, getMainDomains, isValidLocale } from "./config.ts";

// Cache for custom domain lookups
const domainCache = new Map<string, { slug: string; locale: string | null } | null>();
const CACHE_TTL = 5 * 60 * 1000;
const cacheTimestamps = new Map<string, number>();

// Use DEFAULT_LOCALE env var, fallback to config default
function getDefaultLocale(): string {
  return process.env.DEFAULT_LOCALE ?? DEFAULT_LOCALE;
}

function isMainDomain(host: string): boolean {
  const customDomainOverride = process.env.CUSTOM_DOMAIN;
  if (customDomainOverride !== undefined && customDomainOverride !== "") {
    const hostWithoutPort = host.split(":")[0];
    if (hostWithoutPort === "localhost" || hostWithoutPort === "127.0.0.1") {
      return false;
    }
  }

  const hostWithoutPort = host.split(":")[0];
  const mainDomains = getMainDomains();
  return mainDomains.includes(hostWithoutPort);
}

async function fetchCustomDomain(
  host: string,
): Promise<{ slug: string; locale: string | null } | null> {
  const customDomainOverride = process.env.CUSTOM_DOMAIN;
  const effectiveHost = customDomainOverride !== undefined && customDomainOverride !== "" ? customDomainOverride : host;

  const now = Date.now();
  const cached = domainCache.get(effectiveHost);
  const cachedTime = cacheTimestamps.get(effectiveHost);

  if (cached !== undefined && cachedTime !== undefined && now - cachedTime < CACHE_TTL) {
    return cached;
  }

  try {
    const backendUri = process.env.BACKEND_URI ?? process.env.PUBLIC_BACKEND_URI ?? "https://api.aya.is";
    const response = await fetch(
      `${backendUri}/en/site/custom-domains/${encodeURIComponent(effectiveHost)}`,
    );

    if (!response.ok) {
      domainCache.set(effectiveHost, null);
      cacheTimestamps.set(effectiveHost, now);
      return null;
    }

    const responseBody = await response.json();
    const profileData = responseBody.data;

    if (profileData === null || profileData === undefined) {
      domainCache.set(effectiveHost, null);
      cacheTimestamps.set(effectiveHost, now);
      return null;
    }

    const result = { slug: profileData.slug, locale: profileData.locale ?? null };
    domainCache.set(effectiveHost, result);
    cacheTimestamps.set(effectiveHost, now);
    return result;
  } catch (error) {
    console.error(`[custom-domain] Failed to lookup: ${effectiveHost}`, error);
    domainCache.set(effectiveHost, null);
    cacheTimestamps.set(effectiveHost, now);
    return null;
  }
}

function rewriteUrlForCustomDomain(
  url: URL,
  profileSlug: string,
  profileDefaultLocale: string | null,
): URL {
  const pathname = url.pathname;
  const segments = pathname.split("/").filter(Boolean);

  let locale: string;
  let remainingPath: string;

  if (segments.length > 0 && isValidLocale(segments[0])) {
    // Explicit locale: eser.dev/en/blog -> /en/eser/blog
    locale = segments[0];
    remainingPath = segments.slice(1).join("/");
  } else {
    // No locale: eser.dev/blog -> /tr/eser/blog (using profile's default locale)
    locale = profileDefaultLocale ?? getDefaultLocale();
    remainingPath = segments.join("/");
  }

  const newPathname = `/${locale}/${profileSlug}${remainingPath ? "/" + remainingPath : ""}`;

  const newUrl = new URL(url);
  newUrl.pathname = newPathname;
  return newUrl;
}

export default createServerEntry({
  async fetch(request) {
    const url = new URL(request.url);
    const host = url.host;

    // Skip static assets and internal paths
    if (
      url.pathname.startsWith("/_") ||
      url.pathname.startsWith("/api/") ||
      url.pathname.startsWith("/assets/") ||
      url.pathname.includes(".")
    ) {
      return handler.fetch(request);
    }

    // Check if this is a custom domain
    if (!isMainDomain(host)) {
      const customDomain = await fetchCustomDomain(host);

      if (customDomain !== null) {
        const newUrl = rewriteUrlForCustomDomain(url, customDomain.slug, customDomain.locale);

        console.log(`[custom-domain] ${host}${url.pathname} -> ${newUrl.pathname}`);

        // Create a new request with the rewritten URL
        const newRequest = new Request(newUrl.toString(), {
          method: request.method,
          headers: request.headers,
          body: request.body,
          redirect: request.redirect,
          signal: request.signal,
        });

        return handler.fetch(newRequest);
      }
    }

    return handler.fetch(request);
  },
});
