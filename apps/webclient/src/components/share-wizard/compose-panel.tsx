import { useTranslation } from "react-i18next";
import { Hash, Link as LinkIcon } from "lucide-react";
import type { StoryEx } from "@/modules/backend/types";
import styles from "./compose-panel.module.css";

const PLATFORM_LIMITS = [
  { id: "x", name: "X", limit: 280, warning: 250, linkCost: 23 },
  { id: "linkedin", name: "LinkedIn", limit: 3000, warning: 2800, linkCost: 0 },
  { id: "mastodon", name: "Mastodon", limit: 500, warning: 450, linkCost: 23 },
  { id: "bluesky", name: "Bluesky", limit: 300, warning: 270, linkCost: 0 },
] as const;

export type ComposePanelProps = {
  text: string;
  hashtags: string;
  story: StoryEx;
  currentUrl: string;
  onTextChange: (text: string) => void;
  onHashtagsChange: (hashtags: string) => void;
};

function computeCharCount(text: string, hashtags: string, linkCost: number): number {
  const hashtagText = hashtags.trim().length > 0 ? ` ${hashtags.trim()}` : "";
  const fullText = `${text}${hashtagText}`;
  // Link cost: some platforms shorten URLs (X, Mastodon use t.co-style ~23 chars)
  // We add the link cost since the story URL is always appended
  return fullText.length + (linkCost > 0 ? linkCost : 0);
}

function counterStatus(
  charCount: number,
  limit: number,
  warning: number,
): "ok" | "warning" | "exceeded" {
  if (charCount > limit) return "exceeded";
  if (charCount >= warning) return "warning";
  return "ok";
}

const counterStyleMap = {
  ok: styles.counterOk,
  warning: styles.counterWarning,
  exceeded: styles.counterExceeded,
} as const;

export function ComposePanel(props: ComposePanelProps) {
  const { t } = useTranslation();

  const handlePrefillSummary = () => {
    if (props.story.summary !== null && props.story.summary !== undefined) {
      props.onTextChange(props.story.summary);
    }
  };

  const handlePrefillTitle = () => {
    if (props.story.title !== null && props.story.title !== undefined) {
      props.onTextChange(props.story.title);
    }
  };

  return (
    <div className={styles.container}>
      {/* Compose text */}
      <div>
        <div className="flex items-center justify-between mb-1.5">
          <label className={styles.label} htmlFor="share-compose">
            {t("ShareWizard.Compose")}
          </label>
          <div className={styles.prefillRow}>
            {props.story.summary !== null && props.story.summary !== undefined && (
              <button type="button" className={styles.prefillButton} onClick={handlePrefillSummary}>
                {t("ShareWizard.Use Summary")}
              </button>
            )}
            {props.story.title !== null && props.story.title !== undefined && (
              <button type="button" className={styles.prefillButton} onClick={handlePrefillTitle}>
                {t("ShareWizard.Use Title")}
              </button>
            )}
          </div>
        </div>
        <textarea
          id="share-compose"
          className={styles.textarea}
          value={props.text}
          onChange={(e) => props.onTextChange(e.target.value)}
          placeholder={t("ShareWizard.Write your post...")}
        />
      </div>

      {/* Character counters */}
      <div className={styles.counters}>
        {PLATFORM_LIMITS.map((platform) => {
          const charCount = computeCharCount(props.text, props.hashtags, platform.linkCost);
          const status = counterStatus(charCount, platform.limit, platform.warning);
          const remaining = platform.limit - charCount;

          return (
            <span
              key={platform.id}
              className={`${styles.counter} ${counterStyleMap[status]}`}
              title={`${platform.name}: ${charCount}/${platform.limit}`}
            >
              {platform.name} <span>{remaining >= 0 ? remaining : remaining}</span>
            </span>
          );
        })}
      </div>

      {/* Hashtags */}
      <div className={styles.hashtagSection}>
        <label className={styles.label} htmlFor="share-hashtags">
          <Hash className="inline size-3.5 mr-1" />
          {t("ShareWizard.Hashtags")}
        </label>
        <input
          id="share-hashtags"
          type="text"
          className={styles.hashtagInput}
          value={props.hashtags}
          onChange={(e) => props.onHashtagsChange(e.target.value)}
          placeholder="#aya #tech #article"
        />
      </div>

      {/* Story link card */}
      <div>
        <label className={styles.label}>
          <LinkIcon className="inline size-3.5 mr-1" />
          {t("ShareWizard.Story Link")}
        </label>
        <div className={styles.storyLinkCard}>
          {props.story.story_picture_uri !== null && props.story.story_picture_uri !== undefined && (
            <img
              src={props.story.story_picture_uri}
              alt=""
              className={styles.storyLinkImage}
            />
          )}
          <div className={styles.storyLinkMeta}>
            <div className={styles.storyLinkTitle}>{props.story.title}</div>
            <div className={styles.storyLinkUrl}>{props.currentUrl}</div>
          </div>
        </div>
      </div>
    </div>
  );
}
