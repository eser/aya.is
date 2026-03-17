// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Series detail page — shows series info and its stories
import { createFileRoute } from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import { seriesQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import { Story } from "@/components/userland/story";
import i18next from "i18next";

export const Route = createFileRoute("/$locale/series/$seriesSlug/")({
  loader: async ({ params, context }) => {
    const { locale, seriesSlug } = params;

    const [seriesData] = await Promise.all([
      context.queryClient.ensureQueryData(seriesQueryOptions(locale, seriesSlug)),
      i18next.loadLanguages(locale),
    ]);

    return {
      locale,
      seriesSlug,
      seriesTitle: seriesData?.series.title ?? seriesSlug,
      seriesDescription: seriesData?.series.description ?? "",
    };
  },
  head: ({ loaderData }) => {
    const { locale, seriesSlug, seriesTitle, seriesDescription } = loaderData;
    return {
      meta: generateMetaTags({
        title: seriesTitle,
        description: seriesDescription,
        url: buildUrl(locale, `series/${seriesSlug}`),
        locale,
        type: "website",
      }),
      links: [generateCanonicalLink(buildUrl(locale, `series/${seriesSlug}`))],
    };
  },
  errorComponent: QueryError,
  component: SeriesDetailPage,
});

function SeriesDetailPage() {
  const { locale, seriesSlug } = Route.useLoaderData();
  const { data } = useSuspenseQuery(seriesQueryOptions(locale, seriesSlug));
  const { t } = useTranslation();

  if (data === null) {
    return (
      <PageLayout>
        <section className="container px-4 py-8 mx-auto">
          <div className="content">
            <p className="text-muted-foreground">{t("Series.No series yet")}</p>
          </div>
        </section>
      </PageLayout>
    );
  }

  const { series, stories } = data;

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          {series.series_picture_uri !== null && (
            <img
              src={series.series_picture_uri}
              alt={series.title}
              className="w-full h-48 object-cover rounded-lg mb-6"
            />
          )}

          <h1>{series.title}</h1>
          {series.description.length > 0 && <p className="text-muted-foreground mb-6">{series.description}</p>}

          <h2>{t("Series.Stories in this series")}</h2>

          {stories.length === 0
            ? (
              <p className="text-muted-foreground">
                {t("Series.No stories in this series yet")}
              </p>
            )
            : (
              <div>
                {stories.map((story) => <Story key={story.id} story={story} />)}
              </div>
            )}
        </div>
      </section>
    </PageLayout>
  );
}
