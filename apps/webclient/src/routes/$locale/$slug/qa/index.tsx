// Profile Q&A page
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { backend } from "@/modules/backend/backend";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { ChildNotFound } from "../route";
import { buildUrl, generateMetaTags } from "@/lib/seo";
import { compileMdxLite } from "@/lib/mdx";
import type { ProfileQuestion } from "@/modules/backend/types";
import i18next from "i18next";
import { QAPageClient } from "./-components/qa-page-client";

const parentRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/qa/")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const profile = await backend.getProfile(locale, slug);

    if (profile?.feature_qa === "disabled") {
      return { questionsData: null, locale, slug, profileTitle: slug, translatedTitle: "", translatedDescription: "", notFound: true as const };
    }

    const questionsData = await backend.getProfileQuestions(locale, slug);

    if (questionsData === null) {
      return { questionsData: null, locale, slug, profileTitle: slug, translatedTitle: "", translatedDescription: "", notFound: true as const };
    }

    // Pre-compile Q&A content for SSR
    const compiledQuestions: ProfileQuestion[] = await Promise.all(
      questionsData.map(async (q) => {
        let compiledContent: string | null = null;
        let compiledAnswer: string | null = null;
        if (q.content !== "") {
          try { compiledContent = await compileMdxLite(q.content); } catch { /* fallback */ }
        }
        if (q.answer_content !== null && q.answer_content !== "") {
          try { compiledAnswer = await compileMdxLite(q.answer_content); } catch { /* fallback */ }
        }

        return { ...q, compiledContent, compiledAnswer };
      }),
    );

    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    const translatedTitle = `${t("Layout.Q&A")} - ${profile?.title ?? slug}`;
    const translatedDescription = t("QA.Ask a question to this profile.");

    return {
      questionsData: compiledQuestions,
      locale,
      slug,
      profileTitle: profile?.title ?? slug,
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
        url: buildUrl(locale, slug, "qa"),
        locale,
        type: "website",
      }),
    };
  },
  component: QAIndexPage,
  notFoundComponent: ChildNotFound,
});

function QAIndexPage() {
  const loaderData = Route.useLoaderData();
  const { profile, permissions } = parentRoute.useLoaderData();

  if (loaderData.notFound || loaderData.questionsData === null || profile === null) {
    return <ChildNotFound />;
  }

  const { questionsData, locale, slug } = loaderData;

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
