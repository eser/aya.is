import * as React from "react";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { ProfileCard } from "@/components/userland/profile-card/profile-card";
import { Button } from "@/components/ui/button";
import { profilesByKindsQueryOptions } from "@/modules/backend/queries";
import { backend } from "@/modules/backend/backend";
import { ProfileListingFilterBar } from "./filter-bar";
import { ProfileCardSkeletonGrid } from "./profile-card-skeleton";
import styles from "./profile-listing.module.css";

const PAGE_SIZE = 24;
const DEBOUNCE_MS = 300;

export type FilterOption = {
  label: string;
  value: string;
};

export type ProfileListingContentProps = {
  locale: string;
  seed: string;
  /** Base kinds for this listing (e.g., ["product"] or ["individual", "organization"]) */
  baseKinds: string[];
  /** Toggle filter options. First option's value should be "" for "all". Others narrow filter_kind. */
  filterOptions?: FilterOption[];
  /** Label for the filter toggle group */
  filterLabel?: string;
  /** Search input placeholder text */
  searchPlaceholder: string;
  /** Empty state message */
  emptyMessage: string;
};

function resolveKinds(baseKinds: string[], activeFilter: string): string[] {
  if (activeFilter === "" || !baseKinds.includes(activeFilter)) {
    return baseKinds;
  }

  return [activeFilter];
}

export function ProfileListingContent(props: ProfileListingContentProps) {
  const { t } = useTranslation();

  const [searchText, setSearchText] = React.useState("");
  const [debouncedSearch, setDebouncedSearch] = React.useState("");
  const [activeFilter, setActiveFilter] = React.useState("");
  const [additionalPages, setAdditionalPages] = React.useState<
    Array<Array<{ slug: string; [key: string]: unknown }>>
  >([]);
  const [hasMore, setHasMore] = React.useState(true);
  const [isLoadingMore, setIsLoadingMore] = React.useState(false);

  // Debounce search input
  React.useEffect(() => {
    const timer = setTimeout(
      () => setDebouncedSearch(searchText),
      DEBOUNCE_MS,
    );

    return () => clearTimeout(timer);
  }, [searchText]);

  // Reset pagination when search or filter changes
  React.useEffect(() => {
    setAdditionalPages([]);
    setHasMore(true);
  }, [debouncedSearch, activeFilter]);

  const kinds = resolveKinds(props.baseKinds, activeFilter);

  const { data: firstPageData, isLoading, isFetching } = useQuery(
    profilesByKindsQueryOptions(props.locale, kinds, {
      seed: props.seed,
      q: debouncedSearch === "" ? undefined : debouncedSearch,
    }),
  );

  const firstPage = firstPageData ?? [];

  // Update hasMore when first page loads
  React.useEffect(() => {
    if (firstPage.length < PAGE_SIZE) {
      setHasMore(false);
    } else {
      setHasMore(true);
    }
  }, [firstPage]);

  const allProfiles = React.useMemo(() => {
    const profiles = [...firstPage];

    for (const page of additionalPages) {
      for (const profile of page) {
        profiles.push(profile as typeof profiles[number]);
      }
    }

    return profiles;
  }, [firstPage, additionalPages]);

  const loadMore = React.useCallback(async () => {
    setIsLoadingMore(true);

    try {
      const offset = PAGE_SIZE + additionalPages.length * PAGE_SIZE;
      const nextPage = await backend.getProfilesByKinds(
        props.locale,
        kinds,
        props.seed,
        PAGE_SIZE,
        offset,
        debouncedSearch === "" ? undefined : debouncedSearch,
      );

      if (nextPage !== null && nextPage.length > 0) {
        setAdditionalPages((prev) => [...prev, nextPage]);
        setHasMore(nextPage.length === PAGE_SIZE);
      } else {
        setHasMore(false);
      }
    } finally {
      setIsLoadingMore(false);
    }
  }, [props.locale, props.seed, kinds, debouncedSearch, additionalPages.length]);

  const handleFilterChange = React.useCallback((value: string) => {
    setActiveFilter(value);
  }, []);

  const showSkeleton = isLoading && allProfiles.length === 0;
  const showEmpty = !isLoading && allProfiles.length === 0;

  return (
    <>
      <ProfileListingFilterBar
        filterOptions={props.filterOptions}
        activeFilter={activeFilter}
        onFilterChange={handleFilterChange}
        filterLabel={props.filterLabel ?? t("Common.Filter")}
        searchLabel={t("Search.Search")}
        searchText={searchText}
        onSearchTextChange={setSearchText}
        searchPlaceholder={props.searchPlaceholder}
      />

      {showSkeleton && <ProfileCardSkeletonGrid />}

      {showEmpty && (
        <div className={styles.emptyState}>
          <p className={styles.emptyMessage}>{props.emptyMessage}</p>
        </div>
      )}

      {allProfiles.length > 0 && (
        <div className={`${styles.grid} ${isFetching ? "opacity-60 transition-opacity" : ""}`}>
          {allProfiles.map((profile) => (
            <ProfileCard
              key={profile.slug}
              profile={profile}
              variant="cover"
              showKindBadge
            />
          ))}
        </div>
      )}

      {hasMore && allProfiles.length > 0 && (
        <div className={styles.loadMoreContainer}>
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
