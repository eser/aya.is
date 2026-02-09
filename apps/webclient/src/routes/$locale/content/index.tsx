// Content page
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth/auth-context";
import { buildUrl, generateMetaTags } from "@/lib/seo";
import { StoriesPageClient } from "../stories/_components/-stories-page-client";
import i18next from "i18next";

export const Route = createFileRoute("/$locale/content/")({
  validateSearch: (search: Record<string, unknown>) => {
    const offset = Number(search.offset) || 0;
    return offset > 0 ? { offset } : {};
  },
  loader: async ({ params }) => {
    const { locale } = params;
    const content = await backend.getStoriesByKinds(locale, ["content"]);

    // Ensure locale translations are loaded before translating
    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    return {
      content,
      locale,
      translatedTitle: t("Layout.Content"),
      translatedDescription: t("Content.Latest content from the AYA community"),
    };
  },
  head: ({ loaderData }) => {
    const { locale, translatedTitle, translatedDescription } = loaderData;
    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, "content"),
        locale,
        type: "website",
      }),
    };
  },
  component: ContentPage,
});

function ContentPage() {
  const { content, locale } = Route.useLoaderData();
  const { t } = useTranslation();
  const { isAuthenticated } = useAuth();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <div className="flex items-center justify-between mb-4">
            <h1 className="no-margin">{t("Layout.Content")}</h1>
            {isAuthenticated && (
              <Link to="/$locale/content/new" params={{ locale }}>
                <Button variant="default" size="sm">
                  <Plus className="mr-1.5 size-4" />
                  {t("Content.Add Content")}
                </Button>
              </Link>
            )}
          </div>

          <StoriesPageClient
            initialStories={content}
            basePath={`/${locale}/content`}
          />
        </div>
      </section>
    </PageLayout>
  );
}
