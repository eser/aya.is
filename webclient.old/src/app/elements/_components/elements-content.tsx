"use client";

import * as React from "react";

import type { Profile } from "@/shared/components/userland/profile-card/profile-card.tsx";

import { FilterBar, type ProfileKindFilter } from "./filter-bar.tsx";
import { ProfileListDisplay } from "./profile-list-display.tsx";

export type ElementsContentProps = {
  initialProfiles: Profile[];
};

export function ElementsContent(props: ElementsContentProps) {
  const [activeKindFilter, setActiveKindFilter] = React.useState<ProfileKindFilter>("");
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
