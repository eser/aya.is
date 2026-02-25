"use client";

import { useTranslation } from "react-i18next";
import {
  Newspaper,
  PencilLine,
  Megaphone,
  Info,
  Images,
  Presentation,
  Calendar,
} from "lucide-react";
import { LocaleLink } from "@/components/locale-link";
import { SiteAvatar } from "@/components/userland/site-avatar";
import { cn } from "@/lib/utils";
import { formatDateString } from "@/lib/date";
import { stripMarkdown } from "@/lib/strip-markdown";
import { InlineMarkdown } from "@/lib/inline-markdown";
import { LocaleBadge } from "@/components/locale-badge";
import type { Story as StoryType, StoryEx, StoryKind } from "@/modules/backend/types";
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

export type StoryProps = {
  profileSlug?: string;
  story: StoryType | StoryEx;
  /** Whether this is rendered on a custom domain (affects link generation) */
  isCustomDomain?: boolean;
};

export function Story(props: StoryProps) {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;

  // On custom domain, profile slug is implicit (from URL rewriting)
  // So links are relative to root: /stories/slug
  // On main domain, links include profile slug: /profile/stories/slug
  let href: string;
  if (props.isCustomDomain === true) {
    href = `/stories/${props.story.slug}`;
  } else if (props.profileSlug !== undefined) {
    href = `/${props.profileSlug}/stories/${props.story.slug}`;
  } else {
    href = `/stories/${props.story.slug}`;
  }

  return (
    <LocaleLink role="card" to={href} className="no-underline block">
      <article className={styles.story}>
        <div className={cn(styles.imageContainer, "w-[250px] h-[150px]")}>
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
              <div className={styles.imagePlaceholder}>
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
            {(props.story.published_at ?? props.story.created_at) !== null && (
              <span className="flex items-center gap-1.5">
                {(() => {
                  const KindIcon = storyKindIcons[props.story.kind];
                  return KindIcon !== undefined ? <KindIcon className="size-3.5" /> : null;
                })()}
                <time dateTime={props.story.published_at ?? props.story.created_at}>
                  {formatDateString(props.story.published_at ?? props.story.created_at, locale)}
                </time>
              </span>
            )}
            {props.story.author_profile !== null && (
              <span className={styles.author}>
                {props.story.author_profile.title}
              </span>
            )}
          </div>
        </div>
      </article>
    </LocaleLink>
  );
}
