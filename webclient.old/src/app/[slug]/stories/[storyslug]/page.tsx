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

// export const revalidate = 300;

type IndexPageProps = {
  params: Promise<{
    slug: string;
    storyslug: string;
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

  const storiesData = await backend.getProfileStories(locale.code, profileData.slug);

  if (storiesData === null) {
    notFound();
  }

  const storyData = storiesData.find((story) => story.slug === params.storyslug);
  if (storyData === undefined) {
    notFound();
  }

  return {
    title: `${profileData.title} - ${storyData.title}`,
    description: storyData.summary,
  };
}

async function IndexPage(props: IndexPageProps) {
  const params = await props.params;

  const { locale } = await getTranslations();

  const profileData = await backend.getProfile(locale.code, params.slug);

  if (profileData === null) {
    notFound();
  }

  const storyData = await backend.getProfileStory(locale.code, profileData.slug, params.storyslug);

  if (storyData === null) {
    notFound();
  }

  const mdxSource = await mdx(storyData.content, components);

  // Get current URL for sharing
  const baseUrl = siteConfig.host;
  const currentUrl = `${baseUrl}/${params.slug}/stories/${params.storyslug}`;

  return (
    <div className="flex max-w-[980px] flex-col items-start">
      <div className="content">
        <h2>{storyData.title}</h2>

        <StoryMetadata story={storyData} />

        <article>{mdxSource?.content}</article>

        <StoryShareWrapper story={storyData} currentUrl={currentUrl} />

        <StoryFooter story={storyData} />
      </div>
    </div>
  );
}

export { IndexPage as default };
