import { useTranslation } from "react-i18next";
import { LocaleBadge } from "@/components/locale-badge";
import { LocaleLink } from "@/components/locale-link";
import { SiteAvatar } from "@/components/userland/site-avatar";
import { Badge } from "@/components/ui/badge";
import { stripMarkdown } from "@/lib/strip-markdown";
import { InlineMarkdown } from "@/lib/inline-markdown";
import { getProfilePictureUrl } from "@/lib/profile-picture";
import type { Profile } from "@/modules/backend/types";
import styles from "./profile-card.module.css";

export type ProfileCardProps = {
  profile: Profile;
  /** Display mode: "avatar" shows centered avatar, "cover" shows full-width cover image with badge */
  variant?: "avatar" | "cover";
  /** Show the profile kind badge (only used in "cover" variant) */
  showKindBadge?: boolean;
};

export function ProfileCard(props: ProfileCardProps) {
  const { t } = useTranslation();
  const { profile, variant = "avatar", showKindBadge = false } = props;

  if (variant === "cover") {
    return (
      <LocaleLink
        role="card"
        to={`/${profile.slug}`}
        className={styles.cardLink}
      >
        <div className={styles.coverCard}>
          <div className={styles.coverImageContainer}>
            <img
              src={getProfilePictureUrl(profile.profile_picture_uri, profile.title ?? "", profile.slug)}
              alt={profile.title ?? profile.slug}
              className={styles.coverImage}
            />
            {showKindBadge && (
              <div className={styles.badgeContainer}>
                <Badge variant="secondary" className={styles.kindBadge}>
                  {t(`Contributions.${profile.kind}`)}
                </Badge>
              </div>
            )}
          </div>
          <div className={styles.coverInfo}>
            <h3 className={styles.title}>
              {stripMarkdown(profile.title ?? "")}
              <LocaleBadge localeCode={profile.locale_code} className={styles.localeBadge} />
            </h3>
            {profile.description !== null && profile.description !== undefined && (
              <InlineMarkdown content={profile.description} className={styles.description} />
            )}
          </div>
        </div>
      </LocaleLink>
    );
  }

  return (
    <LocaleLink
      role="card"
      to={`/${profile.slug}`}
      className={styles.cardLink}
    >
      <div className={styles.profileCard}>
        <div className={styles.avatarContainer}>
          <SiteAvatar
            src={profile.profile_picture_uri}
            name={profile.title}
            fallbackName={profile.slug}
            className={styles.avatar}
          />
        </div>
        <div className={styles.info}>
          <h3 className={styles.title}>
            {stripMarkdown(profile.title ?? "")}
            <LocaleBadge localeCode={profile.locale_code} className={styles.localeBadge} />
          </h3>
          {profile.description !== null &&
            profile.description !== undefined && (
              <InlineMarkdown content={profile.description} className={styles.description} />
          )}
        </div>
      </div>
    </LocaleLink>
  );
}
