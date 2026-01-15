import type { NextRequest } from "next/server";
import Negotiator from "negotiator";
import picoMatch from "picomatch/posix.js";

export type Locale = {
  code: string;
  matches: string[];
  name: string;
  flag: string;
  dir: "ltr" | "rtl";
};

export function localeMatchFromRequest(
  req: NextRequest,
  supportedLocales: readonly Locale[],
  defaultLocaleCode: string,
) {
  const localeQueryParam = req.nextUrl.searchParams.get("locale");

  if (localeQueryParam !== null) {
    return {
      origin: "queryParam",
      localeCode: localeMatch([localeQueryParam], supportedLocales)?.code ?? defaultLocaleCode,
      updateCookie: true,
    };
  }

  const localeCookie = req.cookies.get("SITE_LOCALE");

  if (localeCookie !== undefined) {
    return {
      origin: "cookie",
      localeCode: localeMatch([localeCookie.value], supportedLocales)?.code ?? defaultLocaleCode,
      updateCookie: false,
    };
  }

  const localeHeader = req.headers.get("accept-language");

  if (localeHeader !== null) {
    const negotiatorRequest = {
      headers: {
        "accept-language": localeHeader,
      },
    };

    const negotiator = new Negotiator(negotiatorRequest);
    const locales = negotiator.languages();

    return {
      origin: "header",
      localeCode: localeMatch(locales, supportedLocales)?.code ?? defaultLocaleCode,
      updateCookie: true,
    };
  }

  return {
    origin: "default",
    localeCode: defaultLocaleCode,
    updateCookie: true,
  };
}

export function localeMatch(locales: readonly string[], supportedLocales: readonly Locale[]) {
  for (const locale of locales) {
    for (const supportedLocale of supportedLocales) {
      if (picoMatch.isMatch(locale, supportedLocale.matches)) {
        return supportedLocale;
      }
    }
  }

  return undefined;
}
