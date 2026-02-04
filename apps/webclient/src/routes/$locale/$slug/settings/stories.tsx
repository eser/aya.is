// Profile stories settings
import * as React from "react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  ExternalLink,
  FileText,
  Globe,
  GlobeLock,
  ImagePlus,
  Images,
  Info,
  Megaphone,
  Newspaper,
  Pencil,
  PencilLine,
  Plus,
  Presentation,
} from "lucide-react";
import { backend, type StoryEx, type StoryKind } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { LocaleLink } from "@/components/locale-link";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { formatDateString } from "@/lib/date";

const ITEMS_PER_PAGE = 10;

const storyKindIcons: Record<StoryKind, React.ElementType> = {
  news: Newspaper,
  article: PencilLine,
  announcement: Megaphone,
  status: Info,
  content: Images,
  presentation: Presentation,
};

export const Route = createFileRoute("/$locale/$slug/settings/stories")({
  validateSearch: (search: Record<string, unknown>) => {
    const page = Number(search.page) || 1;
    return page > 1 ? { page } : {};
  },
  component: StoriesSettingsPage,
});

function StoriesSettingsPage() {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;
  const params = Route.useParams();
  const search = Route.useSearch();
  const currentPage = search.page ?? 1;

  const [stories, setStories] = React.useState<StoryEx[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);

  // Load stories on mount
  React.useEffect(() => {
    loadStories();
  }, [params.locale, params.slug]);

  const loadStories = async () => {
    setIsLoading(true);
    const result = await backend.getProfileAuthoredStories(params.locale, params.slug);
    if (result !== null) {
      // Sort by created_at descending
      const sorted = [...result].sort((a, b) => {
        const dateA = a.created_at !== null ? new Date(a.created_at).getTime() : 0;
        const dateB = b.created_at !== null ? new Date(b.created_at).getTime() : 0;
        return dateB - dateA;
      });
      setStories(sorted);
    }
    setIsLoading(false);
  };

  const totalPages = Math.ceil(stories.length / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const endIndex = startIndex + ITEMS_PER_PAGE;
  const currentStories = stories.slice(startIndex, endIndex);

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
          <h3 className="font-serif text-xl font-semibold text-foreground">{t("Layout.Stories")}</h3>
          <p className="text-muted-foreground text-sm mt-1">
            {t("Profile.Manage your stories and articles.")}
          </p>
        </div>
        <Link to="/$locale/stories/new" params={{ locale: params.locale }}>
          <Button variant="default" size="sm">
            <Plus className="mr-1.5 size-4" />
            {t("ContentEditor.Add Story")}
          </Button>
        </Link>
      </div>

      {stories.length === 0
        ? (
          <div className="text-center py-12 border-2 border-dashed rounded-lg">
            <FileText className="size-12 mx-auto text-muted-foreground mb-4" />
            <p className="text-muted-foreground">{t("Profile.No stories found.")}</p>
          </div>
        )
        : (
          <>
            <div className="space-y-2">
              {currentStories.map((story) => {
                const KindIcon = storyKindIcons[story.kind];
                const isPublished = story.publications.length > 0;

                return (
                  <div
                    key={story.id}
                    className="flex items-center gap-3 p-4 border rounded-lg hover:bg-muted/50 transition-colors"
                  >
                    <div className="flex items-center justify-center size-10 rounded bg-muted shrink-0">
                      {KindIcon !== undefined ? <KindIcon className="size-5" /> : <FileText className="size-5" />}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <p className="font-medium truncate">{story.title}</p>
                        {isPublished
                          ? <Globe className="size-3.5 text-green-600" />
                          : <GlobeLock className="size-3.5 text-yellow-600" />}
                      </div>
                      <p className="text-sm text-muted-foreground">
                        {story.created_at !== null && formatDateString(story.created_at, locale)}
                      </p>
                    </div>
                    <div className="flex items-center gap-1">
                      <LocaleLink to={`/stories/${story.slug}/edit`}>
                        <Button variant="ghost" size="icon" title={t("Common.Edit")}>
                          <Pencil className="size-4" />
                        </Button>
                      </LocaleLink>
                      <LocaleLink to={`/stories/${story.slug}/cover`}>
                        <Button variant="ghost" size="icon" title={t("CoverDesigner.Design Cover")}>
                          <ImagePlus className="size-4" />
                        </Button>
                      </LocaleLink>
                      <LocaleLink to={`/stories/${story.slug}`}>
                        <Button variant="ghost" size="icon" title={t("Common.View")}>
                          <ExternalLink className="size-4" />
                        </Button>
                      </LocaleLink>
                    </div>
                  </div>
                );
              })}
            </div>

            {totalPages > 1 && (
              <div className="mt-6 flex justify-center">
                <Pagination>
                  <PaginationContent>
                    <PaginationItem>
                      <PaginationPrevious
                        className={currentPage <= 1 ? "pointer-events-none opacity-50" : ""}
                        render={(linkProps) => (
                          <Link
                            to="/$locale/$slug/settings/stories"
                            params={{ locale: params.locale, slug: params.slug }}
                            search={currentPage > 2 ? { page: currentPage - 1 } : {}}
                            {...linkProps}
                          />
                        )}
                      >
                        {t("Common.Previous")}
                      </PaginationPrevious>
                    </PaginationItem>

                    {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                      <PaginationItem key={page}>
                        <PaginationLink
                          isActive={currentPage === page}
                          render={(linkProps) => (
                            <Link
                              to="/$locale/$slug/settings/stories"
                              params={{ locale: params.locale, slug: params.slug }}
                              search={page > 1 ? { page } : {}}
                              {...linkProps}
                            />
                          )}
                        >
                          {page}
                        </PaginationLink>
                      </PaginationItem>
                    ))}

                    <PaginationItem>
                      <PaginationNext
                        className={currentPage >= totalPages ? "pointer-events-none opacity-50" : ""}
                        render={(linkProps) => (
                          <Link
                            to="/$locale/$slug/settings/stories"
                            params={{ locale: params.locale, slug: params.slug }}
                            search={{ page: currentPage + 1 }}
                            {...linkProps}
                          />
                        )}
                      >
                        {t("Common.Next")}
                      </PaginationNext>
                    </PaginationItem>
                  </PaginationContent>
                </Pagination>
              </div>
            )}
          </>
        )}
    </Card>
  );
}
