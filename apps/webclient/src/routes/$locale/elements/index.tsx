// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Elements index page
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth/auth-context";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import { getDailySeed } from "@/lib/seed-utils";
import { profilesByKindsQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import { ProfileListingContent } from "@/components/userland/profile-listing";
import i18next from "i18next";

export const Route = createFileRoute("/$locale/elements/")({
  loader: async ({ params, context }) => {
    const { locale } = params;
    const seed = getDailySeed();

    await Promise.all([
      context.queryClient.ensureQueryData(
        profilesByKindsQueryOptions(locale, ["individual", "organization"], { seed }),
      ),
      i18next.loadLanguages(locale),
    ]);

    const t = i18next.getFixedT(locale);
    return {
      locale,
      seed,
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
      links: [generateCanonicalLink(buildUrl(locale, "elements"))],
    };
  },
  errorComponent: QueryError,
  component: ElementsIndexPage,
});

function ElementsIndexPage() {
  const { locale, seed } = Route.useLoaderData();
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

          <ProfileListingContent
            locale={locale}
            seed={seed}
            baseKinds={["individual", "organization"]}
            filterOptions={[
              { label: t("Elements.AllTypes"), value: "" },
              { label: t("Elements.Individuals"), value: "individual" },
              { label: t("Elements.Organizations"), value: "organization" },
            ]}
            filterLabel={t("Elements.FilterByKind")}
            searchPlaceholder={t("Elements.SearchPlaceholder")}
            emptyMessage={t("Elements.NoProfilesFound")}
          />
        </div>
      </section>
    </PageLayout>
  );
}
