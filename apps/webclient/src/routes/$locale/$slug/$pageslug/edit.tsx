// Edit profile page
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
  ssr: false,
  loader: async ({ params }) => {
    const { locale, slug, pageslug } = params;

    // Skip if pageslug matches known routes
    const knownRoutes = ["stories", "settings", "members", "contributions"];
    if (knownRoutes.includes(pageslug)) {
      throw notFound();
    }

    // Get the page (with site locale translation)
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
  // Translation locale is independent from the site locale (params.locale)
  const [translationLocale, setTranslationLocale] = React.useState(params.locale);
  // Translation data for the selected locale (null = not yet loaded, undefined = no translation exists)
  const [translationData, setTranslationData] = React.useState<
    { title: string; summary: string; content: string } | null | undefined
  >(null);

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

  // Load translation data when translationLocale changes
  React.useEffect(() => {
    if (translationLocale === params.locale) {
      // Use the loader data for the site locale
      setTranslationData({
        title: page.title ?? "",
        summary: page.summary ?? "",
        content: page.content ?? "",
      });
      return;
    }

    // Fetch page data for the selected translation locale
    setTranslationData(null);
    backend.getProfilePage(translationLocale, params.slug, params.pageslug).then((data) => {
      if (data === null) {
        // No translation exists for this locale — show empty fields
        setTranslationData(undefined);
      } else {
        setTranslationData({
          title: data.title ?? "",
          summary: data.summary ?? "",
          content: data.content ?? "",
        });
      }
    });
  }, [translationLocale, params.locale, params.slug, params.pageslug, page]);

  // Still checking permissions or loading translation
  if (canEdit === null || translationData === null) {
    return (
      <>
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
      </>
    );
  }

  // No permission
  if (!canEdit) {
    return (
      <>
        <div className="content">
          <h2>Access Denied</h2>
          <p>You don't have permission to edit this page.</p>
        </div>
      </>
    );
  }

  // translationData is undefined when no translation exists for the selected locale
  const isNewTranslation = translationData === undefined;

  const initialData: ContentEditorData = {
    title: isNewTranslation ? "" : translationData.title,
    slug: page.slug ?? "",
    summary: isNewTranslation ? "" : translationData.summary,
    content: isNewTranslation ? "" : translationData.content,
    storyPictureUri: page.cover_picture_uri ?? null,
  };

  const handleLocaleChange = (newLocale: string) => {
    setTranslationLocale(newLocale);
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
        cover_picture_uri: data.storyPictureUri ?? null,
        published_at: null,
      },
    );

    if (updateResult === null) {
      toast.error("Failed to update page");
      return;
    }

    // Update the translation for the selected translation locale
    const translationResult = await backend.updateProfilePageTranslation(
      params.locale,
      params.slug,
      page.id,
      translationLocale,
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
    <>
      <div className="h-[calc(100vh-140px)]">
        <ContentEditor
          key={translationLocale}
          locale={translationLocale}
          profileSlug={params.slug}
          contentType="page"
          initialData={initialData}
          backUrl={`/${params.locale}/${params.slug}/${params.pageslug}`}
          userKind={auth.user?.kind}
          onSave={handleSave}
          onDelete={handleDelete}
          excludeId={page.id}
          onLocaleChange={handleLocaleChange}
        />
      </div>
    </>
  );
}
