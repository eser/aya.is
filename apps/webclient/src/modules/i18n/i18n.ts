import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import resourcesToBackend from "i18next-resources-to-backend";

import { DEFAULT_LOCALE, FALLBACK_LOCALE, SUPPORTED_LOCALES } from "@/config";

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

// Helper to change language and update cookie
export async function changeLanguage(locale: string): Promise<void> {
  if (
    SUPPORTED_LOCALES.includes(locale as (typeof SUPPORTED_LOCALES)[number])
  ) {
    await i18n.changeLanguage(locale);
    // Set cookie for persistence
    if (typeof document !== "undefined") {
      document.cookie = `SITE_LOCALE=${locale};path=/;max-age=${60 * 60 * 24 * 365}`;
    }
  }
}

// Get current language
export function getCurrentLanguage(): string {
  return i18n.language ?? DEFAULT_LOCALE;
}
