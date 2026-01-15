"use client";

import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/components/ui/card.tsx";
import styles from "./page.module.css";

export default function ProfileAdvancedPage() {
  const { t } = useTranslations();

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>
          {t("Profile", "Advanced Settings") ?? "Advanced Settings"}
        </h1>
        <p className={styles.description}>
          {t("Profile", "Advanced configuration and developer options for your profile.") ??
            "Advanced configuration and developer options for your profile."}
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("Profile", "Coming Soon") ?? "Coming Soon"}</CardTitle>
          <CardDescription>
            {t("Profile", "Advanced settings and configuration options are under development.") ??
              "Advanced settings and configuration options are under development."}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className={styles.comingSoon}>
            <span className={styles.icon}>🔧</span>
            <p className={styles.message}>
              {t(
                "Profile",
                "Soon you'll have access to advanced features like API keys, webhooks, and custom integrations.",
              ) || "Soon you'll have access to advanced features like API keys, webhooks, and custom integrations."}
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
