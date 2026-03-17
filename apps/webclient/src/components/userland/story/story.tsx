// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { useTranslation } from "react-i18next";
import { Calendar, Images, Info, Megaphone, Newspaper, PencilLine, Presentation, Tag } from "lucide-react";
import { LocaleLink } from "@/components/locale-link";
import { SiteAvatar } from "@/components/userland/site-avatar";
import { cn } from "@/lib/utils";
import { formatDateString, formatDateTimeRange, formatDateTimeShort } from "@/lib/date";
import { stripMarkdown } from "@/lib/strip-markdown";
import { InlineMarkdown } from "@/lib/inline-markdown";
import { LocaleBadge } from "@/components/locale-badge";
import type { ActivityProperties, Story as StoryType, StoryEx, StoryKind } from "@/modules/backend/types";
import styles from "./story.module.css";

const storyKindIcons: Record<StoryKind, React.ElementType> = {
  news: Newspaper,
  article: PencilLine,
  announcement: Megaphone,
  status: Info,
  content: Images,
  presentation: Presentation,
  activity: Calendar,
};

const activityKindLabels: Record<string, string> = {
  meetup: "Activities.Meetup",
  workshop: "Activities.Workshop",
  conference: "Activities.Conference",
  broadcast: "Activities.Broadcast",
  meeting: "Activities.Meeting",
};

export type StoryProps = {
  profileSlug?: string;
  story: StoryType | StoryEx;
  /** Whether this is rendered on a custom domain (affects link generation) */
  isCustomDomain?: boolean;
};

export function Story(props: StoryProps) {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;
  const isActivity = props.story.kind === "activity";

  let href: string;
  if (props.isCustomDomain === true) {
    href = `/stories/${props.story.slug}`;
  } else if (props.profileSlug !== undefined) {
    href = `/${props.profileSlug}/stories/${props.story.slug}`;
  } else {
    href = `/stories/${props.story.slug}`;
  }

  return (
    <LocaleLink data-slot="card" to={href} className="no-underline block">
      <article className={styles.story}>
        <div className={cn(styles.imageContainer, "w-full h-[150px] md:w-[250px]")}>
          {props.story.story_picture_uri !== null &&
              props.story.story_picture_uri !== undefined
            ? (
              <img
                src={props.story.story_picture_uri}
                alt={stripMarkdown(props.story.title ?? t("News.News item image"))}
                width={250}
                height={150}
                className={styles.image}
              />
            )
            : (
              <div className={styles.imagePlaceholderText}>
                {stripMarkdown(props.story.title ?? t("News.No image available"))}
              </div>
            )}
          {props.story.author_profile !== null &&
            props.story.author_profile !== undefined && (
            <div className={styles.authorAvatarContainer}>
              <SiteAvatar
                src={props.story.author_profile.profile_picture_uri}
                name={props.story.author_profile.title ?? "Author"}
                fallbackName={props.story.author_profile.slug}
                className={styles.authorAvatarImage}
              />
            </div>
          )}
        </div>
        <div className={styles.contentArea}>
          <h3 className={styles.title}>
            {stripMarkdown(props.story.title ?? "")}
            <LocaleBadge localeCode={props.story.locale_code} className={styles.localeBadge} />
          </h3>
          {props.story.summary !== null && props.story.summary !== undefined && (
            <InlineMarkdown content={props.story.summary} className={styles.summary} />
          )}
          <div className={styles.meta}>
            {isActivity
              ? (
                <ActivityMeta
                  properties={props.story.properties}
                  authorProfile={props.story.author_profile}
                  locale={locale}
                />
              )
              : <StoryMeta story={props.story} locale={locale} />}
          </div>
        </div>
      </article>
    </LocaleLink>
  );
}

function StoryMeta(props: { story: StoryType | StoryEx; locale: string }) {
  const dateStr = props.story.published_at ?? props.story.created_at;

  return (
    <>
      {dateStr !== null && (
        <span className="flex items-center gap-1.5">
          {(() => {
            const KindIcon = storyKindIcons[props.story.kind];
            return KindIcon !== undefined ? <KindIcon className="size-3.5" /> : null;
          })()}
          <time dateTime={dateStr}>
            {formatDateString(dateStr, props.locale)}
          </time>
        </span>
      )}
      {props.story.author_profile !== null && (
        <span className={styles.author}>
          {props.story.author_profile.title}
        </span>
      )}
    </>
  );
}

function ActivityMeta(props: {
  properties: Record<string, unknown> | null;
  authorProfile: { title: string | null; slug: string } | null;
  locale: string;
}) {
  const { t } = useTranslation();
  const activityProps = (props.properties ?? {}) as unknown as ActivityProperties;
  const dateMode = activityProps.date_mode ?? "fixed";
  const timeStart = activityProps.activity_time_start !== undefined
    ? new Date(activityProps.activity_time_start)
    : null;
  const timeEnd = activityProps.activity_time_end !== undefined ? new Date(activityProps.activity_time_end) : null;
  const kindLabel = activityKindLabels[activityProps.activity_kind ?? "meetup"] ?? "Activities.Meetup";

  return (
    <>
      {dateMode === "undecided"
        ? (
          <span className="flex items-center gap-1.5">
            <Calendar className="size-3.5" />
            <span className="inline-block px-1.5 py-0.5 rounded-full bg-amber-100 text-amber-800 text-[0.65rem] font-medium dark:bg-amber-900 dark:text-amber-200">
              {t("Activities.Date Undecided")}
            </span>
          </span>
        )
        : timeStart !== null && (
          <span className="flex items-center gap-1.5">
            <Calendar className="size-3.5" />
            {timeEnd !== null
              ? formatDateTimeRange(timeStart, timeEnd, props.locale)
              : formatDateTimeShort(timeStart, props.locale)}
          </span>
        )}
      <span className="flex items-center gap-1.5">
        <Tag className="size-3.5" />
        {t(kindLabel)}
      </span>
      {props.authorProfile !== null && (
        <span className={styles.author}>
          {props.authorProfile.title}
        </span>
      )}
    </>
  );
}
