import { useParams } from "@tanstack/react-router";
import { supportedLocales, type SupportedLocaleCode } from "@/config";

export function LocaleBadge(props: { localeCode?: string; className?: string }) {
  const params = useParams({ strict: false }) as { locale?: string };
  const viewerLocale = params.locale ?? "";

  if (
    props.localeCode === undefined
    || props.localeCode === ""
    || props.localeCode === viewerLocale
  ) {
    return null;
  }

  const data = supportedLocales[props.localeCode as SupportedLocaleCode];
  if (data === undefined) {
    return null;
  }

  return <span className={props.className}>{data.englishName}</span>;
}
