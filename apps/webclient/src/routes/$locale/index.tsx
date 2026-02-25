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
import { backend } from "@/modules/backend/backend";
import type { StoryEx } from "@/modules/backend/types";
import i18next from "i18next";
import styles from "./-home.module.css";

const STORY_KINDS = [
  "article",
  "news",
  "announcement",
  "status",
  "content",
  "presentation",
  "activity",
];

const MAX_STORIES = 20;

type GroupedStories = {
  monthYear: string;
  date: Date;
  stories: StoryEx[];
};

export const Route = createFileRoute("/$locale/")({
  loader: async ({ params }) => {
    const { locale } = params;
    await i18next.loadLanguages(locale);
    const introText = i18next.getFixedT(locale)("Home.IntroText");
    const compiledIntro = await compileMdxLite(introText);

    let allStories: StoryEx[] | null = null;
    try {
      allStories = await backend.getStoriesByKinds(locale, STORY_KINDS);
    } catch {
      // Fetch can fail during HMR — render page without stories
    }

    return { compiledIntro, allStories, locale };
  },
  head: ({ loaderData }) => {
    const { locale } = loaderData;
    return {
      meta: generateMetaTags({
        title: `${siteConfig.name} - Acik Yazilim Agi`,
        description: siteConfig.description,
        url: buildUrl(locale),
        locale,
        type: "website",
      }),
      links: [generateCanonicalLink(buildUrl(locale))],
    };
  },
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
  const { compiledIntro, allStories } = Route.useLoaderData();

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

              <div className="mt-10" />

              {compiledIntro !== null && <MdxContent compiledSource={compiledIntro} />}
            </article>
          </div>
        </div>
      </section>

      {/* Stories section - scrolls over hero for parallax */}
      <section id="latest" className={styles.storiesSection}>
        <div className={styles.storiesContainer}>
          <div className="content">
            <h2 className={styles.storiesHeader}>
              <a href="#latest" className={styles.storiesHeaderLink}>
                {t("Home.Latest Stories")}
              </a>
            </h2>

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
                  {group.stories.map((story) => (
                    <Story key={story.id} story={story} />
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>
    </PageLayout>
  );
}
