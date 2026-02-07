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
  // Translation data tagged with its locale to prevent stale renders.
  // null = loading, data.locale mismatches translationLocale = stale (treat as loading)
  const [translationState, setTranslationState] = React.useState<{
    locale: string;
    data: { title: string; summary: string; content: string } | undefined; // undefined = no translation
  } | null>(null);

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
    let cancelled = false;

    if (translationLocale === params.locale) {
      // Use the loader data for the site locale, but detect fallback via locale_code
      if (page.locale_code !== undefined && page.locale_code.trimEnd() !== translationLocale) {
        // Loader returned a fallback — no translation for the site locale
        setTranslationState({ locale: translationLocale, data: undefined });
      } else {
        setTranslationState({
          locale: translationLocale,
          data: {
            title: page.title ?? "",
            summary: page.summary ?? "",
            content: page.content ?? "",
          },
        });
      }
      return;
    }

    // Fetch page data for the selected translation locale
    setTranslationState(null);
    backend.getProfilePage(translationLocale, params.slug, params.pageslug).then((data) => {
      if (cancelled) return;
      if (data === null || (data.locale_code !== undefined && data.locale_code.trimEnd() !== translationLocale)) {
        // No translation exists for this locale (or fallback was returned) — show empty fields
        setTranslationState({ locale: translationLocale, data: undefined });
      } else {
        setTranslationState({
          locale: translationLocale,
          data: {
            title: data.title ?? "",
            summary: data.summary ?? "",
            content: data.content ?? "",
          },
        });
      }
    });

    return () => { cancelled = true; };
  }, [translationLocale, params.locale, params.slug, params.pageslug, page]);

  // Translation locales state
  const [translationLocales, setTranslationLocales] = React.useState<string[] | null>(null);

  // Load translation locales after page data is loaded
  React.useEffect(() => {
    backend.listProfilePageTranslationLocales(params.locale, params.slug, page.id).then((locales) => {
      setTranslationLocales(locales);
    });
  }, [params.locale, params.slug, page.id]);

  // Show skeleton while loading or when translation data is stale (locale mismatch)
  const translationReady = translationState !== null && translationState.locale === translationLocale;

  if (canEdit === null || !translationReady) {
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
  if (canEdit !== true) {
    return (
      <>
        <div className="content">
          <h2>Access Denied</h2>
          <p>You don't have permission to edit this page.</p>
        </div>
      </>
    );
  }

  // translationState.data is undefined when no translation exists for the selected locale
  const isNewTranslation = translationState.data === undefined;

  const initialData: ContentEditorData = {
    title: isNewTranslation ? "" : translationState.data.title,
    slug: page.slug ?? "",
    summary: isNewTranslation ? "" : translationState.data.summary,
    content: isNewTranslation ? "" : translationState.data.content,
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
      throw new Error("Failed to update page");
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

    if (translationResult === null) {
      throw new Error("Failed to save page translation");
    }

    toast.success("Page saved successfully");
    // Refresh translation locales
    backend.listProfilePageTranslationLocales(params.locale, params.slug, page.id).then((locales) => {
      setTranslationLocales(locales);
    });
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

  const handleAutoTranslate = async (targetLocale: string) => {
    const result = await backend.autoTranslateProfilePage(
      params.locale,
      params.slug,
      page.id,
      targetLocale,
      translationLocale,
    );
    if (result === null) {
      throw new Error("Auto-translate failed");
    }
    // Refresh translation locales
    backend.listProfilePageTranslationLocales(params.locale, params.slug, page.id).then((locales) => {
      setTranslationLocales(locales);
    });
  };

  const handleDeleteTranslation = async (localeToDelete: string) => {
    const result = await backend.deleteProfilePageTranslation(
      params.locale,
      params.slug,
      page.id,
      localeToDelete,
    );
    if (result === null) {
      throw new Error("Delete translation failed");
    }
    // Refresh translation locales
    backend.listProfilePageTranslationLocales(params.locale, params.slug, page.id).then((locales) => {
      setTranslationLocales(locales);
    });
    // If deleting the current locale, switch to the site locale
    if (localeToDelete === translationLocale) {
      setTranslationLocale(params.locale);
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
          translationLocales={translationLocales}
          onAutoTranslate={handleAutoTranslate}
          onDeleteTranslation={handleDeleteTranslation}
        />
      </div>
    </>
  );
}
