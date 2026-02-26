// Locale layout - handles locale validation and i18n sync
import { CatchNotFound, createFileRoute, Outlet, redirect, notFound } from "@tanstack/react-router";
import { isValidLocale, type SupportedLocaleCode } from "@/config";
import { getPreferredLocale } from "@/lib/get-locale";
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
      // Redirect to /{preferredLocale}/{slug}/rest-of-path
      // Uses cookie → Accept-Language → domain default → 'en' detection chain
      const preferredLocale = await getPreferredLocale();
      const newPath = `/${preferredLocale}${location.pathname}`;
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
  // i18n language sync is handled by LocaleSynchronizer in __root.tsx
  return (
    <CatchNotFound fallback={<PageNotFound />}>
      <Outlet />
    </CatchNotFound>
  );
}
