// Profile story page
import * as React from "react";
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { backend } from "@/modules/backend/backend";
import { StoryContent } from "@/components/widgets/story-content";
import { DiscussionThread } from "@/components/widgets/discussion-thread";
import { compileMdx, compileMdxLite } from "@/lib/mdx";
import { siteConfig } from "@/config";
import { useAuth } from "@/lib/auth/auth-context";
import type { DiscussionComment, DiscussionListResponse } from "@/modules/backend/types";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { setResponseHeader } from "@tanstack/react-start/server";
import { computeContentLanguage, computeStoryCanonicalUrl, generateCanonicalLink, generateMetaTags, truncateDescription } from "@/lib/seo";
import { ChildNotFound } from "../../route";

const profileRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/stories/$storyslug/")({
  loader: async ({ params }) => {
    const { locale, slug, storyslug } = params;
    const story = await backend.getProfileStory(locale, slug, storyslug);

    if (story === null || story === undefined) {
      return { story: null, compiledContent: null, currentUrl: null, locale, slug, initialDiscussion: null, notFound: true as const };
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

    // Build current URL for sharing
    const currentUrl = `${siteConfig.host}/${locale}/${slug}/stories/${storyslug}`;

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

    // Set Content-Language header with content locale awareness
    if (import.meta.env.SSR) {
      setResponseHeader("Content-Language", computeContentLanguage(locale, story.locale_code));
    }

    return { story, compiledContent, currentUrl, locale, slug, initialDiscussion, notFound: false as const };
  },
  head: ({ loaderData }) => {
    if (loaderData === undefined || loaderData.notFound || loaderData.story === null) {
      return { meta: [] };
    }
    const { story, currentUrl, locale } = loaderData;
    const canonicalUrl = computeStoryCanonicalUrl(story, locale, "stories");
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
      links: [generateCanonicalLink(canonicalUrl)],
    };
  },
  component: ProfileStoryPage,
  notFoundComponent: ChildNotFound,
});

function ProfileStoryPage() {
  const params = Route.useParams();
  const auth = useAuth();
  const loaderData = Route.useLoaderData();
  const { profile, permissions } = profileRoute.useLoaderData();
  const [canEdit, setCanEdit] = React.useState(false);

  if (loaderData.notFound || loaderData.story === null || profile === null) {
    return <ChildNotFound />;
  }

  const { story, compiledContent, currentUrl, locale, slug, initialDiscussion } = loaderData;

  // Check edit permissions
  React.useEffect(() => {
    if (auth.isAuthenticated && !auth.isLoading) {
      backend
        .getStoryPermissions(params.locale, params.slug, story.id)
        .then((perms) => {
          if (perms !== null) {
            setCanEdit(perms.can_edit);
          }
        });
    } else {
      setCanEdit(false);
    }
  }, [auth.isAuthenticated, auth.isLoading, params.locale, params.slug, story.id]);

  const editUrl = canEdit ? `/${params.locale}/stories/${params.storyslug}/edit` : undefined;

  // Show author info when viewing from a different profile (e.g. a publication),
  // but not when viewing from the author's own profile.
  const isAuthorProfile = story.author_profile?.id === profile.id;

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale} viewerMembershipKind={permissions?.viewer_membership_kind}>
      <StoryContent
        story={story}
        compiledContent={compiledContent}
        currentUrl={currentUrl}
        showAuthor={!isAuthorProfile}
        showPublications={false}
        headingOffset={2}
        editUrl={editUrl}
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
    </ProfileSidebarLayout>
  );
}
