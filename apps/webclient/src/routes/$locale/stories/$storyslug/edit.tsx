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
  const [editData, setEditData] = React.useState<StoryEditData | null>(null);
  const [canEdit, setCanEdit] = React.useState<boolean | null>(null);
  // Translation locale is independent from the site locale (params.locale)
  const [translationLocale, setTranslationLocale] = React.useState(params.locale);

  // Load edit data client-side (auth required), re-fetch when translationLocale changes
  React.useEffect(() => {
    if (auth.isLoading) return;

    if (!auth.isAuthenticated) {
      setCanEdit(false);
      return;
    }

    // Reset while loading new locale
    setEditData(null);
    setCanEdit(null);

    // Use "_" as placeholder profileSlug - the handler doesn't use it
    backend.getStoryForEdit(translationLocale, "_", params.storyslug).then((data) => {
      if (data === null) {
        setCanEdit(false);
      } else {
        setEditData(data);
        setCanEdit(true);
      }
    });
  }, [auth.isAuthenticated, auth.isLoading, translationLocale, params.storyslug]);

  // Get the author's profile slug from edit data
  const authorProfileSlug = editData?.author_profile_slug ?? null;

  // Still checking permissions
  if (canEdit === null) {
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
  if (!canEdit || editData === null || authorProfileSlug === null) {
    return (
      <PageLayout>
        <div className="content">
          <h2>Access Denied</h2>
          <p>You don't have permission to edit this story.</p>
        </div>
      </PageLayout>
    );
  }

  // If the returned locale doesn't match the requested translation locale, this is a new translation
  const isNewTranslation = editData.locale_code !== translationLocale;

  const initialData: ContentEditorData = {
    title: isNewTranslation ? "" : (editData.title ?? ""),
    slug: editData.slug ?? "",
    summary: isNewTranslation ? "" : (editData.summary ?? ""),
    content: isNewTranslation ? "" : editData.content,
    storyPictureUri: editData.story_picture_uri,
    kind: (editData.kind as ContentEditorData["kind"]) ?? "article",
  };

  const handleLocaleChange = (newLocale: string) => {
    setTranslationLocale(newLocale);
  };

  const handleSave = async (data: ContentEditorData) => {
    // Update the story main fields
    const updateResult = await backend.updateStory(
      params.locale,
      authorProfileSlug,
      editData.id,
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
      editData.id,
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
      editData.id,
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
          excludeId={editData.id}
          storyId={editData.id}
          initialPublications={editData.publications ?? []}
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
