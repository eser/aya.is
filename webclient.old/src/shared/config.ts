// import * as flags from "flags/next";

import type { Locale } from "@/shared/modules/i18n/locales.ts";

import * as flags from "./fake-flags-sdk.ts";

export const getBackendUri = () => {
  if (globalThis.localStorage !== undefined) {
    const backendUriFromLocalStorage = localStorage.getItem("backendUri");

    if (backendUriFromLocalStorage !== null) {
      return backendUriFromLocalStorage;
    }
  }

  return siteConfig.backendUri;
};

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
  backendUriProduction: string;
  backendUri: string;

  features: {
    login: flags.Flag<boolean>;
  };
};

export const siteConfig: SiteConfig = {
  name: "AYA",
  fancyName: "𝕒𝕔𝕀𝕜 𝕪𝕒𝕫𝕀𝕃𝕀𝕞 𝕒𝕘𝕀",
  title: "AYA",
  description: "Gönüllü geliştirilen yazılımlarla oluşan bir yazılım vakfı",
  keywords: ["AYA", "Açık Yazılım Ağı", "Açık Kaynak", "Açık Veri"],

  links: {
    x: "https://twitter.com/acikyazilimagi",
    instagram: "https://www.instagram.com/acikyazilimagi/",
    github: "https://github.com/acikyazilimagi",
  },

  // deno-lint-ignore no-process-global
  environment: process.env.NODE_ENV ?? "development",
  // deno-lint-ignore no-process-global
  host: process.env.NEXT_PUBLIC_HOST ?? "https://aya.is",
  // deno-lint-ignore no-process-global
  backendUriProduction: process.env.BACKEND_URI_PRODUCTION ?? process.env.NEXT_PUBLIC_BACKEND_URI_PRODUCTION ??
    "https://api.aya.is",
  // deno-lint-ignore no-process-global
  backendUri: process.env.BACKEND_URI ?? process.env.NEXT_PUBLIC_BACKEND_URI ?? "https://api.aya.is",

  features: {
    login: flags.flag({
      key: "login",

      decide() {
        return true;
      },
    }),
  },
};

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

export const supportedLocales: Record<string, Locale> = {
  tr: {
    code: "tr",
    matches: ["@(tr)?(-*)"],
    name: "Türkçe",
    flag: "🇹🇷",
    dir: "ltr",
  },
  en: {
    code: "en",
    matches: ["@(en)?(-*)"],
    name: "English",
    flag: "🇺🇸",
    dir: "ltr",
  },
  fr: {
    code: "fr",
    matches: ["@(fr)?(-*)"],
    name: "Français",
    flag: "🇫🇷",
    dir: "ltr",
  },
  de: {
    code: "de",
    matches: ["@(de)?(-*)"],
    name: "Deutsch",
    flag: "🇩🇪",
    dir: "ltr",
  },
  es: {
    code: "es",
    matches: ["@(es)?(-*)"],
    name: "Español",
    flag: "🇪🇸",
    dir: "ltr",
  },
  "pt-PT": {
    code: "pt-PT",
    matches: ["@(pt)?(-PT)?(-*)"],
    name: "Português (Portugal)",
    flag: "🇵🇹",
    dir: "ltr",
  },
  it: {
    code: "it",
    matches: ["@(it)?(-*)"],
    name: "Italiano",
    flag: "🇮🇹",
    dir: "ltr",
  },
  nl: {
    code: "nl",
    matches: ["@(nl)?(-*)"],
    name: "Nederlands",
    flag: "🇳🇱",
    dir: "ltr",
  },
  ja: {
    code: "ja",
    matches: ["@(ja)?(-*)"],
    name: "日本語",
    flag: "🇯🇵",
    dir: "ltr",
  },
  ko: {
    code: "ko",
    matches: ["@(ko)?(-*)"],
    name: "한국어",
    flag: "🇰🇷",
    dir: "ltr",
  },
  ru: {
    code: "ru",
    matches: ["@(ru)?(-*)"],
    name: "Русский",
    flag: "🇷🇺",
    dir: "ltr",
  },
  "zh-CN": {
    code: "zh-CN",
    matches: ["@(zh)?(-CN)?(-*)"],
    name: "简体中文",
    flag: "🇨🇳",
    dir: "ltr",
  },
  ar: {
    code: "ar",
    matches: ["@(ar)?(-*)"],
    name: "العربية",
    flag: "🇸🇦",
    dir: "rtl",
  },
};

export type SupportedLocaleCode = keyof typeof supportedLocales;

export const fallbackLocaleCode: SupportedLocaleCode = "tr";
