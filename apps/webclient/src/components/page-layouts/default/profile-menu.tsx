import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
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

type ProfileMenuProps = {
  className?: string;
};

export function ProfileMenu(props: ProfileMenuProps) {
  const { isAuthenticated, user, logout } = useAuth();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const locale = getCurrentLanguage();

  if (!isAuthenticated || user === null) {
    return null;
  }

  // Determine avatar URL with priority: GitHub avatar > null (SiteAvatar will use dicebear fallback)
  const githubAvatarUrl = user.github_handle !== undefined && user.github_handle !== null
    ? `https://github.com/${user.github_handle}.png?size=32`
    : null;

  const handleProfileClick = () => {
    // Navigate to user's profile if they have one, otherwise to create profile page
    if (user.individual_profile_slug !== undefined) {
      navigate({ to: `/${locale}/${user.individual_profile_slug}` });
    } else {
      navigate({ to: `/${locale}/elements/new` });
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
        {user.accessible_profiles !== undefined && user.accessible_profiles.length > 0 ? (
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>
              {t("Auth.My Profile")}
            </DropdownMenuSubTrigger>
            <DropdownMenuSubContent>
              <DropdownMenuItem onClick={handleProfileClick}>
                {user.name}
                <span className="ml-auto text-xs text-muted-foreground">
                  {t("Profile.Individual")}
                </span>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              {user.accessible_profiles.map((profile) => (
                <DropdownMenuItem
                  key={profile.id}
                  onClick={() => navigate({ to: `/${locale}/${profile.slug}` })}
                >
                  {profile.title}
                  <span className="ml-auto text-xs text-muted-foreground">
                    {profile.kind}
                  </span>
                </DropdownMenuItem>
              ))}
            </DropdownMenuSubContent>
          </DropdownMenuSub>
        ) : (
          <DropdownMenuItem onClick={handleProfileClick}>
            {t("Auth.My Profile")}
          </DropdownMenuItem>
        )}
        {user.kind === "admin" && (
          <DropdownMenuItem onClick={() => navigate({ to: `/${locale}/admin` })}>
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
