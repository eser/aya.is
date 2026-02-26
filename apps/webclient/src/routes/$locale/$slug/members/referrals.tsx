// Profile membership referrals page
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { backend } from "@/modules/backend/backend";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import i18next from "i18next";
import { ChildNotFound } from "../route";
import { ReferralsPageClient } from "./-components/referrals-page-client";

const parentRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/members/referrals")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const profile = await backend.getProfile(locale, slug);

    if (profile?.feature_relations === "disabled") {
      return {
        referrals: null,
        teams: null,
        locale,
        slug,
        translatedTitle: "",
        translatedDescription: "",
        notFound: true as const,
      };
    }

    const [referrals, teams] = await Promise.all([
      backend.listReferrals(locale, slug),
      backend.listProfileTeams(locale, slug),
    ]);

    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    const translatedTitle = `${t("Layout.Referrals")} - ${profile?.title ?? slug}`;
    const translatedDescription = t(
      "Referrals.Referral proposals for new members.",
    );

    return {
      referrals: referrals ?? [],
      teams: teams ?? [],
      locale,
      slug,
      translatedTitle,
      translatedDescription,
      notFound: false as const,
    };
  },
  head: ({ loaderData }) => {
    if (loaderData === undefined || loaderData.notFound) {
      return { meta: [] };
    }

    const { locale, slug, translatedTitle, translatedDescription } = loaderData;

    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, slug, "members/referrals"),
        locale,
        type: "website",
      }),
      links: [
        generateCanonicalLink(buildUrl(locale, slug, "members/referrals")),
      ],
    };
  },
  component: ReferralsPage,
  notFoundComponent: ChildNotFound,
});

function ReferralsPage() {
  const loaderData = Route.useLoaderData();
  const { profile } = parentRoute.useLoaderData();

  if (loaderData.notFound || profile === null) {
    return <ChildNotFound />;
  }

  const { referrals, teams, locale, slug } = loaderData;

  return (
    <ReferralsPageClient
      referrals={referrals}
      teams={teams}
      locale={locale}
      slug={slug}
    />
  );
}
