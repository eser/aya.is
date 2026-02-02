// Profile links page
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  Globe,
  Instagram,
  Linkedin,
  Youtube,
  ExternalLink,
  EyeOff,
  type LucideIcon,
} from "lucide-react";
import i18next from "i18next";
import { backend, type ProfileLink, type ProfileLinkKind } from "@/modules/backend/backend";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { buildUrl, generateMetaTags } from "@/lib/seo";
import { Icon, Bsky, Discord, GitHub, Telegram, X } from "@/components/icons";
import { Card } from "@/components/ui/card";

const parentRoute = getRouteApi("/$locale/$slug");

type LinkTypeConfig = {
  kind: ProfileLinkKind;
  label: string;
  icon: Icon | LucideIcon;
};

const LINK_TYPES: LinkTypeConfig[] = [
  { kind: "github", label: "GitHub", icon: GitHub },
  { kind: "x", label: "X (Twitter)", icon: X },
  { kind: "linkedin", label: "LinkedIn", icon: Linkedin },
  { kind: "instagram", label: "Instagram", icon: Instagram },
  { kind: "youtube", label: "YouTube", icon: Youtube },
  { kind: "bsky", label: "Bluesky", icon: Bsky },
  { kind: "discord", label: "Discord", icon: Discord },
  { kind: "telegram", label: "Telegram", icon: Telegram },
  { kind: "website", label: "Website", icon: Globe },
];

function getLinkTypeConfig(kind: ProfileLinkKind): LinkTypeConfig {
  return LINK_TYPES.find((lt) => lt.kind === kind) ?? LINK_TYPES[LINK_TYPES.length - 1];
}

export const Route = createFileRoute("/$locale/$slug/links")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const links = await backend.getProfileLinks(locale, slug);
    const profile = await backend.getProfile(locale, slug);
    return { links, locale, slug, profileTitle: profile?.title ?? slug };
  },
  head: ({ loaderData }) => {
    const { locale, slug, profileTitle } = loaderData;
    const t = i18next.getFixedT(locale);
    return {
      meta: generateMetaTags({
        title: `${t("Layout.Links")} - ${profileTitle}`,
        description: t("Links.All links associated with this profile."),
        url: buildUrl(locale, slug, "links"),
        locale,
        type: "website",
      }),
    };
  },
  component: LinksPage,
});

function LinksPage() {
  const { links, locale, slug } = Route.useLoaderData();
  const { profile } = parentRoute.useLoaderData();
  const { t } = useTranslation();

  if (profile === null) {
    return null;
  }

  // Group links by their group field
  const groupedLinks = groupLinksByGroup(links ?? []);
  const groupNames = Object.keys(groupedLinks).sort((a, b) => {
    // Put "ungrouped" (empty string key) at the TOP
    if (a === "") return -1;
    if (b === "") return 1;
    return a.localeCompare(b);
  });

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale}>
      <div className="space-y-6">
        <div>
          <h2 className="font-serif text-2xl font-bold text-foreground">{t("Layout.Links")}</h2>
          <p className="text-muted-foreground">
            {t("Links.All links associated with this profile.")}
          </p>
        </div>

        {links !== null && links.length > 0
          ? (
            <div className="space-y-8">
              {groupNames.map((groupName) => (
                <div key={groupName || "ungrouped"}>
                  {/* Only show header for named groups, not for ungrouped */}
                  {groupName !== "" && (
                    <h3 className="text-lg font-semibold mb-4">{groupName}</h3>
                  )}
                  <div className="flex flex-col gap-3">
                    {groupedLinks[groupName].map((link) => (
                      <LinkCard key={link.id} link={link} t={t} />
                    ))}
                  </div>
                </div>
              ))}
            </div>
          )
          : (
            <p className="text-muted-foreground">
              {t("Links.No links found.")}
            </p>
          )}
      </div>
    </ProfileSidebarLayout>
  );
}

function groupLinksByGroup(links: ProfileLink[]): Record<string, ProfileLink[]> {
  const groups: Record<string, ProfileLink[]> = {};

  for (const link of links) {
    const groupKey = link.group ?? "";
    if (groups[groupKey] === undefined) {
      groups[groupKey] = [];
    }
    groups[groupKey].push(link);
  }

  return groups;
}

type LinkCardProps = {
  link: ProfileLink;
  t: (key: string) => string;
};

function LinkCard(props: LinkCardProps) {
  const { link, t } = props;
  const config = getLinkTypeConfig(link.kind);
  const IconComponent = config.icon;

  return (
    <Card className="p-4">
      <div className="flex items-center gap-4">
        <div className="flex items-center justify-center size-12 rounded-full bg-muted shrink-0">
          <IconComponent className="size-6" />
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <p className="font-medium">{link.title}</p>
            {link.visibility !== undefined && link.visibility !== "public" && (
              <span className="inline-flex items-center gap-1 text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded">
                <EyeOff className="size-3" />
                {t(`Profile.Visibility.${link.visibility}`)}
              </span>
            )}
          </div>
          {link.description !== null && link.description !== "" && (
            <p className="text-sm text-muted-foreground mt-1">{link.description}</p>
          )}
          {link.uri !== null && link.uri !== "" && (
            <a
              href={link.uri}
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm text-primary hover:underline inline-flex items-center gap-1 mt-1"
            >
              {link.uri}
              <ExternalLink className="size-3" />
            </a>
          )}
        </div>
      </div>
    </Card>
  );
}
