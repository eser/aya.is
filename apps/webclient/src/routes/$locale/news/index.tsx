// News page
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

export const Route = createFileRoute("/$locale/news/")({
  validateSearch: (search: Record<string, unknown>) => {
    const offset = Number(search.offset) || 0;
    return offset > 0 ? { offset } : {};
  },
  loader: async ({ params }) => {
    const { locale } = params;
    const news = await backend.getStoriesByKinds(locale, ["news"]);

    // Ensure locale translations are loaded before translating
    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    return {
      news,
      locale,
      translatedTitle: t("Layout.News"),
      translatedDescription: t("News.Latest news and updates from the AYA community"),
    };
  },
  head: ({ loaderData }) => {
    const { locale, translatedTitle, translatedDescription } = loaderData;
    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, "news"),
        locale,
        type: "website",
      }),
    };
  },
  component: NewsPage,
});

function NewsPage() {
  const { news, locale } = Route.useLoaderData();
  const { t } = useTranslation();
  const { isAuthenticated } = useAuth();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <div className="flex items-center justify-between mb-4">
            <h1 className="no-margin">{t("Layout.News")}</h1>
            {isAuthenticated && (
              <Link to="/$locale/news/new" params={{ locale }}>
                <Button variant="default" size="sm">
                  <Plus className="mr-1.5 size-4" />
                  {t("News.Add News")}
                </Button>
              </Link>
            )}
          </div>

          <StoriesPageClient
            initialStories={news}
            basePath={`/${locale}/news`}
          />
        </div>
      </section>
    </PageLayout>
  );
}
