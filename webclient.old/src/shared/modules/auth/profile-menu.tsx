"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/shared/components/ui/button.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/components/ui/dropdown-menu.tsx";
import { useAuth } from "@/shared/modules/auth/auth-context.tsx";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";

interface ProfileMenuProps {
  className?: string;
}

export function ProfileMenu({ className }: ProfileMenuProps) {
  const { session, logout, isAuthenticated } = useAuth();
  const router = useRouter();
  const { t } = useTranslations();

  if (!isAuthenticated || !session) {
    return null;
  }

  const handleProfileClick = () => {
    if (session.selected_profile) {
      router.push(`/${session.selected_profile.slug}`);
    } else {
      router.push("/elements/create-profile");
    }
  };

  const handleLogout = async () => {
    await logout();
  };

  // Determine avatar URL with priority: profile picture > GitHub avatar > fallback
  const fallbackAvatarUrl = `https://api.dicebear.com/7.x/initials/svg?seed=${encodeURIComponent(session.user?.name)}`;
  const avatarUrl = session.selected_profile?.profile_picture_uri ||
    (session.user?.github_handle ? `https://github.com/${session.user.github_handle}.png?size=32` : null) ||
    fallbackAvatarUrl;

  const [imageError, setImageError] = React.useState(false);

  const handleImageError = () => {
    setImageError(true);
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className={`relative h-8 w-8 rounded-full p-0 ${className}`}>
          <img
            className="h-8 w-8 rounded-full object-cover"
            src={imageError ? fallbackAvatarUrl : avatarUrl}
            alt={session.user?.name}
            onError={handleImageError}
          />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56" align="end" forceMount>
        <DropdownMenuLabel className="font-normal">
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-medium leading-none">
              {session.selected_profile?.title ?? session.user?.name}
            </p>
            <p className="text-xs leading-none text-muted-foreground">
              {session.selected_profile?.slug !== undefined ? `https://aya.is/${session.selected_profile.slug}` : ""}
            </p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        {session.selected_profile !== undefined && (
          <DropdownMenuItem onClick={handleProfileClick}>
            {t("Auth", "My Profile")}
          </DropdownMenuItem>
        )}
        <DropdownMenuItem onClick={handleLogout}>
          {t("Auth", "Logout")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
