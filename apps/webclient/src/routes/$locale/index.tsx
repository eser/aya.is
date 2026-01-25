// Homepage for locale - shows landing page
// On custom domain, server-side URL rewriting redirects to profile page
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { ArrowRight } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { Button } from "@/components/ui/button";
import { LocaleLink } from "@/components/locale-link";
import { Astronaut } from "@/components/widgets/astronaut";
import { MdxContent } from "@/components/userland/mdx-content";
import { compileMdx } from "@/lib/mdx";
import { siteConfig } from "@/config";
import i18next from "i18next";

export const Route = createFileRoute("/$locale/")({
  loader: async ({ params }) => {
    const { locale } = params;
    const introText = i18next.getFixedT(locale)("Home.IntroText");
    const compiledIntro = await compileMdx(introText);
    return { compiledIntro };
  },
  component: LocaleHomePage,
});

function LocaleHomePage() {
  const { t } = useTranslation();
  const { compiledIntro } = Route.useLoaderData();

  return (
    <PageLayout>
      <section className="container items-center px-4 mx-auto grid">
        <div className="flex max-w-[980px] flex-col items-start pt-10">
          {/* Astronaut - hidden on small screens, shown on xl */}
          <div className="absolute hidden xl:inline-block end-0 z-0">
            <Astronaut width={400} height={400} />
          </div>

          <article className="content relative z-10">
            <h1 className="hero">{t("Home.AYA the Open Source Network")}</h1>
            <h2 className="subtitle">
              {t(
                "Home.A software foundation formed by volunteer-developed software",
              )}
            </h2>

            <div className="mt-10" />

            {compiledIntro !== null && <MdxContent compiledSource={compiledIntro} />}
          </article>
        </div>

        <div className="flex gap-4 mt-8">
          <Button size="lg" render={<LocaleLink to="/aya/about" />}>
            {t("Home.Rest of the story")}
            <ArrowRight className="ml-2 h-4 w-4" />
          </Button>
          <Button
            variant="secondary"
            size="lg"
            render={
              <a
                target="_blank"
                rel="noreferrer"
                href={siteConfig.links.github}
              />
            }
          >
            {t("Home.GitHub")}
          </Button>
        </div>
      </section>
    </PageLayout>
  );
}
