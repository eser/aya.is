// Admin layout - handles admin permission check
import { createFileRoute, notFound, Outlet } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { LocaleLink } from "@/components/locale-link";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";

export const Route = createFileRoute("/$locale/admin")({
  ssr: false,
  loader: async ({ params }) => {
    const { locale } = params;

    // Check if user is logged in and is admin
    const session = await backend.getSessionCurrent(locale);

    if (!session.authenticated || session.user === undefined) {
      throw notFound();
    }

    // Check if user is admin
    if (session.user.kind !== "admin") {
      throw notFound();
    }

    return { session, user: session.user };
  },
  component: AdminLayout,
});

function AdminLayout() {
  const { t } = useTranslation();

  return (
    <PageLayout>
      <div className="container mx-auto py-8 px-4">
        <div className="space-y-6">
          <div>
            <LocaleLink to={`/admin`}>
              <h1 className="font-serif text-3xl font-bold text-foreground">
                {t("Admin.Admin Area")}
              </h1>
            </LocaleLink>
            <p className="text-muted-foreground">
              {t("Admin.Manage site content and settings.")}
            </p>
          </div>

          <nav className="flex gap-4 border-b">
            <LocaleLink
              to={`/admin/profiles`}
              className="relative pb-2 text-sm font-medium text-muted-foreground hover:text-foreground"
              activeProps={{
                className:
                  "relative pb-2 text-sm font-medium text-foreground after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-foreground",
              }}
            >
              {t("Admin.Profiles")}
            </LocaleLink>
            <LocaleLink
              to={`/admin/points`}
              className="relative pb-2 text-sm font-medium text-muted-foreground hover:text-foreground"
              activeProps={{
                className:
                  "relative pb-2 text-sm font-medium text-foreground after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-foreground",
              }}
            >
              {t("Admin.Points")}
            </LocaleLink>
          </nav>

          <Outlet />
        </div>
      </div>
    </PageLayout>
  );
}
