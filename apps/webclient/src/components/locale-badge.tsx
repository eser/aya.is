import { useParams } from "@tanstack/react-router";
import { supportedLocales, type SupportedLocaleCode } from "@/config";

export function LocaleBadge(props: { localeCode?: string; className?: string }) {
  const params = useParams({ strict: false }) as { locale?: string };
  const viewerLocale = params.locale ?? "";
  const localeCode = props.localeCode?.trim() ?? "";

  if (localeCode === "" || localeCode === viewerLocale) {
    return null;
  }

  const data = supportedLocales[localeCode as SupportedLocaleCode];
  if (data === undefined) {
    return null;
  }

  return <span className={props.className}>{data.englishName}</span>;
}
