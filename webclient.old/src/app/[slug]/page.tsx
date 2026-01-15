import * as React from "react";
import type { Metadata, ResolvingMetadata } from "next";
import { notFound } from "next/navigation";

import { backend } from "@/shared/modules/backend/backend.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";

import { Timeline } from "./_components/timeline.tsx";
import styles from "./page.module.css";

type IndexPageProps = {
  params: Promise<{
    slug: string;
  }>;
};

// TODO(@eser) add more from https://beta.nextjs.org/docs/api-reference/metadata
export async function generateMetadata(props: IndexPageProps, _parent: ResolvingMetadata): Promise<Metadata> {
  const params = await props.params;

  const { locale } = await getTranslations();

  const profileData = await backend.getProfile(locale.code, params.slug);
  if (profileData === null) {
    notFound();
  }

  return {
    title: profileData.title,
    description: profileData.description,
  };
}

async function IndexPage(props: IndexPageProps) {
  const params = await props.params;

  const { locale } = await getTranslations();

  const storiesData = await backend.getProfileStories(locale.code, params.slug);

  if (storiesData === null) {
    notFound();
  }

  return (
    <div className={styles["timeline-container"]}>
      <Timeline profileSlug={params.slug} stories={storiesData} />
    </div>
  );
}

export { IndexPage as default };
