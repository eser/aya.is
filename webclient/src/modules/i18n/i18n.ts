import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import resourcesToBackend from "i18next-resources-to-backend";

import {
  SUPPORTED_LOCALES,
  FALLBACK_LOCALE,
  DEFAULT_LOCALE,
} from "@/config";

// Dynamic import for translation files
const loadResources = resourcesToBackend(
  (language: string) => import(`@/messages/${language}.json`)
);

i18n
  .use(loadResources)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: FALLBACK_LOCALE,
    supportedLngs: [...SUPPORTED_LOCALES],
    debug: process.env.NODE_ENV === "development",
    lng: DEFAULT_LOCALE, // Default language

    interpolation: {
      escapeValue: false, // React already escapes values
    },

    detection: {
      // Order of detection methods
      order: ["path", "cookie", "navigator"],
      // Look for locale in URL path (e.g., /en/page)
      lookupFromPathIndex: 0,
      // Cookie settings
      lookupCookie: "SITE_LOCALE",
      caches: ["cookie"],
      cookieMinutes: 60 * 24 * 365, // 1 year
    },

    react: {
      useSuspense: true,
    },
  });

export default i18n;

// Helper to change language and update cookie
export async function changeLanguage(locale: string): Promise<void> {
  if (SUPPORTED_LOCALES.includes(locale as (typeof SUPPORTED_LOCALES)[number])) {
    await i18n.changeLanguage(locale);
    // Set cookie for persistence
    if (typeof document !== "undefined") {
      document.cookie = `SITE_LOCALE=${locale};path=/;max-age=${60 * 60 * 24 * 365}`;
    }
  }
}

// Get current language
export function getCurrentLanguage(): string {
  return i18n.language || DEFAULT_LOCALE;
}
