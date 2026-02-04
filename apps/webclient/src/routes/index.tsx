// Main domain root - redirects to preferred locale
import { createFileRoute, redirect } from "@tanstack/react-router";
import { getPreferredLocale } from "@/lib/get-locale";

export const Route = createFileRoute("/")({
  beforeLoad: async () => {
    // Detect preferred locale: cookie → Accept-Language → domain default → DEFAULT_LOCALE
    const locale = await getPreferredLocale();

    throw redirect({
      to: `/${locale}`,
      replace: true,
    });
  },
  component: () => null,
});
