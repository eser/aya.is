// Profile story page
import * as React from "react";
import { createFileRoute, notFound, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { StoryContent } from "@/components/widgets/story-content";
import { compileMdx } from "@/lib/mdx";
import { siteConfig } from "@/config";
import { useAuth } from "@/lib/auth/auth-context";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";

const profileRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/stories/$storyslug/")({
  loader: async ({ params }) => {
    const { locale, slug, storyslug } = params;
    const story = await backend.getProfileStory(locale, slug, storyslug);

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

    // Build current URL for sharing
    const currentUrl = `${siteConfig.host}/${locale}/${slug}/stories/${storyslug}`;

    return { story, compiledContent, currentUrl, locale, slug };
  },
  component: ProfileStoryPage,
  notFoundComponent: StoryNotFound,
});

function ProfileStoryPage() {
  const params = Route.useParams();
  const auth = useAuth();
  const { story, compiledContent, currentUrl, locale, slug } = Route.useLoaderData();
  const { profile } = profileRoute.useLoaderData();
  const [canEdit, setCanEdit] = React.useState(false);

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

  if (profile === null) {
    return <StoryNotFound />;
  }

  const editUrl = canEdit ? `/${params.locale}/stories/${params.storyslug}/edit` : undefined;

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale}>
      <StoryContent
        story={story}
        compiledContent={compiledContent}
        currentUrl={currentUrl}
        showAuthor={false}
        headingOffset={2}
        editUrl={editUrl}
      />
    </ProfileSidebarLayout>
  );
}

function StoryNotFound() {
  const { t } = useTranslation();

  return (
    <div className="content">
      <h2>{t("Layout.Page not found")}</h2>
      <p className="text-muted-foreground">
        {t(
          "Layout.The page you are looking for does not exist. Please check your spelling and try again.",
        )}
      </p>
    </div>
  );
}
