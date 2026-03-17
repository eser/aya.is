// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Application form page — allows non-members to apply to join an organization
import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import {
  applicationFormQueryOptions,
  myApplicationQueryOptions,
  profilePermissionsQueryOptions,
  profileQueryOptions,
} from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import i18next from "i18next";
import { NotFoundContent } from "./route";
import { ApplyPageClient } from "./-components/apply-page-client";

const MEMBER_KINDS = new Set([
  "member",
  "contributor",
  "maintainer",
  "lead",
  "owner",
]);

export const Route = createFileRoute("/$locale/$slug/members/apply")({
  loader: async ({ params, context }) => {
    const { locale, slug } = params;
    const [profile, form] = await Promise.all([
      context.queryClient.ensureQueryData(profileQueryOptions(locale, slug)),
      context.queryClient.ensureQueryData(
        applicationFormQueryOptions(locale, slug),
      ).catch(() => null),
    ]);

    // Applications must be enabled
    if (
      profile === null ||
      profile.feature_applications === "disabled"
    ) {
      return {
        form: null,
        existingApplication: null,
        locale,
        slug,
        translatedTitle: "",
        translatedDescription: "",
        notFound: true as const,
      };
    }

    // Check permissions for membership (works during SSR on aya.is where cookies are forwarded)
    const permissions = await context.queryClient.ensureQueryData(
      profilePermissionsQueryOptions(locale, slug),
    ).catch(() => null);

    const ssrIsMember = permissions?.viewer_membership_kind !== undefined &&
      permissions.viewer_membership_kind !== null &&
      MEMBER_KINDS.has(permissions.viewer_membership_kind);

    // Check if user already applied (optional — null if not authenticated)
    const existingApplication = await context.queryClient.ensureQueryData(
      myApplicationQueryOptions(locale, slug),
    ).catch(() => null);

    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    const translatedTitle = `${t("Applications.Apply to Join")} - ${profile.title ?? slug}`;
    const translatedDescription = t(
      "Applications.Fill out the form below to apply",
    );

    return {
      form,
      existingApplication,
      ssrIsMember,
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
        url: buildUrl(locale, slug, "members/apply"),
        locale,
        type: "website",
      }),
      links: [
        generateCanonicalLink(buildUrl(locale, slug, "members/apply")),
      ],
    };
  },
  errorComponent: QueryError,
  component: ApplyPage,
});

function ApplyPage() {
  const loaderData = Route.useLoaderData();

  if (loaderData.notFound) {
    return <NotFoundContent />;
  }

  const { form, ssrIsMember, locale, slug } = loaderData;

  // Isomorphic queries — read from SSR-hydrated cache, auto-refetch on client
  const { data: permissions } = useQuery(profilePermissionsQueryOptions(locale, slug));
  const { data: existingApplication } = useQuery(myApplicationQueryOptions(locale, slug));

  const queryIsMember = permissions?.viewer_membership_kind !== undefined &&
    permissions.viewer_membership_kind !== null &&
    MEMBER_KINDS.has(permissions.viewer_membership_kind);

  const isMember = ssrIsMember || queryIsMember;

  return (
    <ApplyPageClient
      form={form}
      locale={locale}
      slug={slug}
      existingApplication={existingApplication ?? null}
      isMember={isMember}
    />
  );
}
