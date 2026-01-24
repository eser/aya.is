// Stories index page
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth/auth-context";
import { StoriesPageClient } from "./_components/-stories-page-client";

export const Route = createFileRoute("/$locale/stories/")({
  validateSearch: (search: Record<string, unknown>) => {
    const offset = Number(search.offset) || 0;
    return offset > 0 ? { offset } : {};
  },
  loader: async ({ params }) => {
    const { locale } = params;
    const stories = await backend.getStoriesByKinds(locale, ["article"]);
    return { stories, locale };
  },
  component: StoriesIndexPage,
});

function StoriesIndexPage() {
  const { stories, locale } = Route.useLoaderData();
  const { t } = useTranslation();
  const { isAuthenticated } = useAuth();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <div className="flex items-center justify-between mb-20">
            <h1 className="no-margin">{t("Layout.Articles")}</h1>
            {isAuthenticated && (
              <Link
                to="/$locale/stories/new"
                params={{ locale }}
              >
                <Button variant="default" size="sm">
                  <Plus className="mr-1.5 size-4" />
                  {t("Editor.Add Story")}
                </Button>
              </Link>
            )}
          </div>

          <StoriesPageClient
            initialStories={stories}
            basePath={`/${locale}/stories`}
          />
        </div>
      </section>
    </PageLayout>
  );
}
