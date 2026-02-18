// Profile Q&A page
import { createFileRoute, getRouteApi, notFound } from "@tanstack/react-router";
import { backend } from "@/modules/backend/backend";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { buildUrl, generateMetaTags } from "@/lib/seo";
import i18next from "i18next";
import { QAPageClient } from "./-components/qa-page-client";

const parentRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/qa/")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const profile = await backend.getProfile(locale, slug);

    if (profile?.feature_qa === "disabled") {
      throw notFound();
    }

    const questionsData = await backend.getProfileQuestions(locale, slug);

    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    const translatedTitle = `${t("Layout.Q&A")} - ${profile?.title ?? slug}`;
    const translatedDescription = t("QA.Ask a question to this profile.");

    return {
      questionsData,
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
        url: buildUrl(locale, slug, "qa"),
        locale,
        type: "website",
      }),
    };
  },
  component: QAIndexPage,
});

function QAIndexPage() {
  const { questionsData, locale, slug } = Route.useLoaderData();
  const { profile, permissions } = parentRoute.useLoaderData();

  if (profile === null) {
    return null;
  }

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale} viewerMembershipKind={permissions?.viewer_membership_kind}>
      <QAPageClient
        questions={questionsData ?? []}
        locale={locale}
        slug={slug}
        profileId={profile.id}
        profileKind={profile.kind}
      />
    </ProfileSidebarLayout>
  );
}
