"use client";

import { usePathname, useRouter, useSearchParams } from "next/navigation";
import * as React from "react";

import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/shared/components/ui/pagination.tsx";
import { Story } from "@/shared/components/userland/story/story.tsx";
import { formatMonthYear, parseDateFromSlug } from "@/shared/lib/date.ts";
import type { Story as StoryType, StoryEx } from "@/shared/modules/backend/stories/types.ts";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";

import styles from "../page.module.css";

const STORIES_PER_PAGE = 12;

type GroupedStories = {
  monthYear: string;
  date: Date;
  stories: (StoryType | StoryEx)[];
};

type StoriesPageClientProps = {
  initialStories: (StoryType | StoryEx)[] | null;
};

export function StoriesPageClient(props: StoriesPageClientProps) {
  const { t, locale } = useTranslations();
  const searchParams = useSearchParams();
  const router = useRouter();
  const pathname = usePathname();

  const currentPage = parseInt(searchParams.get("page") || "1", 10);

  const groupedStories = React.useMemo(() => {
    if (!props.initialStories) return [];

    const storiesWithDates = props.initialStories
      .map((story) => ({
        story,
        date: parseDateFromSlug(story.slug),
      }))
      .filter((item): item is { story: StoryType | StoryEx; date: Date } => item.date !== null);

    storiesWithDates.sort((a, b) => b.date.getTime() - a.date.getTime());

    const groups = new Map<string, GroupedStories>();

    storiesWithDates.forEach(({ story, date }) => {
      const monthYear = formatMonthYear(date, locale.code);

      if (!groups.has(monthYear)) {
        groups.set(monthYear, {
          monthYear,
          date,
          stories: [],
        });
      }

      groups.get(monthYear)!.stories.push(story);
    });

    return Array.from(groups.values()).sort((a, b) => b.date.getTime() - a.date.getTime());
  }, [props.initialStories, locale.code]);

  const totalStories = props.initialStories?.length || 0;
  const totalPages = Math.ceil(totalStories / STORIES_PER_PAGE);
  const startIndex = (currentPage - 1) * STORIES_PER_PAGE;
  const endIndex = startIndex + STORIES_PER_PAGE;

  const allStoriesFlat = groupedStories.flatMap((group) =>
    group.stories.map((story) => ({ ...story, groupMonthYear: group.monthYear, groupDate: group.date }))
  );
  const currentPageStories = allStoriesFlat.slice(startIndex, endIndex);

  const currentPageGrouped = React.useMemo(() => {
    const groups = new Map<string, GroupedStories>();

    currentPageStories.forEach((story) => {
      const monthYear = story.groupMonthYear;

      if (!groups.has(monthYear)) {
        groups.set(monthYear, {
          monthYear,
          date: story.groupDate,
          stories: [],
        });
      }

      groups.get(monthYear)!.stories.push(story);
    });

    return Array.from(groups.values()).sort((a, b) => b.date.getTime() - a.date.getTime());
  }, [currentPageStories]);

  const handlePageChange = (page: number) => {
    const params = new URLSearchParams(searchParams);
    if (page === 1) {
      params.delete("page");
    } else {
      params.set("page", page.toString());
    }

    const search = params.toString();
    const url = search ? `${pathname}?${search}` : pathname;
    router.push(url);
  };

  if (!props.initialStories || props.initialStories.length === 0) {
    return (
      <section className="container px-4 py-8 mx-auto">
        <div className="content">
          <h2>{t("Layout", "Articles")}</h2>
          {/* TODO(@eser) Add a specific translation key e.g., t("NewsPage", "NoNewsFound") */}
          <article>{t("Layout", "Content not yet available.")}</article>
        </div>
      </section>
    );
  }

  return (
    <section className="container px-4 py-8 mx-auto">
      <div className="content">
        <h2>{t("Layout", "Articles")}</h2>

        {currentPageGrouped.map((group) => (
          <div key={group.monthYear} className={styles["story-group"]}>
            <h3 className={styles["group-header"]}>{formatMonthYear(group.date, locale.code)}</h3>
            <div className="divide-y divide-border">
              {group.stories.map((story) => <Story key={story.id} story={story} />)}
            </div>
          </div>
        ))}

        {totalPages > 1 && (
          <div className={styles["pagination-container"]}>
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      if (currentPage > 1) {
                        handlePageChange(currentPage - 1);
                      }
                    }}
                    aria-disabled={currentPage <= 1}
                    className={currentPage <= 1 ? "pointer-events-none opacity-50" : ""}
                  />
                </PaginationItem>

                {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                  <PaginationItem key={page}>
                    <PaginationLink
                      href="#"
                      onClick={(e) => {
                        e.preventDefault();
                        handlePageChange(page);
                      }}
                      isActive={currentPage === page}
                    >
                      {page}
                    </PaginationLink>
                  </PaginationItem>
                ))}

                <PaginationItem>
                  <PaginationNext
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      if (currentPage < totalPages) {
                        handlePageChange(currentPage + 1);
                      }
                    }}
                    aria-disabled={currentPage >= totalPages}
                    className={currentPage >= totalPages ? "pointer-events-none opacity-50" : ""}
                  />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </div>
        )}
      </div>
    </section>
  );
}
