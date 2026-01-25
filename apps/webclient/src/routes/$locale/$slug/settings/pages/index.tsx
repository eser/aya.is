// Profile pages settings
import * as React from "react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { FileText, GripVertical, ExternalLink, Pencil, Plus } from "lucide-react";
import { backend, type ProfilePage } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { LocaleLink } from "@/components/locale-link";

export const Route = createFileRoute("/$locale/$slug/settings/pages/")({
  component: PagesSettingsPage,
});

function PagesSettingsPage() {
  const { t } = useTranslation();
  const params = Route.useParams();

  const [pages, setPages] = React.useState<ProfilePage[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);

  // Drag and drop state
  const [draggedId, setDraggedId] = React.useState<string | null>(null);
  const [dragOverId, setDragOverId] = React.useState<string | null>(null);

  // Load pages on mount
  React.useEffect(() => {
    loadPages();
  }, [params.locale, params.slug]);

  const loadPages = async () => {
    setIsLoading(true);
    const result = await backend.listProfilePages(params.locale, params.slug);
    if (result !== null) {
      // Sort by sort_order
      const sorted = [...result].sort((a, b) => a.sort_order - b.sort_order);
      setPages(sorted);
    } else {
      toast.error(t("Profile.Failed to load profile pages"));
    }
    setIsLoading(false);
  };

  // Drag and drop handlers
  const handleDragStart = (e: React.DragEvent, pageId: string) => {
    setDraggedId(pageId);
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/plain", pageId);
    // Add a slight delay to allow the drag image to be captured
    setTimeout(() => {
      const element = document.querySelector(`[data-page-id="${pageId}"]`);
      if (element !== null) {
        element.classList.add("opacity-50");
      }
    }, 0);
  };

  const handleDragEnd = () => {
    // Remove opacity from dragged element
    if (draggedId !== null) {
      const element = document.querySelector(`[data-page-id="${draggedId}"]`);
      if (element !== null) {
        element.classList.remove("opacity-50");
      }
    }
    setDraggedId(null);
    setDragOverId(null);
  };

  const handleDragOver = (e: React.DragEvent, pageId: string) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
    if (pageId !== draggedId) {
      setDragOverId(pageId);
    }
  };

  const handleDragLeave = () => {
    setDragOverId(null);
  };

  const handleDrop = async (e: React.DragEvent, targetId: string) => {
    e.preventDefault();
    setDragOverId(null);

    if (draggedId === null || draggedId === targetId) return;

    const draggedIndex = pages.findIndex((p) => p.id === draggedId);
    const targetIndex = pages.findIndex((p) => p.id === targetId);

    if (draggedIndex === -1 || targetIndex === -1) return;

    // Reorder locally first for immediate feedback
    const newPages = [...pages];
    const [draggedItem] = newPages.splice(draggedIndex, 1);
    newPages.splice(targetIndex, 0, draggedItem);

    // Update local state immediately
    setPages(newPages);

    // Update orders on backend
    const updatePromises = newPages.map((page, index) => {
      const newOrder = index + 1;
      if (page.sort_order !== newOrder) {
        return backend.updateProfilePage(params.locale, params.slug, page.id, {
          slug: page.slug,
          order: newOrder,
          cover_picture_uri: null,
          published_at: null,
        });
      }
      return Promise.resolve(null);
    });

    try {
      await Promise.all(updatePromises);
      toast.success(t("Profile.Pages reordered successfully"));
      // Reload to get fresh data from server
      loadPages();
    } catch {
      toast.error(t("Profile.Failed to reorder pages"));
      // Reload to restore original order
      loadPages();
    }
  };

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="mb-6">
          <Skeleton className="h-7 w-40 mb-2" />
          <Skeleton className="h-4 w-72" />
        </div>
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="flex items-center gap-3 p-4 border rounded-lg"
            >
              <Skeleton className="size-5" />
              <Skeleton className="size-10 rounded" />
              <div className="flex-1">
                <Skeleton className="h-5 w-32 mb-2" />
                <Skeleton className="h-4 w-24" />
              </div>
              <Skeleton className="size-10" />
            </div>
          ))}
        </div>
      </Card>
    );
  }

  return (
    <Card className="p-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-serif text-xl font-semibold text-foreground">{t("Profile.Pages")}</h3>
          <p className="text-muted-foreground text-sm mt-1">
            {t("Profile.Manage and reorder your profile pages.")}
          </p>
        </div>
        <Link
          to="/$locale/$slug/settings/pages/new"
          params={{ locale: params.locale, slug: params.slug }}
        >
          <Button variant="default" size="sm">
            <Plus className="mr-1.5 size-4" />
            {t("Editor.Add Page")}
          </Button>
        </Link>
      </div>

      {pages.length === 0 ? (
        <div className="text-center py-12 border-2 border-dashed rounded-lg">
          <FileText className="size-12 mx-auto text-muted-foreground mb-4" />
          <p className="text-muted-foreground">{t("Profile.No pages found.")}</p>
        </div>
      ) : (
        <div className="space-y-2">
          {pages.map((page) => {
            const isDragOver = dragOverId === page.id;

            return (
              <div
                key={page.id}
                data-page-id={page.id}
                draggable
                onDragStart={(e) => handleDragStart(e, page.id)}
                onDragEnd={handleDragEnd}
                onDragOver={(e) => handleDragOver(e, page.id)}
                onDragLeave={handleDragLeave}
                onDrop={(e) => handleDrop(e, page.id)}
                className={`flex items-center gap-3 p-4 border rounded-lg transition-all cursor-move select-none ${
                  isDragOver
                    ? "border-primary bg-primary/5 border-dashed"
                    : "hover:bg-muted/50"
                }`}
              >
                <div className="flex items-center justify-center text-muted-foreground hover:text-foreground cursor-grab active:cursor-grabbing">
                  <GripVertical className="size-5" />
                </div>
                <div className="flex items-center justify-center size-10 rounded bg-muted shrink-0">
                  <FileText className="size-5" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{page.title}</p>
                  <p className="text-sm text-muted-foreground">/{page.slug}</p>
                </div>
                <div className="flex items-center gap-1">
                  <LocaleLink to={`/${params.slug}/${page.slug}/edit`}>
                    <Button variant="ghost" size="icon">
                      <Pencil className="size-4" />
                    </Button>
                  </LocaleLink>
                  <LocaleLink to={`/${params.slug}/${page.slug}`}>
                    <Button variant="ghost" size="icon">
                      <ExternalLink className="size-4" />
                    </Button>
                  </LocaleLink>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </Card>
  );
}
