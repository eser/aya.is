import type { SupportedLocaleCode } from "@/config";

export type DomainConfiguration = {
  type: "main";
  defaultCulture: SupportedLocaleCode;
  allowsWwwPrefix: boolean;
} | {
  type: "custom-domain";
  defaultCulture: SupportedLocaleCode;
  profileSlug: string;
  profileTitle: string | null;
  allowsWwwPrefix: boolean;
} | {
  type: "not-configured";
};

export type RequestContext = {
  domainConfiguration: DomainConfiguration;
  path: string[];
  originalPath: string[];
  cookieHeader?: string;
  ssrTheme?: string;
};
