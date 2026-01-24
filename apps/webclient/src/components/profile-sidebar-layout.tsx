// Profile sidebar layout wrapper - use this in profile child routes that need the sidebar
import { useTranslation } from "react-i18next";
import { Globe, Instagram, Link, Linkedin, SquarePen, Youtube } from "lucide-react";
import { Bsky, Discord, GitHub, Telegram, X } from "@/components/icons";
import { type Profile } from "@/modules/backend/backend";
import { LocaleLink } from "@/components/locale-link";
import { useProfilePermissions } from "@/lib/hooks/use-profile-permissions";

function findIcon(kind: string) {
  switch (kind) {
    case "github":
      return GitHub;
    case "twitter":
    case "x":
      return X;
    case "linkedin":
      return Linkedin;
    case "instagram":
      return Instagram;
    case "youtube":
      return Youtube;
    case "bsky":
      return Bsky;
    case "discord":
      return Discord;
    case "telegram":
      return Telegram;
    case "website":
      return Globe;
    default:
      return Link;
  }
}

type ProfileSidebarLayoutProps = {
  profile: Profile;
  slug: string;
  locale: string;
  children: React.ReactNode;
};

export function ProfileSidebarLayout(props: ProfileSidebarLayoutProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-[280px_1fr] gap-8 items-start">
      <ProfileSidebar profile={props.profile} slug={props.slug} locale={props.locale} />
      <main>{props.children}</main>
    </div>
  );
}

type ProfileSidebarProps = {
  profile: Profile;
  slug: string;
  locale: string;
};

function ProfileSidebar(props: ProfileSidebarProps) {
  const { t } = useTranslation();
  const { canEdit } = useProfilePermissions(props.locale, props.slug);

  return (
    <aside className="flex flex-col gap-4">
      {/* Profile Picture */}
      {props.profile.profile_picture_uri !== null &&
        props.profile.profile_picture_uri !== undefined && (
        <div className="flex justify-center md:justify-start">
          <img
            src={props.profile.profile_picture_uri}
            alt={`${props.profile.title}'s profile picture`}
            width={280}
            height={280}
            className="border rounded-full"
          />
        </div>
      )}

      {/* Hero Section */}
      <div>
        <h1 className="mt-0 mb-2 font-serif text-base font-semibold leading-none sm:text-lg md:text-xl lg:text-2xl">
          {props.profile.title}
        </h1>

        <div className="mt-0 mb-4 font-sans text-sm font-light leading-none sm:text-base md:text-lg lg:text-xl text-muted-foreground">
          {props.profile.slug}
          {props.profile.pronouns !== null &&
            props.profile.pronouns !== undefined &&
            ` · ${props.profile.pronouns}`}
        </div>

        {props.profile.links !== null &&
          props.profile.links !== undefined &&
          props.profile.links.length > 0 && (
          <div className="flex gap-4 mb-3 text-sm text-muted-foreground">
            {props.profile.links.map((link) => {
              const Icon = findIcon(link.kind);
              return (
                <a
                  key={link.id}
                  href={link.uri ?? undefined}
                  title={link.title !== null && link.title !== undefined ? link.title : link.kind}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="no-underline"
                >
                  <Icon className="transition-colors hover:text-foreground h-5 w-5" />
                </a>
              );
            })}
          </div>
        )}

        {props.profile.description !== null &&
          props.profile.description !== undefined && (
          <p className="mt-0 mb-4 font-sans text-sm font-normal leading-snug text-left">
            {props.profile.description}
            {canEdit && (
              <span className="flex">
                <LocaleLink
                  to={`/${props.slug}/settings`}
                  className="no-underline transition-colors text-muted-foreground hover:text-foreground flex text-sm items-center"
                >
                  <SquarePen size="16" className="mr-1.5" />
                  {t("Profile.Edit Profile")}
                </LocaleLink>
              </span>
            )}
          </p>
        )}
      </div>

      {/* Navigation */}
      <nav className="flex justify-center font-serif md:justify-start">
        <ul className="flex flex-row flex-wrap justify-center p-0 space-y-0 md:space-y-3 lg:space-y-4 list-none md:flex-col">
          <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
            <LocaleLink
              to={`/${props.slug}`}
              className="no-underline transition-colors text-muted-foreground hover:text-foreground"
            >
              {t("Layout.Profile")}
            </LocaleLink>
          </li>

          {props.profile.kind === "individual" && (
            <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
              <LocaleLink
                to={`/${props.slug}/contributions`}
                className="no-underline transition-colors text-muted-foreground hover:text-foreground"
              >
                {t("Layout.Contributions")}
              </LocaleLink>
            </li>
          )}

          {(props.profile.kind === "organization" ||
            props.profile.kind === "project" ||
            props.profile.kind === "product") && (
            <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
              <LocaleLink
                to={`/${props.slug}/members`}
                className="no-underline transition-colors text-muted-foreground hover:text-foreground"
              >
                {t("Layout.Members")}
              </LocaleLink>
            </li>
          )}

          {props.profile.pages?.map((page) => (
            <li
              key={page.slug}
              className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none"
            >
              <LocaleLink
                to={`/${props.slug}/${page.slug}`}
                className="no-underline transition-colors text-muted-foreground hover:text-foreground"
              >
                {page.title}
              </LocaleLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}
