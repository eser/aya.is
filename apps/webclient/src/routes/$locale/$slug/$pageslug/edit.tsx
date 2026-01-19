// Edit page
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

    return { page };
  },
  component: EditPagePage,
});

function EditPagePage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const auth = useAuth();
  const { page } = Route.useLoaderData();
  const [canEdit, setCanEdit] = React.useState<boolean | null>(null);

  // Check permissions client-side
  React.useEffect(() => {
    if (auth.isLoading) return;

    if (!auth.isAuthenticated) {
      setCanEdit(false);
      return;
    }

    backend.getProfilePermissions(params.locale, params.slug).then((perms) => {
      setCanEdit(perms !== null && perms.can_edit);
    });
  }, [auth.isAuthenticated, auth.isLoading, params.locale, params.slug]);

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

  // No permission - redirect to page
  if (!canEdit) {
    return (
      <div className="content">
        <h2>Access Denied</h2>
        <p>You don't have permission to edit this page.</p>
      </div>
    );
  }

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
    <div className="h-[calc(100vh-140px)]">
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
