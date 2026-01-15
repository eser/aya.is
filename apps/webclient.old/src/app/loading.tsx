import * as React from "react";

import { Astronaut } from "@/shared/components/widgets/astronaut/astronaut.tsx";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";

export async function Loading() {
  const { t } = await getTranslations();

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Astronaut width={400} height={400} />
      {t("Loading", "Loading...")}
    </div>
  );
}

export { Loading as default };
