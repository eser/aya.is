import { headers } from "next/headers";

import { fallbackLocaleCode, supportedLocales } from "@/shared/config.ts";
import type { Locale } from "@/shared/modules/i18n/locales.ts";

export type TranslationsStateValues = {
  localeCode: string;
  locale: Locale;
  supportedLocales: Record<string, Locale>;
  messages: Record<string, Record<string, string>>;
};

export type TranslationsState = TranslationsStateValues & {
  t: (category: string, key: string) => string;
};

export async function getTranslations(): Promise<TranslationsState> {
  const headersList = await headers();

  const localeCode = headersList.get("x-locale") ?? fallbackLocaleCode;

  const locale: Locale = supportedLocales[localeCode];

  const messages = (await import(`../../../messages/${localeCode}.json`)).default;

  const t = (category: string, key: string) => {
    return messages[category][key];
  };

  return { localeCode, locale, supportedLocales, messages, t };
}
