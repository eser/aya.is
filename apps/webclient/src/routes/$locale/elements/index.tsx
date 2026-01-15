// Elements index page
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { ElementsContent } from "./_components/-elements-content";

export const Route = createFileRoute("/$locale/elements/")({
  loader: async ({ params }) => {
    const { locale } = params;
    const profiles = await backend.getProfilesByKinds(locale, [
      "individual",
      "organization",
    ]);
    return { profiles: profiles ?? [] };
  },
  component: ElementsIndexPage,
});

function ElementsIndexPage() {
  const { profiles } = Route.useLoaderData();
  const { t } = useTranslation();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <h1>{t("Layout.Elements")}</h1>
          <h2 className="subtitle">{t("Elements.Description")}</h2>

          <ElementsContent initialProfiles={profiles} />
        </div>
      </section>
    </PageLayout>
  );
}
