import * as React from "react";
import { getSessionCurrent } from "@/modules/backend/backend";
import type { SessionPreferences } from "@/modules/backend/types";
import type { AccessibleProfile } from "@/modules/backend/backend";
import { getCurrentLanguage } from "@/modules/i18n/i18n";
import { getBackendUri } from "@/config";
import {
  clearAuthData,
  isTokenExpiringSoon,
  refreshTokenRequest,
  setAuthToken,
} from "@/modules/backend/fetcher";

export type IndividualProfile = {
  id: string;
  slug: string;
  kind: string;
  title: string;
  profile_picture_uri?: string | null;
};

export type User = {
  id: string;
  kind: string;
  name: string;
  email?: string;
  github_handle?: string;
  individual_profile_id?: string;
  individual_profile_slug?: string;
  individual_profile?: IndividualProfile;
  accessible_profiles?: AccessibleProfile[];
};

type AuthContextValue = {
  isLoading: boolean;
  isAuthenticated: boolean;
  user: User | null;
  token: string | null;
  preferences: SessionPreferences | null;
  login: (redirectUri?: string) => void;
  logout: () => Promise<void>;
  refreshAuth: () => Promise<void>;
};

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
    preferences: SessionPreferences | null;
  }>({
    isLoading: true,
    isAuthenticated: false,
    user: null,
    token: null,
    preferences: null,
  });

  const initAuth = React.useCallback(async () => {
    // SSR check
    if (typeof window === "undefined" || !globalThis.localStorage) {
      setState((prev) => ({ ...prev, isLoading: false }));
      return;
    }

    const locale = getCurrentLanguage();

    // Single consolidated call: validates session cookie, returns auth + preferences
    const result = await getSessionCurrent(locale);

    if (result.authenticated && result.token !== undefined) {
      // User is authenticated - update local storage and state
      setAuthToken(
        result.token,
        result.expires_at ?? Date.now() + 24 * 60 * 60 * 1000,
      );

      // Combine user data with profile slug and accessible profiles
      const individualProfile: IndividualProfile | undefined =
        result.selected_profile !== undefined && result.selected_profile !== null
          ? {
            id: result.selected_profile.id,
            slug: result.selected_profile.slug,
            kind: result.selected_profile.kind,
            title: result.selected_profile.title,
            profile_picture_uri: result.selected_profile.profile_picture_uri,
          }
          : undefined;

      const userData: User | null =
        result.user !== undefined && result.user !== null
          ? {
            ...result.user,
            individual_profile_slug: result.selected_profile?.slug,
            individual_profile: individualProfile,
            accessible_profiles: result.accessible_profiles,
          }
          : null;

      localStorage.setItem(
        "auth_session",
        JSON.stringify({ user: userData }),
      );
      setState({
        isLoading: false,
        isAuthenticated: true,
        user: userData,
        token: result.token,
        preferences: result.preferences ?? null,
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
        preferences: result.preferences ?? null,
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

    // Build the callback URL using URL API to avoid double-encoding
    const origin = globalThis.location !== undefined
      ? globalThis.location.origin
      : "";
    const callbackUrlObj = new URL(`${origin}/auth/callback`);
    callbackUrlObj.searchParams.set("redirect", redirectUri ?? `/${locale}`);

    const loginUrlObj = new URL(`${backendUri}/${locale}/auth/github/login`);
    loginUrlObj.searchParams.set("redirect_uri", callbackUrlObj.toString());

    if (globalThis.location !== undefined) {
      globalThis.location.href = loginUrlObj.toString();
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
    try {
      await initAuth();
    } catch {
      // Ensure isLoading is reset even if initAuth fails unexpectedly
      setState((prev) => ({ ...prev, isLoading: false }));
    }
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
