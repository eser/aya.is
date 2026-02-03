import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { updateSessionPreferences } from "@/modules/backend/backend";
import { useAuth } from "@/lib/auth/auth-context";

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

  const auth = useAuth();
  const appliedFromServer = useRef(false);

  const [theme, setThemeState] = useState<Theme>(() => {
    if (typeof globalThis.document === "undefined") {
      return defaultTheme;
    }
    return (localStorage.getItem(storageKey) as Theme) ?? defaultTheme;
  });

  // Sync theme from auth context preferences (loaded by AuthProvider)
  useEffect(() => {
    if (!enableServerSync || auth.isLoading || appliedFromServer.current) {
      return;
    }

    appliedFromServer.current = true;

    if (auth.preferences !== null && auth.preferences.theme !== undefined) {
      setThemeState(auth.preferences.theme);
      localStorage.setItem(storageKey, auth.preferences.theme);
    }
  }, [enableServerSync, auth.isLoading, auth.preferences, storageKey]);

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
          await updateSessionPreferences(locale, { theme: newTheme });
        } catch (error) {
          console.error("[ThemeProvider] Failed to sync theme to server:", error);
        }
      })();
    },
    [storageKey, enableServerSync, locale],
  );

  const value = {
    theme,
    setTheme,
    isLoading: enableServerSync && auth.isLoading,
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
