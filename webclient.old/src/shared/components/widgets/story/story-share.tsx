"use client";

import React, { useState } from "react";
import { FaCheck, FaLinkedin, FaReddit, FaRegCopy, FaTelegram, FaTwitter, FaWhatsapp } from "react-icons/fa6";
import { FiDownload } from "react-icons/fi";

import type { StoryEx } from "@/shared/modules/backend/stories/types.ts";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { cn } from "@/shared/lib/cn.ts";
import styles from "./story-share.module.css";

export type StoryShareProps = {
  story: StoryEx;
  currentUrl: string;
};

export function StoryShare(props: StoryShareProps) {
  const { t } = useTranslations();
  const [copied, setCopied] = useState(false);

  const handleCopyLink = async () => {
    try {
      await navigator.clipboard.writeText(props.currentUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy link:", err);
    }
  };

  const handleWhatsAppShare = () => {
    const text = `${props.story.title} - ${props.story.summary}`;
    const url = `https://wa.me/?text=${encodeURIComponent(`${text} ${props.currentUrl}`)}`;
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const handleTelegramShare = () => {
    const text = `${props.story.title} - ${props.story.summary}`;
    const url = `https://t.me/share/url?url=${encodeURIComponent(props.currentUrl)}&text=${encodeURIComponent(text)}`;
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const handleLinkedInShare = () => {
    const url = `https://www.linkedin.com/sharing/share-offsite/?url=${encodeURIComponent(props.currentUrl)}`;
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const handleRedditShare = () => {
    const url = `https://reddit.com/submit?url=${encodeURIComponent(props.currentUrl)}&title=${
      encodeURIComponent(props.story.title || "")
    }`;
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const handleXShare = () => {
    const text = `${props.story.title} - ${props.story.summary}`;
    const url = `https://twitter.com/intent/tweet?text=${encodeURIComponent(text)}&url=${
      encodeURIComponent(props.currentUrl)
    }`;
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const handleDownloadMarkdown = () => {
    const markdownContent = `# ${props.story.title}\n\n${props.story.summary}\n\n${props.story.content}`;
    const blob = new Blob([markdownContent], { type: "text/markdown" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${props.story.slug || "story"}.md`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div className={styles["story-share"]}>
      <div className={styles["share-container"]}>
        <h3 className={styles["share-title"]}>{t("StoryPage", "Share")}:</h3>

        <div className={styles["share-buttons"]}>
          <button
            type="button"
            onClick={handleCopyLink}
            className={cn(styles["share-button"], copied && styles.copied)}
            aria-label={t("StoryPage", "Copy link")}
          >
            {copied ? <FaCheck className={styles["share-icon"]} /> : <FaRegCopy className={styles["share-icon"]} />}
          </button>

          <button
            type="button"
            onClick={handleWhatsAppShare}
            className={styles["share-button"]}
            aria-label={t("StoryPage", "Share on WhatsApp")}
          >
            <FaWhatsapp className={styles["share-icon"]} />
          </button>

          <button
            type="button"
            onClick={handleTelegramShare}
            className={styles["share-button"]}
            aria-label={t("StoryPage", "Share on Telegram")}
          >
            <FaTelegram className={styles["share-icon"]} />
          </button>

          <button
            type="button"
            onClick={handleLinkedInShare}
            className={styles["share-button"]}
            aria-label={t("StoryPage", "Share on LinkedIn")}
          >
            <FaLinkedin className={styles["share-icon"]} />
          </button>

          <button
            type="button"
            onClick={handleRedditShare}
            className={styles["share-button"]}
            aria-label={t("StoryPage", "Share on Reddit")}
          >
            <FaReddit className={styles["share-icon"]} />
          </button>

          <button
            type="button"
            onClick={handleXShare}
            className={styles["share-button"]}
            aria-label={t("StoryPage", "Share on X")}
          >
            <FaTwitter className={styles["share-icon"]} />
          </button>

          <button
            type="button"
            onClick={handleDownloadMarkdown}
            className={styles["share-button"]}
            aria-label={t("StoryPage", "Download Markdown")}
          >
            <FiDownload className={styles["share-icon"]} />
          </button>
        </div>
      </div>
    </div>
  );
}
