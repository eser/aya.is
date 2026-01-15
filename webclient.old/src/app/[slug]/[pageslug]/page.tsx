import * as React from "react";
import type { Metadata, ResolvingMetadata } from "next";
import { notFound } from "next/navigation";

import { mdx } from "@/shared/lib/mdx.tsx";
import { backend } from "@/shared/modules/backend/backend.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { components } from "@/shared/components/userland/userland.ts";

type IndexPageProps = {
  params: Promise<{
    slug: string;
    pageslug: string;
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

  const pageData = await backend.getProfilePage(locale.code, profileData.slug, params.pageslug);

  if (pageData === null) {
    notFound();
  }

  return {
    title: `${profileData.title} - ${pageData.title}`,
    description: pageData.summary,
  };
}

async function IndexPage(props: IndexPageProps) {
  const params = await props.params;

  const { locale } = await getTranslations();

  const pageData = await backend.getProfilePage(locale.code, params.slug, params.pageslug);

  if (pageData === null) {
    notFound();
  }

  const contentText = `# ${pageData.title}\n\n${pageData.content}`;

  const mdxSource = await mdx(contentText, components);

  return (
    <div className="flex max-w-[980px] flex-col items-start">
      <article className="content">{mdxSource?.content}</article>
    </div>
  );
}

export { IndexPage as default };
