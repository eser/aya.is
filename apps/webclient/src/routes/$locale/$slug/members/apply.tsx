// Application form page — allows non-members to apply to join an organization
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import { applicationFormQueryOptions, myApplicationQueryOptions, profileQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import i18next from "i18next";
import { NotFoundContent } from "./route";
import { ApplyPageClient } from "./-components/apply-page-client";

const parentRoute = getRouteApi("/$locale/$slug");

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

  const { form, existingApplication, locale, slug } = loaderData;
  const { permissions } = parentRoute.useLoaderData();

  const isMember = permissions?.viewer_membership_kind !== undefined &&
    permissions.viewer_membership_kind !== null &&
    MEMBER_KINDS.has(permissions.viewer_membership_kind);

  return (
    <ApplyPageClient
      form={form}
      locale={locale}
      slug={slug}
      existingApplication={existingApplication}
      isMember={isMember}
    />
  );
}
