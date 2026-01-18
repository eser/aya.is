import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import {
  createSession,
  getSessionPreferences,
  updateSessionPreferences,
} from "@/modules/backend/sessions";

type Theme = "dark" | "light" | "system";

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
  locale?: string;
  enableServerSync?: boolean;
};

type ThemeProviderState = {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  isLoading: boolean;
};

const initialState: ThemeProviderState = {
  theme: "system",
  setTheme: () => null,
  isLoading: false,
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider(props: ThemeProviderProps) {
  const {
    children,
    defaultTheme = "system",
    storageKey = "vite-ui-theme",
    locale = "en",
    enableServerSync = false,
  } = props;

  const [theme, setThemeState] = useState<Theme>(() => {
    if (typeof globalThis.document === "undefined") {
      return defaultTheme;
    }
    return (localStorage.getItem(storageKey) as Theme) ?? defaultTheme;
  });

  const [isLoading, setIsLoading] = useState(enableServerSync);
  const [hasSession, setHasSession] = useState(false);

  // Sync from server on mount if enabled
  useEffect(() => {
    if (!enableServerSync) {
      return;
    }

    let mounted = true;

    async function syncFromServer() {
      try {
        const serverPrefs = await getSessionPreferences(locale);
        if (mounted && serverPrefs !== null) {
          setHasSession(true);
          if (serverPrefs.theme !== undefined) {
            setThemeState(serverPrefs.theme);
            localStorage.setItem(storageKey, serverPrefs.theme);
          }
        }
      } catch {
        // Session doesn't exist or error - use localStorage values
      } finally {
        if (mounted) {
          setIsLoading(false);
        }
      }
    }

    syncFromServer();

    return () => {
      mounted = false;
    };
  }, [enableServerSync, locale, storageKey]);

  // Apply theme to DOM
  useEffect(() => {
    const root = globalThis.document.documentElement;

    root.classList.remove("light", "dark");

    if (theme === "system") {
      const systemTheme = globalThis.matchMedia("(prefers-color-scheme: dark)")
          .matches
        ? "dark"
        : "light";

      root.classList.add(systemTheme);
      return;
    }

    root.classList.add(theme);
  }, [theme]);

  // Set theme with optional server sync
  const setTheme = useCallback(
    (newTheme: Theme) => {
      // Update local state immediately
      setThemeState(newTheme);
      localStorage.setItem(storageKey, newTheme);

      // Sync to server if enabled
      if (!enableServerSync) {
        return;
      }

      // Fire and forget - don't block UI
      (async () => {
        try {
          if (hasSession) {
            await updateSessionPreferences(locale, { theme: newTheme });
          } else {
            const result = await createSession(locale, { theme: newTheme });
            if (result !== null) {
              setHasSession(true);
            }
          }
        } catch (error) {
          console.error("[ThemeProvider] Failed to sync theme to server:", error);
        }
      })();
    },
    [storageKey, enableServerSync, hasSession, locale],
  );

  const value = {
    theme,
    setTheme,
    isLoading,
  };

  return (
    <ThemeProviderContext.Provider value={value}>
      {children}
    </ThemeProviderContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeProviderContext);

  if (context === undefined) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }

  return context;
}
