// Profile index - shows profile stories/timeline with date grouping and pagination
import { createFileRoute } from "@tanstack/react-router";
import { backend } from "@/modules/backend/backend";
import { StoriesPageClient } from "@/routes/$locale/stories/_components/-stories-page-client";

export const Route = createFileRoute("/$locale/$slug/")({
  validateSearch: (search: Record<string, unknown>) => {
    const offset = Number(search.offset) || 0;
    return offset > 0 ? { offset } : {};
  },
  loader: async ({ params }) => {
    const { slug, locale } = params;
    const stories = await backend.getProfileStories(locale, slug);
    return { stories, slug, locale };
  },
  component: ProfileIndexPage,
});

function ProfileIndexPage() {
  const { stories, slug, locale } = Route.useLoaderData();

  return (
    <div className="content">
      <StoriesPageClient
        initialStories={stories}
        basePath={`/${locale}/${slug}`}
        profileSlug={slug}
      />
    </div>
  );
}
