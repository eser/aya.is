import { createServerFn } from "@tanstack/react-start";
import { getCookie, getRequestHeader } from "@tanstack/react-start/server";
import { DEFAULT_LOCALE, isValidLocale, SUPPORTED_LOCALES, type SupportedLocaleCode } from "@/config";

/**
 * Parses the Accept-Language header and returns the best matching supported locale.
 * Returns null if no supported locale matches.
 *
 * Example header: "en-US,en;q=0.9,tr;q=0.8,fr;q=0.7"
 */
function parseAcceptLanguage(
  header: string,
): SupportedLocaleCode | null {
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

    // Try exact match first (e.g. "pt-PT" matches "pt-PT")
    const exactMatch = SUPPORTED_LOCALES.find(
      (locale) => locale.toLowerCase() === lang,
    );
    if (exactMatch !== undefined) {
      return exactMatch;
    }

    // Try language-only prefix match (e.g. "en-US" → "en", "pt-BR" → skip if no "pt")
    const langPrefix = lang.split("-")[0];
    if (langPrefix !== undefined) {
      const prefixMatch = SUPPORTED_LOCALES.find(
        (locale) => locale.toLowerCase() === langPrefix,
      );
      if (prefixMatch !== undefined) {
        return prefixMatch;
      }
    }
  }

  return null;
}

/**
 * Server function to get the current locale from cookies.
 * Falls back to DEFAULT_LOCALE if no valid cookie is set.
 */
export const getLocaleFromCookie = createServerFn({ method: "GET" }).handler(
  (): SupportedLocaleCode => {
    const cookieLocale = getCookie("SITE_LOCALE");

    if (cookieLocale !== undefined && isValidLocale(cookieLocale)) {
      return cookieLocale;
    }

    return DEFAULT_LOCALE;
  },
);

/**
 * Server function to detect the user's preferred locale.
 *
 * Priority:
 * 1. SITE_LOCALE cookie (explicit user choice from locale switcher)
 * 2. Accept-Language header (browser language preference)
 * 3. DEFAULT_LOCALE (final fallback)
 */
export const getPreferredLocale = createServerFn({ method: "GET" }).handler(
  (): SupportedLocaleCode => {
    // 1. Check cookie first — explicit user choice takes priority
    const cookieLocale = getCookie("SITE_LOCALE");
    if (cookieLocale !== undefined && isValidLocale(cookieLocale)) {
      return cookieLocale;
    }

    // 2. Check Accept-Language header — browser preference
    const acceptLanguage = getRequestHeader("accept-language");
    if (acceptLanguage !== undefined) {
      const matched = parseAcceptLanguage(acceptLanguage);
      if (matched !== null) {
        return matched;
      }
    }

    // 3. Final fallback
    return DEFAULT_LOCALE;
  },
);
