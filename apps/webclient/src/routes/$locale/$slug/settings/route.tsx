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
        <h2 className="text-2xl font-bold">{t("Profile.Settings")}</h2>
        <p className="text-muted-foreground">
          {t(
            "Profile.Manage your profile information, links, and preferences.",
          )}
        </p>
      </div>

      <nav className="flex gap-4 border-b pb-2">
        <LocaleLink
          to={`/${params.slug}/settings`}
          className="text-sm font-medium hover:text-foreground"
        >
          {t("Profile.General")}
        </LocaleLink>
        <LocaleLink
          to={`/${params.slug}/settings/links`}
          className="text-sm font-medium hover:text-foreground"
        >
          {t("Profile.Social Links")}
        </LocaleLink>
      </nav>

      <Outlet />
    </div>
  );
}
