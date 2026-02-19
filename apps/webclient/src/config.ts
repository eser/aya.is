// Site Configuration

export interface Locale {
  code: string;
  matches: string[];
  name: string;
  asciiName: string;
  englishName: string;
  flag: string;
  dir: "ltr" | "rtl";
}

export type SiteConfig = {
  name: string;
  fancyName: string;
  title: string;
  description: string;
  keywords: string[];

  links: {
    x: string;
    instagram: string;
    github: string;
  };

  environment: string;
  host: string;
  backendUri: string;
  telegramBotUsername: string;
};

export const siteConfig: SiteConfig = {
  name: "AYA",
  fancyName:
    "\u{1D552}\u{1D554}\u{0131}\u{1D55C} \u{1D56A}\u{1D552}\u{1D55D}\u{0131}\u{1D55D}\u{0131}\u{1D55E} \u{1D552}\u{1D558}\u{0131}",
  title: "AYA",
  description: "Gonullu gelistirilen yazilimlarla olusan bir yazilim vakfi",
  keywords: ["AYA", "Acik Yazilim Agi", "Acik Kaynak", "Acik Veri"],

  links: {
    x: "https://twitter.com/acikyazilimagi",
    instagram: "https://www.instagram.com/acikyazilimagi/",
    github: "https://github.com/acikyazilimagi",
  },

  // Vite provides import.meta.env.MODE ("development" | "production")
  environment: import.meta.env.MODE ?? "development",
  // Use VITE_ prefix for client-side variables (Vite bundles them)
  host: import.meta.env.VITE_HOST ?? "https://aya.is",
  backendUri: import.meta.env.VITE_BACKEND_URI ?? "https://api.aya.is",
  telegramBotUsername: import.meta.env.VITE_TELEGRAM_BOT_USERNAME ?? "aya_is_bot",
};

export const getBackendUri = (): string => {
  // Only allow localStorage override in development to prevent SSRF via XSS
  if (siteConfig.environment === "development" && typeof globalThis.localStorage !== "undefined") {
    const backendUriFromLocalStorage = globalThis.localStorage.getItem(
      "backendUri",
    );
    if (backendUriFromLocalStorage !== null) {
      return backendUriFromLocalStorage;
    }
  }
  return siteConfig.backendUri;
};

// Locale Configuration
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

export const predefinedSlugs: readonly string[] = [
  "auth"
];

// Allowed URI prefixes for uploads (must match server-side config)
// These are used to validate URIs for non-admin users
export const allowedURIPrefixes = {
  // Stories: only our upload service
  stories: [
    import.meta.env.VITE_ALLOWED_URI_PREFIXES_STORIES ?? "https://objects.aya.is/",
  ].flatMap((s) => s.split(",").map((p) => p.trim()).filter((p) => p !== "")),

  // Profiles: our upload service + GitHub avatars
  profiles: [
    import.meta.env.VITE_ALLOWED_URI_PREFIXES_PROFILES ?? "https://objects.aya.is/,https://avatars.githubusercontent.com/",
  ].flatMap((s) => s.split(",").map((p) => p.trim()).filter((p) => p !== "")),
};

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
