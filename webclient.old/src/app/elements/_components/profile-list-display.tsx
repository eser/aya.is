"use client";

import * as React from "react";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { type Profile, ProfileCard } from "@/shared/components/userland/profile-card/profile-card.tsx";
import type { ProfileKindFilter } from "./filter-bar.tsx";

export type ProfileListDisplayProps = {
  profiles: Profile[];
  activeKindFilter: ProfileKindFilter;
  searchText: string;
};

export function ProfileListDisplay(props: ProfileListDisplayProps) {
  const { t } = useTranslations();
  const { profiles, activeKindFilter, searchText } = props;

  const filteredProfiles = React.useMemo(() => {
    let profilesFiltered: Profile[] = profiles;

    if (activeKindFilter !== "") {
      profilesFiltered = profiles.filter((profile) => profile.kind === activeKindFilter);
    }

    if (searchText.trim() !== "") {
      const lowerSearchText = searchText.toLowerCase();
      profilesFiltered = profilesFiltered.filter(
        (profile) =>
          profile.title?.toLowerCase().includes(lowerSearchText) ||
          profile.description?.toLowerCase().includes(lowerSearchText),
      );
    }

    return profilesFiltered;
  }, [profiles, activeKindFilter, searchText]);

  if (filteredProfiles.length === 0) {
    return (
      <div className="py-10 text-xl text-center text-muted-foreground">
        {t("Elements", "NoProfilesFound")}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {filteredProfiles.map((profile) => <ProfileCard key={profile.slug} profile={profile} />)}
    </div>
  );
}
