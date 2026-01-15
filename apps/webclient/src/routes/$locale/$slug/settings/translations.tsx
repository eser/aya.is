// Profile translations settings
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Card } from "@/components/ui/card";

export const Route = createFileRoute("/$locale/$slug/settings/translations")({
  component: TranslationsSettingsPage,
});

function TranslationsSettingsPage() {
  const { t } = useTranslation();

  return (
    <Card className="p-6">
      <h3 className="text-xl font-semibold mb-4">
        {t("Profile.Translations")}
      </h3>
      <p className="text-muted-foreground">
        {t("Profile.Manage your profile in multiple languages")}
      </p>
    </Card>
  );
}
