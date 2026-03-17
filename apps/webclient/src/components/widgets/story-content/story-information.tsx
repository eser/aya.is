// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";
import { Calendar, Clock, ExternalLink, PencilLine, Share2, Tag } from "lucide-react";
import { cn } from "@/lib/utils";
import { calculateReadingTime } from "@/lib/reading-time";
import { formatDateTimeLong, formatDateTimeRange } from "@/lib/date";
import type { ActivityProperties, StoryEx } from "@/modules/backend/types";
import styles from "./story-information.module.css";

export type StoryInformationProps = {
  story: StoryEx;
  locale: string;
  editUrl?: string;
  coverUrl?: string;
  shareUrl?: string;
};

const activityKindLabels: Record<string, string> = {
  meetup: "Activities.Meetup",
  workshop: "Activities.Workshop",
  conference: "Activities.Conference",
  broadcast: "Activities.Broadcast",
  meeting: "Activities.Meeting",
};

export function StoryInformation(props: StoryInformationProps) {
  if (props.story.kind === "activity") {
    return (
      <ActivityInformationContent
        story={props.story}
        locale={props.locale}
        editUrl={props.editUrl}
        shareUrl={props.shareUrl}
      />
    );
  }

  return <DefaultInformationContent story={props.story} editUrl={props.editUrl} shareUrl={props.shareUrl} />;
}

// --- Internal: activity metadata ---

function ActivityInformationContent(props: { story: StoryEx; locale: string; editUrl?: string; shareUrl?: string }) {
  const { t } = useTranslation();

  const activityProps = (props.story.properties ?? {}) as unknown as ActivityProperties;
  const dateMode = activityProps.date_mode ?? "fixed";
  const timeStart = activityProps.activity_time_start !== undefined
    ? new Date(activityProps.activity_time_start)
    : null;
  const timeEnd = activityProps.activity_time_end !== undefined ? new Date(activityProps.activity_time_end) : null;
  const kindLabel = activityKindLabels[activityProps.activity_kind ?? "meetup"] ?? "Activities.Meetup";

  return (
    <div className={cn(styles.container, "not-prose")}>
      {dateMode === "undecided"
        ? (
          <span className={styles.item}>
            <Calendar className="size-4" />
            <span className={styles.dateUndecided}>
              {t("Activities.Date Undecided")}
            </span>
          </span>
        )
        : timeStart !== null && (
          <span className={styles.item}>
            <Calendar className="size-4" />
            {timeEnd !== null
              ? formatDateTimeRange(timeStart, timeEnd, props.locale)
              : formatDateTimeLong(timeStart, props.locale)}
          </span>
        )}

      <span className={styles.item}>
        <Tag className="size-4" />
        {t(kindLabel)}
      </span>

      {activityProps.external_activity_uri !== undefined &&
        activityProps.external_activity_uri !== "" && (
        <a
          href={activityProps.external_activity_uri}
          target="_blank"
          rel="noopener noreferrer"
          className={cn(styles.item, "text-primary hover:underline")}
        >
          <ExternalLink className="size-4" />
          {t("Common.View")}
        </a>
      )}

      {(props.editUrl !== undefined || props.shareUrl !== undefined) && (
        <div className={styles.editGroup}>
          {props.shareUrl !== undefined && (
            <Link
              to={props.shareUrl}
              className={cn(styles.editLink, "!no-underline hover:!no-underline hover:text-foreground")}
            >
              <Share2 className="size-3.5" />
              {t("ShareWizard.Share Wizard")}
            </Link>
          )}
          {props.editUrl !== undefined && (
            <Link
              to={props.editUrl}
              className={cn(styles.editLink, "!no-underline hover:!no-underline hover:text-foreground")}
            >
              <PencilLine className="size-3.5" />
              {t("ContentEditor.Edit Story")}
            </Link>
          )}
        </div>
      )}
    </div>
  );
}

// --- Internal: default (article/news/etc.) metadata ---

function DefaultInformationContent(props: { story: StoryEx; editUrl?: string; shareUrl?: string }) {
  const { t } = useTranslation();

  const readingTime = calculateReadingTime(props.story.content);
  const publishedDate = new Date(props.story.published_at ?? props.story.created_at);

  return (
    <div className={cn(styles.container, "not-prose")}>
      <span className={styles.item}>
        <Calendar className="size-4" />
        <time
          dateTime={props.story.published_at ?? props.story.created_at}
          className="text-sm text-foreground"
        >
          {t("Common.DateLong", { date: publishedDate })}
        </time>
      </span>

      <span className={styles.item}>
        <Clock className="size-4" />
        <span className="text-sm text-foreground">
          {readingTime} {t("Stories.min read")}
        </span>
      </span>

      {(props.editUrl !== undefined || props.shareUrl !== undefined) && (
        <div className={styles.editGroup}>
          {props.shareUrl !== undefined && (
            <Link
              to={props.shareUrl}
              className={cn(styles.editLink, "!no-underline hover:!no-underline hover:text-foreground")}
            >
              <Share2 className="size-3.5" />
              {t("ShareWizard.Share Wizard")}
            </Link>
          )}
          {props.editUrl !== undefined && (
            <Link
              to={props.editUrl}
              className={cn(styles.editLink, "!no-underline hover:!no-underline hover:text-foreground")}
            >
              <PencilLine className="size-3.5" />
              {t("ContentEditor.Edit Story")}
            </Link>
          )}
        </div>
      )}
    </div>
  );
}
