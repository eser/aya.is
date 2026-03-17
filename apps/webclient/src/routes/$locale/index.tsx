// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Homepage for locale - shows landing page with parallax hero and latest stories
// On custom domain, server-side URL rewriting redirects to profile page
import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";
import { Astronaut } from "@/components/widgets/astronaut";
import { MdxContent } from "@/components/userland/mdx-content";
import { Story } from "@/components/userland/story";
import { compileMdxLite } from "@/lib/mdx";
import { formatMonthYear, parseDateFromSlug } from "@/lib/date";
import { siteConfig } from "@/config";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import type { StoryEx } from "@/modules/backend/types";
import { liveNowQueryOptions, storiesQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import { LiveNowSection } from "@/components/live-now/live-now";
import { HomeCta } from "@/components/widgets/home-cta";
import { getI18nInstance } from "@/routes/__root";
import styles from "./-home.module.css";

const MAX_STORIES = 20;

type GroupedStories = {
  monthYear: string;
  date: Date;
  stories: StoryEx[];
};

export const Route = createFileRoute("/$locale/")({
  loader: async ({ params, context }) => {
    const { locale } = params;
    const i18nInstance = getI18nInstance();
    await i18nInstance.loadLanguages(locale);
    const t = i18nInstance.getFixedT(locale);
    const introText = t("Home.IntroText");
    const siteSubtitle = t("Home.SiteSubtitle");
    const compiledIntro = await compileMdxLite(introText);

    await Promise.all([
      context.queryClient.ensureQueryData(storiesQueryOptions(locale)),
      context.queryClient.prefetchQuery(liveNowQueryOptions(locale)),
    ]);

    // Fetch GitHub star count for the stargazer CTA (non-blocking fallback)
    let githubStars = 0;
    try {
      const res = await fetch("https://api.github.com/repos/eser/aya.is", {
        headers: { Accept: "application/vnd.github.v3+json" },
        signal: AbortSignal.timeout(3000),
      });
      if (res.ok) {
        const data = await res.json();
        githubStars = data.stargazers_count ?? 0;
      }
    } catch {
      // Use fallback — stargazer button will be hidden
    }

    return { compiledIntro, locale, siteSubtitle, githubStars };
  },
  head: ({ loaderData }) => {
    const { locale, siteSubtitle } = loaderData;
    return {
      meta: generateMetaTags({
        title: `${siteConfig.name} - ${siteSubtitle}`,
        description: siteConfig.description,
        url: buildUrl(locale),
        locale,
        type: "website",
      }),
      links: [generateCanonicalLink(buildUrl(locale))],
    };
  },
  errorComponent: QueryError,
  component: LocaleHomePage,
});

function groupStoriesByMonth(
  stories: StoryEx[],
  locale: string,
): GroupedStories[] {
  const storiesWithDates = stories
    .map((story) => ({
      story,
      date: parseDateFromSlug(story.slug),
    }))
    .filter(
      (item): item is { story: StoryEx; date: Date } => item.date !== null,
    );

  storiesWithDates.sort((a, b) => b.date.getTime() - a.date.getTime());

  const limited = storiesWithDates.slice(0, MAX_STORIES);

  const groups = new Map<string, GroupedStories>();

  limited.forEach(({ story, date }) => {
    const monthYear = formatMonthYear(date, locale);

    if (!groups.has(monthYear)) {
      groups.set(monthYear, {
        monthYear,
        date,
        stories: [],
      });
    }

    groups.get(monthYear)!.stories.push(story);
  });

  return Array.from(groups.values()).sort(
    (a, b) => b.date.getTime() - a.date.getTime(),
  );
}

function LocaleHomePage() {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;
  const { compiledIntro: loaderIntro, locale: loaderLocale, githubStars } = Route.useLoaderData();
  const { data: allStories } = useSuspenseQuery(storiesQueryOptions(loaderLocale));
  const { data: liveStreams } = useQuery(liveNowQueryOptions(loaderLocale));

  // Loader compiles introText at request time for the URL locale.
  // When language changes client-side (e.g. via WebMCP switch-language),
  // the loader doesn't re-run, so we recompile the intro text reactively.
  const [compiledIntro, setCompiledIntro] = React.useState(loaderIntro);

  React.useEffect(() => {
    if (locale === loaderLocale) {
      setCompiledIntro(loaderIntro);
      return;
    }

    const introText = t("Home.IntroText");
    import("@/lib/mdx").then(async ({ compileMdxLite }) => {
      const compiled = await compileMdxLite(introText);
      setCompiledIntro(compiled);
    });
  }, [locale, loaderLocale, loaderIntro, t]);

  React.useEffect(() => {
    document.documentElement.classList.add("scroll-smooth");
    return () => {
      document.documentElement.classList.remove("scroll-smooth");
    };
  }, []);

  const groupedStories = React.useMemo(() => {
    if (allStories === null) return [];
    return groupStoriesByMonth(allStories, locale);
  }, [allStories, locale]);

  return (
    <PageLayout>
      {/* Hero section - sticky for parallax effect */}
      <section className={styles.heroSection}>
        <div className={styles.heroContainer}>
          <div className={styles.heroContent}>
            {/* Astronaut - positioned on the right side */}
            <div className={styles.astronautWrapper}>
              <Astronaut width={400} height={400} />
            </div>

            <article className="content relative z-10">
              <h1 className="hero">{t("Home.AYA the Open Source Network")}</h1>
              <h2 className="subtitle">
                {t(
                  "Home.A platform connecting the elements that produce and develop for the community",
                )}
              </h2>

              {compiledIntro !== null && <MdxContent compiledSource={compiledIntro} />}

              <div className="mt-10">
                <HomeCta githubStars={githubStars} />
              </div>
            </article>
          </div>
        </div>
      </section>

      {/* Live Now section - only visible when streams are active */}
      {liveStreams !== undefined && liveStreams !== null && liveStreams.length > 0 && (
        <LiveNowSection streams={liveStreams} locale={locale} />
      )}

      {/* Stories section - scrolls over hero for parallax */}
      <section id="latest" className={styles.storiesSection}>
        <div className={styles.storiesContainer}>
          <h2 className={styles.storiesHeader}>
            {t("Home.Latest Stories")}
          </h2>
          <div className="content">
            {groupedStories.length === 0 && (
              <p className="text-muted-foreground">
                {t("Layout.Content not yet available.")}
              </p>
            )}

            {groupedStories.map((group) => (
              <div key={group.monthYear} className="mb-8">
                <h3 className="text-lg font-semibold text-muted-foreground mb-4 pb-2 border-b border-border">
                  {formatMonthYear(group.date, locale)}
                </h3>
                <div>
                  {group.stories.map((story) => <Story key={story.id} story={story} />)}
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>
    </PageLayout>
  );
}
