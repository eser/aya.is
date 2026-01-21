// Search page
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { useCallback, useState } from "react";
import { BookOpen, FileText, Search, User } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import type { SearchResult } from "@/modules/backend/search/search";

export const Route = createFileRoute("/$locale/search/")({
  validateSearch: (search: Record<string, unknown>) => {
    const q = typeof search.q === "string" ? search.q : "";
    return q.length > 0 ? { q } : {};
  },
  loaderDeps: ({ search }) => ({ q: search.q }),
  loader: async ({ params, deps }) => {
    const { locale } = params;
    const query = deps.q ?? "";

    if (query.length === 0) {
      return { results: null, query: "", locale };
    }

    const results = await backend.search(locale, query);
    return { results, query, locale };
  },
  component: SearchPage,
});

function SearchPage() {
  const { results, query, locale } = Route.useLoaderData();
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [searchInput, setSearchInput] = useState(query);

  const handleSearch = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      if (searchInput.trim().length > 0) {
        navigate({
          to: "/$locale/search",
          params: { locale },
          search: { q: searchInput.trim() },
        });
      }
    },
    [searchInput, locale, navigate],
  );

  const getIcon = (type: string) => {
    switch (type) {
      case "profile":
        return <User className="size-5 text-muted-foreground" />;
      case "story":
        return <BookOpen className="size-5 text-muted-foreground" />;
      case "page":
        return <FileText className="size-5 text-muted-foreground" />;
      default:
        return null;
    }
  };

  const getLink = (result: SearchResult) => {
    switch (result.type) {
      case "profile":
        return `/${locale}/${result.slug}`;
      case "story":
        return result.profile_slug !== null
          ? `/${locale}/${result.profile_slug}/stories/${result.slug}`
          : `/${locale}/stories/${result.slug}`;
      case "page":
        return result.profile_slug !== null
          ? `/${locale}/${result.profile_slug}/${result.slug}`
          : `/${locale}/${result.slug}`;
      default:
        return "#";
    }
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case "profile":
        return t("Search.Profile");
      case "story":
        return t("Search.Story");
      case "page":
        return t("Search.Page");
      default:
        return type;
    }
  };

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <h1>{t("Search.Title", "Search")}</h1>

          <form onSubmit={handleSearch} className="flex gap-2 mb-8">
            <Input
              type="search"
              placeholder={t("Search.Placeholder", "Search for profiles, stories, pages...")}
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              className="flex-1"
            />
            <Button type="submit">
              <Search className="size-4 mr-2" />
              {t("Search.Button", "Search")}
            </Button>
          </form>

          {query.length > 0 && results !== null && (
            <div className="space-y-4">
              {results.length === 0
                ? (
                  <p className="text-muted-foreground">
                    {t("Search.NoResults", 'No results found for "{{query}}"', { query })}
                  </p>
                )
                : (
                  <>
                    <p className="text-sm text-muted-foreground mb-4">
                      {t("Search.ResultCount", "{{count}} results found", { count: results.length })}
                    </p>
                    <div className="space-y-3">
                      {results.map((result) => (
                        <Link
                          key={`${result.type}-${result.id}`}
                          to={getLink(result)}
                          className="block p-4 border rounded-lg hover:bg-accent transition-colors"
                        >
                          <div className="flex items-start gap-3">
                            <div className="mt-0.5">{getIcon(result.type)}</div>
                            <div className="flex-1 min-w-0">
                              <div className="flex items-center gap-2 mb-1">
                                <span className="text-xs px-2 py-0.5 bg-muted rounded">
                                  {getTypeLabel(result.type)}
                                </span>
                                {result.profile_slug !== null && result.type !== "profile" && (
                                  <span className="text-xs text-muted-foreground">
                                    @{result.profile_slug}
                                  </span>
                                )}
                              </div>
                              <h3 className="font-medium truncate">{result.title}</h3>
                              {result.summary !== null && (
                                <p className="text-sm text-muted-foreground line-clamp-2">
                                  {result.summary}
                                </p>
                              )}
                            </div>
                            {result.image_uri !== null && (
                              <img
                                src={result.image_uri}
                                alt=""
                                className="size-12 rounded object-cover"
                              />
                            )}
                          </div>
                        </Link>
                      ))}
                    </div>
                  </>
                )}
            </div>
          )}

          {query.length === 0 && (
            <p className="text-muted-foreground">
              {t("Search.EnterQuery", "Enter a search term to find profiles, stories, and pages.")}
            </p>
          )}
        </div>
      </section>
    </PageLayout>
  );
}
