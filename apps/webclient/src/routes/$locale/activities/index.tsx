// Activities listing page
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { buildUrl, generateMetaTags } from "@/lib/seo";
import { ActivityCard } from "./_components/-activity-card";
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
    };
  },
  component: ActivitiesPage,
});

function ActivitiesPage() {
  const { activities, locale } = Route.useLoaderData();
  const { t } = useTranslation();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <div className="flex items-center justify-between mb-4">
            <h1 className="no-margin">{t("Layout.Activities")}</h1>
          </div>

          {activities === null || activities.length === 0
            ? (
              <p className="text-muted-foreground">
                {t("Activities.No activities yet")}
              </p>
            )
            : (
              <div>
                {activities.map((activity) => (
                  <ActivityCard key={activity.id} activity={activity} />
                ))}
              </div>
            )}
        </div>
      </section>
    </PageLayout>
  );
}
