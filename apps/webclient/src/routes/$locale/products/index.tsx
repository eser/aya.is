// Products page
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth/auth-context";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import { ProductsContent } from "./_components/-products-content";
import i18next from "i18next";

export const Route = createFileRoute("/$locale/products/")({
  loader: async ({ params }) => {
    const { locale } = params;
    const products = await backend.getProfilesByKinds(locale, ["product"]);

    // Ensure locale translations are loaded before translating
    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    return {
      products: products ?? [],
      locale,
      translatedTitle: t("Layout.Products"),
      translatedDescription: t("Products.Discover open source products and projects"),
    };
  },
  head: ({ loaderData }) => {
    const { locale, translatedTitle, translatedDescription } = loaderData;
    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, "products"),
        locale,
        type: "website",
      }),
      links: [generateCanonicalLink(buildUrl(locale, "products"))],
    };
  },
  component: ProductsPage,
});

function ProductsPage() {
  const { products, locale } = Route.useLoaderData();
  const { t } = useTranslation();
  const { isAuthenticated } = useAuth();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <div className="flex items-center justify-between mb-4">
            <h1 className="no-margin">{t("Layout.Products")}</h1>
            {isAuthenticated && (
              <Link
                to="/$locale/products/new"
                params={{ locale }}
              >
                <Button variant="default" size="sm">
                  <Plus className="mr-1.5 size-4" />
                  {t("Products.Add Product")}
                </Button>
              </Link>
            )}
          </div>

          <ProductsContent initialProfiles={products} />
        </div>
      </section>
    </PageLayout>
  );
}
