// Events page
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";

export const Route = createFileRoute("/$locale/events")({
  component: EventsPage,
});

function EventsPage() {
  const { t } = useTranslation();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <h1>{t("Layout.Events")}</h1>
          <p className="text-muted-foreground">
            {t("Layout.Content not yet available.")}
          </p>
        </div>
      </section>
    </PageLayout>
  );
}
