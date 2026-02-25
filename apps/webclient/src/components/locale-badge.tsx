import { useParams } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { SUPPORTED_LOCALES, type SupportedLocaleCode } from "@/config";

export function LocaleBadge(props: { localeCode?: string; className?: string }) {
  const params = useParams({ strict: false }) as { locale?: string };
  const { t } = useTranslation();
  const viewerLocale = params.locale ?? "";
  const localeCode = props.localeCode?.trim() ?? "";

  if (localeCode === "" || localeCode === viewerLocale) {
    return null;
  }

  if (!SUPPORTED_LOCALES.includes(localeCode as SupportedLocaleCode)) {
    return null;
  }

  return <span className={props.className}>{t(`Locales.${localeCode}`)}</span>;
}
