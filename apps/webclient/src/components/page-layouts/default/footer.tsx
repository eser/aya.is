import * as React from "react";
import { useTranslation } from "react-i18next";
import { Separator } from "@/components/ui/separator";
import { LocaleLink } from "@/components/locale-link";
import { useNavigation } from "@/modules/navigation/navigation-context";
import { LocaleSwitcher } from "./locale-switcher";

export function Footer() {
  const { t } = useTranslation();
  const { isCustomDomain } = useNavigation();

  return (
    <footer className="border-0 border-t border-solid border-t-border">
      <div className="container mx-auto px-4 py-6 flex flex-col sm:flex-row justify-between sm:items-center gap-6 sm:gap-0">
        <div className="flex flex-row flex-wrap justify-center sm:justify-start gap-2 sm:gap-4 text-sm">
          <LocaleLink
            to="/aya"
            className="transition-colors hover:text-foreground/80 text-foreground/60"
          >
            AYA
          </LocaleLink>
          <Separator className="h-auto" orientation="vertical" decorative />
          <LocaleLink
            to="/aya/policies"
            className="transition-colors hover:text-foreground/80 text-foreground/60"
          >
            {t("Layout.Policies")}
          </LocaleLink>
          {!isCustomDomain && (
            <>
              <Separator className="h-auto" orientation="vertical" decorative />
              <LocaleLink
                to="/aya/about"
                className="transition-colors hover:text-foreground/80 text-foreground/60"
              >
                {t("Layout.About")}
              </LocaleLink>
            </>
          )}
        </div>
        <div className="flex flex-row justify-center sm:justify-end gap-2 sm:gap-4 text-sm">
          <LocaleSwitcher />
        </div>
      </div>
    </footer>
  );
}
