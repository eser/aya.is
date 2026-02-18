// Profile contributions page
import { createFileRoute, getRouteApi, notFound } from "@tanstack/react-router";
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
    const profile = await backend.getProfile(locale, slug);

    if (profile?.feature_relations === "disabled") {
      throw notFound();
    }

    const contributions = await backend.getProfileContributions(locale, slug);

    if (contributions === null) {
      throw notFound();
    }

    // Ensure locale translations are loaded before translating
    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    const translatedTitle = `${t("Layout.Contributions")} - ${profile?.title ?? slug}`;
    const translatedDescription = t("Contributions.Organizations and products this person contributes to.");

    return {
      contributions,
      locale,
      slug,
      profileTitle: profile?.title ?? slug,
      translatedTitle,
      translatedDescription,
    };
  },
  head: ({ loaderData }) => {
    const { locale, slug, translatedTitle, translatedDescription } = loaderData;
    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
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
  const { profile, permissions } = parentRoute.useLoaderData();
  const { t } = useTranslation();

  if (profile === null) {
    return null;
  }

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale} viewerMembershipKind={permissions?.viewer_membership_kind}>
      <div className="space-y-6">
        <div>
          <h2 className="font-serif text-2xl font-bold text-foreground">{t("Layout.Contributions")}</h2>
          <p className="text-muted-foreground">
            {t(
              "Contributions.Organizations and products this person contributes to.",
            )}
          </p>
        </div>

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
