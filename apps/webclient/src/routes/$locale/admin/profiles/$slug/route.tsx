// Admin profile edit layout with tabs
import { createFileRoute, Outlet, notFound } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { LocaleLink } from "@/components/locale-link";
import { backend } from "@/modules/backend/backend";
import { SiteAvatar } from "@/components/userland";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ExternalLink } from "lucide-react";

export const Route = createFileRoute("/$locale/admin/profiles/$slug")({
  loader: async ({ params }) => {
    const { locale, slug } = params;

    const profile = await backend.getAdminProfile(locale, slug);

    if (profile === null) {
      throw notFound();
    }

    return { profile, locale, slug };
  },
  component: AdminProfileLayout,
});

function AdminProfileLayout() {
  const { t } = useTranslation();
  const { profile, slug } = Route.useLoaderData();
  const params = Route.useParams();

  const getKindBadge = (profileKind: string) => {
    switch (profileKind) {
      case "individual":
        return <Badge variant="secondary">{t("Admin.Individual")}</Badge>;
      case "organization":
        return <Badge variant="default">{t("Admin.Organization")}</Badge>;
      case "product":
        return (
          <Badge variant="outline" className="border-blue-500 text-blue-500">
            {t("Admin.Product")}
          </Badge>
        );
      default:
        return <Badge variant="outline">{profileKind}</Badge>;
    }
  };

  const displayTitle =
    profile.has_translation === false
      ? t("Admin.no translation found")
      : profile.title;

  return (
    <div className="space-y-6">
      {/* Profile Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <SiteAvatar
            src={profile.profile_picture_uri}
            name={profile.title}
            fallbackName={profile.slug}
            size="lg"
          />
          <div>
            <div className="flex items-center gap-2">
              <h2
                className={`font-serif text-2xl font-bold ${profile.has_translation === false ? "italic text-muted-foreground" : ""}`}
              >
                {displayTitle}
              </h2>
              {getKindBadge(profile.kind)}
            </div>
            <p className="text-muted-foreground font-mono">@{profile.slug}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <LocaleLink to={`/${profile.slug}`}>
            <Button variant="outline" size="sm">
              <ExternalLink className="h-4 w-4 mr-2" />
              {t("Admin.View Public Profile")}
            </Button>
          </LocaleLink>
          <LocaleLink to={`/${profile.slug}/settings`}>
            <Button variant="outline" size="sm">
              {t("Admin.Edit Profile Settings")}
            </Button>
          </LocaleLink>
        </div>
      </div>

      {/* Navigation Tabs */}
      <nav className="flex gap-4 border-b">
        <LocaleLink
          to={`/admin/profiles/${slug}`}
          activeOptions={{ exact: true }}
          className="relative pb-2 text-sm font-medium text-muted-foreground hover:text-foreground"
          activeProps={{
            className:
              "relative pb-2 text-sm font-medium text-foreground after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-foreground",
          }}
        >
          {t("Admin.General")}
        </LocaleLink>
        <LocaleLink
          to={`/admin/profiles/${slug}/points`}
          className="relative pb-2 text-sm font-medium text-muted-foreground hover:text-foreground"
          activeProps={{
            className:
              "relative pb-2 text-sm font-medium text-foreground after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-foreground",
          }}
        >
          {t("Admin.Points")}
        </LocaleLink>
      </nav>

      {/* Content */}
      <Outlet />
    </div>
  );
}
