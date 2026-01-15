"use client";

import * as React from "react";
import { useTranslation } from "react-i18next";
import { ProfileCard } from "@/components/userland/profile-card";
import type { Profile } from "@/modules/backend/types";
import type { ProfileKindFilter } from "./filter-bar";

export type ProfileListDisplayProps = {
  profiles: Profile[];
  activeKindFilter: ProfileKindFilter;
  searchText: string;
};

export function ProfileListDisplay(props: ProfileListDisplayProps) {
  const { t } = useTranslation();
  const { profiles, activeKindFilter, searchText } = props;

  const filteredProfiles = React.useMemo(() => {
    let profilesFiltered: Profile[] = profiles;

    if (activeKindFilter !== "") {
      profilesFiltered = profiles.filter(
        (profile) => profile.kind === activeKindFilter
      );
    }

    if (searchText.trim() !== "") {
      const lowerSearchText = searchText.toLowerCase();
      profilesFiltered = profilesFiltered.filter(
        (profile) =>
          profile.title?.toLowerCase().includes(lowerSearchText) ||
          profile.description?.toLowerCase().includes(lowerSearchText)
      );
    }

    return profilesFiltered;
  }, [profiles, activeKindFilter, searchText]);

  if (filteredProfiles.length === 0) {
    return (
      <div className="py-10 text-xl text-center text-muted-foreground">
        {t("Elements.NoProfilesFound")}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {filteredProfiles.map((profile) => (
        <ProfileCard key={profile.slug} profile={profile} />
      ))}
    </div>
  );
}
