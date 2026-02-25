import { supportedLocales, type SupportedLocaleCode } from "@/config";

export function LocaleBadge(props: { localeCode?: string; viewerLocale: string; className?: string }) {
  if (props.localeCode === undefined || props.localeCode === props.viewerLocale) {
    return null;
  }

  const data = supportedLocales[props.localeCode as SupportedLocaleCode];
  const label = data?.englishName ?? props.localeCode;

  return <span className={props.className}>{label}</span>;
}
