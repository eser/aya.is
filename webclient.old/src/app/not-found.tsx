import * as React from "react";

import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";

export async function NotFoundPage() {
  const { t } = await getTranslations();

  return (
    <section className="container px-4 py-8 mx-auto">
      <div className="content">
        <h2>{t("Layout", "Page not found")}</h2>

        <p>
          {t("Layout", "The page you are looking for does not exist. Please check your spelling and try again.")}
        </p>
      </div>
    </section>
  );
}

export { NotFoundPage as default };
