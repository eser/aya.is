

"use client";

import * as React from "react";
import type { Profile } from "@/shared/components/userland/profile-card/profile-card.tsx";
import { FilterBar } from "./filter-bar.tsx";
import { ProductListDisplay } from "./product-list-display.tsx";
export type ProductsContentProps = {
  initialProfiles: Profile[];
};

export function ProductsContent(props: ProductsContentProps) {
  const [searchText, setSearchText] = React.useState("");
  const [statusFilter, setStatusFilter] = React.useState<ProductStatusFilter>("all");

  return (
    <>
      <FilterBar
        searchText={searchText}
        onSearchTextChange={setSearchText}
        statusFilter={statusFilter}
        onStatusFilterChange={setStatusFilter}
      />
      <ProductListDisplay
        allProducts={props.initialProfiles}
        searchText={searchText}
        statusFilter={statusFilter}
      />
    </>
  );
}
