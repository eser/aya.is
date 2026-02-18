// Profile story page
import * as React from "react";
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { backend } from "@/modules/backend/backend";
import { StoryContent } from "@/components/widgets/story-content";
import { compileMdx } from "@/lib/mdx";
import { siteConfig } from "@/config";
import { useAuth } from "@/lib/auth/auth-context";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { generateMetaTags, truncateDescription } from "@/lib/seo";
import { ChildNotFound } from "../../route";

const profileRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/stories/$storyslug/")({
  loader: async ({ params }) => {
    const { locale, slug, storyslug } = params;
    const story = await backend.getProfileStory(locale, slug, storyslug);

    if (story === null || story === undefined) {
      return { story: null, compiledContent: null, currentUrl: null, locale, slug, notFound: true as const };
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

    return { story, compiledContent, currentUrl, locale, slug, notFound: false as const };
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
        publishedTime: story.created_at,
        modifiedTime: story.updated_at,
        author: story.author_profile?.title ?? null,
      }),
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

  const { story, compiledContent, currentUrl, locale, slug } = loaderData;

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
    </ProfileSidebarLayout>
  );
}
