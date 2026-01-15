// Profile contributions page
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";

export const Route = createFileRoute("/$locale/$slug/contributions")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const contributions = await backend.getProfileContributions(locale, slug);
    return { contributions };
  },
  component: ContributionsPage,
});

function ContributionsPage() {
  const { contributions } = Route.useLoaderData();
  const { t } = useTranslation();

  return (
    <div className="content">
      <h2>{t("Layout.Contributions")}</h2>
      <p className="text-muted-foreground mb-4">
        {t("Contributions.A collection of open source projects and organizations.")}
      </p>

      {contributions && contributions.length > 0 ? (
        <div className="space-y-4">
          {/* Contributions will be rendered here */}
        </div>
      ) : (
        <p className="text-muted-foreground">
          {t("Contributions.No contributions found.")}
        </p>
      )}
    </div>
  );
}
