// Profile links page
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { ChildNotFound } from "./route";
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
import { Icon, Bsky, Discord, GitHub, SpeakerDeck, Telegram, X } from "@/components/icons";
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
  { kind: "speakerdeck", label: "SpeakerDeck", icon: SpeakerDeck },
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
    const profile = await backend.getProfile(locale, slug);

    if (profile?.feature_links === "disabled") {
      return { links: null, locale, slug, profileTitle: slug, translatedTitle: "", translatedDescription: "", notFound: true as const };
    }

    const links = await backend.getProfileLinks(locale, slug);

    if (links === null) {
      return { links: null, locale, slug, profileTitle: slug, translatedTitle: "", translatedDescription: "", notFound: true as const };
    }

    // Ensure locale translations are loaded before translating
    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    const translatedTitle = `${t("Layout.Links")} - ${profile?.title ?? slug}`;
    const translatedDescription = t("Links.All links associated with this profile.");

    return {
      links,
      locale,
      slug,
      profileTitle: profile?.title ?? slug,
      translatedTitle,
      translatedDescription,
      notFound: false as const,
    };
  },
  head: ({ loaderData }) => {
    if (loaderData === undefined || loaderData.notFound) {
      return { meta: [] };
    }

    const { locale, slug, translatedTitle, translatedDescription } = loaderData;

    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, slug, "links"),
        locale,
        type: "website",
      }),
    };
  },
  component: LinksPage,
  notFoundComponent: ChildNotFound,
});

function LinksPage() {
  const loaderData = Route.useLoaderData();
  const { profile, permissions } = parentRoute.useLoaderData();
  const { t } = useTranslation();

  if (loaderData.notFound || loaderData.links === null || profile === null) {
    return <ChildNotFound />;
  }

  const { links, locale, slug } = loaderData;

  // Group links by their group field
  const groupedLinks = groupLinksByGroup(links ?? []);
  const groupNames = Object.keys(groupedLinks).sort((a, b) => {
    // Put "ungrouped" (empty string key) at the TOP
    if (a === "") return -1;
    if (b === "") return 1;
    return a.localeCompare(b);
  });

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale} viewerMembershipKind={permissions?.viewer_membership_kind}>
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
                  <div className="grid grid-cols-1 xl:grid-cols-2 gap-2">
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

  // Use custom icon if specified, otherwise fall back to kind-based icon
  const hasCustomIcon = link.icon !== undefined && link.icon !== null && link.icon !== "";

  const cardContent = (
    <div className="flex items-center gap-3">
      <div className="flex items-center justify-center size-9 rounded-full bg-muted shrink-0">
        {hasCustomIcon
          ? <span className="text-base leading-none">{link.icon}</span>
          : <IconComponent className="size-4" />}
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <p className="font-medium text-sm">{link.title}</p>
          {link.visibility !== undefined && link.visibility !== "public" && (
            <span className="inline-flex items-center gap-1 text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded">
              <EyeOff className="size-3" />
              {t(`Profile.Visibility.${link.visibility}`)}
            </span>
          )}
        </div>
        {link.description !== null && link.description !== "" && (
          <p className="text-xs text-muted-foreground mt-0.5 line-clamp-1">{link.description}</p>
        )}
      </div>
      {link.uri !== null && link.uri !== "" && (
        <ExternalLink className="size-4 text-muted-foreground shrink-0" />
      )}
    </div>
  );

  if (link.uri !== null && link.uri !== "") {
    return (
      <a
        href={link.uri}
        target="_blank"
        rel="noopener noreferrer"
        className="block"
      >
        <Card className="px-3 py-2.5 hover:bg-muted/50 transition-colors cursor-pointer">
          {cardContent}
        </Card>
      </a>
    );
  }

  return (
    <Card className="px-3 py-2.5">
      {cardContent}
    </Card>
  );
}
