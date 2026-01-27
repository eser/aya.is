// Admin points layout
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { LocaleLink } from "@/components/locale-link";

export const Route = createFileRoute("/$locale/admin/points")({
  component: AdminPointsLayout,
});

function AdminPointsLayout() {
  const { t } = useTranslation();

  return (
    <div className="space-y-6">
      <div className="flex gap-4">
        <LocaleLink
          to="/admin/points"
          activeOptions={{ exact: true }}
          className="text-sm text-muted-foreground hover:text-foreground"
          activeProps={{ className: "text-sm font-medium text-foreground" }}
        >
          {t("Admin.Dashboard")}
        </LocaleLink>
        <LocaleLink
          to="/admin/points/pending"
          className="text-sm text-muted-foreground hover:text-foreground"
          activeProps={{ className: "text-sm font-medium text-foreground" }}
        >
          {t("Admin.Pending Awards")}
        </LocaleLink>
      </div>

      <Outlet />
    </div>
  );
}
