"use client";

import * as React from "react";
import { useTranslation } from "react-i18next";
import type { Profile } from "@/modules/backend/types";
import { ProfileCard } from "@/components/userland/profile-card/profile-card";
import type { ProductStatusFilter } from "./-filter-bar";

export type ProductListDisplayProps = {
  allProducts: Profile[];
  searchText: string;
  statusFilter: ProductStatusFilter;
};

export function ProductListDisplay(props: ProductListDisplayProps) {
  const { t } = useTranslation();
  const { allProducts, searchText, statusFilter } = props;

  const filteredProducts = React.useMemo(() => {
    let products = allProducts;

    if (searchText.trim() !== "") {
      const lowerSearchText = searchText.toLowerCase();
      products = products.filter(
        (product) =>
          product.title?.toLowerCase().includes(lowerSearchText) ||
          product.description?.toLowerCase().includes(lowerSearchText),
      );
    }

    // Filter by status - when API supports tags, this can be updated
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
        <p className="text-xl text-muted-foreground">
          {t("Products.NoProductsFound")}
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {filteredProducts.map((product) => (
        <ProfileCard key={product.slug} profile={product} variant="cover" showKindBadge />
      ))}
    </div>
  );
}
