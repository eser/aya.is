// Profile route - loads profile data and passes through to children
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { PageLayout } from "@/components/page-layouts/default";

export const Route = createFileRoute("/$locale/$slug")({
  beforeLoad: ({ params }) => {
    const { slug } = params;
    return { profileSlug: slug };
  },
  loader: async ({ params }) => {
    const { slug, locale } = params;

    const profile = await backend.getProfile(locale, slug);

    if (profile === null) {
      return { profile: null, notFound: true };
    }

    return { profile, notFound: false };
  },
  component: ProfileRoute,
  notFoundComponent: ProfileNotFound,
});

function ProfileNotFound() {
  const { t } = useTranslation();

  return (
    <PageLayout>
      <div className="container mx-auto py-16 px-4 text-center">
        <h1 className="font-serif text-4xl font-bold mb-4">{t("Layout.Page not found")}</h1>
        <p className="text-muted-foreground">
          {t("Layout.The page you are looking for does not exist. Please check your spelling and try again.")}
        </p>
      </div>
    </PageLayout>
  );
}

function ProfileRoute() {
  const loaderData = Route.useLoaderData();

  // If notFound flag is set, render 404 page
  if (loaderData.notFound || loaderData.profile === null) {
    return <ProfileNotFound />;
  }

  // Just pass through - children handle their own layout
  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <Outlet />
      </section>
    </PageLayout>
  );
}
