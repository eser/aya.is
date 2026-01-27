"use client";

import * as React from "react";
import { useNavigate, useSearch } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Story } from "@/components/userland/story";
import { formatMonthYear, parseDateFromSlug } from "@/lib/date";
import type { StoryEx } from "@/modules/backend/types";

const ITEMS_PER_PAGE = 12;

type GroupedStories = {
  monthYear: string;
  date: Date;
  stories: StoryEx[];
};

type StoriesPageClientProps = {
  initialStories: StoryEx[] | null;
  basePath: string;
  /** Profile slug for generating story links (e.g., "eser" for /eser/stories/...) */
  profileSlug?: string;
};

export function StoriesPageClient(props: StoriesPageClientProps) {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;
  const navigate = useNavigate();

  // Get offset from URL search params
  const searchParams = useSearch({ strict: false }) as { offset?: number };
  const currentOffset = searchParams.offset ?? 0;

  const groupedStories = React.useMemo(() => {
    if (props.initialStories === null) return [];

    const storiesWithDates = props.initialStories
      .map((story) => ({
        story,
        date: parseDateFromSlug(story.slug),
      }))
      .filter(
        (item): item is { story: StoryEx; date: Date } => item.date !== null,
      );

    storiesWithDates.sort((a, b) => b.date.getTime() - a.date.getTime());

    const groups = new Map<string, GroupedStories>();

    storiesWithDates.forEach(({ story, date }) => {
      const monthYear = formatMonthYear(date, locale);

      if (!groups.has(monthYear)) {
        groups.set(monthYear, {
          monthYear,
          date,
          stories: [],
        });
      }

      groups.get(monthYear)!.stories.push(story);
    });

    return Array.from(groups.values()).sort(
      (a, b) => b.date.getTime() - a.date.getTime(),
    );
  }, [props.initialStories, locale]);

  const totalStories = props.initialStories !== null ? props.initialStories.length : 0;
  const totalPages = Math.ceil(totalStories / ITEMS_PER_PAGE);
  const currentPage = Math.floor(currentOffset / ITEMS_PER_PAGE) + 1;
  const endIndex = currentOffset + ITEMS_PER_PAGE;

  const allStoriesFlat = groupedStories.flatMap((group) =>
    group.stories.map((story) => ({
      ...story,
      groupMonthYear: group.monthYear,
      groupDate: group.date,
    }))
  );
  const currentPageStories = allStoriesFlat.slice(currentOffset, endIndex);

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

    return Array.from(groups.values()).sort(
      (a, b) => b.date.getTime() - a.date.getTime(),
    );
  }, [currentPageStories]);

  const handleOffsetChange = (newOffset: number) => {
    const search = newOffset === 0 ? {} : { offset: newOffset };
    navigate({ to: props.basePath, search });
  };

  if (props.initialStories === null || props.initialStories.length === 0) {
    return (
      <p className="text-muted-foreground">
        {t("Layout.Content not yet available.")}
      </p>
    );
  }

  return (
    <>
      {currentPageGrouped.map((group) => (
        <div key={group.monthYear} className="mb-8">
          <h2 className="text-lg font-semibold text-muted-foreground mb-4 pb-2 border-b border-border">
            {formatMonthYear(group.date, locale)}
          </h2>
          <div>
            {group.stories.map((story) => (
              <Story
                key={story.id}
                story={story}
                profileSlug={props.profileSlug}
              />
            ))}
          </div>
        </div>
      ))}

      {totalPages > 1 && (
        <div className="mt-8 flex justify-center">
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    if (currentOffset > 0) {
                      handleOffsetChange(
                        Math.max(0, currentOffset - ITEMS_PER_PAGE),
                      );
                    }
                  }}
                  aria-disabled={currentOffset <= 0}
                  className={currentOffset <= 0 ? "pointer-events-none opacity-50" : ""}
                >
                  {t("Common.Previous")}
                </PaginationPrevious>
              </PaginationItem>

              {Array.from({ length: totalPages }, (_, i) => i + 1).map(
                (page) => {
                  const pageOffset = (page - 1) * ITEMS_PER_PAGE;
                  return (
                    <PaginationItem key={page}>
                      <PaginationLink
                        href="#"
                        onClick={(e) => {
                          e.preventDefault();
                          handleOffsetChange(pageOffset);
                        }}
                        isActive={currentPage === page}
                      >
                        {page}
                      </PaginationLink>
                    </PaginationItem>
                  );
                },
              )}

              <PaginationItem>
                <PaginationNext
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    const nextOffset = currentOffset + ITEMS_PER_PAGE;
                    if (nextOffset < totalStories) {
                      handleOffsetChange(nextOffset);
                    }
                  }}
                  aria-disabled={currentOffset + ITEMS_PER_PAGE >= totalStories}
                  className={currentOffset + ITEMS_PER_PAGE >= totalStories ? "pointer-events-none opacity-50" : ""}
                >
                  {t("Common.Next")}
                </PaginationNext>
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </div>
      )}
    </>
  );
}
