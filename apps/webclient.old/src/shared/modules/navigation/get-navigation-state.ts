import { headers } from "next/headers";

import { fallbackLocaleCode, supportedLocales } from "@/shared/config.ts";
import type { Locale } from "@/shared/modules/i18n/locales.ts";

export type NavigationState = {
  host: string | null;
  profile: string | null;
  locale: Locale;
};

export async function getNavigationState() {
  const headersList = await headers();

  const host = headersList.get("x-custom-domain-host");
  const profile = headersList.get("x-custom-domain-profile");
  const localeCode = headersList.get("x-locale") ?? fallbackLocaleCode;

  const selectedLocale: Locale = supportedLocales[localeCode];

  const navigationState: NavigationState = {
    host: host,
    profile: profile,
    locale: selectedLocale,
  };

  return navigationState;
}
