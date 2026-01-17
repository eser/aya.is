// Profile stories index - shows all profile stories
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend, type StoryEx } from "@/modules/backend/backend";
import { Story } from "@/components/userland/story";

export const Route = createFileRoute("/$locale/$slug/stories/")({
  loader: async ({ params }) => {
    const { slug, locale } = params;
    const stories = await backend.getProfileStories(locale, slug);
    return { stories };
  },
  component: ProfileStoriesIndexPage,
});

function ProfileStoriesIndexPage() {
  const { stories } = Route.useLoaderData();
  const params = Route.useParams();
  const { t } = useTranslation();

  if (stories.length === 0) {
    return (
      <div className="content">
        <p className="text-muted-foreground">
          {t("Layout.Content not yet available.")}
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col">
      {stories.map((story: StoryEx) => (
        <Story key={story.id} story={story} profileSlug={params.slug} />
      ))}
    </div>
  );
}
