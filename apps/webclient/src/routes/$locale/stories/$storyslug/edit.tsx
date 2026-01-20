// Edit story page
import * as React from "react";
import { createFileRoute, useNavigate, notFound } from "@tanstack/react-router";
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
  loader: async ({ params }) => {
    const { locale, storyslug } = params;

    // Get the story (public data)
    const story = await backend.getStory(locale, storyslug);
    if (story === null) {
      throw notFound();
    }

    return { story };
  },
  component: EditStoryPage,
});

function EditStoryPage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const auth = useAuth();
  const { story } = Route.useLoaderData();
  const [editData, setEditData] = React.useState<StoryEditData | null>(null);
  const [canEdit, setCanEdit] = React.useState<boolean | null>(null);

  // Get the author's profile slug for API calls
  const authorProfileSlug = story.author_profile?.slug ?? null;

  // Check permissions and load edit data client-side
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

    // Load edit data (which also checks permissions)
    backend.getStoryForEdit(params.locale, authorProfileSlug, story.id).then((data) => {
      if (data === null) {
        setCanEdit(false);
      } else {
        setEditData(data);
        setCanEdit(true);
      }
    });
  }, [auth.isAuthenticated, auth.isLoading, params.locale, authorProfileSlug, story.id]);

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

  const initialData: ContentEditorData = {
    title: editData.title ?? "",
    slug: editData.slug ?? "",
    summary: editData.summary ?? "",
    content: editData.content,
    coverImageUrl: editData.story_picture_uri,
    status: editData.status === "published" ? "published" : "draft",
    kind: (editData.kind as ContentEditorData["kind"]) ?? "article",
    isFeatured: editData.is_featured,
  };

  const handleSave = async (data: ContentEditorData) => {
    // Update the story main fields
    const updateResult = await backend.updateStory(
      params.locale,
      authorProfileSlug,
      story.id,
      {
        slug: data.slug,
        status: data.status,
        is_featured: data.isFeatured ?? editData.is_featured,
        story_picture_uri: data.coverImageUrl,
        kind: data.kind,
      },
    );

    if (updateResult === null) {
      toast.error("Failed to update story");
      return;
    }

    // Update the translation
    const translationResult = await backend.updateStoryTranslation(
      params.locale,
      authorProfileSlug,
      story.id,
      params.locale,
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
      story.id,
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
          locale={params.locale}
          profileSlug={authorProfileSlug}
          contentType="story"
          initialData={initialData}
          backUrl={`/${params.locale}/stories/${params.storyslug}`}
          onSave={handleSave}
          onDelete={handleDelete}
        />
      </div>
    </PageLayout>
  );
}
