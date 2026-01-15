import React from "react";
// import Image from "next/image";

import type { StoryEx } from "@/shared/modules/backend/stories/types.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { SiteLink } from "@/shared/components/userland/site-link/site-link.tsx";
import styles from "./story-footer.module.css";

export type StoryFooterProps = {
  story: StoryEx;
};

export async function StoryFooter(props: StoryFooterProps) {
  if (!props.story.author_profile) {
    return null;
  }

  const { t } = await getTranslations();

  // Filter out publications where the id matches the author profile id
  const filteredPublications = props.story.publications.filter(
    (publication) => publication.id !== props.story.author_profile.id,
  );

  return (
    <div className={styles["story-footer"]}>
      <div className={styles["author-details"]}>
        {props.story.author_profile.profile_picture_uri && (
          <div className={styles["author-avatar-container"]}>
            <img
              src={props.story.author_profile.profile_picture_uri}
              alt={props.story.author_profile.title ?? "Author image"}
              width={72}
              height={72}
              className={styles["author-avatar-image"]}
            />
          </div>
        )}
        <div className={styles["author-text-container"]}>
          <div className={styles["written-by-text"]}>{t("Stories", "Written by")}</div>
          <div className={styles["author-name"]}>
            <SiteLink href={`/${props.story.author_profile.slug}`}>
              {props.story.author_profile.title}
            </SiteLink>
          </div>
          <div className={styles["author-description"]}>{props.story.author_profile.description}</div>
        </div>
      </div>
      {filteredPublications.length > 0 && (
        <div className={styles["publications-details"]}>
          <div className={styles["publications-label"]}>{t("Stories", "Publications")}</div>
          <div className={styles["publications-list"]}>
            {filteredPublications.map((publication) => (
              <div key={publication.id} className={styles["publication-item"]}>
                <SiteLink href={`/${publication.slug}`}>
                  {publication.title}
                </SiteLink>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
