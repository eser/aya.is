import { useTranslation } from "react-i18next";
import { LocaleLink } from "@/components/locale-link";
import { SiteAvatar } from "@/components/userland/site-avatar";
import type { ProfileMembership } from "@/modules/backend/types";
import styles from "./member-card.module.css";

export type MemberCardProps = {
  membership: ProfileMembership;
};

export function MemberCard(props: MemberCardProps) {
  const { t } = useTranslation();
  const memberProfile = props.membership.member_profile;
  const properties = props.membership.properties;
  const githubStats = properties !== null ? properties.github : undefined;
  const videoStats = properties !== null ? properties.videos : undefined;

  return (
    <LocaleLink
      role="card"
      to={`/${memberProfile.slug}`}
      className={styles.cardLink}
    >
      <div className={styles.memberCard}>
        <div className={styles.avatarContainer}>
          <SiteAvatar
            src={memberProfile.profile_picture_uri}
            name={memberProfile.title}
            fallbackName={memberProfile.slug}
            className={styles.avatar}
          />
        </div>

        <div className={styles.content}>
          <div className={styles.titleRow}>
            <h3 className={styles.title}>{memberProfile.title}</h3>
            <span className={styles.role}>
              {t(`Contributions.${props.membership.kind}`)}
            </span>
          </div>
          {memberProfile.description !== null &&
            memberProfile.description !== undefined && (
            <p className={styles.description}>{memberProfile.description}</p>
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
