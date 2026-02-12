const DICEBEAR_BASE_URL = "https://api.dicebear.com/7.x/initials/svg";

export function getDisplayName(name: string, fallbackName?: string): string {
  if (name !== "") {
    return name;
  }
  if (fallbackName !== undefined && fallbackName !== "") {
    return fallbackName;
  }
  return "?";
}

export function getInitials(displayName: string): string {
  return displayName.charAt(0).toUpperCase();
}

export function getProfilePictureUrl(
  src: string | null | undefined,
  name: string,
  fallbackName?: string,
): string {
  if (src !== null && src !== undefined && src !== "") {
    return src;
  }
  const seed = getDisplayName(name, fallbackName);
  return `${DICEBEAR_BASE_URL}?seed=${encodeURIComponent(seed)}`;
}
