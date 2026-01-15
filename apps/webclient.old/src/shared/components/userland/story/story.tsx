"use client";

import * as React from "react";
import NextImage from "next/image";

import { SiteLink } from "@/shared/components/userland/site-link/site-link.tsx";
import type { Story, StoryEx } from "@/shared/modules/backend/stories/types.ts";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { cn } from "@/shared/lib/utils.ts";
import { formatDateString } from "@/shared/lib/date.ts";

import styles from "./story.module.css";

const Image = NextImage;

export type StoryProps = {
  profileSlug?: string;
  story: Story | StoryEx;
};

export function Story(props: StoryProps) {
  const { t, locale } = useTranslations();

  let href: string;
  if (props.profileSlug !== undefined) {
    href = `/${props.profileSlug}/stories/${props.story.slug}`;
  } else {
    href = `/stories/${props.story.slug}`;
  }

  return (
    <SiteLink role="card" href={href}>
      <article className={styles.story}>
        <div className={cn(styles["image-container"], "w-[250px] h-[150px]")}>
          {props.story.story_picture_uri
            ? (
              <Image
                src={props.story.story_picture_uri}
                alt={props.story.title ?? t("News", "News item image")}
                width={250}
                height={150}
                className={styles.image}
              />
            )
            : (
              <div className={styles["image-placeholder"]}>
                {props.story.title ?? t("News", "No image available")}
              </div>
            )}
          {props.story.author_profile?.profile_picture_uri && (
            <div className={styles["author-avatar-container"]}>
              <Image
                src={props.story.author_profile.profile_picture_uri}
                alt={props.story.author_profile.title ?? "Author image"}
                width={60}
                height={60}
                className={styles["author-avatar-image"]}
              />
            </div>
          )}
        </div>
        <div className={styles["content-area"]}>
          <h3 className={styles.title}>{props.story.title}</h3>
          <p className={styles.summary}>{props.story.summary}</p>
          <div className={styles.meta}>
            {props.story.created_at !== null && (
              <time dateTime={props.story.created_at} className={styles.date}>
                {formatDateString(props.story.created_at, locale.code)}
              </time>
            )}
            {props.story.author_profile !== null && (
              <span className={styles.author}>
                {props.story.author_profile.title}
              </span>
            )}
          </div>
        </div>
      </article>
    </SiteLink>
  );
}
