import * as React from "react";

import { SiteLink } from "@/shared/components/userland/site-link/site-link.tsx";
import { Card } from "@/shared/components/ui/card.tsx";
import { Button } from "@/shared/components/ui/button.tsx";
import { Icons } from "@/shared/components/icons.tsx";
import type { Profile, ProfileMembership } from "@/shared/modules/backend/profiles/types.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";

import styles from "./profile-list-item.module.css";

export type ProfileListItemProps = {
  focus?: "profile" | "member";
  profile: Profile;
  membership?: ProfileMembership;
};

export async function ProfileListItem(props: ProfileListItemProps) {
  const { t } = await getTranslations();

  const information = props.focus === "member" && props.membership
    ? {
      title: props.membership.member_profile.title,
      description: props.membership.member_profile.description,
      link: `/${props.membership.member_profile.slug}`,
      linkText: t("Contributions", "Profile Details"),
    }
    : {
      title: props.profile.title,
      description: props.profile.description,
      link: `/${props.profile.slug}`,
      linkText: t("Contributions", "Project Details"),
    };

  return (
    <Card className={styles.card}>
      <div className={styles.content}>
        <SiteLink href={information.link} className={styles.header}>
          <h3 className={styles.title}>{information.title}</h3>
          <div className={styles.stats}>
            {props.membership?.properties?.github?.commits !== undefined && (
              <span className={styles.stat}>
                <Icons.gitCommit className={styles.icon} />
                {props.membership?.properties?.github?.commits}
              </span>
            )}
            {props.membership?.properties?.github?.prs && (
              <span
                className={styles.stat}
                title={`${props.membership?.properties?.github?.prs.resolved} resolved out of ${props.membership?.properties?.github?.prs.total} PRs`}
              >
                <Icons.gitPullRequest className={styles.icon} />
                {props.membership?.properties?.github?.prs.resolved}/{props.membership?.properties?.github?.prs.total}
              </span>
            )}
            {props.membership?.properties?.github?.issues !== undefined && (
              <span className={styles.stat}>
                <Icons.gitFork className={styles.icon} />
                {props.membership?.properties?.github?.issues.total}
              </span>
            )}
            {props.membership?.properties?.github?.stars !== undefined && (
              <span className={styles.stat}>
                <Icons.star className={styles.icon} />
                {props.membership?.properties?.github?.stars}
              </span>
            )}
          </div>
        </SiteLink>
        <div className={styles.description}>{information.description}</div>
        <div className={styles.attributes}>
          <div className={styles.section}>
            <div className={styles.heading}>{t("Contributions", "My role in the project")}</div>
            <div className={styles.text}>{t("Contributions", props.membership?.kind ?? "-")}</div>
            <div className={styles.information}>
              <div>
                <Icons.gitCommit className={styles.icon} />
                <div>
                  <div className={styles.key}>{t("Contributions", "My first contribution")}</div>
                  <div className={styles.value}>01 Mayıs 2022</div>
                </div>
              </div>
              <div>
                <Icons.gitCommit className={styles.icon} />
                <div>
                  <div className={styles.key}>{t("Contributions", "My last contribution")}</div>
                  <div className={styles.value}>23 Mayıs 2025</div>
                </div>
              </div>
            </div>
          </div>
          {props.membership?.properties?.github !== undefined && (
            <>
              <div className={styles.section}>
                <div className={styles.heading}>Out of 50 total</div>
                <div className={styles.text}>{t("Contributions", "PRs")}</div>
                <div className={styles.information}>
                  <div>
                    <Icons.gitPullRequest className={styles.icon} />
                    <div>
                      <div className={styles.key}>{t("Contributions", "Opened PRs by me")}</div>
                      <div className={styles.value}>20</div>
                    </div>
                  </div>
                  <div>
                    <Icons.gitPullRequest className={styles.icon} />
                    <div>
                      <div className={styles.key}>{t("Contributions", "Merged PRs by me")}</div>
                      <div className={styles.value}>10</div>
                    </div>
                  </div>
                </div>
              </div>
              <div className={styles.section}>
                <div className={styles.heading}>Out of 200 total</div>
                <div className={styles.text}>{t("Contributions", "Issues")}</div>
                <div className={styles.information}>
                  <div>
                    <Icons.gitPullRequest className={styles.icon} />
                    <div>
                      <div className={styles.key}>{t("Contributions", "Issues opened by me")}</div>
                      <div className={styles.value}>20</div>
                    </div>
                  </div>
                  <div>
                    <Icons.gitPullRequest className={styles.icon} />
                    <div>
                      <div className={styles.key}>{t("Contributions", "Issues contributed by me")}</div>
                      <div className={styles.value}>10</div>
                    </div>
                  </div>
                </div>
              </div>
            </>
          )}
          {props.membership?.properties?.videos !== undefined && (
            <>
              <div className={styles.section}>
                <div className={styles.heading}>Out of 80 total</div>
                <div className={styles.text}>{t("Contributions", "Videos")}</div>
                <div className={styles.information}>
                  <div>
                    <Icons.video className={styles.icon} />
                    <div>
                      <div className={styles.key}>{t("Contributions", "Contributed by me")}</div>
                      <div className={styles.value}>{props.membership?.properties?.videos}</div>
                    </div>
                  </div>
                </div>
              </div>
            </>
          )}
          {
            /* <div className={styles.section}>
            <div>{t("Contributions", "My contributions")}</div>
            <div>~</div>
          </div> */
          }
        </div>
        <div className={styles.footer}>
          <SiteLink href={information.link} className={styles.header}>
            <Button variant="default">
              <Icons.star className={styles.icon} />
              {information.linkText}
              <Icons.arrowRight />
            </Button>
          </SiteLink>
        </div>
      </div>
    </Card>
  );
}
