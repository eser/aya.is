// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Profile membership candidates page
import { createFileRoute } from "@tanstack/react-router";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import { candidatesQueryOptions, profilePermissionsQueryOptions, profileQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import i18next from "i18next";
import { NotFoundContent } from "./route";
import { CandidatesPageClient } from "./-components/candidates-page-client";

const MEMBER_PLUS_KINDS = new Set([
  "member",
  "contributor",
  "maintainer",
  "lead",
  "owner",
]);

export const Route = createFileRoute("/$locale/$slug/members/candidates")({
  loader: async ({ params, context }) => {
    const { locale, slug } = params;
    const [profile, permissions] = await Promise.all([
      context.queryClient.ensureQueryData(profileQueryOptions(locale, slug)),
      context.queryClient.ensureQueryData(profilePermissionsQueryOptions(locale, slug)).catch(() => null),
    ]);

    const isMemberPlus = permissions?.viewer_membership_kind !== undefined &&
      permissions.viewer_membership_kind !== null &&
      MEMBER_PLUS_KINDS.has(permissions.viewer_membership_kind);

    if (profile?.feature_relations === "disabled" || !isMemberPlus) {
      return {
        candidates: null,
        locale,
        slug,
        translatedTitle: "",
        translatedDescription: "",
        notFound: true as const,
      };
    }

    const candidates = await context.queryClient.ensureQueryData(candidatesQueryOptions(locale, slug));

    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    const translatedTitle = `${t("Layout.Candidates")} - ${profile?.title ?? slug}`;
    const translatedDescription = t(
      "Candidates.Candidate proposals for new members.",
    );

    return {
      candidates: candidates ?? [],
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
        url: buildUrl(locale, slug, "members/candidates"),
        locale,
        type: "website",
      }),
      links: [
        generateCanonicalLink(buildUrl(locale, slug, "members/candidates")),
      ],
    };
  },
  errorComponent: QueryError,
  component: CandidatesPage,
});

function CandidatesPage() {
  const loaderData = Route.useLoaderData();

  if (loaderData.notFound) {
    return <NotFoundContent />;
  }

  const { candidates, locale, slug, viewerMembershipKind } = loaderData;

  return (
    <CandidatesPageClient
      candidates={candidates}
      locale={locale}
      slug={slug}
      viewerMembershipKind={viewerMembershipKind}
    />
  );
}
