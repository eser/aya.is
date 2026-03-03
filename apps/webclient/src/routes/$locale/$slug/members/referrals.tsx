// Profile membership referrals page
import { createFileRoute } from "@tanstack/react-router";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import { backend } from "@/modules/backend/backend";
import i18next from "i18next";
import { NotFoundContent } from "./route";
import { ReferralsPageClient } from "./-components/referrals-page-client";

const MEMBER_PLUS_KINDS = new Set([
  "member",
  "contributor",
  "maintainer",
  "lead",
  "owner",
]);

export const Route = createFileRoute("/$locale/$slug/members/referrals")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const [profile, permissions] = await Promise.all([
      backend.getProfile(locale, slug),
      backend.getProfilePermissions(locale, slug),
    ]);

    const isMemberPlus = permissions?.viewer_membership_kind !== undefined &&
      permissions.viewer_membership_kind !== null &&
      MEMBER_PLUS_KINDS.has(permissions.viewer_membership_kind);

    if (profile?.feature_relations === "disabled" || !isMemberPlus) {
      return {
        referrals: null,
        locale,
        slug,
        translatedTitle: "",
        translatedDescription: "",
        notFound: true as const,
      };
    }

    const referrals = await backend.listReferrals(locale, slug);

    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    const translatedTitle = `${t("Layout.Referrals")} - ${profile?.title ?? slug}`;
    const translatedDescription = t(
      "Referrals.Referral proposals for new members.",
    );

    return {
      referrals: referrals ?? [],
      locale,
      slug,
      viewerMembershipKind: permissions?.viewer_membership_kind ?? null,
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
});

function ReferralsPage() {
  const loaderData = Route.useLoaderData();

  if (loaderData.notFound) {
    return <NotFoundContent />;
  }

  const { referrals, locale, slug, viewerMembershipKind } = loaderData;

  return (
    <ReferralsPageClient
      referrals={referrals}
      locale={locale}
      slug={slug}
      viewerMembershipKind={viewerMembershipKind}
    />
  );
}
