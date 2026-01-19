// Profile stories index - shows all profile stories with date grouping and pagination
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { backend } from "@/modules/backend/backend";
import { StoriesPageClient } from "@/routes/$locale/stories/_components/-stories-page-client";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";

const profileRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/stories/")({
  validateSearch: (search: Record<string, unknown>) => {
    const offset = Number(search.offset) || 0;
    return offset > 0 ? { offset } : {};
  },
  loader: async ({ params }) => {
    const { slug, locale } = params;
    const stories = await backend.getProfileStories(locale, slug);
    return { stories, slug, locale };
  },
  component: ProfileStoriesIndexPage,
});

function ProfileStoriesIndexPage() {
  const { stories, slug, locale } = Route.useLoaderData();
  const { profile } = profileRoute.useLoaderData();

  if (profile === null) {
    return null;
  }

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale}>
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
