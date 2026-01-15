// Main domain root - redirects to default locale
import { createFileRoute, redirect } from "@tanstack/react-router";
import { DEFAULT_LOCALE } from "@/config";

export const Route = createFileRoute("/")({
  beforeLoad: () => {
    // Redirect root to default locale
    throw redirect({
      to: `/${DEFAULT_LOCALE}`,
      replace: true,
    });
  },
  component: () => null,
});
