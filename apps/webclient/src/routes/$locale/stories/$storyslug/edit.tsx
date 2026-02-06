// Edit story page
import * as React from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { backend, type StoryEditData } from "@/modules/backend/backend";
import {
  ContentEditor,
  type ContentEditorData,
} from "@/components/content-editor";
import { useAuth } from "@/lib/auth/auth-context";
import { Skeleton } from "@/components/ui/skeleton";
import { PageLayout } from "@/components/page-layouts/default";

export const Route = createFileRoute("/$locale/stories/$storyslug/edit")({
  ssr: false,
  component: EditStoryPage,
  notFoundComponent: StoryNotFound,
});

function EditStoryPage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const auth = useAuth();
  // Story-level data (shared across all locales) — loaded once
  const [storyData, setStoryData] = React.useState<StoryEditData | null>(null);
  const [canEdit, setCanEdit] = React.useState<boolean | null>(null);
  // Translation locale is independent from the site locale (params.locale)
  const [translationLocale, setTranslationLocale] = React.useState(params.locale);
  // Translation data tagged with its locale to prevent stale renders.
  // null = loading, data.locale mismatches translationLocale = stale (treat as loading)
  const [translationState, setTranslationState] = React.useState<{
    locale: string;
    data: { title: string; summary: string; content: string } | undefined; // undefined = no translation
  } | null>(null);

  // Load story data and check permissions (once)
  React.useEffect(() => {
    if (auth.isLoading) return;

    if (!auth.isAuthenticated) {
      setCanEdit(false);
      return;
    }

    // Use "_" as placeholder profileSlug - the handler doesn't use it
    backend.getStoryForEdit(params.locale, "_", params.storyslug).then((data) => {
      if (data === null) {
        setCanEdit(false);
      } else {
        setStoryData(data);
        setCanEdit(true);
      }
    });
  }, [auth.isAuthenticated, auth.isLoading, params.locale, params.storyslug]);

  // Load translation data when translationLocale changes
  React.useEffect(() => {
    if (storyData === null) return;

    if (translationLocale === params.locale) {
      // Use the initial data from storyData for the site locale
      const isFallback = storyData.is_fallback;
      setTranslationState({
        locale: translationLocale,
        data: isFallback ? undefined : {
          title: storyData.title ?? "",
          summary: storyData.summary ?? "",
          content: storyData.content,
        },
      });
      return;
    }

    // Fetch translation for a different locale
    setTranslationState(null);
    backend.getStoryForEdit(translationLocale, "_", params.storyslug).then((data) => {
      if (data === null || data.is_fallback) {
        // No translation exists for this locale — show empty fields
        setTranslationState({ locale: translationLocale, data: undefined });
      } else {
        setTranslationState({
          locale: translationLocale,
          data: {
            title: data.title ?? "",
            summary: data.summary ?? "",
            content: data.content,
          },
        });
      }
    });
  }, [translationLocale, storyData, params.locale, params.storyslug]);

  // Get the author's profile slug from story data
  const authorProfileSlug = storyData?.author_profile_slug ?? null;

  // Show skeleton while loading or when translation data is stale (locale mismatch)
  const translationReady = translationState !== null && translationState.locale === translationLocale;

  if (canEdit === null || !translationReady) {
    return (
      <PageLayout>
        <div className="flex h-[calc(100vh-140px)] flex-col">
          {/* Header skeleton */}
          <div className="flex items-center justify-between border-b p-4">
            <div className="flex items-center gap-3">
              <Skeleton className="size-10 rounded-full" />
              <Skeleton className="h-6 w-24" />
            </div>
            <div className="flex gap-2">
              <Skeleton className="h-9 w-20" />
              <Skeleton className="h-9 w-9" />
            </div>
          </div>
          {/* Main content skeleton */}
          <div className="flex flex-1 overflow-hidden">
            {/* Sidebar skeleton */}
            <div className="w-80 shrink-0 border-r p-4 space-y-4">
              <div className="flex items-center justify-between mb-4">
                <Skeleton className="h-4 w-16" />
                <Skeleton className="size-8" />
              </div>
              <div className="space-y-2">
                <Skeleton className="h-4 w-12" />
                <Skeleton className="h-10 w-full" />
              </div>
              <div className="space-y-2">
                <Skeleton className="h-4 w-10" />
                <Skeleton className="h-10 w-full" />
              </div>
              <div className="space-y-2">
                <Skeleton className="h-4 w-10" />
                <Skeleton className="h-10 w-full" />
              </div>
              <div className="space-y-2">
                <Skeleton className="h-4 w-16" />
                <Skeleton className="h-20 w-full" />
              </div>
            </div>
            {/* Editor content skeleton */}
            <div className="flex flex-1 flex-col overflow-hidden">
              {/* Toolbar skeleton */}
              <div className="flex items-center justify-between border-b px-4 py-2">
                <div className="flex gap-1">
                  {[1, 2, 3, 4, 5].map((i) => (
                    <Skeleton key={i} className="size-8" />
                  ))}
                </div>
                <div className="flex gap-1">
                  {[1, 2, 3].map((i) => (
                    <Skeleton key={i} className="size-8" />
                  ))}
                </div>
              </div>
              {/* Panels skeleton */}
              <div className="flex flex-1 overflow-hidden">
                <Skeleton className="flex-1 m-4" />
              </div>
            </div>
          </div>
        </div>
      </PageLayout>
    );
  }

  // No permission
  if (canEdit !== true || storyData === null || authorProfileSlug === null) {
    return (
      <PageLayout>
        <div className="content">
          <h2>Access Denied</h2>
          <p>You don't have permission to edit this story.</p>
        </div>
      </PageLayout>
    );
  }

  // translationState.data is undefined when no translation exists for the selected locale
  const isNewTranslation = translationState.data === undefined;

  const initialData: ContentEditorData = {
    title: isNewTranslation ? "" : translationState.data.title,
    slug: storyData.slug ?? "",
    summary: isNewTranslation ? "" : translationState.data.summary,
    content: isNewTranslation ? "" : translationState.data.content,
    storyPictureUri: storyData.story_picture_uri,
    kind: (storyData.kind as ContentEditorData["kind"]) ?? "article",
  };

  const handleLocaleChange = (newLocale: string) => {
    setTranslationLocale(newLocale);
  };

  const handleSave = async (data: ContentEditorData) => {
    // Update the story main fields
    const updateResult = await backend.updateStory(
      params.locale,
      authorProfileSlug,
      storyData.id,
      {
        slug: data.slug,
        story_picture_uri: data.storyPictureUri,
      },
    );

    if (updateResult === null) {
      toast.error("Failed to update story");
      return;
    }

    // Update the translation for the selected translation locale
    const translationResult = await backend.updateStoryTranslation(
      params.locale,
      authorProfileSlug,
      storyData.id,
      translationLocale,
      {
        title: data.title,
        summary: data.summary,
        content: data.content,
      },
    );

    if (translationResult !== null) {
      toast.success("Story saved successfully");
      // If slug changed, navigate to new URL
      if (data.slug !== params.storyslug) {
        navigate({
          to: "/$locale/stories/$storyslug",
          params: {
            locale: params.locale,
            storyslug: data.slug,
          },
        });
      }
    } else {
      toast.error("Failed to save story translation");
    }
  };

  const handleDelete = async () => {
    const result = await backend.removeStory(
      params.locale,
      authorProfileSlug,
      storyData.id,
    );

    if (result !== null) {
      toast.success("Story deleted successfully");
      navigate({
        to: "/$locale/stories",
        params: { locale: params.locale },
      });
    } else {
      toast.error("Failed to delete story");
    }
  };

  return (
    <PageLayout>
      <div className="h-[calc(100vh-140px)]">
        <ContentEditor
          key={translationLocale}
          locale={translationLocale}
          profileSlug={authorProfileSlug}
          contentType="story"
          initialData={initialData}
          backUrl={`/${params.locale}/stories/${params.storyslug}`}
          userKind={auth.user?.kind}
          validateSlugDatePrefix
          onSave={handleSave}
          onDelete={handleDelete}
          excludeId={storyData.id}
          storyId={storyData.id}
          initialPublications={storyData.publications ?? []}
          accessibleProfiles={auth.user?.accessible_profiles ?? []}
          individualProfile={auth.user?.individual_profile}
          onLocaleChange={handleLocaleChange}
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
