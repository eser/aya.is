import {
  DEFAULT_LOCALE,
  isValidLocale,
  SUPPORTED_LOCALES,
  type SupportedLocaleCode,
} from "@/config";

interface UrlOptions {
  locale?: string;
  isCustomDomain?: boolean;
  currentLocale?: string;
}

/**
 * Build URL with locale prefix
 *
 * With the new URL-based locale system, locale is ALWAYS in the URL:
 * - Main domain: /tr/page, /en/page, etc.
 * - Custom domain: /tr/page, /en/page (but no profile slug since it's implicit)
 */
export function localizedUrl(path: string, options: UrlOptions = {}): string {
  const { locale, isCustomDomain: _isCustomDomain, currentLocale } = options;
  const targetLocale = locale ?? currentLocale ?? DEFAULT_LOCALE;

  // Ensure path starts with /
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;

  // Always add locale prefix
  const cleanPath = normalizedPath === "/" ? "" : normalizedPath;
  return `/${targetLocale}${cleanPath}`;
}

/**
 * Parse locale from URL path
 *
 * Returns the locale if found in the first path segment, otherwise null
 */
export function parseLocaleFromPath(pathname: string): {
  locale: SupportedLocaleCode | null;
  restPath: string;
} {
  const segments = pathname.split("/").filter(Boolean);
  const firstSegment = segments[0];

  if (firstSegment && isValidLocale(firstSegment)) {
    return {
      locale: firstSegment,
      restPath: "/" + (segments.slice(1).join("/") || ""),
    };
  }

  return { locale: null, restPath: pathname };
}

/**
 * Build alternate URLs for all supported locales (for SEO hreflang)
 */
export function getAlternateUrls(
  basePath: string,
  isCustomDomain: boolean = false,
): Array<{ locale: SupportedLocaleCode; url: string }> {
  return SUPPORTED_LOCALES.map((locale: SupportedLocaleCode) => ({
    locale,
    url: localizedUrl(basePath, { locale, isCustomDomain }),
  }));
}
