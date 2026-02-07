// Individual story page
import * as React from "react";
import { createFileRoute, notFound } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { StoryContent } from "@/components/widgets/story-content";
import { compileMdx } from "@/lib/mdx";
import { siteConfig } from "@/config";
import { useAuth } from "@/lib/auth/auth-context";
import { generateMetaTags, truncateDescription } from "@/lib/seo";

export const Route = createFileRoute("/$locale/stories/$storyslug/")({
  loader: async ({ params }) => {
    const { locale, storyslug } = params;
    const story = await backend.getStory(locale, storyslug);

    if (story === null || story === undefined) {
      throw notFound();
    }

    // Compile MDX content on the server
    let compiledContent: string | null = null;
    if (story.content !== null && story.content !== undefined) {
      try {
        compiledContent = await compileMdx(story.content);
      } catch (error) {
        console.error("Failed to compile MDX:", error);
        compiledContent = null;
      }
    }

    // Get current URL for sharing
    const currentUrl = `${siteConfig.host}/${locale}/stories/${storyslug}`;

    return { story, compiledContent, currentUrl, locale };
  },
  head: ({ loaderData }) => {
    if (loaderData === undefined) {
      return { meta: [] };
    }
    const { story, currentUrl, locale } = loaderData;
    return {
      meta: generateMetaTags({
        title: story.title ?? "Story",
        description: truncateDescription(story.summary),
        url: currentUrl,
        image: story.story_picture_uri,
        locale,
        type: "article",
        publishedTime: story.created_at,
        modifiedTime: story.updated_at,
        author: story.author_profile?.title ?? null,
      }),
    };
  },
  component: StoryPage,
  notFoundComponent: StoryNotFound,
});

function StoryPage() {
  const params = Route.useParams();
  const auth = useAuth();
  const { story, compiledContent, currentUrl } = Route.useLoaderData();
  const [canEdit, setCanEdit] = React.useState(false);

  // Get author profile slug for permissions check
  const authorProfileSlug = story.author_profile?.slug ?? null;

  // Check edit permissions
  React.useEffect(() => {
    if (auth.isAuthenticated && !auth.isLoading && authorProfileSlug !== null) {
      backend
        .getStoryPermissions(params.locale, authorProfileSlug, story.id)
        .then((perms) => {
          if (perms !== null) {
            setCanEdit(perms.can_edit);
          }
        });
    } else {
      setCanEdit(false);
    }
  }, [auth.isAuthenticated, auth.isLoading, params.locale, authorProfileSlug, story.id]);

  const editUrl = canEdit ? `/${params.locale}/stories/${params.storyslug}/edit` : undefined;
  const coverUrl = canEdit ? `/${params.locale}/stories/${params.storyslug}/cover` : undefined;

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto max-w-4xl">
        <StoryContent
          story={story}
          compiledContent={compiledContent}
          currentUrl={currentUrl}
          showAuthor
          editUrl={editUrl}
          coverUrl={coverUrl}
        />
      </section>
    </PageLayout>
  );
}

function StoryNotFound() {
  const { t } = useTranslation();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto text-center">
        <h1 className="font-serif text-3xl font-bold mb-4">
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
