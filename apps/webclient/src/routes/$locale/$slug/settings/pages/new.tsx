// Create new profile page
import * as React from "react";
import { createFileRoute, useNavigate, getRouteApi } from "@tanstack/react-router";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import {
  ContentEditor,
  type ContentEditorData,
} from "@/components/content-editor";
import { useAuth } from "@/lib/auth/auth-context";
import { Skeleton } from "@/components/ui/skeleton";

const profileRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/settings/pages/new")({
  component: NewPagePage,
});

function NewPagePage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const auth = useAuth();
  const { profile } = profileRoute.useLoaderData();
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

  // No permission
  if (!canEdit || profile === null) {
    return (
      <div className="content">
        <h2>Access Denied</h2>
        <p>You don't have permission to create pages for this profile.</p>
      </div>
    );
  }

  const initialData: ContentEditorData = {
    title: "",
    slug: "",
    summary: "",
    content: "",
    storyPictureUri: null,
    status: "draft",
  };

  const handleSave = async (data: ContentEditorData) => {
    const result = await backend.createProfilePage(
      params.locale,
      params.slug,
      {
        slug: data.slug,
        title: data.title,
        summary: data.summary,
        content: data.content,
        cover_picture_uri: data.storyPictureUri,
        published_at: data.publishedAt,
      },
    );

    if (result !== null) {
      toast.success("Page created successfully");
      navigate({
        to: "/$locale/$slug/$pageslug",
        params: {
          locale: params.locale,
          slug: params.slug,
          pageslug: data.slug,
        },
      });
    } else {
      toast.error("Failed to create page");
    }
  };

  return (
    <div className="h-[calc(100vh-140px)]">
      <ContentEditor
        locale={params.locale}
        profileSlug={params.slug}
        contentType="page"
        initialData={initialData}
        backUrl={`/${params.locale}/${params.slug}/settings/pages`}
        userKind={auth.user?.kind}
        onSave={handleSave}
        isNew
      />
    </div>
  );
}
