// Site Configuration

export interface Locale {
  code: string;
  matches: string[];
  name: string;
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
};

export const siteConfig: SiteConfig = {
  name: "AYA",
  fancyName: "\u{1D552}\u{1D554}\u{0131}\u{1D55C} \u{1D56A}\u{1D552}\u{1D55D}\u{0131}\u{1D55D}\u{0131}\u{1D55E} \u{1D552}\u{1D558}\u{0131}",
  title: "AYA",
  description: "Gonullu gelistirilen yazilimlarla olusan bir yazilim vakfi",
  keywords: ["AYA", "Acik Yazilim Agi", "Acik Kaynak", "Acik Veri"],

  links: {
    x: "https://twitter.com/acikyazilimagi",
    instagram: "https://www.instagram.com/acikyazilimagi/",
    github: "https://github.com/acikyazilimagi",
  },

  environment: process.env.NODE_ENV ?? "development",
  host: process.env.PUBLIC_HOST ?? "https://aya.is",
  backendUri: process.env.BACKEND_URI ?? process.env.PUBLIC_BACKEND_URI ?? "https://api.aya.is",
};

export const getBackendUri = (): string => {
  if (typeof window !== "undefined" && window.localStorage) {
    const backendUriFromLocalStorage = localStorage.getItem("backendUri");
    if (backendUriFromLocalStorage !== null) {
      return backendUriFromLocalStorage;
    }
  }
  return siteConfig.backendUri;
};

// Locale Configuration
export const SUPPORTED_LOCALES = [
  "tr",
  "en",
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

export const DEFAULT_LOCALE: SupportedLocaleCode = "tr"; // Default for main domain
export const CUSTOM_DOMAIN_DEFAULT_LOCALE: SupportedLocaleCode = "en"; // Default for custom domains
export const FALLBACK_LOCALE: SupportedLocaleCode = "tr"; // Fallback for missing translations

export const supportedLocales: Record<SupportedLocaleCode, Locale> = {
  tr: {
    code: "tr",
    matches: ["@(tr)?(-*)"],
    name: "Turkce",
    flag: "\u{1F1F9}\u{1F1F7}",
    dir: "ltr",
  },
  en: {
    code: "en",
    matches: ["@(en)?(-*)"],
    name: "English",
    flag: "\u{1F1FA}\u{1F1F8}",
    dir: "ltr",
  },
  fr: {
    code: "fr",
    matches: ["@(fr)?(-*)"],
    name: "Francais",
    flag: "\u{1F1EB}\u{1F1F7}",
    dir: "ltr",
  },
  de: {
    code: "de",
    matches: ["@(de)?(-*)"],
    name: "Deutsch",
    flag: "\u{1F1E9}\u{1F1EA}",
    dir: "ltr",
  },
  es: {
    code: "es",
    matches: ["@(es)?(-*)"],
    name: "Espanol",
    flag: "\u{1F1EA}\u{1F1F8}",
    dir: "ltr",
  },
  "pt-PT": {
    code: "pt-PT",
    matches: ["@(pt)?(-PT)?(-*)"],
    name: "Portugues (Portugal)",
    flag: "\u{1F1F5}\u{1F1F9}",
    dir: "ltr",
  },
  it: {
    code: "it",
    matches: ["@(it)?(-*)"],
    name: "Italiano",
    flag: "\u{1F1EE}\u{1F1F9}",
    dir: "ltr",
  },
  nl: {
    code: "nl",
    matches: ["@(nl)?(-*)"],
    name: "Nederlands",
    flag: "\u{1F1F3}\u{1F1F1}",
    dir: "ltr",
  },
  ja: {
    code: "ja",
    matches: ["@(ja)?(-*)"],
    name: "\u65E5\u672C\u8A9E",
    flag: "\u{1F1EF}\u{1F1F5}",
    dir: "ltr",
  },
  ko: {
    code: "ko",
    matches: ["@(ko)?(-*)"],
    name: "\uD55C\uAD6D\uC5B4",
    flag: "\u{1F1F0}\u{1F1F7}",
    dir: "ltr",
  },
  ru: {
    code: "ru",
    matches: ["@(ru)?(-*)"],
    name: "\u0420\u0443\u0441\u0441\u043A\u0438\u0439",
    flag: "\u{1F1F7}\u{1F1FA}",
    dir: "ltr",
  },
  "zh-CN": {
    code: "zh-CN",
    matches: ["@(zh)?(-CN)?(-*)"],
    name: "\u7B80\u4F53\u4E2D\u6587",
    flag: "\u{1F1E8}\u{1F1F3}",
    dir: "ltr",
  },
  ar: {
    code: "ar",
    matches: ["@(ar)?(-*)"],
    name: "\u0627\u0644\u0639\u0631\u0628\u064A\u0629",
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

// Forbidden slugs (reserved for system routes)
export const forbiddenSlugs: readonly string[] = [
  "about",
  "admin",
  "api",
  "auth",
  "communities",
  "community",
  "config",
  "contact",
  "contributions",
  "dashboard",
  "element",
  "elements",
  "events",
  "faq",
  "feed",
  "guide",
  "help",
  "home",
  "impressum",
  "imprint",
  "jobs",
  "legal",
  "login",
  "logout",
  "news",
  "null",
  "organizations",
  "orgs",
  "people",
  "policies",
  "policy",
  "privacy",
  "product",
  "products",
  "profile",
  "profiles",
  "projects",
  "register",
  "root",
  "search",
  "services",
  "settings",
  "signin",
  "signout",
  "signup",
  "stories",
  "story",
  "support",
  "tag",
  "tags",
  "terms",
  "tos",
  "undefined",
  "user",
  "users",
  "verify",
  "wiki",
];
