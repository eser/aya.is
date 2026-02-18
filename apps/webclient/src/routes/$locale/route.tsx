// Locale layout - handles locale validation and i18n sync
import { CatchNotFound, createFileRoute, Outlet, redirect, notFound } from "@tanstack/react-router";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { DEFAULT_LOCALE, isValidLocale, type SupportedLocaleCode } from "@/config";
import { PageNotFound } from "@/components/page-not-found";

// Paths that should NOT be matched by the $locale route
// These are handled by dedicated API/server routes
const EXCLUDED_PATHS = ["api", "llms.txt", "llms-full.txt"];

export const Route = createFileRoute("/$locale")({
  beforeLoad: async ({ params, location }) => {
    const { locale } = params;

    // Don't match paths that should be handled by API routes
    // This allows /api/* routes to take precedence
    if (EXCLUDED_PATHS.includes(locale)) {
      throw notFound();
    }

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
  notFoundComponent: PageNotFound,
});

function LocaleLayout() {
  const { locale } = Route.useRouteContext();
  const { i18n } = useTranslation();

  // Sync i18n language with URL locale
  useEffect(() => {
    if (locale !== undefined && i18n.language !== locale) {
      i18n.changeLanguage(locale);
    }
  }, [locale, i18n]);

  return (
    <CatchNotFound fallback={<PageNotFound />}>
      <Outlet />
    </CatchNotFound>
  );
}
