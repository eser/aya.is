// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { useTranslation } from "react-i18next";
import { Sparkles } from "lucide-react";
import { Button } from "@/components/ui/button";
import styles from "./ai-assist-panel.module.css";

export type AiAssistPanelProps = {
  text: string;
  storyContent: string | null;
  storySlug: string;
  locale: string;
  activePlatform: string;
};

export function AiAssistPanel(_props: AiAssistPanelProps) {
  const { t } = useTranslation();

  // AI features are temporarily disabled — they will be enabled once
  // profile points integration is in place to avoid spending AI credits.

  return (
    <div className={styles.container}>
      <div className={styles.title}>
        <Sparkles className="size-4" />
        {t("ShareWizard.AI Assist")}
      </div>

      <div className={styles.actions}>
        <Button variant="outline" size="sm" className={styles.actionButton} disabled>
          {t("ShareWizard.Summarize")}
        </Button>

        <Button variant="outline" size="sm" className={styles.actionButton} disabled>
          {t("ShareWizard.Optimize for Platform")}
        </Button>

        <Button variant="outline" size="sm" className={styles.actionButton} disabled>
          {t("ShareWizard.Adjust Tone")}
        </Button>

        <Button variant="outline" size="sm" className={styles.actionButton} disabled>
          {t("ShareWizard.Translate")}
        </Button>
      </div>

      <p className="text-xs text-muted-foreground mt-2">
        {t("Common.Coming soon")}
      </p>
    </div>
  );
}
