"use client";

import * as React from "react";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { type Profile } from "@/shared/components/userland/profile-card/profile-card.tsx";
import { ProductCard } from "./product-card.tsx";
import { type ProductStatusFilter } from "./filter-bar.tsx";

export type ProductListDisplayProps = {
  allProducts: Profile[]; // Assuming products also use the Profile type structure
  searchText: string;
  statusFilter: ProductStatusFilter;
};

export function ProductListDisplay(props: ProductListDisplayProps) {
  const { t } = useTranslations();
  const { allProducts, searchText, statusFilter } = props;

  const filteredProducts = React.useMemo(() => {
    let products = allProducts;

    if (searchText.trim() !== "") {
      const lowerSearchText = searchText.toLowerCase();
      products = products.filter((product) =>
        product.title?.toLowerCase().includes(lowerSearchText) ||
        product.description?.toLowerCase().includes(lowerSearchText)
      );
    }

    //TODO(@sameterkanboz): implement filtering by tags when available instead of kind
    if (statusFilter !== "all") {
      products = products.filter((product) => {
        switch (statusFilter) {
          case "help-needed":
            return product.kind === "help-needed";
          case "looking-for-participants":
            return product.kind === "looking-for-participants";
          default:
            return true;
        }
      });
    }

    return products;
  }, [allProducts, searchText, statusFilter]);

  if (filteredProducts.length === 0) {
    return (
      <div className="py-10 text-center">
        <p className="text-xl text-muted-foreground">{t("Products", "NoProductsFound")}</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {filteredProducts.map((product) => <ProductCard key={product.slug} product={product} />)}
    </div>
  );
}
