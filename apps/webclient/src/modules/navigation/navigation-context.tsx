import * as React from "react";
import type { Locale, SupportedLocaleCode } from "@/config";

export interface NavigationState {
  /**
   * Current locale code (e.g., "tr", "en", "de")
   */
  locale: SupportedLocaleCode;

  /**
   * Locale data (name, flag, dir, etc.)
   */
  localeData: Locale;

  /**
   * Whether the current request is from a custom domain
   */
  isCustomDomain: boolean;

  /**
   * The custom domain host (e.g., "eser.aya.is"), null if not custom domain
   */
  customDomainHost: string | null;

  /**
   * The profile slug for custom domains (determined by the domain), null otherwise
   */
  customDomainProfileSlug: string | null;

  /**
   * The profile title for custom domains, null otherwise
   */
  customDomainProfileTitle: string | null;
}

const NavigationContext = React.createContext<NavigationState | undefined>(
  undefined,
);

interface NavigationProviderProps {
  state: NavigationState;
  children: React.ReactNode;
}

export function NavigationProvider(props: NavigationProviderProps) {
  return (
    <NavigationContext.Provider value={props.state}>
      {props.children}
    </NavigationContext.Provider>
  );
}

export function useNavigation(): NavigationState {
  const context = React.useContext(NavigationContext);

  if (context === undefined) {
    throw new Error("useNavigation must be used within a NavigationProvider");
  }

  return context;
}

/**
 * Hook to get just the current locale
 */
export function useLocale(): SupportedLocaleCode {
  const { locale } = useNavigation();
  return locale;
}

/**
 * Hook to check if we're on a custom domain
 */
export function useIsCustomDomain(): boolean {
  const { isCustomDomain } = useNavigation();
  return isCustomDomain;
}
