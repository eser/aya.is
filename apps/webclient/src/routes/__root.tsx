import type { ReactNode } from "react";
import { useEffect } from "react";
import { createRootRoute, HeadContent, Link, Outlet, Scripts, useRouterState } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { I18nextProvider, useTranslation } from "react-i18next";
import { ThemeProvider } from "@/components/theme-provider";
import { AuthProvider } from "@/lib/auth/auth-context";
import { Button } from "@/components/ui/button";
import { NavigationProvider, type NavigationState } from "@/modules/navigation/navigation-context";
import { EasterEgg, ResponsiveIndicator } from "@/components/page-layouts/default";
import {
  CUSTOM_DOMAIN_DEFAULT_LOCALE,
  DEFAULT_LOCALE,
  isValidLocale,
  type SupportedLocaleCode,
  supportedLocales,
} from "@/config";
import { parseLocaleFromPath } from "@/lib/url";
import i18n from "@/modules/i18n/i18n";
import "@/styles.css";

const themeScript = `
  (function() {
    const storageKey = 'vite-ui-theme';
    const theme = localStorage.getItem(storageKey);
    const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    const resolvedTheme = theme === 'system' || !theme ? systemTheme : theme;
    document.documentElement.classList.add(resolvedTheme);
  })();
`;

// Helper to detect custom domain from host
function isCustomDomain(host: string): boolean {
  // Main domain patterns
  const mainDomains = ["aya.is", "localhost:3000", "127.0.0.1:3000"];
  return !mainDomains.some(
    (domain) => host === domain || host === `www.${domain}`,
  );
}

// Helper to extract profile from custom domain
function getProfileFromCustomDomain(host: string): string | null {
  // Pattern: {profile}.aya.is or {profile}.localhost:3000
  const patterns = [/^([^.]+)\.aya\.is$/, /^([^.]+)\.localhost:3000$/];
  for (const pattern of patterns) {
    const match = host.match(pattern);
    if (match && match[1] !== "www" && match[1] !== "api") {
      return match[1];
    }
  }
  return null;
}

// Detect navigation state from URL and host
function detectNavigationState(
  pathname: string,
  host?: string,
): NavigationState {
  const isCustom = host ? isCustomDomain(host) : false;
  const customProfile = host ? getProfileFromCustomDomain(host) : null;

  // Parse locale from URL path
  const { locale: urlLocale } = parseLocaleFromPath(pathname);

  // Determine effective locale
  let localeCode: SupportedLocaleCode;
  if (urlLocale && isValidLocale(urlLocale)) {
    localeCode = urlLocale;
  } else {
    localeCode = isCustom ? CUSTOM_DOMAIN_DEFAULT_LOCALE : DEFAULT_LOCALE;
  }

  return {
    locale: localeCode,
    localeData: supportedLocales[localeCode],
    isCustomDomain: isCustom,
    customDomainHost: isCustom ? host || null : null,
    customDomainProfile: customProfile,
  };
}

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: "AYA - Acik Yazilim Agi" },
      {
        name: "description",
        content: "Gonullu gelistirilen yazilimlarla olusan bir yazilim vakfi",
      },
    ],
  }),
  component: RootComponent,
  notFoundComponent: NotFoundComponent,
});

function RootComponent() {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;

  // Get host from window on client side
  const host = typeof window !== "undefined" ? globalThis.location.host : undefined;

  // Detect navigation state from URL and host
  const navigationState = detectNavigationState(pathname, host);

  return (
    <I18nextProvider i18n={i18n}>
      <NavigationProvider state={navigationState}>
        <RootDocument
          locale={navigationState.locale}
          dir={navigationState.localeData.dir}
        >
          <AuthProvider>
            <ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
              <LocaleSynchronizer />
              <Outlet />
            </ThemeProvider>
          </AuthProvider>
        </RootDocument>
      </NavigationProvider>
    </I18nextProvider>
  );
}

// Component to sync i18n language with navigation state
function LocaleSynchronizer() {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;
  const { i18n: i18nInstance } = useTranslation();
  const host = typeof window !== "undefined" ? globalThis.location.host : undefined;

  useEffect(() => {
    const navigationState = detectNavigationState(pathname, host);
    if (i18nInstance.language !== navigationState.locale) {
      i18nInstance.changeLanguage(navigationState.locale);
    }
  }, [pathname, host, i18nInstance]);

  return null;
}

interface RootDocumentProps {
  children: ReactNode;
  locale?: string;
  dir?: "ltr" | "rtl";
}

function RootDocument({
  children,
  locale = DEFAULT_LOCALE,
  dir = "ltr",
}: Readonly<RootDocumentProps>) {
  return (
    <html lang={locale} dir={dir} suppressHydrationWarning>
      <head>
        <HeadContent />
        <script dangerouslySetInnerHTML={{ __html: themeScript }} />
      </head>
      <body className="min-h-screen bg-background font-sans antialiased">
        {children}
        <EasterEgg />
        <ResponsiveIndicator />
        <TanStackRouterDevtools position="bottom-right" />
        <Scripts />
      </body>
    </html>
  );
}

function NotFoundComponent() {
  const { t } = useTranslation();

  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-4">
      <h1 className="text-4xl font-bold">404</h1>
      <p className="text-muted-foreground">{t("Error.Page not found")}</p>
      <Button render={<Link to="/">{t("Layout.Go Home")}</Link>} />
    </div>
  );
}
