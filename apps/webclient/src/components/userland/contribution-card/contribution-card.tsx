import { useTranslation } from "react-i18next";
import { LocaleLink } from "@/components/locale-link";
import { SiteAvatar } from "@/components/userland/site-avatar";
import type { ProfileMembership } from "@/modules/backend/types";
import styles from "./contribution-card.module.css";

export type ContributionCardProps = {
  membership: ProfileMembership;
};

export function ContributionCard(props: ContributionCardProps) {
  const { t } = useTranslation();
  const profile = props.membership.profile;
  const properties = props.membership.properties;
  const githubStats = properties !== null ? properties.github : undefined;
  const videoStats = properties !== null ? properties.videos : undefined;

  return (
    <LocaleLink
      role="card"
      to={`/${profile.slug}`}
      className={styles.cardLink}
    >
      <div className={styles.contributionCard}>
        <div className={styles.avatarContainer}>
          <SiteAvatar
            src={profile.profile_picture_uri}
            name={profile.title}
            fallbackName={profile.slug}
            className={styles.avatar}
          />
        </div>

        <div className={styles.content}>
          <div className={styles.titleRow}>
            <h3 className={styles.title}>{profile.title}</h3>
            <span className={styles.role}>
              {t(`Contributions.${props.membership.kind}`)}
            </span>
            <span className={styles.profileKind}>
              {t(`Contributions.${profile.kind}`)}
            </span>
          </div>
          {profile.description !== null &&
            profile.description !== undefined &&
            profile.description !== "" && (
            <p className={styles.description}>{profile.description}</p>
          )}
        </div>

        {githubStats !== undefined && (
          <div className={styles.stats}>
            <div className={styles.statsGrid}>
              <div className={styles.statItem}>
                <span className={styles.statValue}>{githubStats.commits}</span>
                <span className={styles.statLabel}>{t("Members.Commits")}</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statValue}>
                  {githubStats.prs.resolved}/{githubStats.prs.total}
                </span>
                <span className={styles.statLabel}>{t("Contributions.PRs")}</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statValue}>
                  {githubStats.issues.resolved}/{githubStats.issues.total}
                </span>
                <span className={styles.statLabel}>{t("Contributions.Issues")}</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statValue}>{githubStats.stars}</span>
                <span className={styles.statLabel}>{t("Members.Stars")}</span>
              </div>
            </div>
          </div>
        )}

        {videoStats !== undefined && (
          <div className={styles.stats}>
            <div className={styles.statsGrid}>
              <div className={styles.statItem}>
                <span className={styles.statValue}>{videoStats}</span>
                <span className={styles.statLabel}>{t("Contributions.Videos")}</span>
              </div>
            </div>
          </div>
        )}
      </div>
    </LocaleLink>
  );
}
