// Profile custom page
import * as React from "react";
import { createFileRoute, Link, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { PencilLine } from "lucide-react";
import { backend } from "@/modules/backend/backend";
import { TextContent } from "@/components/widgets/text-content";
import { compileMdx } from "@/lib/mdx";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth/auth-context";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";

const profileRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/$pageslug/")({
  loader: async ({ params }) => {
    const { locale, slug, pageslug } = params;

    // Skip if pageslug matches known routes
    const knownRoutes = ["stories", "settings", "members", "contributions"];
    if (knownRoutes.includes(pageslug)) {
      return { page: null, compiledContent: null, notFound: true, locale, slug };
    }

    const page = await backend.getProfilePage(locale, slug, pageslug);

    if (page === null || page === undefined) {
      return { page: null, compiledContent: null, notFound: true, locale, slug };
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

    return { page, compiledContent, notFound: false, locale, slug };
  },
  component: ProfileCustomPage,
  notFoundComponent: PageNotFound,
});

function ProfileCustomPage() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const auth = useAuth();
  const loaderData = Route.useLoaderData();
  const { profile } = profileRoute.useLoaderData();
  const [canEdit, setCanEdit] = React.useState(false);

  // If notFound flag is set, render 404 page
  if (loaderData.notFound || loaderData.page === null || profile === null) {
    return <PageNotFound />;
  }

  const { page, compiledContent, locale, slug } = loaderData;

  // Check edit permissions (pages use profile permissions)
  React.useEffect(() => {
    if (auth.isAuthenticated && !auth.isLoading) {
      backend
        .getProfilePermissions(params.locale, params.slug)
        .then((perms) => {
          if (perms !== null) {
            setCanEdit(perms.can_edit);
          }
        });
    } else {
      setCanEdit(false);
    }
  }, [auth.isAuthenticated, auth.isLoading, params.locale, params.slug]);

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale}>
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
              <PencilLine className="mr-1.5 size-4" />
              {t("Editor.Edit Page")}
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

function PageNotFound() {
  const { t } = useTranslation();

  return (
    <div className="content">
      <h2>{t("Layout.Page not found")}</h2>
      <p className="text-muted-foreground">
        {t(
          "Layout.The page you are looking for does not exist. Please check your spelling and try again.",
        )}
      </p>
    </div>
  );
}
