"use client";

import { useTranslation } from "react-i18next";
import {
  Newspaper,
  PencilLine,
  Megaphone,
  Info,
  Images,
  Presentation,
} from "lucide-react";
import { LocaleLink } from "@/components/locale-link";
import { cn } from "@/lib/utils";
import { formatDateString } from "@/lib/date";
import type { Story as StoryType, StoryEx, StoryKind } from "@/modules/backend/types";
import styles from "./story.module.css";

const storyKindIcons: Record<StoryKind, React.ElementType> = {
  news: Newspaper,
  article: PencilLine,
  announcement: Megaphone,
  status: Info,
  content: Images,
  presentation: Presentation,
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
                alt={props.story.title ?? t("News.News item image")}
                width={250}
                height={150}
                className={styles.image}
              />
            )
            : (
              <div className={styles.imagePlaceholder}>
                {props.story.title ?? t("News.No image available")}
              </div>
            )}
          {props.story.author_profile?.profile_picture_uri !== null &&
            props.story.author_profile?.profile_picture_uri !== undefined && (
            <div className={styles.authorAvatarContainer}>
              <img
                src={props.story.author_profile.profile_picture_uri}
                alt={props.story.author_profile.title ?? "Author image"}
                width={60}
                height={60}
                className={styles.authorAvatarImage}
              />
            </div>
          )}
        </div>
        <div className={styles.contentArea}>
          <h3 className={styles.title}>{props.story.title}</h3>
          <p className={styles.summary}>{props.story.summary}</p>
          <div className={styles.meta}>
            {props.story.created_at !== null && (
              <span className="flex items-center gap-1.5">
                {(() => {
                  const KindIcon = storyKindIcons[props.story.kind];
                  return KindIcon !== undefined ? <KindIcon className="size-3.5" /> : null;
                })()}
                <time dateTime={props.story.created_at}>
                  {formatDateString(props.story.created_at, locale)}
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
