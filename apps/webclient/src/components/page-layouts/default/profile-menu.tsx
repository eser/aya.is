import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { siteConfig } from "@/config";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { SiteAvatar } from "@/components/userland";
import { useAuth } from "@/lib/auth/auth-context";
import { getCurrentLanguage } from "@/modules/i18n/i18n";
import { useNavigation } from "@/modules/navigation/navigation-context";

type ProfileMenuProps = {
  className?: string;
};

export function ProfileMenu(props: ProfileMenuProps) {
  const { isAuthenticated, user, logout } = useAuth();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const locale = getCurrentLanguage();
  const { isCustomDomain } = useNavigation();

  if (!isAuthenticated || user === null) {
    return null;
  }

  // Determine avatar URL with priority: GitHub avatar > null (SiteAvatar will use dicebear fallback)
  const githubAvatarUrl = user.github_handle !== undefined && user.github_handle !== null
    ? `https://github.com/${user.github_handle}.png?size=32`
    : null;

  const accessibleProfiles = user.accessible_profiles !== undefined
    ? user.accessible_profiles
    : [];

  const navigateToMain = (path: string) => {
    if (isCustomDomain) {
      globalThis.location.href = `${siteConfig.host}${path}`;
    } else {
      navigate({ to: path });
    }
  };

  const handleProfileClick = () => {
    // Navigate to user's profile if they have one, otherwise to create profile page
    if (user.individual_profile_slug !== undefined) {
      navigateToMain(`/${locale}/${user.individual_profile_slug}`);
    } else {
      navigateToMain(`/${locale}/elements/new`);
    }
  };

  const handleLogout = async () => {
    await logout();
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        render={(triggerProps) => (
          <Button
            {...triggerProps}
            variant="ghost"
            className={`relative h-8 w-8 rounded-full p-0 ${props.className !== undefined ? props.className : ""}`}
          >
            <SiteAvatar
              src={githubAvatarUrl}
              name={user.name ?? "User"}
              size="default"
            />
          </Button>
        )}
      />
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuGroup>
          <DropdownMenuLabel className="font-normal">
            <div className="flex flex-col space-y-1">
              <p className="text-sm font-medium leading-none">
                {user.name ?? user.email ?? t("Layout.Profile")}
                {user.individual_profile?.points !== undefined && (
                  <span className="ml-1 text-xs font-normal text-muted-foreground">
                    ({user.individual_profile.points.toLocaleString()}xp)
                  </span>
                )}
              </p>
              {user.github_handle !== undefined &&
                user.github_handle !== null && (
                <p className="text-xs leading-none text-muted-foreground">
                  {user.github_handle}@github
                </p>
              )}
            </div>
          </DropdownMenuLabel>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={handleProfileClick}>
          {t("Auth.My Profile")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => navigateToMain(`/${locale}/mailbox`)}>
          {t("Layout.Mailbox")}
          {user.total_pending_envelopes !== undefined && (
            <span className="ml-auto text-xs text-muted-foreground">
              ({user.total_pending_envelopes})
            </span>
          )}
        </DropdownMenuItem>
        {accessibleProfiles.length > 0 && (
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>
              {t("Auth.Profiles")}
            </DropdownMenuSubTrigger>
            <DropdownMenuSubContent>
              {accessibleProfiles.map((profile) => (
                <DropdownMenuItem
                  key={profile.id}
                  onClick={() => navigateToMain(`/${locale}/${profile.slug}`)}
                >
                  {profile.title}
                  <span className="ml-auto text-xs text-muted-foreground">
                    {t(`Profile.MembershipKind.${profile.membership_kind}`)}
                  </span>
                </DropdownMenuItem>
              ))}
            </DropdownMenuSubContent>
          </DropdownMenuSub>
        )}
        {user.kind === "admin" && (
          <DropdownMenuItem onClick={() => navigateToMain(`/${locale}/admin`)}>
            {t("Admin.Admin Area")}
          </DropdownMenuItem>
        )}
        <DropdownMenuItem onClick={handleLogout}>
          {t("Auth.Logout")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
