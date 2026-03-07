import { useState } from "react";
import { useTranslation } from "react-i18next";
import { ExternalLink, Settings2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import type { StoryEx } from "@/modules/backend/types";
import styles from "./platform-preview.module.css";

type Platform = {
  id: string;
  name: string;
  limit: number;
  linkCost: number;
  buildIntentUrl: (text: string, storyUrl: string) => string;
};

const PLATFORMS: Platform[] = [
  {
    id: "x",
    name: "X",
    limit: 280,
    linkCost: 23,
    buildIntentUrl: (text, url) =>
      `https://twitter.com/intent/tweet?text=${encodeURIComponent(text)}&url=${encodeURIComponent(url)}`,
  },
  {
    id: "linkedin",
    name: "LinkedIn",
    limit: 3000,
    linkCost: 0,
    buildIntentUrl: (_text, url) => `https://www.linkedin.com/sharing/share-offsite/?url=${encodeURIComponent(url)}`,
  },
  {
    id: "mastodon",
    name: "Mastodon",
    limit: 500,
    linkCost: 23,
    buildIntentUrl: (text, url) => `https://mastodon.social/share?text=${encodeURIComponent(`${text}\n${url}`)}`,
  },
  {
    id: "bluesky",
    name: "Bluesky",
    limit: 300,
    linkCost: 0,
    buildIntentUrl: (text, url) => `https://bsky.app/intent/compose?text=${encodeURIComponent(`${text}\n${url}`)}`,
  },
];

export type PlatformPreviewProps = {
  text: string;
  story: StoryEx;
  currentUrl: string;
  platformOverrides: Map<string, string>;
  activePlatform: string;
  onActivePlatformChange: (platformId: string) => void;
  onPlatformOverrideChange: (platformId: string, text: string | null) => void;
};

function getDisplayText(
  platformId: string,
  mainText: string,
  overrides: Map<string, string>,
): string {
  return overrides.get(platformId) ?? mainText;
}

function getEffectiveLength(text: string, linkCost: number): number {
  return text.length + (linkCost > 0 ? linkCost : 0);
}

export function PlatformPreview(props: PlatformPreviewProps) {
  const { t } = useTranslation();
  const [customizing, setCustomizing] = useState<Set<string>>(new Set());

  const activeTab = props.activePlatform;
  const activePlatformDef = PLATFORMS.find((p) => p.id === activeTab) ?? PLATFORMS[0];
  const displayText = getDisplayText(activeTab, props.text, props.platformOverrides);
  const effectiveLength = getEffectiveLength(displayText, activePlatformDef.linkCost);
  const isExceeded = effectiveLength > activePlatformDef.limit;
  const isCustomized = props.platformOverrides.has(activeTab);

  const handleToggleCustomize = (platformId: string) => {
    const next = new Set(customizing);
    if (next.has(platformId)) {
      next.delete(platformId);
      props.onPlatformOverrideChange(platformId, null);
    } else {
      next.add(platformId);
      props.onPlatformOverrideChange(platformId, props.text);
    }
    setCustomizing(next);
  };

  const handleCustomTextChange = (platformId: string, value: string) => {
    props.onPlatformOverrideChange(platformId, value);
  };

  const handleOpenIntent = (platform: Platform) => {
    const text = getDisplayText(platform.id, props.text, props.platformOverrides);
    const url = platform.buildIntentUrl(text, props.currentUrl);
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const authorName = props.story.author_profile?.title ?? "Author";

  return (
    <div className={styles.container}>
      <div className={styles.sectionTitle}>{t("ShareWizard.Platform Previews")}</div>

      {/* Tabs */}
      <div className={styles.tabs}>
        {PLATFORMS.map((platform) => (
          <button
            key={platform.id}
            type="button"
            className={cn(styles.tab, activeTab === platform.id && styles.tabActive)}
            onClick={() => props.onActivePlatformChange(platform.id)}
          >
            {platform.name}
          </button>
        ))}
      </div>

      {/* Preview card */}
      <div className={styles.previewCard}>
        {/* Header */}
        <div className={styles.previewHeader}>
          {props.story.author_profile?.profile_picture_uri !== null &&
              props.story.author_profile?.profile_picture_uri !== undefined
            ? (
              <img
                src={props.story.author_profile.profile_picture_uri}
                alt=""
                className={styles.previewAvatar}
              />
            )
            : <div className={styles.previewAvatar} />}
          <div className={styles.previewNameGroup}>
            <span className={styles.previewName}>{authorName}</span>
            <span className={styles.previewHandle}>
              @{props.story.author_profile?.slug ?? "user"}
            </span>
          </div>
        </div>

        {/* Text */}
        <div className={styles.previewText}>
          {displayText.length > 0 ? displayText : t("ShareWizard.Write your post...")}
        </div>

        {/* Truncation warning */}
        {isExceeded && (
          <div className={styles.previewTruncated}>
            {t("ShareWizard.Exceeds limit")} ({effectiveLength}/{activePlatformDef.limit})
          </div>
        )}

        {/* Link preview card */}
        <div className={styles.previewLinkCard}>
          {props.story.story_picture_uri !== null && props.story.story_picture_uri !== undefined && (
            <img
              src={props.story.story_picture_uri}
              alt=""
              className={styles.previewLinkImage}
            />
          )}
          <div className={styles.previewLinkMeta}>
            <div className={styles.previewLinkTitle}>{props.story.title}</div>
            <div className={styles.previewLinkUrl}>{props.currentUrl}</div>
          </div>
        </div>

        {/* Customize toggle */}
        <button
          type="button"
          className={styles.customizeToggle}
          onClick={() => handleToggleCustomize(activeTab)}
        >
          <Settings2 className="size-3.5" />
          {isCustomized ? t("ShareWizard.Remove customization") : t("ShareWizard.Customize for this platform")}
        </button>

        {/* Per-platform text override */}
        {customizing.has(activeTab) && (
          <textarea
            className={styles.customizeTextarea}
            value={props.platformOverrides.get(activeTab) ?? props.text}
            onChange={(e) => handleCustomTextChange(activeTab, e.target.value)}
            placeholder={t("ShareWizard.Custom text for platform", { platform: activePlatformDef.name })}
          />
        )}

        {/* Action buttons */}
        <div className={styles.previewActions}>
          <Button size="sm" onClick={() => handleOpenIntent(activePlatformDef)}>
            <ExternalLink className="size-3.5 mr-1.5" />
            {t("ShareWizard.Open in platform", { platform: activePlatformDef.name })}
          </Button>
        </div>
      </div>
    </div>
  );
}
