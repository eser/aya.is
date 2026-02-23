// Individual story page
import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { StoryContent } from "@/components/widgets/story-content";
import { DiscussionThread } from "@/components/widgets/discussion-thread";
import { compileMdx, compileMdxLite } from "@/lib/mdx";
import { siteConfig } from "@/config";
import { useAuth } from "@/lib/auth/auth-context";
import { generateMetaTags, truncateDescription } from "@/lib/seo";
import { PageNotFound } from "@/components/page-not-found";
import type { DiscussionComment, DiscussionListResponse } from "@/modules/backend/types";

export const Route = createFileRoute("/$locale/stories/$storyslug/")({
  loader: async ({ params }) => {
    const { locale, storyslug } = params;
    const story = await backend.getStory(locale, storyslug);

    if (story === null || story === undefined) {
      return { story: null, compiledContent: null, currentUrl: null, locale, initialDiscussion: null, notFound: true as const };
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

    // Pre-fetch discussion data for SSR
    let initialDiscussion: { thread: DiscussionListResponse["thread"]; comments: DiscussionComment[] } | null = null;
    if (story.feat_discussions === true) {
      try {
        const discussion = await backend.getStoryDiscussion(locale, storyslug, "hot");
        if (discussion !== null && discussion !== undefined) {
          const compiledComments = await Promise.all(
            discussion.comments.map(async (comment) => {
              if (comment.content === "") {
                return { ...comment, compiledContent: null };
              }
              try {
                const compiled = await compileMdxLite(comment.content);
                return { ...comment, compiledContent: compiled };
              } catch {
                return { ...comment, compiledContent: null };
              }
            }),
          );
          initialDiscussion = { thread: discussion.thread, comments: compiledComments };
        }
      } catch {
        // Discussion fetch failed — component will fetch client-side
      }
    }

    return { story, compiledContent, currentUrl, locale, initialDiscussion, notFound: false as const };
  },
  head: ({ loaderData }) => {
    if (loaderData === undefined || loaderData.notFound || loaderData.story === null) {
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
        publishedTime: story.published_at ?? story.created_at,
        modifiedTime: story.updated_at,
        author: story.author_profile?.title ?? null,
      }),
    };
  },
  component: StoryPage,
  notFoundComponent: PageNotFound,
});

function StoryPage() {
  const params = Route.useParams();
  const auth = useAuth();
  const loaderData = Route.useLoaderData();
  const [canEdit, setCanEdit] = React.useState(false);

  if (loaderData.notFound || loaderData.story === null) {
    return <PageNotFound />;
  }

  const { story, compiledContent, currentUrl, initialDiscussion } = loaderData;

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

        {story.feat_discussions === true && (
          <DiscussionThread
            storySlug={params.storyslug}
            locale={params.locale}
            profileId={story.author_profile?.id ?? ""}
            profileKind={story.author_profile?.kind ?? "individual"}
            initialData={initialDiscussion}
          />
        )}
      </section>
    </PageLayout>
  );
}
