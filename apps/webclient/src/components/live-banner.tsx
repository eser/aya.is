import { useTranslation } from "react-i18next";
import { Radio } from "lucide-react";
import { ExternalLink } from "@/components/external-link";
import type { ProfileLink } from "@/modules/backend/types";
import styles from "./live-banner.module.css";

type LiveBannerProps = {
  links: ProfileLink[] | undefined | null;
};

function buildLiveUrl(link: ProfileLink): string {
  // Use broadcast URL from properties if available, otherwise build from channel URL
  if (link.properties !== null && link.properties !== undefined) {
    const onlineInfo = link.properties.online_information;
    if (onlineInfo !== null && onlineInfo !== undefined && typeof onlineInfo === "object") {
      const broadcastUrl = (onlineInfo as Record<string, unknown>).broadcast_url;
      if (typeof broadcastUrl === "string" && broadcastUrl !== "") {
        return broadcastUrl;
      }
    }
  }

  if (link.uri !== null && link.uri !== undefined && link.uri !== "") {
    // YouTube channel URLs like https://youtube.com/@channel — append /live
    return link.uri.replace(/\/+$/, "") + "/live";
  }

  return "https://www.youtube.com";
}

function getOnlineTitle(link: ProfileLink): string | null {
  if (link.properties === null || link.properties === undefined) {
    return null;
  }

  const onlineInfo = link.properties.online_information;
  if (onlineInfo === null || onlineInfo === undefined || typeof onlineInfo !== "object") {
    return null;
  }

  const title = (onlineInfo as Record<string, unknown>).title;
  if (typeof title === "string" && title !== "") {
    return title;
  }

  return null;
}

export function LiveBanner(props: LiveBannerProps) {
  const { t } = useTranslation();

  if (props.links === null || props.links === undefined) {
    return null;
  }

  const liveLink = props.links.find(
    (link) => link.is_online === true && link.kind === "youtube",
  );

  if (liveLink === undefined) {
    return null;
  }

  const liveUrl = buildLiveUrl(liveLink);
  const onlineTitle = getOnlineTitle(liveLink);

  return (
    <div className={styles.banner}>
      <div className={styles.indicator}>
        <div className={styles.dot} />
        <span className={styles.label}>{t("Profile.Live")}</span>
      </div>
      <span className={styles.title}>
        {liveLink.title}
        {onlineTitle !== null && <span className={styles.streamTitle}>· {onlineTitle}</span>}
      </span>
      <ExternalLink
        href={liveUrl}
        className={styles.watchLink}
      >
        <Radio className="size-3.5" />
        {t("Profile.Watch Live")}
      </ExternalLink>
    </div>
  );
}
