"use client";

import * as React from "react";
import { useTheme } from "next-themes";

import { cn } from "@/shared/lib/cn.ts";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { Icons } from "@/shared/components/icons.tsx";
import { Button } from "@/shared/components/ui/button.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/shared/components/ui/dropdown-menu.tsx";

import styles from "./theme-switcher.module.css";

export function ThemeSwitcher() {
  const { t } = useTranslations();

  const { theme, setTheme } = useTheme();

  return (
    <div className="hidden lg:block">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="sm" className={styles["theme-selector"]} asChild>
            <div>
              <Icons.sun className={cn(styles.sun, "rotate-0 scale-100 dark:-rotate-90 dark:scale-0")} />
              <Icons.moon className={cn(styles.moon, "rotate-90 scale-0 dark:rotate-0 dark:scale-100")} />
              <span className={styles.description}>{t("Layout", "Change theme")}</span>
            </div>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" sideOffset={15} forceMount>
          <DropdownMenuRadioGroup value={theme} onValueChange={setTheme}>
            <DropdownMenuRadioItem value="default">{t("Layout", "System")}</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="light">{t("Layout", "Light")}</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="midnight">{t("Layout", "Midnight")}</DropdownMenuRadioItem>
          </DropdownMenuRadioGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
