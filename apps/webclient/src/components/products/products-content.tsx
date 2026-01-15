"use client";

import * as React from "react";
import type { Profile } from "@/modules/backend/types";
import { FilterBar, type ProductStatusFilter } from "./filter-bar";
import { ProductListDisplay } from "./product-list-display";

export type ProductsContentProps = {
  initialProfiles: Profile[];
};

export function ProductsContent(props: ProductsContentProps) {
  const [searchText, setSearchText] = React.useState("");
  const [statusFilter, setStatusFilter] = React.useState<ProductStatusFilter>(
    "all",
  );

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
