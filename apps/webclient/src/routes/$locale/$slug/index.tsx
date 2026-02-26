// Profile index - shows profile stories/timeline with date grouping and pagination
import { createFileRoute, Link, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { backend } from "@/modules/backend/backend";
import { StoriesPageClient } from "@/routes/$locale/stories/_components/-stories-page-client";
import { Button } from "@/components/ui/button";
import { useProfilePermissions } from "@/lib/hooks/use-profile-permissions";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { buildUrl, generateCanonicalLink } from "@/lib/seo";

const parentRoute = getRouteApi("/$locale/$slug");

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
  head: ({ loaderData }) => {
    const { locale, slug } = loaderData;
    return {
      links: [generateCanonicalLink(buildUrl(locale, slug))],
    };
  },
  component: ProfileIndexPage,
});

function ProfileIndexPage() {
  const { t } = useTranslation();
  const { stories, slug, locale } = Route.useLoaderData();
  const { profile, permissions } = parentRoute.useLoaderData();
  const { canEdit } = useProfilePermissions(profile?.id ?? "");

  if (profile === null) {
    return null;
  }

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale} viewerMembershipKind={permissions?.viewer_membership_kind}>
      <div className="content relative">
        {canEdit && (
          <Link
            to="/$locale/stories/new"
            params={{ locale }}
            className="absolute right-0 top-0 z-10"
          >
            <Button variant="default" size="sm">
              <Plus className="mr-1.5 size-4" />
              {t("ContentEditor.Add Story")}
            </Button>
          </Link>
        )}
        <StoriesPageClient
          initialStories={stories}
          basePath={`/${locale}/${slug}`}
          profileSlug={slug}
        />
      </div>
    </ProfileSidebarLayout>
  );
}
