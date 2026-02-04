// Main domain root - redirects to preferred locale
import { createFileRoute, redirect } from "@tanstack/react-router";
import { getPreferredLocale } from "@/lib/get-locale";

export const Route = createFileRoute("/")({
  beforeLoad: async ({ context }) => {
    // Determine domain-specific default locale (if custom domain has one configured)
    const domainConfig = context.requestContext?.domainConfiguration;
    const domainDefault = domainConfig?.type === "custom-domain" || domainConfig?.type === "main"
      ? domainConfig.defaultCulture
      : undefined;

    // Detect preferred locale: cookie → Accept-Language → domain default → DEFAULT_LOCALE
    const locale = await getPreferredLocale({ data: domainDefault });

    throw redirect({
      to: `/${locale}`,
      replace: true,
    });
  },
  component: () => null,
});
