// Create new story page
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import {
  ContentEditor,
  type ContentEditorData,
} from "@/components/content-editor";
import { useAuth } from "@/lib/auth/auth-context";
import { Skeleton } from "@/components/ui/skeleton";
import { PageLayout } from "@/components/page-layouts/default";

export const Route = createFileRoute("/$locale/stories/new")({
  component: NewStoryPage,
});

// Helper to get current date as YYYYMMDD- prefix
function getDatePrefix(): string {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  return `${year}${month}${day}-`;
}

function NewStoryPage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const auth = useAuth();

  // Get user's profile slug directly from auth context
  const userProfileSlug = auth.user?.individual_profile_slug ?? null;

  if (auth.isLoading) {
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

  if (!auth.isAuthenticated || userProfileSlug === null) {
    return (
      <PageLayout>
        <div className="content">
          <h2>Access Denied</h2>
          <p>You need to be logged in with a profile to create stories.</p>
        </div>
      </PageLayout>
    );
  }

  const initialData: ContentEditorData = {
    title: "",
    slug: getDatePrefix(),
    summary: "",
    content: "",
    storyPictureUri: null,
    status: "draft",
    kind: "article",
  };

  const handleSave = async (data: ContentEditorData) => {
    const result = await backend.insertStory(params.locale, userProfileSlug, {
      slug: data.slug,
      kind: data.kind ?? "article",
      title: data.title,
      summary: data.summary,
      content: data.content,
      story_picture_uri: data.storyPictureUri,
      status: data.status,
      is_featured: data.isFeatured ?? false,
      published_at: data.publishedAt,
    });

    if (result !== null) {
      toast.success("Story created successfully");
      navigate({
        to: "/$locale/stories/$storyslug",
        params: {
          locale: params.locale,
          storyslug: data.slug,
        },
      });
    } else {
      toast.error("Failed to create story");
    }
  };

  return (
    <PageLayout>
      <div className="h-[calc(100vh-140px)]">
        <ContentEditor
          locale={params.locale}
          profileSlug={userProfileSlug}
          contentType="story"
          initialData={initialData}
          backUrl={`/${params.locale}/stories`}
          userKind={auth.user?.kind}
          validateSlugDatePrefix
          onSave={handleSave}
          isNew
        />
      </div>
    </PageLayout>
  );
}
