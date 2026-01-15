import * as React from "react";
import type { Metadata } from "next";

import { backend } from "@/shared/modules/backend/backend.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { ElementsContent } from "./_components/elements-content.tsx";

export async function generateMetadata(): Promise<Metadata> {
  const { t } = await getTranslations();
  return {
    title: t("Layout", "Elements"),
  };
}

async function IndexPage() {
  const { t, locale } = await getTranslations();

  const profiles = await backend.getProfilesByKinds(locale.code, ["individual", "organization"]);

  return (
    <section className="container px-4 py-8 mx-auto">
      <div className="content">
        <h2>{t("Layout", "Elements")}</h2>

        <ElementsContent initialProfiles={profiles!} />
      </div>
    </section>
  );
}

export { IndexPage as default };
