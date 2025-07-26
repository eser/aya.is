import type { Metadata } from "next";

import { backend } from "@/shared/modules/backend/backend.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";

import { StoriesPageClient } from "./_components/stories-page-client.tsx";

export async function generateMetadata(): Promise<Metadata> {
  const { t } = await getTranslations();

  return {
    title: t("Layout", "Articles"),
  };
}

async function IndexPage() {
  const { locale } = await getTranslations();

  const stories = await backend.getStoriesByKinds(locale.code, ["article"]);

  return <StoriesPageClient initialStories={stories} />;
}

export { IndexPage as default };
