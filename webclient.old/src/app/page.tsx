import * as React from "react";

import { siteConfig } from "@/shared/config.ts";
import { mdx } from "@/shared/lib/mdx.tsx";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { Button } from "@/shared/components/ui/button.tsx";
import { Astronaut } from "@/shared/components/widgets/astronaut/astronaut.tsx";
import { SiteLink } from "@/shared/components/userland/site-link/site-link.tsx";
import { components } from "@/shared/components/userland/userland.ts";
import { Icons } from "@/shared/components/icons.tsx";

import styles from "./page.module.css";

async function IndexPage() {
  const { t } = await getTranslations();

  const contentText = t("Home", "IntroText");

  const mdxSource = await mdx(contentText, components);

  return (
    <section className="container items-center px-4 mx-auto grid">
      <div className="flex max-w-[980px] flex-col items-start pt-10">
        <div className={styles["astronaut-layer"]}>
          <Astronaut width={400} height={400} />
        </div>
        <article className="content">
          <h1>{t("Home", "AYA the Open Source Network")}</h1>
          <h3>{t("Home", "A software foundation formed by volunteer-developed software")}</h3>

          <div className="mt-10" />

          {mdxSource?.content}
        </article>
      </div>
      <div className="flex gap-4">
        <Button variant="default" size="lg" asChild>
          <SiteLink href="/aya/about" rel="noreferrer">
            {t("Home", "Rest of the story")}
            <Icons.arrowRight />
          </SiteLink>
        </Button>
        <Button variant="secondary" size="lg" asChild>
          <SiteLink target="_blank" rel="noreferrer" href={siteConfig.links.github}>
            {t("Home", "GitHub")}
          </SiteLink>
        </Button>
      </div>
    </section>
  );
}

export { IndexPage as default };
