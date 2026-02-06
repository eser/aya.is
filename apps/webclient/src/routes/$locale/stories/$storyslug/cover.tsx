// Cover generator page for stories
import * as React from "react";
import { createFileRoute, useNavigate, notFound } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend.ts";
import { CoverGenerator } from "@/components/cover-generator/index.ts";
import { useAuth } from "@/lib/auth/auth-context.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { Button } from "@/components/ui/button.tsx";
import { PageLayout } from "@/components/page-layouts/default";
import { ShieldX } from "lucide-react";
import styles from "@/components/cover-generator/cover-generator.module.css";

export const Route = createFileRoute("/$locale/stories/$storyslug/cover")({
  ssr: false,
  loader: async ({ params }) => {
    const { locale, storyslug } = params;

    // Get the story (public data)
    const story = await backend.getStory(locale, storyslug);
    if (story === null) {
      throw notFound();
    }

    return { story };
  },
  component: StoryCoverPage,
  notFoundComponent: StoryNotFound,
});

function StoryCoverPage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const auth = useAuth();
  const { t } = useTranslation();
  const { story } = Route.useLoaderData();
  const [canEdit, setCanEdit] = React.useState<boolean | null>(null);

  // Get the author's profile slug for API calls
  const authorProfileSlug = story.author_profile?.slug ?? null;

  // Check permissions client-side
  React.useEffect(() => {
    if (auth.isLoading) return;

    if (!auth.isAuthenticated) {
      setCanEdit(false);
      return;
    }

    if (authorProfileSlug === null) {
      setCanEdit(false);
      return;
    }

    // Check story permissions
    backend.getStoryPermissions(params.locale, authorProfileSlug, story.id).then((perms) => {
      setCanEdit(perms?.can_edit ?? false);
    });
  }, [auth.isAuthenticated, auth.isLoading, params.locale, authorProfileSlug, story.id]);

  // Handle back navigation
  const handleBack = () => {
    navigate({
      to: "/$locale/stories/$storyslug",
      params: {
        locale: params.locale,
        storyslug: params.storyslug,
      },
    });
  };

  // Handle cover set success
  const handleCoverSet = () => {
    // Navigate back to the story page after setting cover
    navigate({
      to: "/$locale/stories/$storyslug",
      params: {
        locale: params.locale,
        storyslug: params.storyslug,
      },
    });
  };

  // Still checking permissions - show skeleton
  if (canEdit === null) {
    return (
      <PageLayout>
        <div className="flex h-[calc(100vh-140px)] flex-col">
          {/* Header skeleton */}
          <div className={styles.skeletonHeader}>
            <div className="flex items-center justify-between p-4">
              <div className="flex items-center gap-3">
                <Skeleton className="h-9 w-20" />
                <Skeleton className="h-6 w-40" />
              </div>
              <div className="flex gap-2">
                <Skeleton className="h-9 w-28" />
                <Skeleton className="h-9 w-32" />
              </div>
            </div>
          </div>
          {/* Content skeleton */}
          <div className="flex flex-1 overflow-hidden">
            {/* Preview skeleton */}
            <div className={styles.skeletonPreview}>
              <Skeleton className={styles.skeletonPreviewBox} />
            </div>
            {/* Controls skeleton */}
            <div className={styles.skeletonControls}>
              <Skeleton className={styles.skeletonControl} />
              <Skeleton className={styles.skeletonControl} />
              <Skeleton className={styles.skeletonControl} />
              <Skeleton className={styles.skeletonControl} />
            </div>
          </div>
        </div>
      </PageLayout>
    );
  }

  // No permission - show access denied
  if (canEdit !== true || authorProfileSlug === null) {
    return (
      <PageLayout>
        <div className={styles.accessDenied}>
          <ShieldX className={styles.accessDeniedIcon} />
          <h2 className={styles.accessDeniedTitle}>
            {t("CoverDesigner.Access Denied")}
          </h2>
          <p className={styles.accessDeniedMessage}>
            {t("CoverDesigner.You don't have permission to generate covers for this story.")}
          </p>
          <Button variant="outline" onClick={handleBack}>
            {t("CoverDesigner.Back to Story")}
          </Button>
        </div>
      </PageLayout>
    );
  }

  // Render cover generator
  return (
    <PageLayout>
      <div className="h-[calc(100vh-140px)]">
        <CoverGenerator
          story={story}
          locale={params.locale}
          onBack={handleBack}
          onCoverSet={handleCoverSet}
        />
      </div>
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
