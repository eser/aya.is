import type { ReactNode } from "react";
import { useEffect } from "react";
import { createRootRouteWithContext, HeadContent, Link, Outlet, Scripts, useRouterState } from "@tanstack/react-router";
import { TanStackDevtools } from "@tanstack/react-devtools";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { I18nextProvider, useTranslation } from "react-i18next";
import { ThemeProvider } from "@/components/theme-provider";
import { AuthProvider } from "@/lib/auth/auth-context";
import { Button } from "@/components/ui/button";
import { NavigationProvider, type NavigationState } from "@/modules/navigation/navigation-context";
import { EasterEgg, ResponsiveIndicator } from "@/components/page-layouts/default";
import {
  DEFAULT_LOCALE,
  isValidLocale,
  type SupportedLocaleCode,
  supportedLocales,
} from "@/config";
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
  const customProfile = isCustom ? domainConfig.profileSlug : null;

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
    customDomainProfile: customProfile,
  };
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  beforeLoad: ({ context }) => {
    // Pass the context from middleware
    return { requestContext: context.requestContext };
  },
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
  const { requestContext } = Route.useRouteContext();
  const host = typeof window !== "undefined" ? globalThis.location.host : undefined;

  useEffect(() => {
    const navigationState = detectNavigationState(pathname, requestContext, host);
    if (i18nInstance.language !== navigationState.locale) {
      i18nInstance.changeLanguage(navigationState.locale);
    }
  }, [pathname, requestContext, host, i18nInstance]);

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

  const themeScript = `
(function() {
  const storageKey = 'vite-ui-theme';
  const theme = localStorage.getItem(storageKey);
  const systemTheme = matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  const resolvedTheme = theme === 'system' || !theme ? systemTheme : theme;
  document.documentElement.classList.add(resolvedTheme);
})();

globalThis.__REQUEST_CONTEXT__ = ${JSON.stringify(requestContext)};
`;

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
        <TanStackDevtools
          config={{
            position: "bottom-right",
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
      <h1 className="text-4xl font-bold">404</h1>
      <p className="text-muted-foreground">{t("Error.Page not found")}</p>
      <Button render={<Link to="/">{t("Layout.Go Home")}</Link>} />
    </div>
  );
}
