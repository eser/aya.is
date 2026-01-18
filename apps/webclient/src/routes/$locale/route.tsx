// Locale layout - handles locale validation and i18n sync
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { DEFAULT_LOCALE, isValidLocale, type SupportedLocaleCode } from "@/config";
import { PageLayout } from "@/components/page-layouts/default";

export const Route = createFileRoute("/$locale")({
  beforeLoad: async ({ params, location }) => {
    const { locale } = params;

    // Check if this is a valid locale
    if (!isValidLocale(locale)) {
      // Not a valid locale - treat as profile slug
      // Redirect to /{DEFAULT_LOCALE}/{slug}/rest-of-path
      const newPath = `/${DEFAULT_LOCALE}${location.pathname}`;
      throw redirect({ to: newPath, replace: true });
    }

    return {
      locale: locale as SupportedLocaleCode,
    };
  },
  component: LocaleLayout,
  notFoundComponent: LocaleNotFound,
});

function LocaleNotFound() {
  const { t } = useTranslation();

  return (
    <PageLayout>
      <div className="container mx-auto py-16 px-4 text-center">
        <h1 className="text-4xl font-bold mb-4">{t("Layout.Page not found")}</h1>
        <p className="text-muted-foreground">
          {t("Layout.The page you are looking for does not exist. Please check your spelling and try again.")}
        </p>
      </div>
    </PageLayout>
  );
}

function LocaleLayout() {
  const { locale } = Route.useRouteContext();
  const { i18n } = useTranslation();

  // Sync i18n language with URL locale
  useEffect(() => {
    if (locale !== undefined && i18n.language !== locale) {
      i18n.changeLanguage(locale);
    }
  }, [locale, i18n]);

  return <Outlet />;
}
