// News page
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { StoriesPageClient } from "./stories/_components/-stories-page-client";

export const Route = createFileRoute("/$locale/news")({
  validateSearch: (search: Record<string, unknown>) => {
    const offset = Number(search.offset) || 0;
    return offset > 0 ? { offset } : {};
  },
  loader: async ({ params }) => {
    const { locale } = params;
    const news = await backend.getStoriesByKinds(locale, ["news"]);
    return { news, locale };
  },
  component: NewsPage,
});

function NewsPage() {
  const { news, locale } = Route.useLoaderData();
  const { t } = useTranslation();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <h1>{t("Layout.News")}</h1>

          <StoriesPageClient
            initialStories={news}
            basePath={`/${locale}/news`}
          />
        </div>
      </section>
    </PageLayout>
  );
}
