// Activities listing page
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { Button } from "@/components/ui/button";
import { backend } from "@/modules/backend/backend";
import { useAuth } from "@/lib/auth/auth-context";
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
  const { isAuthenticated } = useAuth();

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
