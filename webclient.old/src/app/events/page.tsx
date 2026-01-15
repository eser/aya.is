import * as React from "react";
import type { Metadata, ResolvingMetadata } from "next";

import { mdx } from "@/shared/lib/mdx.tsx";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { components } from "@/shared/components/userland/userland.ts";

type IndexPageProps = {
  params: Promise<never>;
};

// TODO(@eser) add more from https://beta.nextjs.org/docs/api-reference/metadata
export async function generateMetadata(_props: IndexPageProps, _parent: ResolvingMetadata): Promise<Metadata> {
  const { t } = await getTranslations();

  return {
    title: t("Layout", "Events"),
    description: "",
  };
}

async function IndexPage(_props: IndexPageProps) {
  const { t } = await getTranslations();

  const contentText = `${t("Layout", "Content not yet available.")}`;

  const mdxSource = await mdx(contentText, components);

  return (
    <section className="container px-4 py-8 mx-auto">
      <div className="content">
        <h2>{t("Layout", "Events")}</h2>

        <article>{mdxSource?.content}</article>
      </div>
    </section>
  );
}

export { IndexPage as default };
