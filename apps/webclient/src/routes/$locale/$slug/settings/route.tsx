// Profile settings layout
import { createFileRoute, Outlet, redirect, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { LocaleLink } from "@/components/locale-link";
import { backend } from "@/modules/backend/backend";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";

const profileRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/settings")({
  loader: async ({ params }) => {
    const { locale, slug } = params;

    // Fetch permissions
    const permissions = await backend.getProfilePermissions(locale, slug);

    // Redirect if user can't edit
    if (permissions === null || !permissions.can_edit) {
      throw redirect({ to: `/${locale}/${slug}` });
    }

    return { permissions, locale, slug };
  },
  component: SettingsLayout,
});

function SettingsLayout() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const { locale, slug } = Route.useLoaderData();
  const { profile } = profileRoute.useLoaderData();

  if (profile === null) {
    return null;
  }

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale}>
      <div className="space-y-6">
        <div>
          <h2 className="font-serif text-2xl font-bold text-foreground">{t("Profile.Settings")}</h2>
          <p className="text-muted-foreground">
            {t(
              "Profile.Manage your profile information, links, and preferences.",
            )}
          </p>
        </div>

        <nav className="flex gap-4 border-b">
          <LocaleLink
            to={`/${params.slug}/settings`}
            activeOptions={{ exact: true }}
            className="relative pb-2 text-sm font-medium text-muted-foreground hover:text-foreground"
            activeProps={{
              className:
                "relative pb-2 text-sm font-medium text-foreground after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-foreground",
            }}
          >
            {t("Profile.General")}
          </LocaleLink>
          <LocaleLink
            to={`/${params.slug}/settings/pages`}
            className="relative pb-2 text-sm font-medium text-muted-foreground hover:text-foreground"
            activeProps={{
              className:
                "relative pb-2 text-sm font-medium text-foreground after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-foreground",
            }}
          >
            {t("Profile.Pages")}
          </LocaleLink>
          <LocaleLink
            to={`/${params.slug}/settings/links`}
            className="relative pb-2 text-sm font-medium text-muted-foreground hover:text-foreground"
            activeProps={{
              className:
                "relative pb-2 text-sm font-medium text-foreground after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-foreground",
            }}
          >
            {t("Profile.Social Links")}
          </LocaleLink>
        </nav>

        <Outlet />
      </div>
    </ProfileSidebarLayout>
  );
}
