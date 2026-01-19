// Profile index - shows profile stories/timeline with date grouping and pagination
import * as React from "react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { backend } from "@/modules/backend/backend";
import { StoriesPageClient } from "@/routes/$locale/stories/_components/-stories-page-client";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth/auth-context";

export const Route = createFileRoute("/$locale/$slug/")({
  validateSearch: (search: Record<string, unknown>) => {
    const offset = Number(search.offset) || 0;
    return offset > 0 ? { offset } : {};
  },
  loader: async ({ params }) => {
    const { slug, locale } = params;
    const stories = await backend.getProfileStories(locale, slug);
    return { stories, slug, locale };
  },
  component: ProfileIndexPage,
});

function ProfileIndexPage() {
  const { t } = useTranslation();
  const auth = useAuth();
  const { stories, slug, locale } = Route.useLoaderData();
  const [canEdit, setCanEdit] = React.useState(false);

  // Check edit permissions
  React.useEffect(() => {
    if (auth.isAuthenticated && !auth.isLoading) {
      backend.getProfilePermissions(locale, slug).then((perms) => {
        if (perms !== null) {
          setCanEdit(perms.can_edit);
        }
      });
    } else {
      setCanEdit(false);
    }
  }, [auth.isAuthenticated, auth.isLoading, locale, slug]);

  return (
    <div className="content">
      {canEdit && (
        <div className="flex justify-end mb-4">
          <Link
            to="/$locale/$slug/stories/new"
            params={{ locale, slug }}
          >
            <Button variant="default" size="sm">
              <Plus className="mr-1.5 size-4" />
              {t("Editor.Add Story")}
            </Button>
          </Link>
        </div>
      )}
      <StoriesPageClient
        initialStories={stories}
        basePath={`/${locale}/${slug}`}
        profileSlug={slug}
      />
    </div>
  );
}
