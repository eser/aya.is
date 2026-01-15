// Locale layout - handles both valid locales and profile slugs
import {
  createFileRoute,
  Outlet,
  redirect,
} from "@tanstack/react-router";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import {
  isValidLocale,
  DEFAULT_LOCALE,
  type SupportedLocaleCode,
} from "@/config";

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

    return { locale: locale as SupportedLocaleCode };
  },
  component: LocaleLayout,
});

function LocaleLayout() {
  const { locale } = Route.useRouteContext();
  const { i18n } = useTranslation();

  // Sync i18n language with URL locale
  useEffect(() => {
    if (locale && i18n.language !== locale) {
      i18n.changeLanguage(locale);
    }
  }, [locale, i18n]);

  return <Outlet />;
}
