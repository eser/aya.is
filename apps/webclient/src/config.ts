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
  name: "aya.is",
  fancyName:
    "\u{1D400}\u{1D418}\u{1D400} \u{1D402}\u{1D428}\u{1D426}\u{1D426}\u{1D42E}\u{1D427}\u{1D422}\u{1D42D}\u{1D432} \u{1D407}\u{1D42E}\u{1D41B}",
  title: "AYA Community Hub",
  description: "The content and collaboration platform of our open software community",
  keywords: [
    "AYA",
    "Community Hub",
    "Community Platform",
    "Community Site",
    "Open Source",
    "Open Data",
    "Content",
  ],

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

  // Content links: safe URL protocols for user-facing content (ButtonCTA, etc.)
  content: ["http://", "https://", "/", "#"],
};
