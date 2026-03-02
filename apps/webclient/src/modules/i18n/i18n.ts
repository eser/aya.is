import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import resourcesToBackend from "i18next-resources-to-backend";
import { useNavigate } from "@tanstack/react-router";
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
const loadResources = resourcesToBackend(
  (language: string) => import(`@/messages/${language}.json`),
);

i18n
  .use(loadResources)
  .use(initReactI18next)
  .init({
    fallbackLng: FALLBACK_LOCALE,
    supportedLngs: [...SUPPORTED_LOCALES],
    debug: false,
    showSupportNotice: false,
    lng: DEFAULT_LOCALE, // Default language

    interpolation: {
      escapeValue: false, // React already escapes values
      // Custom format function for date/time formatting
      format: (value, format, lng) => {
        if (value instanceof Date) {
          const locale = lng ?? DEFAULT_LOCALE;

          switch (format) {
            case "date-long":
              return value.toLocaleDateString(locale, {
                year: "numeric",
                month: "long",
                day: "numeric",
              });
            case "date-short":
              return value.toLocaleDateString(locale, {
                year: "numeric",
                month: "short",
                day: "numeric",
              });
            case "date-only":
              return value.toLocaleDateString(locale);
            case "datetime-long":
              return value.toLocaleDateString(locale, {
                year: "numeric",
                month: "long",
                day: "numeric",
                hour: "2-digit",
                minute: "2-digit",
              });
            case "datetime-short":
              return value.toLocaleDateString(locale, {
                year: "numeric",
                month: "short",
                day: "numeric",
                hour: "2-digit",
                minute: "2-digit",
              });
            case "month-year":
              return value.toLocaleDateString(locale, {
                year: "numeric",
                month: "long",
              });
            case "month-short":
              return value.toLocaleDateString(locale, {
                month: "short",
              });
            case "relative":
              return formatRelativeTime(value, locale);
            default:
              return value.toLocaleDateString(locale);
          }
        }
        return String(value);
      },
    },

    react: {
      useSuspense: true,
    },
  });

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
