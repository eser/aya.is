// Create new story page
import * as React from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import { ContentEditor, type ContentEditorData } from "@/components/content-editor";
import { useAuth } from "@/lib/auth/auth-context";
import { Skeleton } from "@/components/ui/skeleton";
import { PageLayout } from "@/components/page-layouts/default";
import { getDatePrefix } from "@/lib/slugify";

export const Route = createFileRoute("/$locale/stories/new")({
  component: NewStoryPage,
});

function NewStoryPage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const auth = useAuth();

  // Content locale — owned by the route, passed to ContentEditor as locale prop
  const [contentLocale, setContentLocale] = React.useState(params.locale);

  // Get user's profile slug directly from auth context
  const userProfileSlug = auth.user?.individual_profile_slug ?? null;

  // Fetch profile to get discussion default setting
  const [discussionsByDefault, setDiscussionsByDefault] = React.useState(false);
  React.useEffect(() => {
    if (userProfileSlug === null) return;
    backend.getProfile(params.locale, userProfileSlug).then((profile) => {
      if (profile !== null && profile.option_story_discussions_by_default === true) {
        setDiscussionsByDefault(true);
      }
    });
  }, [params.locale, userProfileSlug]);

  if (auth.isLoading) {
    return (
      <PageLayout fullHeight>
        <div className="flex h-full flex-col">
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
                  {[1, 2, 3, 4, 5].map((i) => <Skeleton key={i} className="size-8" />)}
                </div>
                <div className="flex gap-1">
                  {[1, 2, 3].map((i) => <Skeleton key={i} className="size-8" />)}
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
          <h2>{t("Auth.Access Denied")}</h2>
          <p>{t("Auth.You need to be logged in with a profile to create stories.")}</p>
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
    kind: "article",
    visibility: "public",
    featDiscussions: discussionsByDefault,
  };

  const handleSave = async (data: ContentEditorData) => {
    const result = await backend.insertStory(contentLocale, userProfileSlug, {
      slug: data.slug,
      kind: data.kind ?? "article",
      title: data.title,
      summary: data.summary,
      content: data.content,
      story_picture_uri: data.storyPictureUri,
      visibility: data.visibility ?? "public",
      feat_discussions: data.featDiscussions,
    });

    if (result !== null) {
      toast.success(t("ContentEditor.Story created successfully"));
      navigate({
        to: "/$locale/stories/$storyslug/edit",
        params: {
          locale: params.locale,
          storyslug: data.slug,
        },
        search: { lang: contentLocale },
      });
    } else {
      toast.error(t("ContentEditor.Failed to create story"));
    }
  };

  return (
    <PageLayout fullHeight>
      <ContentEditor
        key={`${contentLocale}-${discussionsByDefault}`}
        locale={contentLocale}
        profileSlug={userProfileSlug}
        contentType="story"
        initialData={initialData}
        backUrl={`/${params.locale}/stories`}
        userKind={auth.user?.kind}
        onSave={handleSave}
        isNew
        accessibleProfiles={auth.user?.accessible_profiles ?? []}
        individualProfile={auth.user?.individual_profile}
        onLocaleChange={setContentLocale}
      />
    </PageLayout>
  );
}
