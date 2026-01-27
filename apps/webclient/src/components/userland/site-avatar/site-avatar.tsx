import * as React from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";

const DICEBEAR_BASE_URL = "https://api.dicebear.com/7.x/initials/svg";

export type SiteAvatarProps = {
  /** Profile picture URL */
  src?: string | null;
  /** Name or title to use for alt text and fallback initials */
  name: string;
  /** Optional fallback name (e.g., slug) if name is empty */
  fallbackName?: string;
  /** Size of the avatar - can use predefined sizes or custom class */
  size?: "sm" | "default" | "lg" | "xl" | "2xl";
  /** Additional class name for custom sizing or styling */
  className?: string;
  /** Handler for image load errors */
  onError?: () => void;
};

function getDisplayName(name: string, fallbackName?: string): string {
  if (name !== "") {
    return name;
  }
  if (fallbackName !== undefined && fallbackName !== "") {
    return fallbackName;
  }
  return "?";
}

function getInitials(displayName: string): string {
  return displayName.charAt(0).toUpperCase();
}

function getDicebearUrl(seed: string): string {
  return `${DICEBEAR_BASE_URL}?seed=${encodeURIComponent(seed)}`;
}

function getSizeClass(size: SiteAvatarProps["size"]): string {
  switch (size) {
    case "xl":
      return "size-16";
    case "2xl":
      return "size-20";
    default:
      return "";
  }
}

export function SiteAvatar(props: SiteAvatarProps) {
  const [imageError, setImageError] = React.useState(false);

  const displayName = getDisplayName(props.name, props.fallbackName);
  const initials = getInitials(displayName);
  const dicebearUrl = getDicebearUrl(displayName);

  // Determine the image source
  const hasCustomImage = props.src !== null && props.src !== undefined && props.src !== "";
  const imageSrc = imageError || !hasCustomImage ? dicebearUrl : props.src;

  const handleImageError = () => {
    setImageError(true);
    if (props.onError !== undefined) {
      props.onError();
    }
  };

  // For predefined sizes (sm, default, lg), use the Avatar component's size prop
  const isPredefinedSize = props.size === "sm" || props.size === "default" || props.size === "lg";
  const avatarSize = isPredefinedSize ? props.size : undefined;
  const customSizeClass = !isPredefinedSize ? getSizeClass(props.size) : "";

  return (
    <Avatar
      size={avatarSize}
      className={cn(customSizeClass, props.className)}
    >
      <AvatarImage
        src={imageSrc}
        alt={displayName}
        onError={handleImageError}
      />
      <AvatarFallback>{initials}</AvatarFallback>
    </Avatar>
  );
}
