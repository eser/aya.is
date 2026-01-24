import { useCallback } from "react";
import {
  updateSessionPreferences,
} from "@/modules/backend/sessions";
import type { SessionPreferences } from "@/modules/backend/sessions/types.ts";
import { useAuth } from "@/lib/auth/auth-context";

type Theme = "dark" | "light" | "system";

interface UseSessionPreferencesOptions {
  locale: string;
  defaultTheme?: Theme;
}

interface UseSessionPreferencesResult {
  preferences: SessionPreferences;
  isLoading: boolean;
  setTheme: (theme: Theme) => Promise<void>;
  setLocale: (locale: string) => Promise<void>;
}

const LOCAL_STORAGE_THEME_KEY = "vite-ui-theme";
const LOCAL_STORAGE_LOCALE_KEY = "preferred-locale";

/**
 * Hook for managing session preferences with server-side sync.
 * Reads initial preferences from AuthProvider (consolidated GET /sessions/_current).
 * Writes preference changes via PATCH /sessions/_current.
 */
export function useSessionPreferences(
  options: UseSessionPreferencesOptions,
): UseSessionPreferencesResult {
  const { locale, defaultTheme = "system" } = options;
  const auth = useAuth();

  const preferences: SessionPreferences = {
    theme: auth.preferences?.theme ?? defaultTheme,
    locale: auth.preferences?.locale ?? locale,
  };

  // Set theme with server sync
  const setTheme = useCallback(
    async (theme: Theme) => {
      // Save to localStorage
      if (typeof globalThis.localStorage !== "undefined") {
        globalThis.localStorage.setItem(LOCAL_STORAGE_THEME_KEY, theme);
      }

      // Sync to server
      try {
        await updateSessionPreferences(locale, { theme });
      } catch (error) {
        console.error(
          "[SessionPreferences] Failed to sync theme to server:",
          error,
        );
      }
    },
    [locale],
  );

  // Set locale with server sync
  const setLocale = useCallback(
    async (newLocale: string) => {
      // Save to localStorage
      if (typeof globalThis.localStorage !== "undefined") {
        globalThis.localStorage.setItem(LOCAL_STORAGE_LOCALE_KEY, newLocale);
      }

      // Sync to server
      try {
        await updateSessionPreferences(locale, { locale: newLocale });
      } catch (error) {
        console.error(
          "[SessionPreferences] Failed to sync locale to server:",
          error,
        );
      }
    },
    [locale],
  );

  return {
    preferences,
    isLoading: auth.isLoading,
    setTheme,
    setLocale,
  };
}
