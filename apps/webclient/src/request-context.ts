export type DomainConfiguration = {
  type: "main";
  allowsWwwPrefix: boolean;
} | {
  type: "custom-domain";
  profileSlug: string;
  allowsWwwPrefix: boolean;
} | {
  type: "not-configured";
};

export type RequestContext = {
  domainConfiguration: DomainConfiguration;
  path: string[];
  originalPath: string[];
};
