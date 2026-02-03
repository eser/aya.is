// Elements index page
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth/auth-context";
import { buildUrl, generateMetaTags } from "@/lib/seo";
import { ElementsContent } from "./_components/-elements-content";
import i18next from "i18next";

export const Route = createFileRoute("/$locale/elements/")({
  loader: async ({ params }) => {
    const { locale } = params;
    const profiles = await backend.getProfilesByKinds(locale, [
      "individual",
      "organization",
    ]);

    // Ensure locale translations are loaded before translating
    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    return {
      profiles: profiles ?? [],
      locale,
      translatedTitle: t("Layout.Elements"),
      translatedDescription: t("Elements.Discover individuals and organizations in the AYA community"),
    };
  },
  head: ({ loaderData }) => {
    const { locale, translatedTitle, translatedDescription } = loaderData;
    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, "elements"),
        locale,
        type: "website",
      }),
    };
  },
  component: ElementsIndexPage,
});

function ElementsIndexPage() {
  const { profiles, locale } = Route.useLoaderData();
  const { t } = useTranslation();
  const { isAuthenticated } = useAuth();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <div className="flex items-center justify-between mb-4">
            <h1 className="no-margin">{t("Layout.Elements")}</h1>
            {isAuthenticated && (
              <Link
                to="/$locale/elements/new"
                params={{ locale }}
              >
                <Button variant="default" size="sm">
                  <Plus className="mr-1.5 size-4" />
                  {t("Elements.Add Element")}
                </Button>
              </Link>
            )}
          </div>

          <ElementsContent initialProfiles={profiles} />
        </div>
      </section>
    </PageLayout>
  );
}
