// Profile stories index - shows all profile stories with date grouping and pagination
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";
import { StoriesPageClient } from "@/routes/$locale/stories/_components/-stories-page-client";
import { profileQueryOptions, profileStoriesQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import i18next from "i18next";

const profileRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/stories/")({
  validateSearch: (search: Record<string, unknown>) => {
    const offset = Number(search.offset) || 0;
    return offset > 0 ? { offset } : {};
  },
  loader: async ({ params, context }) => {
    const { slug, locale } = params;
    await context.queryClient.ensureQueryData(profileStoriesQueryOptions(locale, slug));
    const profile = await context.queryClient.ensureQueryData(profileQueryOptions(locale, slug));
    const profileTitle = profile?.title ?? slug;

    // Ensure locale translations are loaded before translating
    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    return {
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
      links: [generateCanonicalLink(buildUrl(locale, slug, "stories"))],
    };
  },
  errorComponent: QueryError,
  component: ProfileStoriesIndexPage,
});

function ProfileStoriesIndexPage() {
  const { slug, locale } = Route.useLoaderData();
  const { data: stories } = useSuspenseQuery(profileStoriesQueryOptions(locale, slug));
  const { profile, permissions } = profileRoute.useLoaderData();

  if (profile === null) {
    return null;
  }

  return (
    <ProfileSidebarLayout
      profile={profile}
      slug={slug}
      locale={locale}
      viewerMembershipKind={permissions?.viewer_membership_kind}
    >
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
