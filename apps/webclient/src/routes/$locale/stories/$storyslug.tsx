// Individual story page
import { createFileRoute, notFound } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { StoryFooter, StoryMetadata, StoryShare } from "@/components/widgets/story";
import { MdxContent } from "@/components/userland/mdx-content";
import { compileMdx } from "@/lib/mdx";
import { siteConfig } from "@/config";

export const Route = createFileRoute("/$locale/stories/$storyslug")({
  loader: async ({ params }) => {
    const { locale, storyslug } = params;
    const story = await backend.getStory(locale, storyslug);

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
        // Fall back to raw content (will be rendered as HTML)
        compiledContent = null;
      }
    }

    // Get current URL for sharing
    const currentUrl = `${siteConfig.host}/stories/${storyslug}`;

    return { story, compiledContent, currentUrl };
  },
  component: StoryPage,
  notFoundComponent: StoryNotFound,
});

function StoryPage() {
  const { story, compiledContent, currentUrl } = Route.useLoaderData();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto max-w-4xl">
        <div className="content">
          <h1>{story.title}</h1>

          <StoryMetadata story={story} />

          <article>
            {compiledContent ? <MdxContent compiledSource={compiledContent} /> : (
              story.content && <div dangerouslySetInnerHTML={{ __html: story.content }} />
            )}
          </article>

          <StoryShare story={story} currentUrl={currentUrl} />

          <StoryFooter story={story} />
        </div>
      </section>
    </PageLayout>
  );
}

function StoryNotFound() {
  const { t } = useTranslation();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto text-center">
        <h1 className="text-3xl font-bold mb-4">
          {t("Layout.Page not found")}
        </h1>
        <p className="text-muted-foreground">
          {t(
            "Layout.The page you are looking for does not exist. Please check your spelling and try again.",
          )}
        </p>
      </section>
    </PageLayout>
  );
}
