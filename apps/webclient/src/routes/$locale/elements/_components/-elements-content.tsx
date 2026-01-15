"use client";

import * as React from "react";
import type { Profile } from "@/modules/backend/types";
import { FilterBar, type ProfileKindFilter } from "./-filter-bar";
import { ProfileListDisplay } from "./-profile-list-display";

export type ElementsContentProps = {
  initialProfiles: Profile[];
};

export function ElementsContent(props: ElementsContentProps) {
  const [activeKindFilter, setActiveKindFilter] = React.useState<
    ProfileKindFilter
  >("");
  const [searchText, setSearchText] = React.useState("");

  return (
    <>
      <FilterBar
        activeKindFilter={activeKindFilter}
        onKindChange={setActiveKindFilter}
        searchText={searchText}
        onSearchTextChange={setSearchText}
      />
      <ProfileListDisplay
        profiles={props.initialProfiles}
        activeKindFilter={activeKindFilter}
        searchText={searchText}
      />
    </>
  );
}
