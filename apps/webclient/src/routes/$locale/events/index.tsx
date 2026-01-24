// Events page
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth/auth-context";

export const Route = createFileRoute("/$locale/events")({
  loader: async ({ params }) => {
    const { locale } = params;
    // const events = await backend.getEvents(locale);
    return { events: [], locale };
  },
  component: EventsPage,
});

function EventsPage() {
  const { events, locale } = Route.useLoaderData();
  const { t } = useTranslation();
  const { isAuthenticated } = useAuth();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <div className="flex items-center justify-between mb-20">
            <h1 className="mb-0">{t("Layout.Events")}</h1>
            {isAuthenticated && (
              <Link
                to="/$locale/events/new"
                params={{ locale }}
                disabled
              >
                <Button variant="default" size="sm">
                  <Plus className="mr-1.5 size-4" />
                  {t("Events.Add Event")}
                </Button>
              </Link>
            )}
          </div>

          <p className="text-muted-foreground">
            {t("Layout.Content not yet available.")}
          </p>
        </div>
      </section>
    </PageLayout>
  );
}
