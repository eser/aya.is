// Profile custom page
import { createFileRoute, Link, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { SquarePen } from "lucide-react";
import { backend } from "@/modules/backend/backend";
import { TextContent } from "@/components/widgets/text-content";
import { compileMdx } from "@/lib/mdx";
import { Button } from "@/components/ui/button";
import { useProfilePermissions } from "@/lib/hooks/use-profile-permissions";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { buildUrl, generateMetaTags, truncateDescription } from "@/lib/seo";
import { ChildNotFound } from "../route";

const profileRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/$pageslug/")({
  loader: async ({ params }) => {
    const { locale, slug, pageslug } = params;

    // Skip if pageslug matches known routes
    const knownRoutes = ["stories", "settings", "members", "contributions"];
    if (knownRoutes.includes(pageslug)) {
      return { page: null, compiledContent: null, notFound: true, locale, slug, pageslug };
    }

    const page = await backend.getProfilePage(locale, slug, pageslug);

    if (page === null || page === undefined) {
      return { page: null, compiledContent: null, notFound: true, locale, slug, pageslug };
    }

    // Compile MDX content on the server
    let compiledContent: string | null = null;
    if (page.content !== null && page.content !== undefined) {
      try {
        compiledContent = await compileMdx(page.content);
      } catch (error) {
        console.error("Failed to compile MDX:", error);
        compiledContent = null;
      }
    }

    return { page, compiledContent, notFound: false, locale, slug, pageslug };
  },
  head: ({ loaderData }) => {
    const { page, locale, slug, pageslug } = loaderData;
    if (page === null) {
      return { meta: [] };
    }
    return {
      meta: generateMetaTags({
        title: page.title,
        description: truncateDescription(page.summary),
        url: buildUrl(locale, slug, pageslug),
        image: page.cover_picture_uri,
        locale,
        type: "article",
        publishedTime: page.published_at,
      }),
    };
  },
  component: ProfileCustomPage,
  notFoundComponent: ChildNotFound,
});

function ProfileCustomPage() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const loaderData = Route.useLoaderData();
  const { profile, permissions } = profileRoute.useLoaderData();
  const { canEdit } = useProfilePermissions(profile?.id ?? "");

  // If notFound flag is set, render 404 page
  if (loaderData.notFound || loaderData.page === null || profile === null) {
    return <ChildNotFound />;
  }

  const { page, compiledContent, locale, slug } = loaderData;

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale} viewerMembershipKind={permissions?.viewer_membership_kind}>
      <div className="relative">
        {canEdit && (
          <Link
            to="/$locale/$slug/$pageslug/edit"
            params={{
              locale: params.locale,
              slug: params.slug,
              pageslug: params.pageslug,
            }}
            className="absolute right-0 top-0 z-10"
          >
            <Button variant="outline" size="sm">
              <SquarePen className="mr-1.5 size-4" />
              {t("ContentEditor.Edit Page")}
            </Button>
          </Link>
        )}
        <TextContent
          title={page.title}
          compiledContent={compiledContent}
          rawContent={page.content}
          headingOffset={2}
        />
      </div>
    </ProfileSidebarLayout>
  );
}
