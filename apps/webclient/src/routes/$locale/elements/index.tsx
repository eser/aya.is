// Elements index page
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth/auth-context";
import { ElementsContent } from "./_components/-elements-content";

export const Route = createFileRoute("/$locale/elements/")({
  loader: async ({ params }) => {
    const { locale } = params;
    const profiles = await backend.getProfilesByKinds(locale, [
      "individual",
      "organization",
    ]);
    return { profiles: profiles ?? [], locale };
  },
  component: ElementsIndexPage,
});

function ElementsIndexPage() {
  const { profiles, locale } = Route.useLoaderData();
  const { t } = useTranslation();
  const { isAuthenticated } = useAuth();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <div className="flex items-center justify-between mb-20">
            <h1 className="no-margin">{t("Layout.Elements")}</h1>
            {isAuthenticated && (
              <Link
                to="/$locale/elements/new"
                params={{ locale }}
              >
                <Button variant="default" size="sm">
                  <Plus className="mr-1.5 size-4" />
                  {t("Elements.Add Element")}
                </Button>
              </Link>
            )}
          </div>

          <ElementsContent initialProfiles={profiles} />
        </div>
      </section>
    </PageLayout>
  );
}
