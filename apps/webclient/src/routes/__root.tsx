import type { ReactNode } from "react";
import { useEffect } from "react";
import { createRootRouteWithContext, HeadContent, Link, Outlet, Scripts, useRouterState } from "@tanstack/react-router";
import { TanStackDevtools } from "@tanstack/react-devtools";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { I18nextProvider, useTranslation } from "react-i18next";
import { ThemeProvider } from "@/components/theme-provider";
import { Toaster } from "@/components/ui/sonner";
import { AuthProvider } from "@/lib/auth/auth-context";
import { Button } from "@/components/ui/button";
import { NavigationProvider, type NavigationState } from "@/modules/navigation/navigation-context";
import { EasterEgg, ResponsiveIndicator } from "@/components/page-layouts/default";
import {
  DEFAULT_LOCALE,
  isValidLocale,
  siteConfig,
  type SupportedLocaleCode,
  supportedLocales,
} from "@/config";
import { generateMetaTags } from "@/lib/seo";
import { parseLocaleFromPath } from "@/lib/url";
import type { RequestContext } from "@/request-context";
import i18n from "@/modules/i18n/i18n";
import "@/styles.css";

type MyRouterContext = {
  requestContext: RequestContext | undefined;
};

// Detect navigation state from URL, host, and request context
function detectNavigationState(
  pathname: string,
  requestContext: RequestContext | undefined,
  host?: string,
): NavigationState {
  // Check if we have custom domain info from middleware context
  const domainConfig = requestContext?.domainConfiguration;
  const isCustom = domainConfig?.type === "custom-domain";
  const customProfileSlug = isCustom ? domainConfig.profileSlug : null;
  const customProfileTitle = isCustom ? domainConfig.profileTitle : null;

  // Parse locale from URL path
  const { locale: urlLocale } = parseLocaleFromPath(pathname);

  // Determine effective locale
  let localeCode: SupportedLocaleCode;
  if (urlLocale !== null && isValidLocale(urlLocale)) {
    localeCode = urlLocale;
  } else {
    localeCode = DEFAULT_LOCALE;
  }

  return {
    locale: localeCode,
    localeData: supportedLocales[localeCode],
    isCustomDomain: isCustom,
    customDomainHost: isCustom ? host ?? null : null,
    customDomainProfileSlug: customProfileSlug,
    customDomainProfileTitle: customProfileTitle,
  };
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  beforeLoad: ({ context }) => {
    // Pass the context from middleware
    return { requestContext: context.requestContext };
  },
  head: () => ({
    meta: generateMetaTags({
      title: siteConfig.name,
      description: siteConfig.description,
      url: siteConfig.host,
      type: "website",
    }),
  }),
  component: RootComponent,
  notFoundComponent: NotFoundComponent,
});

function RootComponent() {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;
  const { requestContext } = Route.useRouteContext();

  // Get host from window on client side
  const host = typeof window !== "undefined" ? globalThis.location.host : undefined;

  // Detect navigation state from URL, context, and host
  const navigationState = detectNavigationState(pathname, requestContext, host);

  return (
    <I18nextProvider i18n={i18n}>
      <NavigationProvider state={navigationState}>
        <RootDocument
          locale={navigationState.locale}
          dir={navigationState.localeData.dir}
          requestContext={requestContext}
        >
          <AuthProvider>
            <ThemeProvider
              defaultTheme="system"
              storageKey="vite-ui-theme"
              locale={navigationState.locale}
              enableServerSync
            >
              <LocaleSynchronizer />
              <Outlet />
              <Toaster />
            </ThemeProvider>
          </AuthProvider>
        </RootDocument>
      </NavigationProvider>
    </I18nextProvider>
  );
}

// Component to sync i18n language with navigation state.
// The URL pathname is the single source of truth for locale.
function LocaleSynchronizer() {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;
  const { requestContext } = Route.useRouteContext();
  const host = typeof window !== "undefined" ? globalThis.location.host : undefined;

  useEffect(() => {
    const navigationState = detectNavigationState(pathname, requestContext, host);
    if (i18n.language !== navigationState.locale) {
      i18n.changeLanguage(navigationState.locale);
    }
  }, [pathname, requestContext, host]);

  return null;
}

interface RootDocumentProps {
  children: ReactNode;
  locale?: string;
  dir?: "ltr" | "rtl";
  requestContext?: RequestContext;
}

function RootDocument(props: Readonly<RootDocumentProps>) {
  const { children, locale = DEFAULT_LOCALE, dir = "ltr", requestContext } = props;

  // Resolve SSR theme class for flicker-free rendering
  const ssrTheme = requestContext?.ssrTheme;
  let ssrThemeClass: string | undefined;
  if (ssrTheme === "light" || ssrTheme === "dark") {
    ssrThemeClass = ssrTheme;
  }
  // "system" or undefined — let the inline script handle it

  const themeScript = `
(function() {
  var storageKey = 'vite-ui-theme';
  // Read theme cookie (cross-domain), then localStorage, then system
  var cookieTheme = (document.cookie.match(/(?:^|;\\s*)site_theme=([^;]+)/) || [])[1] || null;
  var localTheme = localStorage.getItem(storageKey);
  var theme = cookieTheme || localTheme;
  var systemTheme = matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  var resolvedTheme = theme === 'system' || !theme ? systemTheme : theme;
  document.documentElement.classList.remove('light', 'dark');
  document.documentElement.classList.add(resolvedTheme);
  // Sync cookie theme to localStorage for consistency
  if (cookieTheme && cookieTheme !== localTheme) {
    localStorage.setItem(storageKey, cookieTheme);
  }
})();

globalThis.__REQUEST_CONTEXT__ = ${JSON.stringify(requestContext).replace(/</g, "\\u003c").replace(/>/g, "\\u003e").replace(/\//g, "\\u002f")};
`;

  return (
    <html lang={locale} dir={dir} className={ssrThemeClass} suppressHydrationWarning>
      <head>
        <HeadContent />
        <script dangerouslySetInnerHTML={{ __html: themeScript }} />
      </head>
      <body className="min-h-screen bg-background font-sans antialiased">
        {children}
        <EasterEgg />
        <ResponsiveIndicator />
        <TanStackDevtools
          config={{
            position: "bottom-left",
          }}
          plugins={[
            {
              name: "Tanstack Router",
              render: <TanStackRouterDevtoolsPanel />,
            },
          ]}
        />
        <Scripts />
      </body>
    </html>
  );
}

function NotFoundComponent() {
  const { t } = useTranslation();

  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-4">
      <h1 className="font-serif text-4xl font-bold">404</h1>
      <p className="text-muted-foreground">{t("Error.Page not found")}</p>
      <Button render={<Link to="/">{t("Layout.Go Home")}</Link>} />
    </div>
  );
}
