import * as React from "react";
import type { Metadata, ResolvingMetadata } from "next";
import { notFound } from "next/navigation";

import { mdx } from "@/shared/lib/mdx.tsx";
import { backend } from "@/shared/modules/backend/backend.ts";
import { siteConfig } from "@/shared/config.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";

import { components } from "@/shared/components/userland/userland.ts";
import { StoryMetadata } from "@/shared/components/widgets/story/story-metadata.tsx";
import { StoryFooter } from "@/shared/components/widgets/story/story-footer.tsx";
import { StoryShareWrapper } from "@/shared/components/widgets/story/story-share-wrapper.tsx";

export const revalidate = 300;

type IndexPageProps = {
  params: Promise<{
    storyslug: string;
  }>;
};

// TODO(@eser) add more from https://beta.nextjs.org/docs/api-reference/metadata
export async function generateMetadata(props: IndexPageProps, _parent: ResolvingMetadata): Promise<Metadata> {
  const params = await props.params;

  const { locale } = await getTranslations();

  const storyData = await backend.getStory(locale.code, params.storyslug);

  if (storyData === null) {
    notFound();
  }

  return {
    title: storyData.title,
    description: storyData.summary,
  };
}

async function IndexPage(props: IndexPageProps) {
  const params = await props.params;

  const { locale } = await getTranslations();

  const storyData = await backend.getStory(locale.code, params.storyslug);

  if (storyData === null) {
    notFound();
  }

  const mdxSource = await mdx(storyData.content, components);

  // Get current URL for sharing
  const baseUrl = siteConfig.host;
  const currentUrl = `${baseUrl}/stories/${params.storyslug}`;

  return (
    <section className="container px-4 py-8 mx-auto">
      <div className="content">
        <h2>{storyData.title}</h2>

        <StoryMetadata story={storyData} />

        <article>{mdxSource?.content}</article>

        <StoryShareWrapper story={storyData} currentUrl={currentUrl} />

        <StoryFooter story={storyData} />
      </div>
    </section>
  );
}

export { IndexPage as default };
