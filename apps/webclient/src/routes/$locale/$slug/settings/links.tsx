// Profile links settings
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Card } from "@/components/ui/card";

export const Route = createFileRoute("/$locale/$slug/settings/links")({
  component: LinksSettingsPage,
});

function LinksSettingsPage() {
  const { t } = useTranslation();

  return (
    <Card className="p-6">
      <h3 className="text-xl font-semibold mb-4">
        {t("Profile.Social Links")}
      </h3>
      <p className="text-muted-foreground">
        {t("Profile.Add social media and website links")}
      </p>
    </Card>
  );
}
