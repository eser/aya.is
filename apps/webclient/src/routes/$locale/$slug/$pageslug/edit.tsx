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
