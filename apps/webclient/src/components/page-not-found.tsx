import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";

export function PageNotFound() {
  const { t } = useTranslation();

  return (
    <PageLayout>
      <div className="container mx-auto py-16 px-4 text-center">
        <h1 className="font-serif text-4xl font-bold mb-4">{t("Layout.Page not found")}</h1>
        <p className="text-muted-foreground">
          {t("Layout.The page you are looking for does not exist. Please check your spelling and try again.")}
        </p>
      </div>
    </PageLayout>
  );
}
