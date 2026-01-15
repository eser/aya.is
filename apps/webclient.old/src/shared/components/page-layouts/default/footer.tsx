import * as React from "react";

import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { Separator } from "@/shared/components/ui/separator.tsx";
import { SiteLink } from "@/shared/components/userland/site-link/site-link.tsx";

import { LocaleSwitcher } from "./locale-switcher.tsx";
import styles from "./footer.module.css";

export async function Footer() {
  const { t } = await getTranslations();

  return (
    <footer className={styles.footer}>
      <div>
        <div className={styles.links}>
          <SiteLink href="/aya/">AYA</SiteLink>
          <Separator className={styles.div} orientation="vertical" decorative />
          <SiteLink href="/aya/policies/">{t("Layout", "Policies")}</SiteLink>
          <Separator className={styles.div} orientation="vertical" decorative />
          <SiteLink href="/aya/about/">{t("Layout", "About")}</SiteLink>
        </div>
        <div className={styles.locales}>
          <LocaleSwitcher />
        </div>
      </div>
    </footer>
  );
}
