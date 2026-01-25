// Profile contributions page
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { ContributionCard } from "@/components/userland/contribution-card/contribution-card";
import { buildUrl, generateMetaTags } from "@/lib/seo";
import i18next from "i18next";

const parentRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/contributions")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const contributions = await backend.getProfileContributions(locale, slug);
    const profile = await backend.getProfile(locale, slug);
    return { contributions, locale, slug, profileTitle: profile?.title ?? slug };
  },
  head: ({ loaderData }) => {
    const { locale, slug, profileTitle } = loaderData;
    const t = i18next.getFixedT(locale);
    return {
      meta: generateMetaTags({
        title: `${t("Layout.Contributions")} - ${profileTitle}`,
        description: t("Contributions.Organizations and products this person contributes to."),
        url: buildUrl(locale, slug, "contributions"),
        locale,
        type: "website",
      }),
    };
  },
  component: ContributionsPage,
});

function ContributionsPage() {
  const { contributions, locale, slug } = Route.useLoaderData();
  const { profile } = parentRoute.useLoaderData();
  const { t } = useTranslation();

  if (profile === null) {
    return null;
  }

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale}>
      <div className="content">
        <h2>{t("Layout.Contributions")}</h2>
        <p className="text-muted-foreground mb-6">
          {t(
            "Contributions.Organizations and products this person contributes to.",
          )}
        </p>

        {contributions !== null && contributions.length > 0
          ? (
            <div className="flex flex-col gap-4">
              {contributions.map((membership) => (
                <ContributionCard key={membership.id} membership={membership} />
              ))}
            </div>
          )
          : (
            <p className="text-muted-foreground">
              {t("Contributions.No contributions found.")}
            </p>
          )}
      </div>
    </ProfileSidebarLayout>
  );
}
