// Profile settings layout
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { LocaleLink } from "@/components/locale-link";

export const Route = createFileRoute("/$locale/$slug/settings")({
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
          {t("Profile.Manage your profile information, links, and preferences.")}
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
          to={`/${params.slug}/settings/translations`}
          className="text-sm font-medium hover:text-foreground"
        >
          {t("Profile.Translations")}
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
