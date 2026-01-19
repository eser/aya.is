// Profile settings layout
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { LocaleLink } from "@/components/locale-link";
import { backend } from "@/modules/backend/backend";

export const Route = createFileRoute("/$locale/$slug/settings")({
  loader: async ({ params }) => {
    const { locale, slug } = params;

    // Fetch profile and permissions in parallel
    const [profile, permissions] = await Promise.all([
      backend.getProfile(locale, slug),
      backend.getProfilePermissions(locale, slug),
    ]);

    // Redirect if profile doesn't exist or user can't edit
    if (profile === null || permissions === null || !permissions.can_edit) {
      throw redirect({ to: `/${locale}/${slug}` });
    }

    return { profile, permissions };
  },
  component: SettingsLayout,
});

function SettingsLayout() {
  const { t } = useTranslation();
  const params = Route.useParams();

  return (
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
  );
}
