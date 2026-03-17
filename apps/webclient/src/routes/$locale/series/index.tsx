// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Series listing page
import { createFileRoute } from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import { seriesListQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import { LocaleLink } from "@/components/locale-link";
import type { StorySeries } from "@/modules/backend/types";
import i18next from "i18next";

export const Route = createFileRoute("/$locale/series/")({
  loader: async ({ params, context }) => {
    const { locale } = params;

    await Promise.all([
      context.queryClient.ensureQueryData(seriesListQueryOptions(locale)),
      i18next.loadLanguages(locale),
    ]);

    const t = i18next.getFixedT(locale);
    return {
      locale,
      translatedTitle: t("Layout.Series"),
      translatedDescription: t("Series.Browse story series from the AYA community"),
    };
  },
  head: ({ loaderData }) => {
    const { locale, translatedTitle, translatedDescription } = loaderData;
    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, "series"),
        locale,
        type: "website",
      }),
      links: [generateCanonicalLink(buildUrl(locale, "series"))],
    };
  },
  errorComponent: QueryError,
  component: SeriesListPage,
});

function SeriesListPage() {
  const { locale } = Route.useLoaderData();
  const { data: seriesList } = useSuspenseQuery(seriesListQueryOptions(locale));
  const { t } = useTranslation();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <h1>{t("Layout.Series")}</h1>

          {seriesList === null || seriesList.length === 0
            ? (
              <p className="text-muted-foreground">
                {t("Series.No series yet")}
              </p>
            )
            : (
              <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
                {seriesList.map((series) => <SeriesCard key={series.id} series={series} />)}
              </div>
            )}
        </div>
      </section>
    </PageLayout>
  );
}

function SeriesCard(props: { series: StorySeries }) {
  return (
    <LocaleLink
      data-slot="card"
      to={`/series/${props.series.slug}`}
      className="no-underline block"
    >
      <div className="border rounded-lg overflow-hidden hover:border-primary transition-colors">
        {props.series.series_picture_uri !== null && (
          <img
            src={props.series.series_picture_uri}
            alt={props.series.title}
            className="w-full h-40 object-cover"
          />
        )}
        <div className="p-4">
          <h3 className="text-lg font-semibold mb-1">{props.series.title}</h3>
          {props.series.description.length > 0 && (
            <p className="text-sm text-muted-foreground line-clamp-2">
              {props.series.description}
            </p>
          )}
        </div>
      </div>
    </LocaleLink>
  );
}
