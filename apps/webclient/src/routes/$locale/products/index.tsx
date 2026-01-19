// Products page
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { ProductsContent } from "./_components/-products-content";

export const Route = createFileRoute("/$locale/products/")({
  loader: async ({ params }) => {
    const { locale } = params;
    const products = await backend.getProfilesByKinds(locale, ["product"]);
    return { products: products ?? [] };
  },
  component: ProductsPage,
});

function ProductsPage() {
  const { products } = Route.useLoaderData();
  const { t } = useTranslation();

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <h1>{t("Layout.Products")}</h1>

          <ProductsContent initialProfiles={products} />
        </div>
      </section>
    </PageLayout>
  );
}
