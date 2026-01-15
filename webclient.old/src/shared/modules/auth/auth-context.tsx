"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { getBackendUri } from "@/shared/config.ts";
import { clearAuthData, getTokenExpirationTime, refreshToken } from "./token-refresh.ts";
import { useNavigationClient } from "@/shared/modules/navigation/use-navigation-client.tsx";
import { getCurrentSession } from "@/shared/modules/backend/sessions/get-current-session.ts";
import type { Session } from "@/shared/modules/backend/sessions/types.ts";

interface AuthContextType {
  session: Session | null;
  token: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  authError: string | null;
  login: (redirectUri?: string) => void;
  logout: () => Promise<void>;
  refreshAuth: () => void;
  clearError: () => void;
}

const AuthContext = React.createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const { state } = useNavigationClient();
  const [session, setSession] = React.useState<Session | null>(null);
  const [token, setToken] = React.useState<string | null>(null);
  const [isLoading, setIsLoading] = React.useState(true);
  const [authError, setAuthError] = React.useState<string | null>(null);

  // Load auth state from localStorage on mount
  React.useEffect(() => {
    const loadAuthState = () => {
      if (globalThis.localStorage === undefined) {
        setIsLoading(false);
        return;
      }

      try {
        const storedToken = localStorage.getItem("auth_token");
        const storedSession = localStorage.getItem("auth_session");

        if (storedToken && storedSession) {
          setToken(storedToken);
          setSession(JSON.parse(storedSession));
        }
      } catch (error) {
        console.error("Failed to load auth state:", error);
        // Clear invalid data
        clearAuthData();
      } finally {
        setIsLoading(false);
      }
    };

    loadAuthState();
  }, []);

  // Check for auth_token in query string (OAuth callback)
  React.useEffect(() => {
    const handleAuthTokenFromQuery = async () => {
      if (globalThis.location === undefined) {
        return;
      }

      const urlParams = new URLSearchParams(globalThis.location.search);
      const authToken = urlParams.get("auth_token");

      if (!authToken) {
        return;
      }

      try {
        // Store token in localStorage
        localStorage.setItem("auth_token", authToken);
        setToken(authToken);

        // Parse JWT to get token expiration
        const tokenParts = authToken.split(".");
        if (tokenParts.length === 3) {
          const payload = JSON.parse(atob(tokenParts[1]));

          // Fetch session details from authenticated backend endpoint
          const locale = state.locale.code;
          const sessionPayload = await getCurrentSession(locale);

          if (sessionPayload !== null) {
            localStorage.setItem("auth_session", JSON.stringify(sessionPayload));
            setSession(sessionPayload);
            console.log("✅ Session state updated with:", sessionPayload);

            // Store token expiration
            if (payload.exp) {
              localStorage.setItem("auth_token_expires_at", String(payload.exp * 1000));
            }
          } else {
            console.log("❌ No session payload received");
          }
        }

        // Clean URL by removing auth_token from query string
        urlParams.delete("auth_token");
        const newUrl = urlParams.toString()
          ? `${globalThis.location.pathname}?${urlParams.toString()}`
          : globalThis.location.pathname;
        globalThis.history.replaceState({}, "", newUrl);
      } catch (error) {
        console.error("Failed to handle auth token from query:", error);
        setAuthError("Failed to complete authentication. Please try again.");
      }
    };

    handleAuthTokenFromQuery();
  }, [state.locale.code]);

  // Set up automatic token refresh
  React.useEffect(() => {
    if (token === null) {
      return;
    }

    const checkAndRefresh = async () => {
      const expirationTime = getTokenExpirationTime();
      if (expirationTime === null) {
        return 0;
      }

      const now = Date.now();
      const timeUntilExpiry = expirationTime - now;

      // Refresh 5 minutes before expiration
      const refreshTime = Math.max(0, timeUntilExpiry - (5 * 60 * 1000));

      if (refreshTime === 0) {
        // Token needs immediate refresh
        const newToken = await refreshToken();

        if (newToken) {
          setToken(newToken);

          return 0;
        }

        // Failed to refresh, clear auth state
        clearAuthData();
        setSession(null);
        setToken(null);
        setAuthError("Your session has expired. Please log in again.");

        return 0;
      }

      return refreshTime;
    };

    let timeoutId: ReturnType<typeof setTimeout>;

    const setupRefreshTimer = async () => {
      const refreshTime = await checkAndRefresh();

      if (refreshTime && refreshTime > 0) {
        timeoutId = setTimeout(async () => {
          const newToken = await refreshToken();

          if (newToken !== null) {
            setToken(newToken);
            // Set up next refresh
            setupRefreshTimer();

            return;
          }

          // Failed to refresh, clear auth state
          clearAuthData();
          setSession(null);
          setToken(null);
          setAuthError("Your session has expired. Please log in again.");
        }, refreshTime);
      }
    };

    setupRefreshTimer();

    return () => {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    };
  }, [token]);

  const login = React.useCallback((redirectUri?: string) => {
    const backendUri = getBackendUri();
    const currentUrl = globalThis.location?.href ?? "";
    const finalRedirectUri = redirectUri || currentUrl;
    const locale = state.locale.code;

    // Redirect to backend OAuth endpoint
    const loginUrl = `${backendUri}/${locale}/auth/github/login${
      finalRedirectUri ? `?redirect_uri=${encodeURIComponent(finalRedirectUri)}` : ""
    }`;
    if (globalThis.location !== undefined) {
      globalThis.location.href = loginUrl;
    }
  }, [state.locale.code]);

  const logout = React.useCallback(async () => {
    try {
      // Call backend logout endpoint
      const backendUri = getBackendUri();
      const locale = state.locale.code;
      await fetch(`${backendUri}/${locale}/auth/logout`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
      });
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      // Clear local state regardless of backend response
      clearAuthData();
      setSession(null);
      setToken(null);
      router.push("/");
    }
  }, [token, router, state.locale.code]);

  const refreshAuth = React.useCallback(() => {
    // Re-load auth state from localStorage
    const storedToken = localStorage.getItem("auth_token");
    const storedSession = localStorage.getItem("auth_session");

    if (storedToken && storedSession) {
      try {
        setToken(storedToken);
        setSession(JSON.parse(storedSession));
      } catch (error) {
        console.error("Failed to refresh auth:", error);
      }
    }
  }, []);

  const clearError = React.useCallback(() => {
    setAuthError(null);
  }, []);

  const value = React.useMemo(
    () => ({
      session,
      token,
      isLoading,
      isAuthenticated: !!session && !!token,
      authError,
      login,
      logout,
      refreshAuth,
      clearError,
    }),
    [session, token, isLoading, authError, login, logout, refreshAuth, clearError],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = React.useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
