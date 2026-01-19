// Create new story page
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import {
  ContentEditor,
  type ContentEditorData,
} from "@/components/content-editor";

export const Route = createFileRoute("/$locale/$slug/stories/new")({
  component: NewStoryPage,
});

function NewStoryPage() {
  const params = Route.useParams();
  const navigate = useNavigate();

  const initialData: ContentEditorData = {
    title: "",
    slug: "",
    summary: "",
    content: "",
    coverImageUrl: null,
    status: "draft",
  };

  const handleSave = async (data: ContentEditorData) => {
    const result = await backend.insertStory(params.locale, params.slug, {
      slug: data.slug,
      kind: "article",
      title: data.title,
      summary: data.summary,
      content: data.content,
      story_picture_uri: data.coverImageUrl,
      status: data.status,
      is_featured: false,
    });

    if (result !== null) {
      toast.success("Story created successfully");
      navigate({
        to: "/$locale/$slug/stories/$storyslug",
        params: {
          locale: params.locale,
          slug: params.slug,
          storyslug: data.slug,
        },
      });
    } else {
      toast.error("Failed to create story");
    }
  };

  return (
    <div className="content h-[calc(100vh-200px)]">
      <ContentEditor
        locale={params.locale}
        profileSlug={params.slug}
        contentType="story"
        initialData={initialData}
        backUrl={`/${params.locale}/${params.slug}`}
        backLabel="Back to Profile"
        onSave={handleSave}
        isNew
      />
    </div>
  );
}
