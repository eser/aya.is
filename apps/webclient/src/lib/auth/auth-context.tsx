import * as React from "react";
import { checkSessionViaCookie } from "@/modules/backend/auth/session-check";
import { getCurrentLanguage } from "@/modules/i18n/i18n";
import { getBackendUri } from "@/config";
import {
  clearAuthData,
  isTokenExpiringSoon,
  refreshTokenRequest,
  setAuthToken,
} from "@/modules/backend/fetcher";

interface User {
  id: string;
  name: string;
  email?: string;
  github_handle?: string;
  individual_profile_id?: string;
}

interface AuthContextValue {
  isLoading: boolean;
  isAuthenticated: boolean;
  user: User | null;
  token: string | null;
  login: (redirectUri?: string) => void;
  logout: () => Promise<void>;
  refreshAuth: () => Promise<void>;
}

const AuthContext = React.createContext<AuthContextValue | null>(null);

type AuthProviderProps = {
  children: React.ReactNode;
};

export function AuthProvider(props: AuthProviderProps) {
  const [state, setState] = React.useState<{
    isLoading: boolean;
    isAuthenticated: boolean;
    user: User | null;
    token: string | null;
  }>({
    isLoading: true,
    isAuthenticated: false,
    user: null,
    token: null,
  });

  const initAuth = React.useCallback(async () => {
    // SSR check
    if (typeof window === "undefined" || !globalThis.localStorage) {
      setState((prev) => ({ ...prev, isLoading: false }));
      return;
    }

    const locale = getCurrentLanguage();

    // Always validate session via cookie-based check (works cross-domain)
    // This is the source of truth for authentication state
    const result = await checkSessionViaCookie(locale);

    if (result.authenticated && result.token !== undefined) {
      // User is authenticated - update local storage and state
      setAuthToken(
        result.token,
        result.expires_at ?? Date.now() + 24 * 60 * 60 * 1000,
      );
      localStorage.setItem(
        "auth_session",
        JSON.stringify({ user: result.user }),
      );
      setState({
        isLoading: false,
        isAuthenticated: true,
        user: result.user || null,
        token: result.token,
      });
    } else {
      // User is NOT authenticated - clear any stale local data
      // This ensures logout on one domain propagates to all domains
      clearAuthData();
      setState({
        isLoading: false,
        isAuthenticated: false,
        user: null,
        token: null,
      });
    }
  }, []);

  React.useEffect(() => {
    initAuth();
  }, [initAuth]);

  // Auto-refresh token before expiration
  React.useEffect(() => {
    if (!state.isAuthenticated || state.token === null) {
      return;
    }

    const checkAndRefresh = async () => {
      if (isTokenExpiringSoon()) {
        const locale = getCurrentLanguage();
        const result = await refreshTokenRequest(locale);
        if (result !== null) {
          setState((prev) => ({
            ...prev,
            token: result.token,
          }));
        }
      }
    };

    // Check every minute
    const interval = setInterval(checkAndRefresh, 60 * 1000);
    return () => clearInterval(interval);
  }, [state.isAuthenticated, state.token]);

  const login = React.useCallback((redirectUri?: string) => {
    const locale = getCurrentLanguage();
    const backendUri = getBackendUri();

    // Build the callback URL - use locale-aware path with optional redirect
    const callbackUrl = typeof window !== "undefined"
      ? `${globalThis.location.origin}/${locale}/auth/callback`
      : `/${locale}/auth/callback`;

    // If a redirectUri was provided, append it as a query param to the callback
    const finalCallbackUrl = redirectUri !== undefined
      ? `${callbackUrl}?redirect=${encodeURIComponent(redirectUri)}`
      : callbackUrl;

    const loginUrl = `${backendUri}/${locale}/auth/github/login?redirect_uri=${encodeURIComponent(finalCallbackUrl)}`;

    if (typeof window !== "undefined") {
      globalThis.location.href = loginUrl;
    }
  }, []);

  const logout = React.useCallback(async () => {
    const locale = getCurrentLanguage();
    const backendUri = getBackendUri();

    // Call backend logout endpoint
    try {
      await fetch(`${backendUri}/${locale}/auth/logout`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          ...(state.token !== null ? { Authorization: `Bearer ${state.token}` } : {}),
        },
      });
    } catch {
      // Ignore logout errors
    }

    // Clear local storage using centralized function
    clearAuthData();

    // Navigate to home - force page reload to clear state
    if (typeof window !== "undefined") {
      globalThis.location.href = `/${locale}`;
    }
  }, [state.token]);

  const refreshAuth = React.useCallback(async () => {
    setState((prev) => ({ ...prev, isLoading: true }));
    await initAuth();
  }, [initAuth]);

  const value = React.useMemo(
    () => ({
      ...state,
      login,
      logout,
      refreshAuth,
    }),
    [state, login, logout, refreshAuth],
  );

  return <AuthContext.Provider value={value}>{props.children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = React.useContext(AuthContext);
  if (ctx === null) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
