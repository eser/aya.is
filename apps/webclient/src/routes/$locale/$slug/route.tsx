// Profile route - loads profile data and passes through to children
import { CatchNotFound, createFileRoute, ErrorComponent, getRouteApi, Outlet, useMatches } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { PageLayout } from "@/components/page-layouts/default";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { buildUrl, generateMetaTags, truncateDescription } from "@/lib/seo";
import { PageNotFound } from "@/components/page-not-found";

export const Route = createFileRoute("/$locale/$slug")({
  beforeLoad: ({ params }) => {
    const { slug } = params;
    return { profileSlug: slug };
  },
  loader: async ({ params }) => {
    const { slug, locale } = params;

    const profile = await backend.getProfile(locale, slug);

    if (profile === null) {
      return { profile: null, notFound: true, locale, slug, permissions: null };
    }

    // Returns null for unauthenticated users (401 → null via fetcher)
    const permissions = await backend.getProfilePermissions(locale, slug).catch(() => null);

    return { profile, notFound: false, locale, slug, permissions };
  },
  head: ({ loaderData }) => {
    const { profile, locale, slug } = loaderData;
    if (profile === null) {
      return { meta: [] };
    }
    return {
      meta: generateMetaTags({
        title: profile.title,
        description: truncateDescription(profile.description),
        url: buildUrl(locale, slug),
        image: profile.profile_picture_uri,
        locale,
        type: "profile",
      }),
    };
  },
  component: ProfileRoute,
  errorComponent: ErrorComponent,
  notFoundComponent: PageNotFound,
});

const parentRoute = getRouteApi("/$locale/$slug");

export function ChildNotFound() {
  const { profile, permissions, locale, slug } = parentRoute.useLoaderData();
  const { t } = useTranslation();

  const notFoundContent = (
    <div className="content">
      <h2>{t("Layout.Page not found")}</h2>
      <p className="text-muted-foreground">
        {t("Layout.The page you are looking for does not exist. Please check your spelling and try again.")}
      </p>
    </div>
  );

  if (profile === null) {
    return notFoundContent;
  }

  return (
    <ProfileSidebarLayout
      profile={profile}
      slug={slug}
      locale={locale}
      viewerMembershipKind={permissions?.viewer_membership_kind}
    >
      {notFoundContent}
    </ProfileSidebarLayout>
  );
}

function ProfileRoute() {
  const loaderData = Route.useLoaderData();
  const matches = useMatches();

  // If notFound flag is set, render 404 page
  if (loaderData.notFound || loaderData.profile === null) {
    return <PageNotFound />;
  }

  // Check if we're on an edit/new route - these need full-width layout without section wrapper
  const isFullWidthRoute = matches.some((match) =>
    match.pathname.endsWith("/edit") || match.pathname.endsWith("/settings/pages/new")
  );

  if (isFullWidthRoute) {
    return (
      <PageLayout fullHeight>
        <CatchNotFound fallback={<ChildNotFound />}>
          <Outlet />
        </CatchNotFound>
      </PageLayout>
    );
  }

  // Just pass through - children handle their own layout
  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <CatchNotFound fallback={<ChildNotFound />}>
          <Outlet />
        </CatchNotFound>
      </section>
    </PageLayout>
  );
}
