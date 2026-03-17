// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Download, FileText, ImageIcon, Megaphone, Quote } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import type { StoryEx } from "@/modules/backend/types";
import styles from "./image-composer.module.css";

type Template = {
  id: string;
  labelKey: string;
  icon: React.ComponentType<{ className?: string }>;
};

type Dimension = {
  id: string;
  labelKey: string;
  width: number;
  height: number;
};

const TEMPLATES: Template[] = [
  { id: "quote-card", labelKey: "ShareWizard.Quote Card", icon: Quote },
  { id: "article-card", labelKey: "ShareWizard.Article Card", icon: FileText },
  { id: "announcement-card", labelKey: "ShareWizard.Announcement Card", icon: Megaphone },
];

const DIMENSIONS: Dimension[] = [
  { id: "og-image", labelKey: "ShareWizard.OG Image", width: 1200, height: 630 },
  { id: "instagram-post", labelKey: "ShareWizard.Instagram Post", width: 1080, height: 1350 },
  { id: "instagram-story", labelKey: "ShareWizard.Instagram Story", width: 1080, height: 1920 },
];

export type ImageComposerProps = {
  story: StoryEx;
  locale: string;
};

export function ImageComposer(props: ImageComposerProps) {
  const { t } = useTranslation();
  const [selectedTemplate, setSelectedTemplate] = useState("article-card");
  const [selectedDimension, setSelectedDimension] = useState("og-image");

  const dimension = DIMENSIONS.find((d) => d.id === selectedDimension) ?? DIMENSIONS[0];

  const handleDownload = () => {
    // Phase 1: Use story cover image as download if available
    if (props.story.story_picture_uri !== null && props.story.story_picture_uri !== undefined) {
      const link = document.createElement("a");
      link.href = props.story.story_picture_uri;
      link.download = `${props.story.slug ?? "share"}-${selectedDimension}.png`;
      link.target = "_blank";
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.title}>
        <ImageIcon className="size-4" />
        {t("ShareWizard.Generate Image")}
      </div>

      {/* Template selection */}
      <div className={styles.templates}>
        {TEMPLATES.map((tmpl) => {
          const Icon = tmpl.icon;
          return (
            <button
              key={tmpl.id}
              type="button"
              className={cn(styles.templateCard, selectedTemplate === tmpl.id && styles.templateCardActive)}
              onClick={() => setSelectedTemplate(tmpl.id)}
            >
              <Icon className={styles.templateIcon} />
              <div className={styles.templateLabel}>{t(tmpl.labelKey)}</div>
            </button>
          );
        })}
      </div>

      {/* Dimension selection */}
      <div className={styles.dimensions}>
        {DIMENSIONS.map((dim) => (
          <button
            key={dim.id}
            type="button"
            className={cn(styles.dimensionButton, selectedDimension === dim.id && styles.dimensionButtonActive)}
            onClick={() => setSelectedDimension(dim.id)}
          >
            {t(dim.labelKey)} ({dim.width}x{dim.height})
          </button>
        ))}
      </div>

      {/* Preview */}
      <div className={styles.preview} style={{ aspectRatio: `${dimension.width}/${dimension.height}` }}>
        {props.story.story_picture_uri !== null && props.story.story_picture_uri !== undefined
          ? (
            <img
              src={props.story.story_picture_uri}
              alt={props.story.title ?? ""}
              className={styles.previewImage}
            />
          )
          : <span className={styles.previewPlaceholder}>{t("ShareWizard.No cover image")}</span>}
      </div>

      {/* Download */}
      <Button
        variant="outline"
        size="sm"
        className={styles.downloadButton}
        onClick={handleDownload}
        disabled={props.story.story_picture_uri === null || props.story.story_picture_uri === undefined}
      >
        <Download className="size-3.5 mr-1.5" />
        {t("ShareWizard.Download Image")}
      </Button>
    </div>
  );
}
