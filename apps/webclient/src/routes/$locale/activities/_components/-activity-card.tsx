"use client";

import { useTranslation } from "react-i18next";
import { Calendar, Clock, MapPin, User } from "lucide-react";
import { LocaleLink } from "@/components/locale-link";
import { cn } from "@/lib/utils";
import { formatDateTimeShort } from "@/lib/date";
import { stripMarkdown } from "@/lib/strip-markdown";
import { InlineMarkdown } from "@/lib/inline-markdown";
import type { StoryEx, ActivityProperties } from "@/modules/backend/types";
import styles from "./-activity-card.module.css";

const activityKindLabels: Record<string, string> = {
  meetup: "Activities.Meetup",
  workshop: "Activities.Workshop",
  conference: "Activities.Conference",
  broadcast: "Activities.Broadcast",
  meeting: "Activities.Meeting",
};

export type ActivityCardProps = {
  activity: StoryEx;
};

export function ActivityCard(props: ActivityCardProps) {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;

  const activityProps = (props.activity.properties ?? {}) as unknown as ActivityProperties;
  const activityKind = activityProps.activity_kind ?? "meetup";
  const kindLabelKey = activityKindLabels[activityKind] ?? "Activities.Meetup";

  const timeStart = activityProps.activity_time_start !== undefined
    ? new Date(activityProps.activity_time_start)
    : null;
  const timeEnd = activityProps.activity_time_end !== undefined
    ? new Date(activityProps.activity_time_end)
    : null;

  const href = `/activities/${props.activity.slug}`;

  return (
    <LocaleLink role="card" to={href} className="no-underline block">
      <article className={styles.card}>
        <div className={cn(styles.imageContainer, "w-[200px] h-[120px]")}>
          {props.activity.story_picture_uri !== null &&
              props.activity.story_picture_uri !== undefined
            ? (
              <img
                src={props.activity.story_picture_uri}
                alt={stripMarkdown(props.activity.title ?? t("Activities.Activity"))}
                width={200}
                height={120}
                className={styles.image}
              />
            )
            : (
              <div className={styles.imagePlaceholder}>
                <Calendar className="size-8 text-muted-foreground/50" />
              </div>
            )}
        </div>
        <div className={styles.contentArea}>
          <h3 className={styles.title}>
            <span className={styles.badge}>{t(kindLabelKey)}</span>
            {stripMarkdown(props.activity.title ?? "")}
          </h3>
          {props.activity.summary !== null && props.activity.summary !== undefined && (
            <InlineMarkdown content={props.activity.summary} className={styles.summary} />
          )}
          <div className={styles.meta}>
            {timeStart !== null && (
              <span className={styles.timeRange}>
                <Clock className="size-3.5" />
                {formatDateTimeShort(timeStart, locale)}
                {timeEnd !== null && (
                  <> – {formatDateTimeShort(timeEnd, locale)}</>
                )}
              </span>
            )}
            {props.activity.author_profile !== null &&
              props.activity.author_profile !== undefined && (
              <span className={styles.organizer}>
                <User className="size-3.5" />
                {props.activity.author_profile.title}
              </span>
            )}
          </div>
        </div>
      </article>
    </LocaleLink>
  );
}
