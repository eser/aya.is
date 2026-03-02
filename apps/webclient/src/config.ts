// Site Configuration

// Re-export pure locale utilities (no import.meta.env dependency)
export {
  DEFAULT_LOCALE,
  FALLBACK_LOCALE,
  getLocaleData,
  isAllowedURI,
  isValidLocale,
  SUPPORTED_LOCALES,
  supportedLocales,
} from "@/lib/locale-utils.ts";
export type { Locale, SupportedLocaleCode } from "@/lib/locale-utils.ts";

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

  authGithubEnabled: boolean;
  authAppleEnabled: boolean;
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

  authGithubEnabled: import.meta.env.VITE_AUTH_GITHUB_ENABLED !== "false",
  authAppleEnabled: import.meta.env.VITE_AUTH_APPLE_ENABLED === "true",
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

export const predefinedSlugs: readonly string[] = [
  "auth",
  "mailbox",
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
    import.meta.env.VITE_ALLOWED_URI_PREFIXES_PROFILES ??
      "https://objects.aya.is/,https://avatars.githubusercontent.com/",
  ].flatMap((s) => s.split(",").map((p) => p.trim()).filter((p) => p !== "")),
};
