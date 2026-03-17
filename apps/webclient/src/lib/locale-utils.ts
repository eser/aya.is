// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Pure locale utilities — no import.meta.env dependency
// Extracted from config.ts for testability with Deno's test runner

export interface Locale {
  code: string;
  matches: string[];
  name: string;
  asciiName: string;
  englishName: string;
  flag: string;
  dir: "ltr" | "rtl";
}

export const SUPPORTED_LOCALES = [
  "en",
  "tr",
  "fr",
  "de",
  "es",
  "pt-PT",
  "it",
  "nl",
  "ja",
  "ko",
  "ru",
  "zh-CN",
  "ar",
] as const;

export type SupportedLocaleCode = (typeof SUPPORTED_LOCALES)[number];

export const DEFAULT_LOCALE: SupportedLocaleCode = "en"; // Default for main domain
export const FALLBACK_LOCALE: SupportedLocaleCode = "en"; // Fallback for missing translations

export const supportedLocales: Record<SupportedLocaleCode, Locale> = {
  tr: {
    code: "tr",
    matches: ["@(tr)?(-*)"],
    name: "Türkçe",
    asciiName: "Turkce",
    englishName: "Turkish",
    flag: "\u{1F1F9}\u{1F1F7}",
    dir: "ltr",
  },
  en: {
    code: "en",
    matches: ["@(en)?(-*)"],
    name: "English",
    asciiName: "English",
    englishName: "English",
    flag: "\u{1F1FA}\u{1F1F8}",
    dir: "ltr",
  },
  fr: {
    code: "fr",
    matches: ["@(fr)?(-*)"],
    name: "Français",
    asciiName: "Francais",
    englishName: "French",
    flag: "\u{1F1EB}\u{1F1F7}",
    dir: "ltr",
  },
  de: {
    code: "de",
    matches: ["@(de)?(-*)"],
    name: "Deutsch",
    asciiName: "Deutsch",
    englishName: "German",
    flag: "\u{1F1E9}\u{1F1EA}",
    dir: "ltr",
  },
  es: {
    code: "es",
    matches: ["@(es)?(-*)"],
    name: "Español",
    asciiName: "Espanol",
    englishName: "Spanish",
    flag: "\u{1F1EA}\u{1F1F8}",
    dir: "ltr",
  },
  "pt-PT": {
    code: "pt-PT",
    matches: ["@(pt)?(-PT)?(-*)"],
    name: "Português (Portugal)",
    asciiName: "Portugues (Portugal)",
    englishName: "Portuguese (Portugal)",
    flag: "\u{1F1F5}\u{1F1F9}",
    dir: "ltr",
  },
  it: {
    code: "it",
    matches: ["@(it)?(-*)"],
    name: "Italiano",
    asciiName: "Italiano",
    englishName: "Italian",
    flag: "\u{1F1EE}\u{1F1F9}",
    dir: "ltr",
  },
  nl: {
    code: "nl",
    matches: ["@(nl)?(-*)"],
    name: "Nederlands",
    asciiName: "Nederlands",
    englishName: "Dutch",
    flag: "\u{1F1F3}\u{1F1F1}",
    dir: "ltr",
  },
  ja: {
    code: "ja",
    matches: ["@(ja)?(-*)"],
    name: "\u65E5\u672C\u8A9E",
    asciiName: "Nihongo",
    englishName: "Japanese",
    flag: "\u{1F1EF}\u{1F1F5}",
    dir: "ltr",
  },
  ko: {
    code: "ko",
    matches: ["@(ko)?(-*)"],
    name: "\uD55C\uAD6D\uC5B4",
    asciiName: "Hangugeo",
    englishName: "Korean",
    flag: "\u{1F1F0}\u{1F1F7}",
    dir: "ltr",
  },
  ru: {
    code: "ru",
    matches: ["@(ru)?(-*)"],
    name: "\u0420\u0443\u0441\u0441\u043A\u0438\u0439",
    asciiName: "Russkiy",
    englishName: "Russian",
    flag: "\u{1F1F7}\u{1F1FA}",
    dir: "ltr",
  },
  "zh-CN": {
    code: "zh-CN",
    matches: ["@(zh)?(-CN)?(-*)"],
    name: "\u7B80\u4F53\u4E2D\u6587",
    asciiName: "Zhongwen",
    englishName: "Chinese (Simplified)",
    flag: "\u{1F1E8}\u{1F1F3}",
    dir: "ltr",
  },
  ar: {
    code: "ar",
    matches: ["@(ar)?(-*)"],
    name: "\u0627\u0644\u0639\u0631\u0628\u064A\u0629",
    asciiName: "Al-Arabiyyah",
    englishName: "Arabic",
    flag: "\u{1F1F8}\u{1F1E6}",
    dir: "rtl",
  },
};

export function isValidLocale(locale: string): locale is SupportedLocaleCode {
  return SUPPORTED_LOCALES.includes(locale as SupportedLocaleCode);
}

export function getLocaleData(locale: string): Locale | undefined {
  if (isValidLocale(locale)) {
    return supportedLocales[locale];
  }
  return undefined;
}

// Validate if a URI starts with one of the allowed prefixes
export function isAllowedURI(uri: string | null | undefined, prefixes: string[]): boolean {
  if (uri === null || uri === undefined || uri === "") {
    return true; // Empty/null URIs are allowed
  }
  if (prefixes.length === 0) {
    return true; // No restrictions if no prefixes configured
  }
  return prefixes.some((prefix) => uri.startsWith(prefix));
}
