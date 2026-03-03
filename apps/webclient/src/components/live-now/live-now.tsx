import { useTranslation } from "react-i18next";
import { Radio } from "lucide-react";
import { SiteAvatar } from "@/components/userland";
import { LocaleLink } from "@/components/locale-link";
import type { LiveStreamInfo } from "@/modules/backend/types";
import styles from "./live-now.module.css";

type LiveNowSectionProps = {
  streams: LiveStreamInfo[];
  locale: string;
};

function buildBroadcastUrl(stream: LiveStreamInfo): string {
  if (stream.properties !== null && stream.properties !== undefined) {
    const onlineInfo = stream.properties.online_information;
    if (onlineInfo !== null && onlineInfo !== undefined && typeof onlineInfo === "object") {
      const broadcastUrl = (onlineInfo as Record<string, unknown>).broadcast_url;
      if (typeof broadcastUrl === "string" && broadcastUrl !== "") {
        return broadcastUrl;
      }
    }
  }

  if (stream.uri !== null && stream.uri !== undefined && stream.uri !== "") {
    return stream.uri.replace(/\/+$/, "") + "/live";
  }

  return "https://www.youtube.com";
}

function getStreamTitle(stream: LiveStreamInfo): string {
  if (stream.properties !== null && stream.properties !== undefined) {
    const onlineInfo = stream.properties.online_information;
    if (onlineInfo !== null && onlineInfo !== undefined && typeof onlineInfo === "object") {
      const title = (onlineInfo as Record<string, unknown>).title;
      if (typeof title === "string" && title !== "") {
        return title;
      }
    }
  }

  return stream.link_title;
}

function getThumbnailUrl(stream: LiveStreamInfo): string | null {
  if (stream.properties === null || stream.properties === undefined) {
    return null;
  }

  const onlineInfo = stream.properties.online_information;
  if (onlineInfo === null || onlineInfo === undefined || typeof onlineInfo !== "object") {
    return null;
  }

  const thumbnailUrl = (onlineInfo as Record<string, unknown>).thumbnail_url;
  if (typeof thumbnailUrl === "string" && thumbnailUrl !== "") {
    return thumbnailUrl;
  }

  return null;
}

type LiveStreamCardProps = {
  stream: LiveStreamInfo;
  locale: string;
};

function LiveStreamCard(props: LiveStreamCardProps) {
  const { t } = useTranslation();
  const broadcastUrl = buildBroadcastUrl(props.stream);
  const streamTitle = getStreamTitle(props.stream);
  const thumbnailUrl = getThumbnailUrl(props.stream);

  return (
    <div className={styles.card}>
      <div className={styles.thumbnailArea}>
        {thumbnailUrl !== null && (
          <img
            src={thumbnailUrl}
            alt={streamTitle}
            className={styles.thumbnailImage}
          />
        )}
        <div className={styles.liveBadge}>
          <span className={styles.liveBadgeDot} />
          {t("Home.Live")}
        </div>
      </div>

      <div className={styles.cardBody}>
        <div className={styles.avatarRow}>
          <div className={styles.avatarRing}>
            <SiteAvatar
              src={props.stream.profile_picture_uri}
              name={props.stream.profile_title}
              fallbackName={props.stream.profile_slug}
              className="size-12"
            />
          </div>
        </div>

        <div className={styles.streamTitle}>{streamTitle}</div>
        <div className={styles.profileName}>
          <LocaleLink
            to={`/${props.stream.profile_slug}`}
            className={styles.profileLink}
          >
            {props.stream.profile_title}
          </LocaleLink>
        </div>

        <a
          href={broadcastUrl}
          target="_blank"
          rel="noopener noreferrer"
          className={styles.watchButton}
        >
          <Radio className="size-4" />
          {t("Home.Watch Live")}
        </a>
      </div>
    </div>
  );
}

export function LiveNowSection(props: LiveNowSectionProps) {
  const { t } = useTranslation();

  if (props.streams.length === 0) {
    return null;
  }

  return (
    <section className={styles.section}>
      <div className={styles.container}>
        <div className={styles.header}>
          <span className={styles.headerDot} />
          <h2 className={styles.headerTitle}>{t("Home.Live Now")}</h2>
        </div>

        <div className={styles.scrollContainer}>
          {props.streams.map((stream) => (
            <LiveStreamCard
              key={stream.link_id}
              stream={stream}
              locale={props.locale}
            />
          ))}
        </div>
      </div>
    </section>
  );
}
