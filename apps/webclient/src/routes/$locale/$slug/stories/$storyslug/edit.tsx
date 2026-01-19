// Edit story page
import * as React from "react";
import { createFileRoute, useNavigate, notFound } from "@tanstack/react-router";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import {
  ContentEditor,
  type ContentEditorData,
} from "@/components/content-editor";
import { useAuth } from "@/lib/auth/auth-context";
import { Skeleton } from "@/components/ui/skeleton";

export const Route = createFileRoute("/$locale/$slug/stories/$storyslug/edit")({
  loader: async ({ params }) => {
    const { locale, slug, storyslug } = params;

    // Get the story (public data)
    const story = await backend.getProfileStory(locale, slug, storyslug);
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
  const [editData, setEditData] = React.useState<{
    title: string;
    slug: string;
    summary: string;
    content: string;
    story_picture_uri: string | null;
    status: string;
    is_featured: boolean;
  } | null>(null);
  const [canEdit, setCanEdit] = React.useState<boolean | null>(null);

  // Check permissions and load edit data client-side
  React.useEffect(() => {
    if (auth.isLoading) return;

    if (!auth.isAuthenticated) {
      setCanEdit(false);
      return;
    }

    // Load edit data (which also checks permissions)
    backend.getStoryForEdit(params.locale, params.slug, story.id).then((data) => {
      if (data === null) {
        setCanEdit(false);
      } else {
        setEditData(data);
        setCanEdit(true);
      }
    });
  }, [auth.isAuthenticated, auth.isLoading, params.locale, params.slug, story.id]);

  // Still checking permissions
  if (canEdit === null) {
    return (
      <div className="h-[calc(100vh-140px)]">
        <div className="flex flex-col gap-4">
          <div className="flex items-center justify-between">
            <Skeleton className="h-8 w-32" />
            <div className="flex gap-2">
              <Skeleton className="h-9 w-24" />
              <Skeleton className="h-9 w-24" />
            </div>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 flex-1">
            <div className="flex flex-col gap-2">
              <Skeleton className="h-10 w-full" />
              <Skeleton className="h-10 w-full" />
              <Skeleton className="h-[calc(100vh-400px)] w-full" />
            </div>
            <Skeleton className="h-[calc(100vh-320px)] w-full" />
          </div>
        </div>
      </div>
    );
  }

  // No permission
  if (!canEdit || editData === null) {
    return (
      <div className="content">
        <h2>Access Denied</h2>
        <p>You don't have permission to edit this story.</p>
      </div>
    );
  }

  const initialData: ContentEditorData = {
    title: editData.title ?? "",
    slug: editData.slug ?? "",
    summary: editData.summary ?? "",
    content: editData.content,
    coverImageUrl: editData.story_picture_uri,
    status: editData.status === "published" ? "published" : "draft",
  };

  const handleSave = async (data: ContentEditorData) => {
    // Update the story main fields
    const updateResult = await backend.updateStory(
      params.locale,
      params.slug,
      story.id,
      {
        slug: data.slug,
        status: data.status,
        is_featured: editData.is_featured,
        story_picture_uri: data.coverImageUrl,
      },
    );

    if (updateResult === null) {
      toast.error("Failed to update story");
      return;
    }

    // Update the translation
    const translationResult = await backend.updateStoryTranslation(
      params.locale,
      params.slug,
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
          to: "/$locale/$slug/stories/$storyslug",
          params: {
            locale: params.locale,
            slug: params.slug,
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
      params.slug,
      story.id,
    );

    if (result !== null) {
      toast.success("Story deleted successfully");
      navigate({
        to: "/$locale/$slug",
        params: { locale: params.locale, slug: params.slug },
      });
    } else {
      toast.error("Failed to delete story");
    }
  };

  return (
    <div className="h-[calc(100vh-140px)]">
      <ContentEditor
        locale={params.locale}
        profileSlug={params.slug}
        contentType="story"
        initialData={initialData}
        backUrl={`/${params.locale}/${params.slug}/stories/${params.storyslug}`}
        backLabel="Back to Story"
        onSave={handleSave}
        onDelete={handleDelete}
      />
    </div>
  );
}
