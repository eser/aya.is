"use client";

import * as React from "react";
// import NextImage from "next/image";

// import { formatDateString } from "@/shared/lib/date.ts";
import type { Story, StoryEx } from "@/shared/modules/backend/stories/types.ts";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { Story as StoryComponent } from "@/shared/components/userland/story/story.tsx";
// import { SiteLink } from "@/shared/components/userland/site-link/site-link.tsx";

import { FilterBar, type FilterKeyType } from "./filter-bar.tsx"; // Import the type
// import styles from "./timeline.module.css";

// const Image = NextImage; // Simple reassignment, will add @ts-expect-error at usage

export type TimelineProps = {
  profileSlug?: string;
  stories: (Story | StoryEx)[];
};

export function Timeline(props: TimelineProps) {
  const { t } = useTranslations();

  const [activeFilter, setActiveFilter] = React.useState<FilterKeyType>("");

  const filteredStories = props.stories.filter((story) => activeFilter === "" || story.kind === activeFilter);

  const filterMapping: Record<FilterKeyType, string> = {
    "": t("Stories", "All"),
    "status": t("Stories", "Status"),
    "announcement": t("Stories", "Announcement"),
    "article": t("Stories", "Article"),
    "content": t("Stories", "Content"),
    "presentation": t("Stories", "Presentation"),
  };

  return (
    <>
      <FilterBar activeFilter={activeFilter} setActiveFilter={setActiveFilter} mapping={filterMapping} />

      <div className="py-10">
        {filteredStories !== null && filteredStories.length > 0
          ? (
            <div className="divide-y divide-border">
              {filteredStories.map((story) => (
                <StoryComponent
                  key={story.id}
                  profileSlug={props.profileSlug}
                  story={story}
                />
              ))}
            </div>
          )
          : (
            <article>
              {/* TODO(@eser) Add a specific translation key e.g., t("NewsPage", "NoNewsFound") */}
              {t("Layout", "Content not yet available.")}
            </article>
          )}
      </div>
    </>
  );
}
