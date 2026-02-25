// Profile sidebar layout wrapper - use this in profile child routes that need the sidebar
import * as React from "react";
import { useTranslation } from "react-i18next";
import { Coins, Ellipsis, Globe, Instagram, Link, Linkedin, SquarePen, UserMinus, UserPlus, Youtube } from "lucide-react";
import { toast } from "sonner";
import { Bsky, Discord, GitHub, SpeakerDeck, Telegram, X } from "@/components/icons";
import { backend, type Profile } from "@/modules/backend/backend";
import { LocaleBadge } from "@/components/locale-badge";
import { LocaleLink } from "@/components/locale-link";
import { SiteAvatar } from "@/components/userland";
import { useAuth } from "@/lib/auth/auth-context";
import { InlineMarkdown } from "@/lib/inline-markdown";
import { useProfilePermissions } from "@/lib/hooks/use-profile-permissions";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

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
    case "speakerdeck":
      return SpeakerDeck;
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
  viewerMembershipKind?: string | null;
  children: React.ReactNode;
};

export function ProfileSidebarLayout(props: ProfileSidebarLayoutProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-[280px_1fr] gap-8 items-start">
      <ProfileSidebar profile={props.profile} slug={props.slug} locale={props.locale} viewerMembershipKind={props.viewerMembershipKind} />
      <main className="min-w-0">{props.children}</main>
    </div>
  );
}

type ProfileSidebarProps = {
  profile: Profile;
  slug: string;
  locale: string;
  viewerMembershipKind?: string | null;
};

// Membership kinds at or above sponsor level (cannot self-unfollow, shown as badge)
const ADVANCED_KINDS = new Set(["sponsor", "contributor", "maintainer", "lead", "owner"]);

function ProfileSidebar(props: ProfileSidebarProps) {
  const { t } = useTranslation();
  const { canEdit } = useProfilePermissions(props.profile.id);
  const { isAuthenticated, user } = useAuth();

  const [isUnfollowDialogOpen, setIsUnfollowDialogOpen] = React.useState(false);
  const [isProcessing, setIsProcessing] = React.useState(false);
  const [localMembershipKind, setLocalMembershipKind] = React.useState<string | null>(
    props.viewerMembershipKind ?? null,
  );

  // Sync when prop changes (e.g. navigation between profiles)
  React.useEffect(() => {
    setLocalMembershipKind(props.viewerMembershipKind ?? null);
  }, [props.viewerMembershipKind]);

  // Determine viewer's membership relationship to this profile
  const isOwnProfile = user?.individual_profile_id === props.profile.id;

  const isFollower = localMembershipKind === "follower";
  const isAdvanced = localMembershipKind !== null && ADVANCED_KINDS.has(localMembershipKind);
  const showFollowButton = isAuthenticated && !isOwnProfile && localMembershipKind === null;
  const showUnfollowButton = isAuthenticated && !isOwnProfile && isFollower;
  const showBadge = isAuthenticated && !isOwnProfile && isAdvanced;

  const handleFollow = async () => {
    setIsProcessing(true);
    try {
      const success = await backend.followProfile(props.locale, props.slug);
      if (success) {
        toast.success(t("Profile.Followed successfully"));
        setLocalMembershipKind("follower");
      } else {
        toast.error(t("Profile.Failed to follow"));
      }
    } catch {
      toast.error(t("Profile.Failed to follow"));
    }
    setIsProcessing(false);
  };

  const handleUnfollow = async () => {
    setIsProcessing(true);
    try {
      const success = await backend.unfollowProfile(props.locale, props.slug);
      if (success) {
        toast.success(t("Profile.Unfollowed successfully"));
        setIsUnfollowDialogOpen(false);
        setLocalMembershipKind(null);
      } else {
        toast.error(t("Profile.Failed to unfollow"));
      }
    } catch {
      toast.error(t("Profile.Failed to unfollow"));
    }
    setIsProcessing(false);
  };

  const buttonClass = "no-underline inline-flex items-center gap-1.5 rounded-md bg-background/80 backdrop-blur-sm px-2.5 py-1.5 text-xs font-medium text-muted-foreground shadow-sm border border-border/50 transition-colors hover:text-foreground hover:bg-background";

  return (
    <aside className="flex flex-col gap-4">
      {/* Profile Picture */}
      <div className="relative flex justify-center md:justify-start">
        <SiteAvatar
          src={props.profile.profile_picture_uri}
          name={props.profile.title}
          fallbackName={props.slug}
          className="size-[280px] border"
        />
        <div className="absolute bottom-0 right-0 flex gap-1">
          {canEdit && (
            <LocaleLink
              to={`/${props.slug}/settings`}
              className={buttonClass}
            >
              <SquarePen size="14" />
              {t("Profile.Edit Profile")}
            </LocaleLink>
          )}
          {showFollowButton && (
            <button
              type="button"
              onClick={handleFollow}
              disabled={isProcessing}
              className={buttonClass}
            >
              <UserPlus size="14" />
              {isProcessing ? t("Common.Saving...") : t("Profile.Follow")}
            </button>
          )}
          {showUnfollowButton && (
            <DropdownMenu>
              <DropdownMenuTrigger className={buttonClass}>
                <Ellipsis size="14" />
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => setIsUnfollowDialogOpen(true)}>
                  <UserMinus size="14" />
                  {t("Profile.Unfollow")}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          )}
          {showBadge && localMembershipKind !== null && (
            <span className={`${buttonClass} cursor-default`}>
              {t(`Profile.MembershipKind.${localMembershipKind}`)}
            </span>
          )}
        </div>
      </div>

      {/* Unfollow Confirmation Dialog */}
      <AlertDialog open={isUnfollowDialogOpen} onOpenChange={setIsUnfollowDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("Profile.Unfollow")}</AlertDialogTitle>
            <AlertDialogDescription>
              {t("Profile.Are you sure you want to unfollow this profile?")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t("Common.Cancel")}</AlertDialogCancel>
            <AlertDialogAction onClick={handleUnfollow} disabled={isProcessing}>
              {isProcessing ? t("Common.Saving...") : t("Profile.Unfollow")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

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
          <InlineMarkdown
            content={props.profile.description}
            className="mt-0 mb-4 font-sans text-sm font-normal leading-snug text-left"
          />
        )}
      </div>

      {/* Points Display */}
      {/* {props.profile.points > 0 && (
        <div className="flex items-center gap-2 text-muted-foreground">
          <Coins className="size-4" />
          <span className="font-semibold text-foreground">{props.profile.points.toLocaleString()}</span>
          <span>{t("Profile.points")}</span>
        </div>
      )} */}

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

          {props.profile.kind === "individual" && props.profile.feature_relations === "public" && (
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
            props.profile.kind === "product") && (
            <>
              {props.profile.feature_relations === "public" && (
                <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
                  <LocaleLink
                    to={`/${props.slug}/members`}
                    className="no-underline transition-colors text-muted-foreground hover:text-foreground"
                  >
                    {t("Layout.Members")}
                  </LocaleLink>
                </li>
              )}
              {props.profile.feature_links === "public" && (
                <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
                  <LocaleLink
                    to={`/${props.slug}/links`}
                    className="no-underline transition-colors text-muted-foreground hover:text-foreground"
                  >
                    {t("Layout.Links")}
                  </LocaleLink>
                </li>
              )}
            </>
          )}

          {props.profile.feature_qa === "public" && (
            <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
              <LocaleLink
                to={`/${props.slug}/qa`}
                className="no-underline transition-colors text-muted-foreground hover:text-foreground"
              >
                {t("Layout.Q&A")}
              </LocaleLink>
            </li>
          )}

          {props.profile.pages?.filter((page) => page.visibility === "public").map((page) => (
            <li
              key={page.slug}
              className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none"
            >
              <LocaleLink
                to={`/${props.slug}/${page.slug}`}
                className="no-underline transition-colors text-muted-foreground hover:text-foreground"
              >
                {page.title}
                <LocaleBadge
                  localeCode={page.locale_code}
                  className="text-xs font-medium px-2 py-0.5 rounded-full bg-primary/10 text-primary ml-2 align-middle"
                />
              </LocaleLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}
