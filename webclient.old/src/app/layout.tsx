import "server-only";
import * as React from "react";
import type { Metadata } from "next";

import { siteConfig } from "@/shared/config.ts";
import { NavigationProvider } from "@/shared/modules/navigation/navigation-provider.tsx";
import { TranslationsProvider } from "@/shared/modules/i18n/translations-provider.tsx";
import { getNavigationState } from "@/shared/modules/navigation/get-navigation-state.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { PageLayout } from "@/shared/components/page-layouts/default/page-layout.tsx";
import { AuthProvider } from "@/shared/modules/auth/auth-context.tsx";
import { AuthErrorBanner } from "@/shared/modules/auth/auth-error-banner.tsx";
import { RegisterGlobals } from "./register-globals.tsx";
import { Analytics } from "./analytics.tsx";
import "../shared/globals.css";

// TODO(@eser) add more from https://beta.nextjs.org/docs/api-reference/metadata
export const metadata: Metadata = {
  title: {
    default: siteConfig.title,
    template: `%s | ${siteConfig.title}`,
  },
  description: siteConfig.description,

  icons: {
    icon: "/favicon.ico",
  },
};

export const viewport = {
  width: "device-width",
  initialScale: 1,
  // maximumScale: 1,
};

type LayoutProps = {
  children: React.ReactNode;
};

async function Layout(props: LayoutProps) {
  const navigationState = await getNavigationState();

  const translationsState = await getTranslations();
  const translationsStateValues = {
    localeCode: translationsState.localeCode,
    locale: translationsState.locale,
    supportedLocales: translationsState.supportedLocales,
    messages: translationsState.messages,
  };

  const placeholders: Record<string, string> = {
    locale: translationsState.locale.name,
  };

  return (
    <html lang={translationsState.locale.code} dir={translationsState.locale.dir} suppressHydrationWarning>
      <head />
      <body>
        <TranslationsProvider state={translationsStateValues}>
          <NavigationProvider state={navigationState}>
            <AuthProvider>
              <AuthErrorBanner />
              <PageLayout placeholders={placeholders}>
                {props.children}
              </PageLayout>
            </AuthProvider>
          </NavigationProvider>
        </TranslationsProvider>
        <RegisterGlobals />
        <Analytics />
      </body>
    </html>
  );
}

export { Layout as default };
