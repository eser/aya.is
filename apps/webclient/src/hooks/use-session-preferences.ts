import { useCallback, useEffect, useState } from "react";
import {
  createSession,
  getSessionPreferences,
  updateSessionPreferences,
} from "@/modules/backend/sessions";
import type { SessionPreferences } from "@/modules/backend/sessions/types.ts";

type Theme = "dark" | "light" | "system";

interface UseSessionPreferencesOptions {
  locale: string;
  defaultTheme?: Theme;
  defaultLocale?: string;
}

interface UseSessionPreferencesResult {
  preferences: SessionPreferences;
  isLoading: boolean;
  isSessionReady: boolean;
  setTheme: (theme: Theme) => Promise<void>;
  setLocale: (locale: string) => Promise<void>;
  syncFromServer: () => Promise<void>;
}

const LOCAL_STORAGE_THEME_KEY = "vite-ui-theme";
const LOCAL_STORAGE_LOCALE_KEY = "preferred-locale";

/**
 * Hook for managing session preferences with server-side sync
 *
 * On mount:
 * 1. Load from localStorage for immediate display
 * 2. Try to fetch from server (if session exists)
 * 3. If server returns preferences, update local state
 *
 * On change:
 * 1. Update local state immediately
 * 2. Save to localStorage for persistence
 * 3. Sync to server (create session if needed)
 */
export function useSessionPreferences(
  options: UseSessionPreferencesOptions,
): UseSessionPreferencesResult {
  const { locale, defaultTheme = "system" } = options;

  const [preferences, setPreferences] = useState<SessionPreferences>(() => {
    // Initialize from localStorage on client
    if (typeof globalThis.localStorage !== "undefined") {
      const savedTheme =
        globalThis.localStorage.getItem(LOCAL_STORAGE_THEME_KEY) as
          | Theme
          | null;
      const savedLocale = globalThis.localStorage.getItem(
        LOCAL_STORAGE_LOCALE_KEY,
      );
      return {
        theme: savedTheme ?? defaultTheme,
        locale: savedLocale ?? locale,
      };
    }
    return { theme: defaultTheme, locale };
  });

  const [isLoading, setIsLoading] = useState(true);
  const [isSessionReady, setIsSessionReady] = useState(false);
  const [hasSession, setHasSession] = useState(false);

  // Sync preferences from server
  const syncFromServer = useCallback(async () => {
    try {
      const serverPrefs = await getSessionPreferences(locale);
      if (serverPrefs !== null) {
        setHasSession(true);
        setPreferences((prev) => ({
          ...prev,
          ...serverPrefs,
        }));
        // Update localStorage to match server
        if (
          typeof globalThis.localStorage !== "undefined" &&
          serverPrefs.theme !== undefined
        ) {
          globalThis.localStorage.setItem(
            LOCAL_STORAGE_THEME_KEY,
            serverPrefs.theme,
          );
        }
        if (
          typeof globalThis.localStorage !== "undefined" &&
          serverPrefs.locale !== undefined
        ) {
          globalThis.localStorage.setItem(
            LOCAL_STORAGE_LOCALE_KEY,
            serverPrefs.locale,
          );
        }
      }
    } catch (error) {
      console.error("[SessionPreferences] Failed to sync from server:", error);
    }
  }, [locale]);

  // Initial load from server
  useEffect(() => {
    let mounted = true;

    async function loadFromServer() {
      setIsLoading(true);
      try {
        const serverPrefs = await getSessionPreferences(locale);
        if (mounted && serverPrefs !== null) {
          setHasSession(true);
          setPreferences((prev) => ({
            ...prev,
            ...serverPrefs,
          }));
          // Update localStorage to match server
          if (
            typeof globalThis.localStorage !== "undefined" &&
            serverPrefs.theme !== undefined
          ) {
            globalThis.localStorage.setItem(
              LOCAL_STORAGE_THEME_KEY,
              serverPrefs.theme,
            );
          }
        }
      } catch {
        // Session doesn't exist or error - use localStorage values
      } finally {
        if (mounted) {
          setIsLoading(false);
          setIsSessionReady(true);
        }
      }
    }

    // Only run on client
    if (typeof globalThis.localStorage !== "undefined") {
      loadFromServer();
    } else {
      setIsLoading(false);
      setIsSessionReady(true);
    }

    return () => {
      mounted = false;
    };
  }, [locale]);

  // Set theme with server sync
  const setTheme = useCallback(
    async (theme: Theme) => {
      // Update local state immediately
      setPreferences((prev) => ({ ...prev, theme }));

      // Save to localStorage
      if (typeof globalThis.localStorage !== "undefined") {
        globalThis.localStorage.setItem(LOCAL_STORAGE_THEME_KEY, theme);
      }

      // Sync to server
      try {
        if (hasSession) {
          // Update existing session
          await updateSessionPreferences(locale, { theme });
        } else {
          // Create new session with preferences
          const result = await createSession(locale, { theme });
          if (result !== null) {
            setHasSession(true);
          }
        }
      } catch (error) {
        console.error(
          "[SessionPreferences] Failed to sync theme to server:",
          error,
        );
        // Don't fail - local preference is still saved
      }
    },
    [locale, hasSession],
  );

  // Set locale with server sync
  const setLocale = useCallback(
    async (newLocale: string) => {
      // Update local state immediately
      setPreferences((prev) => ({ ...prev, locale: newLocale }));

      // Save to localStorage
      if (typeof globalThis.localStorage !== "undefined") {
        globalThis.localStorage.setItem(LOCAL_STORAGE_LOCALE_KEY, newLocale);
      }

      // Sync to server
      try {
        if (hasSession) {
          // Update existing session
          await updateSessionPreferences(locale, { locale: newLocale });
        } else {
          // Create new session with preferences
          const result = await createSession(locale, { locale: newLocale });
          if (result !== null) {
            setHasSession(true);
          }
        }
      } catch (error) {
        console.error(
          "[SessionPreferences] Failed to sync locale to server:",
          error,
        );
        // Don't fail - local preference is still saved
      }
    },
    [locale, hasSession],
  );

  return {
    preferences,
    isLoading,
    isSessionReady,
    setTheme,
    setLocale,
    syncFromServer,
  };
}
