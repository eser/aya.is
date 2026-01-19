// Edit story page
import { createFileRoute, useNavigate, notFound } from "@tanstack/react-router";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import {
  ContentEditor,
  type ContentEditorData,
} from "@/components/content-editor";

export const Route = createFileRoute("/$locale/$slug/stories/$storyslug/edit")({
  loader: async ({ params }) => {
    const { locale, slug, storyslug } = params;

    // First, get the story to find its ID
    const story = await backend.getProfileStory(locale, slug, storyslug);
    if (story === null) {
      throw notFound();
    }

    // Then get the editable version with permissions check
    const editData = await backend.getStoryForEdit(locale, slug, story.id);
    if (editData === null) {
      throw notFound();
    }

    return { story, editData };
  },
  component: EditStoryPage,
});

function EditStoryPage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const { story, editData } = Route.useLoaderData();

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
    <div className="content h-[calc(100vh-200px)]">
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
