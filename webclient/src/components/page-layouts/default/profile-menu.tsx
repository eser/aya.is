import * as React from "react";
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
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
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
  const [imageError, setImageError] = React.useState(false);

  if (!isAuthenticated || user === null) {
    return null;
  }

  // Determine avatar URL with priority: GitHub avatar > DiceBear fallback
  const fallbackAvatarUrl = `https://api.dicebear.com/7.x/initials/svg?seed=${encodeURIComponent(user.name || "User")}`;
  const githubAvatarUrl = user.github_handle !== undefined && user.github_handle !== null
    ? `https://github.com/${user.github_handle}.png?size=32`
    : null;
  const avatarUrl = githubAvatarUrl !== null ? githubAvatarUrl : fallbackAvatarUrl;

  const handleProfileClick = () => {
    // Navigate to user's profile if they have one, otherwise to create profile page
    if (user.individual_profile_id !== undefined) {
      // TODO: Get slug from profile data when available
      navigate({ to: `/${locale}/elements/create-profile` });
    } else {
      navigate({ to: `/${locale}/elements/create-profile` });
    }
  };

  const handleLogout = async () => {
    await logout();
  };

  const handleImageError = () => {
    setImageError(true);
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
            <Avatar size="default">
              <AvatarImage
                src={imageError ? fallbackAvatarUrl : avatarUrl}
                alt={user.name || "User avatar"}
                onError={handleImageError}
              />
              <AvatarFallback>
                {(user.name || "U").charAt(0).toUpperCase()}
              </AvatarFallback>
            </Avatar>
          </Button>
        )}
      />
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuGroup>
          <DropdownMenuLabel className="font-normal">
            <div className="flex flex-col space-y-1">
              <p className="text-sm font-medium leading-none">
                {user.name || user.email || t("Layout.Profile")}
              </p>
              {user.github_handle !== undefined && user.github_handle !== null && (
                <p className="text-xs leading-none text-muted-foreground">
                  @{user.github_handle}
                </p>
              )}
            </div>
          </DropdownMenuLabel>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={handleProfileClick}>
          {t("Auth.My Profile")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={handleLogout}>
          {t("Auth.Logout")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
