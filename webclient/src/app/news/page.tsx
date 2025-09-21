import type { Metadata } from "next";

import { backend } from "@/shared/modules/backend/backend.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { StoriesPageClient } from "../stories/_components/stories-page-client.tsx";

export async function generateMetadata(): Promise<Metadata> {
  const { t } = await getTranslations();

  return {
    title: t("Layout", "News"),
  };
}

async function IndexPage() {
  const { t, locale } = await getTranslations();

  const news = await backend.getStoriesByKinds(locale.code, ["news"]);
  return (
    <section className="container px-4 py-8 mx-auto">
      <div className="content">
        <h2>{t("Layout", "News")}</h2>

        {news !== null && news.length > 0
          ? (
            <div className="divide-y divide-border">
              <StoriesPageClient initialStories={news} />
            </div>
          )
          : (
            <article>
              {/* TODO(@eser) Add a specific translation key e.g., t("NewsPage", "NoNewsFound") */}
              {t("Layout", "Content not yet available.")}
            </article>
          )}
      </div>
    </section>
  );
}

export { IndexPage as default };
