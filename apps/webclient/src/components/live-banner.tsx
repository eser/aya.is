import { useTranslation } from "react-i18next";
import { Radio } from "lucide-react";
import type { ProfileLink } from "@/modules/backend/types";
import styles from "./live-banner.module.css";

type LiveBannerProps = {
  links: ProfileLink[] | undefined | null;
};

function buildLiveUrl(link: ProfileLink): string {
  if (link.uri !== null && link.uri !== undefined && link.uri !== "") {
    // YouTube channel URLs like https://youtube.com/@channel — append /live
    return link.uri.replace(/\/+$/, "") + "/live";
  }

  return "https://www.youtube.com";
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

  return (
    <div className={styles.banner}>
      <div className={styles.indicator}>
        <div className={styles.dot} />
        <span className={styles.label}>{t("Profile.Live")}</span>
      </div>
      <span className={styles.title}>
        {liveLink.title}
      </span>
      <a
        href={liveUrl}
        target="_blank"
        rel="noopener noreferrer"
        className={styles.watchLink}
      >
        <Radio className="size-3.5" />
        {t("Profile.Watch Live")}
      </a>
    </div>
  );
}
