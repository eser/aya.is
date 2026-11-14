// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import resourcesToBackend from "i18next-resources-to-backend";
import { useNavigate } from "@tanstack/react-router";
import { normalizeSpaces } from "@/lib/date";
import { localizedUrl, parseLocaleFromPath } from "@/lib/url";

import { DEFAULT_LOCALE, FALLBACK_LOCALE, isValidLocale, SUPPORTED_LOCALES, supportedLocales } from "@/config";

// Helper function to format relative time using Intl.RelativeTimeFormat
function formatRelativeTime(date: Date, locale: string): string {
  const now = new Date();
  const diffInSeconds = Math.floor((date.getTime() - now.getTime()) / 1000);
  const absDiff = Math.abs(diffInSeconds);

  const rtf = new Intl.RelativeTimeFormat(locale, { numeric: "auto" });

  if (absDiff < 60) {
    return rtf.format(diffInSeconds, "second");
  }
  if (absDiff < 3600) {
    return rtf.format(Math.floor(diffInSeconds / 60), "minute");
  }
  if (absDiff < 86400) {
    return rtf.format(Math.floor(diffInSeconds / 3600), "hour");
  }
  if (absDiff < 2592000) {
    return rtf.format(Math.floor(diffInSeconds / 86400), "day");
  }
  if (absDiff < 31536000) {
    return rtf.format(Math.floor(diffInSeconds / 2592000), "month");
  }
  return rtf.format(Math.floor(diffInSeconds / 31536000), "year");
}

// Dynamic import for translation files
const loadResources = resourcesToBackend((language: string) => import(`@/messages/${language}.json`));

// Named date/time formatters, referenced from translation strings as `{{value, date-long}}`.
const dateFormatters: Record<string, (value: Date, locale: string) => string> = {
  "date-long": (value, locale) => value.toLocaleDateString(locale, { year: "numeric", month: "long", day: "numeric" }),
  "date-short": (value, locale) =>
    value.toLocaleDateString(locale, { year: "numeric", month: "short", day: "numeric" }),
  "date-only": (value, locale) => value.toLocaleDateString(locale),
  "datetime-long": (value, locale) =>
    value.toLocaleDateString(locale, {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }),
  "datetime-short": (value, locale) =>
    value.toLocaleDateString(locale, {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }),
  "month-year": (value, locale) => value.toLocaleDateString(locale, { year: "numeric", month: "long" }),
  "month-short": (value, locale) => value.toLocaleDateString(locale, { month: "short" }),
  relative: (value, locale) => formatRelativeTime(value, locale),
};

i18n
  .use(loadResources)
  .use(initReactI18next)
  .init({
    fallbackLng: FALLBACK_LOCALE,
    supportedLngs: [...SUPPORTED_LOCALES],
    debug: false,
    lng: DEFAULT_LOCALE, // Default language

    interpolation: {
      escapeValue: false, // React already escapes values
    },

    react: {
      useSuspense: true,
    },
  });

// i18next 26 removed `interpolation.format`; custom formats are registered on the
// formatter service by name instead.
for (const [name, formatDate] of Object.entries(dateFormatters)) {
  i18n.services.formatter?.add(name, (value, lng) => {
    if (value instanceof Date) {
      return normalizeSpaces(formatDate(value, lng ?? DEFAULT_LOCALE));
    }

    return String(value);
  });
}

export default i18n;

// Helper to change locale: updates i18n language, document dir/lang, and persistence cookie
export function changeLocale(locale: string, isCustomDomain: boolean, navigate: ReturnType<typeof useNavigate>): void {
  if (!isValidLocale(locale)) {
    return;
  }

  // await i18n.changeLanguage(locale);

  if (typeof document !== "undefined") {
    const localeData = supportedLocales[locale];
    document.documentElement.dir = localeData.dir;
    document.documentElement.lang = locale;
    document.cookie = `SITE_LOCALE=${locale};path=/;max-age=${60 * 60 * 24 * 365}`;
  }

  // Get the current path without locale prefix
  const { restPath } = parseLocaleFromPath(location.pathname);

  // Build new URL with the new locale
  const newPath = localizedUrl(restPath ?? "/", {
    locale: locale,
    isCustomDomain,
    currentLocale: locale,
  });

  // Navigate to the new localized URL
  navigate({ to: newPath });
}

// Get current locale
export function getCurrentLocale(): string {
  return i18n.language ?? DEFAULT_LOCALE;
}
