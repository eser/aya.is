// Edit page
import { createFileRoute, useNavigate, notFound } from "@tanstack/react-router";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import {
  ContentEditor,
  type ContentEditorData,
} from "@/components/content-editor";

export const Route = createFileRoute("/$locale/$slug/$pageslug/edit")({
  loader: async ({ params }) => {
    const { locale, slug, pageslug } = params;

    // Skip if pageslug matches known routes
    const knownRoutes = ["stories", "settings", "members", "contributions"];
    if (knownRoutes.includes(pageslug)) {
      throw notFound();
    }

    // Get the page
    const page = await backend.getProfilePage(locale, slug, pageslug);
    if (page === null) {
      throw notFound();
    }

    // Check permissions
    const permissions = await backend.getProfilePermissions(locale, slug);
    if (permissions === null || !permissions.can_edit) {
      throw notFound();
    }

    return { page };
  },
  component: EditPagePage,
});

function EditPagePage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const { page } = Route.useLoaderData();

  const initialData: ContentEditorData = {
    title: page.title ?? "",
    slug: page.slug ?? "",
    summary: "",
    content: page.content ?? "",
    coverImageUrl: null,
    status: "published",
  };

  const handleSave = async (data: ContentEditorData) => {
    // Update the page main fields
    const updateResult = await backend.updateProfilePage(
      params.locale,
      params.slug,
      page.id,
      {
        slug: data.slug,
        order: page.sort_order,
        cover_picture_uri: data.coverImageUrl,
        published_at: null,
      },
    );

    if (updateResult === null) {
      toast.error("Failed to update page");
      return;
    }

    // Update the translation
    const translationResult = await backend.updateProfilePageTranslation(
      params.locale,
      params.slug,
      page.id,
      params.locale,
      {
        title: data.title,
        summary: data.summary,
        content: data.content,
      },
    );

    if (translationResult !== null) {
      toast.success("Page saved successfully");
      // If slug changed, navigate to new URL
      if (data.slug !== params.pageslug) {
        navigate({
          to: "/$locale/$slug/$pageslug",
          params: {
            locale: params.locale,
            slug: params.slug,
            pageslug: data.slug,
          },
        });
      }
    } else {
      toast.error("Failed to save page translation");
    }
  };

  const handleDelete = async () => {
    const result = await backend.deleteProfilePage(
      params.locale,
      params.slug,
      page.id,
    );

    if (result !== null) {
      toast.success("Page deleted successfully");
      navigate({
        to: "/$locale/$slug",
        params: { locale: params.locale, slug: params.slug },
      });
    } else {
      toast.error("Failed to delete page");
    }
  };

  return (
    <div className="content h-[calc(100vh-200px)]">
      <ContentEditor
        locale={params.locale}
        profileSlug={params.slug}
        contentType="page"
        initialData={initialData}
        backUrl={`/${params.locale}/${params.slug}/${params.pageslug}`}
        backLabel="Back to Page"
        onSave={handleSave}
        onDelete={handleDelete}
      />
    </div>
  );
}
