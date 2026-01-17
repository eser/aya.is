import process from "node:process";
import { CUSTOM_DOMAIN_DEFAULT_LOCALE } from "@/config";
import { isMainDomain } from "@/shared.ts";

import type { DomainConfiguration } from "@/request-context";

// Cache for custom domain lookups
const domainCache = new Map<string, { slug: string; locale: string | null } | null>();
const CACHE_TTL = 5 * 60 * 1000;
const cacheTimestamps = new Map<string, number>();

export const defaultDomainConfiguration: DomainConfiguration = {
  type: "main",
  defaultCulture: CUSTOM_DOMAIN_DEFAULT_LOCALE,
  allowsWwwPrefix: false,
};

export async function fetchCustomDomainFromBackend(
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

export async function getDomainConfiguration(address: string | undefined): Promise<DomainConfiguration> {
  if (address === undefined) {
    return defaultDomainConfiguration;
  }

  const hostname = process.env.CUSTOM_DOMAIN ?? address.split(":")[0];

  if (isMainDomain(hostname)) {
    return defaultDomainConfiguration;
  }

  // For non-main domains, try backend lookup
  const backendResult = await fetchCustomDomainFromBackend(hostname);

  if (backendResult !== null) {
    return {
      type: "custom-domain",
      profileSlug: backendResult.slug,
      defaultCulture: CUSTOM_DOMAIN_DEFAULT_LOCALE,
      allowsWwwPrefix: false,
    };
  }

  return {
    type: "not-configured",
  };
}
