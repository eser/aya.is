// Application form page — allows non-members to apply to join an organization
import { createFileRoute } from "@tanstack/react-router";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import { applicationFormQueryOptions, myApplicationQueryOptions, profileQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import i18next from "i18next";
import { NotFoundContent } from "./route";
import { ApplyPageClient } from "./-components/apply-page-client";

export const Route = createFileRoute("/$locale/$slug/members/apply")({
  loader: async ({ params, context }) => {
    const { locale, slug } = params;
    const [profile, form] = await Promise.all([
      context.queryClient.ensureQueryData(profileQueryOptions(locale, slug)),
      context.queryClient.ensureQueryData(
        applicationFormQueryOptions(locale, slug),
      ).catch(() => null),
    ]);

    // Applications must be enabled and a form must exist
    if (
      profile === null ||
      profile.feature_applications === "disabled" ||
      form === null
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

  if (loaderData.notFound || loaderData.form === null) {
    return <NotFoundContent />;
  }

  const { form, existingApplication, locale, slug } = loaderData;

  return (
    <ApplyPageClient
      form={form}
      locale={locale}
      slug={slug}
      existingApplication={existingApplication}
    />
  );
}
