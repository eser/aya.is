"use client";

import * as React from "react";
import { useTranslation } from "react-i18next";
import type { Profile } from "@/modules/backend/types";
import { backend } from "@/modules/backend/backend";
import { Button } from "@/components/ui/button";
import { FilterBar, type ProductStatusFilter } from "./-filter-bar";
import { ProductListDisplay } from "./-product-list-display";

const PAGE_SIZE = 24;

export type ProductsContentProps = {
  initialProfiles: Profile[];
  locale: string;
  seed: string;
  kinds: string[];
};

export function ProductsContent(props: ProductsContentProps) {
  const { t } = useTranslation();
  const [searchText, setSearchText] = React.useState("");
  const [statusFilter, setStatusFilter] = React.useState<ProductStatusFilter>(
    "all",
  );
  const [allProfiles, setAllProfiles] = React.useState<Profile[]>(
    props.initialProfiles ?? [],
  );
  const [offset, setOffset] = React.useState(PAGE_SIZE);
  const [hasMore, setHasMore] = React.useState(
    (props.initialProfiles?.length ?? 0) === PAGE_SIZE,
  );
  const [isLoadingMore, setIsLoadingMore] = React.useState(false);

  React.useEffect(() => {
    setAllProfiles(props.initialProfiles ?? []);
    setOffset(PAGE_SIZE);
    setHasMore((props.initialProfiles?.length ?? 0) === PAGE_SIZE);
  }, [props.initialProfiles]);

  const loadMore = React.useCallback(async () => {
    setIsLoadingMore(true);

    try {
      const nextPage = await backend.getProfilesByKinds(
        props.locale,
        props.kinds,
        props.seed,
        PAGE_SIZE,
        offset,
      );

      if (nextPage !== null && nextPage.length > 0) {
        setAllProfiles((prev) => [...prev, ...nextPage]);
        setOffset((prev) => prev + PAGE_SIZE);
        setHasMore(nextPage.length === PAGE_SIZE);
      } else {
        setHasMore(false);
      }
    } finally {
      setIsLoadingMore(false);
    }
  }, [props.locale, props.kinds, props.seed, offset]);

  return (
    <>
      <FilterBar
        searchText={searchText}
        onSearchTextChange={setSearchText}
        statusFilter={statusFilter}
        onStatusFilterChange={setStatusFilter}
      />
      <ProductListDisplay
        allProducts={allProfiles}
        searchText={searchText}
        statusFilter={statusFilter}
      />
      {hasMore && (
        <div className="flex justify-center mt-8">
          <Button
            variant="outline"
            size="lg"
            onClick={loadMore}
            disabled={isLoadingMore}
          >
            {isLoadingMore ? t("Common.Loading...") : t("Common.Load more")}
          </Button>
        </div>
      )}
    </>
  );
}
