// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// SEO Meta Tag Utilities
// Thin wrapper around seo-utils.ts that injects siteConfig values

import { siteConfig } from "@/config";
import {
  buildUrl as buildUrlPure,
  computeContentLanguage,
  computeStoryCanonicalUrl as computeStoryCanonicalUrlPure,
  generateCanonicalLink,
  generateMetaTags as generateMetaTagsPure,
  truncateDescription as truncateDescriptionPure,
} from "@/lib/seo-utils.ts";

// Re-export types
export type { MetaTag, SeoMeta } from "@/lib/seo-utils.ts";

// Re-export pure functions that don't need siteConfig
export { computeContentLanguage, generateCanonicalLink };

const SITE_DEFAULTS = {
  host: siteConfig.host,
  name: siteConfig.name,
  defaultImage: `${siteConfig.host}/og-image.png`,
};

export function generateMetaTags(
  meta: Parameters<typeof generateMetaTagsPure>[0],
): ReturnType<typeof generateMetaTagsPure> {
  return generateMetaTagsPure(meta, SITE_DEFAULTS);
}

export function buildUrl(locale: string, ...pathSegments: string[]): string {
  return buildUrlPure(siteConfig.host, locale, ...pathSegments);
}

export function computeStoryCanonicalUrl(
  story: { is_managed: boolean; properties: Record<string, unknown> | null; locale_code?: string; slug: string | null },
  viewerLocale: string,
  ...pathPrefix: string[]
): string {
  return computeStoryCanonicalUrlPure(siteConfig.host, story, viewerLocale, ...pathPrefix);
}

export function truncateDescription(text: string | null | undefined, maxLength = 160): string {
  return truncateDescriptionPure(text, siteConfig.description, maxLength);
}
