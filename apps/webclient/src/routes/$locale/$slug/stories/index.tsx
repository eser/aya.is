// Profile stories index - shows all profile stories with date grouping and pagination
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { backend } from "@/modules/backend/backend";
import { StoriesPageClient } from "@/routes/$locale/stories/_components/-stories-page-client";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { buildUrl, generateMetaTags } from "@/lib/seo";
import i18next from "i18next";

const profileRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/stories/")({
  validateSearch: (search: Record<string, unknown>) => {
    const offset = Number(search.offset) || 0;
    return offset > 0 ? { offset } : {};
  },
  loader: async ({ params }) => {
    const { slug, locale } = params;
    const stories = await backend.getProfileStories(locale, slug);
    const profile = await backend.getProfile(locale, slug);
    const profileTitle = profile?.title ?? slug;

    // Ensure locale translations are loaded before translating
    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    return {
      stories,
      slug,
      locale,
      profileTitle,
      translatedTitle: `${t("Layout.Stories")} - ${profileTitle}`,
      translatedDescription: t("Stories.Browse stories from {{profile}}", { profile: profileTitle }),
    };
  },
  head: ({ loaderData }) => {
    const { locale, slug, translatedTitle, translatedDescription } = loaderData;
    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, slug, "stories"),
        locale,
        type: "website",
      }),
    };
  },
  component: ProfileStoriesIndexPage,
});

function ProfileStoriesIndexPage() {
  const { stories, slug, locale } = Route.useLoaderData();
  const { profile, permissions } = profileRoute.useLoaderData();

  if (profile === null) {
    return null;
  }

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale} viewerMembershipKind={permissions?.viewer_membership_kind}>
      <div className="content">
        <StoriesPageClient
          initialStories={stories}
          basePath={`/${locale}/${slug}/stories`}
          profileSlug={slug}
        />
      </div>
    </ProfileSidebarLayout>
  );
}
