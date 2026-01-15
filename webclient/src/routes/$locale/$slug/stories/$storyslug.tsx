// Profile story page
import { createFileRoute, notFound } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { MdxContent } from "@/components/userland/mdx-content";
import { compileMdx } from "@/lib/mdx";

export const Route = createFileRoute("/$locale/$slug/stories/$storyslug")({
  loader: async ({ params }) => {
    const { locale, slug, storyslug } = params;
    const story = await backend.getProfileStory(locale, slug, storyslug);

    if (!story) {
      throw notFound();
    }

    // Compile MDX content on the server
    let compiledContent: string | null = null;
    if (story.content) {
      try {
        compiledContent = await compileMdx(story.content);
      } catch (error) {
        console.error("Failed to compile MDX:", error);
        compiledContent = null;
      }
    }

    return { story, compiledContent };
  },
  component: ProfileStoryPage,
  notFoundComponent: StoryNotFound,
});

function ProfileStoryPage() {
  const { story, compiledContent } = Route.useLoaderData();

  return (
    <article className="content">
      <h2>{story.title}</h2>
      {compiledContent ? (
        <MdxContent compiledSource={compiledContent} />
      ) : (
        story.content && (
          <div dangerouslySetInnerHTML={{ __html: story.content }} />
        )
      )}
    </article>
  );
}

function StoryNotFound() {
  const { t } = useTranslation();

  return (
    <div className="content">
      <h2>{t("Layout.Page not found")}</h2>
      <p className="text-muted-foreground">
        {t("Layout.The page you are looking for does not exist. Please check your spelling and try again.")}
      </p>
    </div>
  );
}
