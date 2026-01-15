// Profile settings index (general settings)
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Card } from "@/components/ui/card";

export const Route = createFileRoute("/$locale/$slug/settings/")({
  component: SettingsIndexPage,
});

function SettingsIndexPage() {
  const { t } = useTranslation();

  return (
    <Card className="p-6">
      <h3 className="text-xl font-semibold mb-4">
        {t("Profile.Profile Settings")}
      </h3>
      <p className="text-muted-foreground">
        {t("Profile.Basic profile information and settings")}
      </p>
    </Card>
  );
}
