import { createServerFn } from "@tanstack/react-start";
import { getCookie } from "@tanstack/react-start/server";
import { DEFAULT_LOCALE, isValidLocale, type SupportedLocaleCode } from "@/config";

/**
 * Server function to get the current locale from cookies
 * Falls back to DEFAULT_LOCALE if no valid cookie is set
 */
export const getLocaleFromCookie = createServerFn({ method: "GET" }).handler(
  async (): Promise<SupportedLocaleCode> => {
    const cookieLocale = getCookie("SITE_LOCALE");

    if (cookieLocale && isValidLocale(cookieLocale)) {
      return cookieLocale;
    }

    return DEFAULT_LOCALE;
  }
);
