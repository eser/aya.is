// Catch-all route for any unmatched paths under profile
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

export const Route = createFileRoute("/$locale/$slug/$")({
  component: CatchAllNotFound,
});

function CatchAllNotFound() {
  const { t } = useTranslation();

  return (
    <div className="content">
      <h2>{t("Layout.Page not found")}</h2>
      <p className="text-muted-foreground">
        {t(
          "Layout.The page you are looking for does not exist. Please check your spelling and try again.",
        )}
      </p>
    </div>
  );
}
