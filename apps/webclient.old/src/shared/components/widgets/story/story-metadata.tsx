import React from "react";

import type { StoryEx } from "@/shared/modules/backend/stories/types.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { calculateReadingTime } from "@/shared/lib/reading-time.ts";
import styles from "./story-metadata.module.css";

export type StoryMetadataProps = {
  story: StoryEx;
};

export async function StoryMetadata(props: StoryMetadataProps) {
  const { t, locale } = await getTranslations();

  const readingTime = calculateReadingTime(props.story.content);

  // Format the published date
  const publishedDate = new Date(props.story.created_at);
  const formattedDate = publishedDate.toLocaleDateString(locale.code, {
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  return (
    <div className={styles["story-metadata"]}>
      <div className={styles["publish-date"]}>
        <span className={styles.label}>{t("StoryPage", "Published on")}:</span>
        <time dateTime={props.story.created_at} className={styles.date}>
          {formattedDate}
        </time>
      </div>

      <div className={styles["reading-time"]}>
        <span className={styles.label}>{t("StoryPage", "Reading time")}:</span>
        <span className={styles.time}>
          {readingTime} {t("StoryPage", "min read")}
        </span>
      </div>
    </div>
  );
}
