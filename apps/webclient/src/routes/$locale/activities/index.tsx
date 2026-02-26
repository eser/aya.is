// Activities listing page
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { Button } from "@/components/ui/button";
import { backend } from "@/modules/backend/backend";
import { useAuth } from "@/lib/auth/auth-context";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import { formatMonthYear } from "@/lib/date";
import { ActivityCard } from "./_components/-activity-card";
import type { ActivityProperties, StoryEx } from "@/modules/backend/types";
import i18next from "i18next";

export const Route = createFileRoute("/$locale/activities/")({
  loader: async ({ params }) => {
    const { locale } = params;
    const activities = await backend.getActivities(locale);

    // Ensure locale translations are loaded before translating
    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    return {
      activities,
      locale,
      translatedTitle: t("Layout.Activities"),
      translatedDescription: t("Activities.Discover upcoming activities and meetups"),
    };
  },
  head: ({ loaderData }) => {
    const { locale, translatedTitle, translatedDescription } = loaderData;
    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, "activities"),
        locale,
        type: "website",
      }),
      links: [generateCanonicalLink(buildUrl(locale, "activities"))],
    };
  },
  component: ActivitiesPage,
});

type GroupedActivities = {
  monthYear: string;
  date: Date;
  activities: StoryEx[];
};

function getActivityDate(activity: StoryEx): Date {
  const activityProps = (activity.properties ?? {}) as unknown as ActivityProperties;

  if (activityProps.activity_time_start !== undefined) {
    const parsed = new Date(activityProps.activity_time_start);
    if (!Number.isNaN(parsed.getTime())) {
      return parsed;
    }
  }

  const fallback = activity.published_at ?? activity.created_at;
  const parsed = new Date(fallback);
  return Number.isNaN(parsed.getTime()) || parsed.getFullYear() < 1900 ? new Date() : parsed;
}

function groupActivitiesByMonth(activities: StoryEx[], locale: string): GroupedActivities[] {
  const groups = new Map<string, GroupedActivities>();

  const sorted = [...activities].sort((a, b) => getActivityDate(b).getTime() - getActivityDate(a).getTime());

  for (const activity of sorted) {
    const date = getActivityDate(activity);
    const monthYear = formatMonthYear(date, locale);

    if (!groups.has(monthYear)) {
      groups.set(monthYear, { monthYear, date, activities: [] });
    }

    groups.get(monthYear)!.activities.push(activity);
  }

  return Array.from(groups.values()).sort(
    (a, b) => b.date.getTime() - a.date.getTime(),
  );
}

function ActivitiesPage() {
  const { activities, locale } = Route.useLoaderData();
  const { t, i18n } = useTranslation();
  const { isAuthenticated } = useAuth();
  const currentLocale = i18n.language;

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <div className="flex items-center justify-between mb-4">
            <h1 className="no-margin">{t("Layout.Activities")}</h1>
            {isAuthenticated && (
              <Link
                to="/$locale/activities/new"
                params={{ locale }}
              >
                <Button variant="default" size="sm">
                  <Plus className="mr-1.5 size-4" />
                  {t("Activities.Add Activity")}
                </Button>
              </Link>
            )}
          </div>

          {activities === null || activities.length === 0
            ? (
              <p className="text-muted-foreground">
                {t("Activities.No activities yet")}
              </p>
            )
            : (
              <>
                {groupActivitiesByMonth(activities, currentLocale).map((group) => (
                  <div key={group.monthYear} className="mb-8">
                    <h2 className="text-lg font-semibold text-muted-foreground mb-4 pb-2 border-b border-border">
                      {formatMonthYear(group.date, currentLocale)}
                    </h2>
                    <div>
                      {group.activities.map((activity) => <ActivityCard key={activity.id} activity={activity} />)}
                    </div>
                  </div>
                ))}
              </>
            )}
        </div>
      </section>
    </PageLayout>
  );
}
