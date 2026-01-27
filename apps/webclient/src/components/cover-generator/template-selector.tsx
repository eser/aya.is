// Template Selector Component
// Grid of template thumbnails for selection

import * as React from "react";
import { useTranslation } from "react-i18next";
import type { TemplateId } from "@/lib/cover-generator/types.ts";
import { getTemplateList, type TemplateConfig } from "@/lib/cover-generator/templates.ts";
import { cn } from "@/lib/utils.ts";
import styles from "./cover-generator.module.css";

interface TemplateSelectorProps {
  selectedTemplate: TemplateId;
  onSelect: (templateId: TemplateId) => void;
}

export function TemplateSelector(props: TemplateSelectorProps) {
  const { t } = useTranslation();
  const templates = getTemplateList();

  return (
    <div className={styles.templateSelector}>
      <h3 className={styles.sectionTitle}>{t("CoverDesigner.Templates")}</h3>
      <div className={styles.templateGrid}>
        {templates.map((template) => (
          <TemplateCard
            key={template.id}
            template={template}
            isSelected={props.selectedTemplate === template.id}
            onSelect={() => props.onSelect(template.id)}
          />
        ))}
      </div>
    </div>
  );
}

interface TemplateCardProps {
  template: TemplateConfig;
  isSelected: boolean;
  onSelect: () => void;
}

function TemplateCard(props: TemplateCardProps) {
  const { template, isSelected, onSelect } = props;

  // Generate a simple preview based on template defaults
  const bgColor = template.defaults.backgroundColor ?? "#1a1a2e";
  const accentColor = template.defaults.accentColor ?? "#e94560";
  const textColor = template.defaults.textColor ?? "#ffffff";

  return (
    <button
      type="button"
      className={cn(styles.templateCard, isSelected && styles.templateCardSelected)}
      onClick={onSelect}
      aria-pressed={isSelected}
    >
      <div
        className={styles.templatePreview}
        style={{ backgroundColor: bgColor }}
      >
        {/* Mini preview representation */}
        <div className={styles.templatePreviewContent}>
          {/* Title placeholder */}
          <div
            className={styles.templatePreviewTitle}
            style={{ backgroundColor: textColor }}
          />
          <div
            className={styles.templatePreviewTitle}
            style={{ backgroundColor: textColor, width: "60%" }}
          />
          {/* Accent line for bold template */}
          {template.id === "bold" && (
            <div
              className={styles.templatePreviewAccent}
              style={{ backgroundColor: accentColor }}
            />
          )}
          {/* Author placeholder */}
          {template.defaults.showAuthor !== false && (
            <div className={styles.templatePreviewAuthor}>
              <div
                className={styles.templatePreviewAvatar}
                style={{ backgroundColor: textColor }}
              />
              <div
                className={styles.templatePreviewName}
                style={{ backgroundColor: textColor }}
              />
            </div>
          )}
        </div>
        {/* Pattern overlay */}
        {template.defaults.backgroundPattern === "diagonal" && (
          <div className={styles.templatePreviewPatternDiagonal} />
        )}
        {template.defaults.backgroundPattern === "dots" && (
          <div className={styles.templatePreviewPatternDots} />
        )}
      </div>
      <span className={styles.templateName}>{template.name}</span>
    </button>
  );
}
